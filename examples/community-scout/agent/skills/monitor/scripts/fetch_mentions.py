#!/usr/bin/env python3
"""Fetch and persist monitored Reddit mentions (deterministic — no model in the loop).

Monitoring is done as a script so the bulk of Reddit data never round-trips
through the agent's context, which would be slow and expensive. The script reads
the watchlist, searches Reddit directly, normalizes the hits, and dedups them into
the mention store. The agent only ever sees the small summary this prints — plus
the specific threads it later chooses to judge.

It queries Reddit's public Atom (RSS) search feed with a realistic, rotating
User-Agent — the same approach the platform's Reddit ability uses, because the
plain JSON endpoint and bot-looking User-Agents are blocked.

Stdlib only. Usage (from the workspace root):
    python .skills/monitor/scripts/fetch_mentions.py

Requires outbound HTTPS to reddit.com from the shell sandbox. If the sandbox has
no network egress, fall back to the read-only reddit tools for fetching.
"""
import glob
import html
import json
import os
import random
import re
import time
import urllib.parse
import urllib.request
import xml.etree.ElementTree as ET
from datetime import datetime, timezone

ROOT = os.getcwd()
WATCHLIST = os.path.join(ROOT, "watchlist.md")
MENTIONS_DIR = os.path.join(ROOT, "mentions")

ATOM = "http://www.w3.org/2005/Atom"
LIMIT = 25
SLEEP_S = 2.0  # be polite to the unauthenticated endpoint

USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
]

_TAG_RE = re.compile(r"<[^>]+>")


def parse_section(text, header):
    """Collect '- item' bullets under a '## header' until the next '## '."""
    items = []
    in_section = False
    for line in text.splitlines():
        if line.strip().lower().startswith("## "):
            in_section = line.strip().lower() == f"## {header.lower()}"
            continue
        if in_section:
            match = re.match(r"\s*-\s+(.+)", line)
            if match:
                items.append(match.group(1).strip())
    return items


def load_watchlist():
    if not os.path.exists(WATCHLIST):
        return [], []
    text = open(WATCHLIST, encoding="utf-8").read()
    keywords = parse_section(text, "Keywords")
    subs = []
    for raw in parse_section(text, "Subreddits to watch"):
        name = raw.split()[0] if raw.split() else ""
        if name.startswith("r/"):
            name = name[2:]
        if name:
            subs.append(name)
    return keywords, subs


def build_query(keywords):
    parts = []
    for kw in keywords:
        kw = kw.strip()
        if kw:
            parts.append(f'"{kw}"' if " " in kw else kw)
    return " OR ".join(parts)


def search(subreddit, query):
    qs = urllib.parse.urlencode(
        {"q": query, "restrict_sr": "on", "sort": "new", "limit": LIMIT}
    )
    url = f"https://www.reddit.com/r/{urllib.parse.quote(subreddit)}/search.rss?{qs}"
    req = urllib.request.Request(
        url, headers={"User-Agent": random.choice(USER_AGENTS)}
    )
    with urllib.request.urlopen(req, timeout=20) as resp:
        root = ET.fromstring(resp.read())
    return root.findall(f"{{{ATOM}}}entry")


def _text(entry, tag):
    el = entry.find(f"{{{ATOM}}}{tag}")
    return el.text or "" if el is not None else ""


def normalize(entry, matched):
    link = entry.find(f"{{{ATOM}}}link")
    author = entry.find(f"{{{ATOM}}}author/{{{ATOM}}}name")
    category = entry.find(f"{{{ATOM}}}category")
    snippet = html.unescape(_TAG_RE.sub(" ", _text(entry, "content")))
    snippet = re.sub(r"\s+", " ", snippet).strip()[:280]
    return {
        "id": _text(entry, "id").strip(),  # Reddit fullname, e.g. t3_abc123
        "source": "reddit",
        "subreddit": category.get("term") if category is not None else "",
        "title": html.unescape(_text(entry, "title").strip()),
        "url": link.get("href") if link is not None else "",
        "author": author.text if author is not None else "",
        "snippet": snippet,
        "matched": matched,
    }


def main():
    os.makedirs(MENTIONS_DIR, exist_ok=True)
    keywords, subs = load_watchlist()

    if not keywords or not subs:
        print(json.dumps({"error": "no keywords or subreddits found in watchlist.md"}))
        return

    query = build_query(keywords)

    pulled = 0
    new = 0
    duplicates = 0
    errors = 0

    for sub in subs:
        try:
            entries = search(sub, query)
        except Exception as exc:  # network/HTTP error (403/429/...) — skip this sub
            errors += 1
            print(f"warning: search r/{sub} failed: {exc}", flush=True)
            time.sleep(SLEEP_S)
            continue

        for entry in entries:
            mention = normalize(entry, query)
            mid = mention["id"]
            if not mid:
                continue

            pulled += 1
            path = os.path.join(MENTIONS_DIR, f"{mid}.json")
            if os.path.exists(path):
                duplicates += 1
                continue

            mention.update(
                {
                    "discovered_at": datetime.now(timezone.utc).isoformat(),
                    "status": "new",  # new -> scored -> drafted -> suggested | skipped
                    "relevance": None,
                    "draft": None,
                }
            )
            with open(path, "w") as handle:
                json.dump(mention, handle, indent=2)
            new += 1

        time.sleep(SLEEP_S)

    total = len(glob.glob(os.path.join(MENTIONS_DIR, "*.json")))

    print(
        json.dumps(
            {
                "subreddits_searched": len(subs),
                "pulled": pulled,
                "new": new,
                "duplicates": duplicates,
                "errors": errors,
                "mentions_total": total,
            },
            indent=2,
        )
    )


if __name__ == "__main__":
    main()
