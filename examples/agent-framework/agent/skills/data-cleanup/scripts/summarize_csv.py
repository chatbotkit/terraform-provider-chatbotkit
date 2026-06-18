#!/usr/bin/env python3
"""Profile a CSV file: row/column counts, fill rate, inferred type, samples.

Usage: python summarize_csv.py <path-to-csv>
Standard library only, so it runs anywhere in the workspace.
"""

import csv
import sys
from collections import Counter


def infer_type(value):
    v = value.strip()
    if v == "":
        return "empty"
    try:
        int(v)
        return "int"
    except ValueError:
        pass
    try:
        float(v)
        return "float"
    except ValueError:
        pass
    return "str"


def main(path):
    with open(path, newline="") as f:
        reader = csv.reader(f)
        try:
            header = next(reader)
        except StopIteration:
            print("empty file")
            return
        rows = list(reader)

    total = len(rows)
    print(f"file: {path}")
    print(f"rows: {total}")
    print(f"columns: {len(header)}")
    print()

    for i, name in enumerate(header):
        values = [row[i] for row in rows if i < len(row)]
        filled = [v for v in values if v.strip() != ""]
        types = Counter(infer_type(v) for v in filled)
        dominant = types.most_common(1)[0][0] if types else "empty"
        fill = (len(filled) / total * 100) if total else 0
        samples = ", ".join(filled[:3]) if filled else "-"
        print(f"[{i}] {name}")
        print(f"    fill: {fill:.0f}%  type: {dominant}  samples: {samples}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("usage: python summarize_csv.py <path-to-csv>", file=sys.stderr)
        sys.exit(1)
    main(sys.argv[1])
