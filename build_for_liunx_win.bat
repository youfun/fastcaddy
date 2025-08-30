@echo off


powershell -Command "Set-Item -Path Env:GOOS -Value 'linux'; Set-Item -Path Env:GOARCH -Value 'amd64'; go build -o fastcaddy ./cmd/fastcaddy"
REM 可选：使用 upx 进一步压缩（需安装 upx）
REM upx fastcaddy
echo Linux build completed: fastcaddy
