pkgs=$(go list -f '{{ if .TestGoFiles }}{{.ImportPath}}{{end}}' ./pkg/... | xargs)

cov_tmp="$(mktemp)"
echo "mode: count" >"${cov_tmp}"

go test \
    -covermode=count \
    -coverprofile="${cov_tmp}" \
    ${pkgs}

cat $cov_tmp