---
description: Profile and sanity-check a tabular CSV file.
---

# Data Cleanup

Use this when you're handed a CSV and need a quick, reliable profile before
working with it.

## Steps

1. Put the CSV somewhere in the workspace (e.g. `data/input.csv`).
2. Run the bundled profiler with the shell tools:

   ```bash
   python .skills/data-cleanup/scripts/summarize_csv.py data/input.csv
   ```

3. Read the report: row/column counts, per-column fill rate, inferred type, and
   a few sample values.
4. Flag columns that are mostly empty or have inconsistent types, and propose
   concrete cleanup steps before transforming anything.

## Notes

- The profiler uses only the Python standard library, so it runs anywhere in the
  workspace.
