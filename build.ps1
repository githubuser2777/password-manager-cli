Write-Host "Building Password Manager CLI..."
go build -o passmgr.exe
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: passmgr.exe" -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
}
