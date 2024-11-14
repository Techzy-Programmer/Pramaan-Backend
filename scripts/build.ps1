Write-Host "[>>]: Building for production..." -ForegroundColor Cyan

$appName = "pch-backend"
Write-Host "Building for Unix..." -ForegroundColor Cyan

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$appName" ./app/

if ($LASTEXITCODE -eq 0) {
  Write-Host "Build successful: $appName" -ForegroundColor Green
} else {
  Write-Host "Build failed" -ForegroundColor Red
}
