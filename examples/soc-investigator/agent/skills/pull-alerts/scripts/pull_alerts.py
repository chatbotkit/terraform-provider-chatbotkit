#!/usr/bin/env python3
"""Deterministic alert puller for the SOC investigator example.

This is the *deterministic* half of the SOC loop — no LLM involved. It polls the
(mock) SIEM, normalizes raw alerts into a common shape, correlates them into
cases by a stable correlation UID (rule + asset + hourly time bucket), and writes
open cases as JSON files under cases/.

Re-running is idempotent: an alert that maps to an existing case is appended to
it, not turned into a duplicate case. That is what makes a *polling* pull safe to
run on a cycle.

Stdlib only, so it runs in any sandbox with no pip install. To use a real SIEM,
replace read_alerts() with an ELK/Splunk query (see the SIEM seam below).

Usage (from the workspace root):
    python .skills/pull-alerts/scripts/pull_alerts.py
"""
import glob
import hashlib
import json
import os
from datetime import datetime

ROOT = os.getcwd()
ALERTS_FILE = os.path.join(ROOT, "data", "sample-alerts.json")
CASES_DIR = os.path.join(ROOT, "cases")

SEVERITY_ORDER = ["low", "medium", "high", "critical"]
SEVERITY_BY_SCORE = [(90, "critical"), (70, "high"), (40, "medium"), (0, "low")]


def severity_for(score):
    for threshold, label in SEVERITY_BY_SCORE:
        if score >= threshold:
            return label
    return "low"


def max_severity(a, b):
    return a if SEVERITY_ORDER.index(a) >= SEVERITY_ORDER.index(b) else b


def time_bucket(ts):
    # @note round down to the hour: alerts of the same rule on the same asset
    # within an hour are treated as one case.
    try:
        dt = datetime.fromisoformat(ts.replace("Z", "+00:00"))
    except (ValueError, AttributeError):
        dt = datetime.utcnow()
    return dt.strftime("%Y-%m-%dT%H")


def correlation_uid(alert):
    key = "|".join(
        [
            alert["rule"],
            alert["host"],
            alert["user"],
            time_bucket(alert["timestamp"]),
        ]
    )
    return hashlib.sha1(key.encode()).hexdigest()[:12]


def alert_fingerprint(alert):
    # @note a stable per-alert id so re-polling the same window does not append
    # the same alert twice (true idempotency, not just case-level dedup).
    key = "|".join(
        [
            alert["timestamp"],
            alert["rule"],
            alert["host"],
            alert["user"],
            alert["process"],
        ]
    )
    return hashlib.sha1(key.encode()).hexdigest()[:16]


def read_alerts():
    # @note SIEM seam: swap this for a real query against your SIEM
    # (e.g. an adaptive query / SPL / ES|QL against ELK or Splunk). Everything
    # downstream only depends on the normalized shape from normalize().
    with open(ALERTS_FILE) as handle:
        return json.load(handle)


def normalize(raw):
    score = float(raw.get("risk_score", 50))
    alert = {
        "timestamp": raw.get("@timestamp") or raw.get("timestamp") or "",
        "rule": raw.get("rule") or "unknown",
        "host": raw.get("host") or "",
        "user": raw.get("user") or "",
        "process": raw.get("process") or "",
        "risk_score": score,
        "severity": severity_for(score),
        "iocs": raw.get("iocs", []),
        "raw": raw,
    }
    alert["alert_id"] = alert_fingerprint(alert)
    return alert


def main():
    os.makedirs(CASES_DIR, exist_ok=True)
    alerts = read_alerts()

    new_cases = 0
    correlated = 0
    duplicates = 0

    for raw in alerts:
        alert = normalize(raw)
        uid = correlation_uid(alert)
        path = os.path.join(CASES_DIR, f"{uid}.json")

        if os.path.exists(path):
            with open(path) as handle:
                case = json.load(handle)
            seen = {a.get("alert_id") for a in case["alerts"]}
            if alert["alert_id"] in seen:
                duplicates += 1
                continue
            case["alerts"].append(alert)
            case["alert_count"] = len(case["alerts"])
            case["severity"] = max_severity(case["severity"], alert["severity"])
            with open(path, "w") as handle:
                json.dump(case, handle, indent=2)
            correlated += 1
        else:
            asset = alert["host"] or "unknown asset"
            case = {
                "id": uid,
                "status": "open",
                "severity": alert["severity"],
                "title": f'{alert["rule"]} on {asset}',
                "created_at": alert["timestamp"],
                "alert_count": 1,
                "alerts": [alert],
                "report": None,
                "enrichment": {},
            }
            with open(path, "w") as handle:
                json.dump(case, handle, indent=2)
            new_cases += 1

    open_total = len(glob.glob(os.path.join(CASES_DIR, "*.json")))

    print(
        json.dumps(
            {
                "pulled": len(alerts),
                "new_cases": new_cases,
                "correlated_into_existing": correlated,
                "duplicates_skipped": duplicates,
                "open_cases_total": open_total,
            },
            indent=2,
        )
    )


if __name__ == "__main__":
    main()
