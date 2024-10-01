$platform = Read-Host "Build for Windows or Unix? (Enter 'w' or 'u')" # Ask user for the target platform
$appName = "pi-server"

function Build-Windows {
  Write-Host "Building for Windows..." -ForegroundColor Cyan

  $env:GOOS = "windows"
  $env:GOARCH = "amd64"
  go build -o "../builds/$appName.exe" ./app/cmd
  
  if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $appName.exe" -ForegroundColor Green
  } else {
    Write-Host "Build failed" -ForegroundColor Red
  }
}

function Build-Unix {
  Write-Host "Building for Unix..." -ForegroundColor Cyan
  
  $env:GOOS = "linux"
  $env:GOARCH = "arm64"
  go build -o "../builds/$appName" ./app/cmd

  if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $appName" -ForegroundColor Green
  } else {
    Write-Host "Build failed" -ForegroundColor Red
  }
}

switch ($platform.ToLower()) {
  "u" { Build-Unix }
  "w" { Build-Windows }
  default { Write-Host "Invalid input. Please enter 'w' or 'u'" -ForegroundColor Red }
}
