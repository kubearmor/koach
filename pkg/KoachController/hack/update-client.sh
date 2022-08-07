#!/bin/bash

dir="$(cd "$(dirname "$0")" || exit; pwd)"
topdir="$(cd "$(dirname "$0")/.." || exit; pwd)"
project_module="$(awk '/^module / { print $2 }' "$topdir/go.mod")"
evtax_module="$(echo "$project_mode" | sed -e 's,/taxman/,/event-taxonomy/,')"
apis_dir=pkg/apis
output_dir=pkg/client
clientset_name="clientset"
go run k8s.io/code-generator/cmd/client-gen \
  --clientset-name "$clientset_name" \
  --input-base "gitlab.com/zedge-oss/zeppa/event-taxonomy" \
  --input "$evtax_module/$apis_dir/dwh/v1beta1" \
  --output-package "$project_module/$output_dir" \
  --output-base "" \
  --go-header-file "$dir/boilerplate.go.txt" \
  "$@"
# shellcheck disable=SC2115
rm -rf "$output_dir/$clientset_name"
# Since client-gen is still incapable of generating working code outside of GOPATH,
# we need to do this hack:
mv "$topdir/$project_module/$output_dir/$clientset_name" "$topdir/$output_dir/$clientset_name"
rm -rf "$topdir/gitlab.com"
