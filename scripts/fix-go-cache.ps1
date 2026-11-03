# Fix "corrupt index" / bad module zip on Windows.
# go.mod must stay minimal (module + go 1.27.3 only); run go mod tidy after this.

Write-Host "GOMODCACHE: $(go env GOMODCACHE)"

go clean -modcache
go clean -cache

Set-Location $PSScriptRoot\..

go mod tidy
go test ./...

Write-Host "Done."
