$appName = "sec-surv-server"

function Build-Unix {
  Write-Host "Building for Unix..." -ForegroundColor Cyan
  
  $env:GOOS = "linux"
  $env:GOARCH = "arm64"
  go build -o "$appName" ./app/cmd

  if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $appName" -ForegroundColor Green
  } else {
    Write-Host "Build failed" -ForegroundColor Red
  }
}

Build-Unix
