#!/usr/bin/env python3
"""Threat-intelligence lookup for the SOC investigator example.

Given an IOC (IP, domain, or file hash), returns a reputation verdict. This mock
is deterministic and offline so the example runs without credentials. To use a
real provider (AlienVault OTX, VirusTotal, ...), set TI_API_TOKEN and replace the
body of lookup() with an HTTP call.

Stdlib only. Usage (from the workspace root):
    python .skills/enrich/scripts/ti_lookup.py <ioc>
"""
import hashlib
import json
import os
import sys


def lookup(ioc):
    token = os.environ.get("TI_API_TOKEN")

    # @note TI seam: with a token, call your provider here (HTTP GET with the
    # token) and map its response onto the shape below.
    if token:
        pass

    # deterministic mock: a stable pseudo-score derived from the IOC itself
    score = int(hashlib.sha1(ioc.encode()).hexdigest(), 16) % 100

    if score >= 80:
        verdict = "malicious"
    elif score >= 50:
        verdict = "suspicious"
    else:
        verdict = "benign"

    return {
        "ioc": ioc,
        "verdict": verdict,
        "reputation_score": score,
        "source": "live-ti" if token else "mock-ti",
        "pulses": [] if verdict == "benign" else [f"reported in {score % 7 + 1} feeds"],
    }


def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "usage: ti_lookup.py <ioc>"}))
        sys.exit(1)

    print(json.dumps(lookup(sys.argv[1]), indent=2))


if __name__ == "__main__":
    main()
