#!/usr/bin/env bash
#
# Validates every Terraform example against the locally-built provider using a
# dev override. This needs no registry download and no API credentials —
# `terraform validate` checks each example's resources and attributes against the
# provider schema, catching regressions such as renamed resources or removed
# fields (e.g. a stale `chatbotkit_mcp_server_integration`).
#
# Usage: ./scripts/validate-examples.sh   (run from the provider root, or anywhere)

set -euo pipefail

cd "$(dirname "$0")/.."

bindir="$(mktemp -d)"
tfrc="$(mktemp)"
trap 'rm -rf "$bindir" "$tfrc"' EXIT

echo "==> building provider"
go build -o "$bindir/terraform-provider-chatbotkit" .

cat > "$tfrc" <<EOF
provider_installation {
  dev_overrides {
    "chatbotkit/chatbotkit" = "$bindir"
  }
  direct {}
}
EOF
export TF_CLI_CONFIG_FILE="$tfrc"

echo "==> terraform fmt -check"
terraform fmt -check -recursive examples/

failed=0
for dir in examples/*/; do
  ls "$dir"*.tf >/dev/null 2>&1 || continue
  echo "==> validating ${dir}"
  if ! terraform -chdir="$dir" validate; then
    failed=1
  fi
done

if [ "$failed" -ne 0 ]; then
  echo "FAILED: one or more examples did not validate"
  exit 1
fi

echo "OK: all examples validated"
