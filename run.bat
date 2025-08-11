@echo off
echo Chrome-Stable uTLS Template Generator
echo =====================================
echo.
echo Available commands:
echo   test --template templates\chrome-120.0.6099.109.json
echo   --help
echo   generate --help
echo   monitor --help
echo.
echo Type your command (or just press Enter to exit):
set /p command="chrome-utls-gen.exe "

if not "%command%"=="" (
    chrome-utls-gen.exe %command%
    echo.
    echo Command completed. Press any key to exit...
    pause >nul
)
