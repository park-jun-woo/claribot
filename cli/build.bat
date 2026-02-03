@echo off
REM Claritask CLI Build Script for Windows (Batch)
REM Usage: build.bat [command]
REM Commands: build, install, clean, all (default: build)

setlocal enabledelayedexpansion

set BIN_NAME=clari.exe
set INSTALL_DIR=%USERPROFILE%\bin
set COMMAND=%1

if "%COMMAND%"=="" set COMMAND=build

cd /d "%~dp0"

if "%COMMAND%"=="build" goto :build
if "%COMMAND%"=="install" goto :install
if "%COMMAND%"=="clean" goto :clean
if "%COMMAND%"=="all" goto :all

echo Unknown command: %COMMAND%
echo Usage: build.bat [build^|install^|clean^|all]
exit /b 1

:check_prereq
echo Checking prerequisites...

where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ERROR: Go is not installed. Install from https://go.dev/dl/
    exit /b 1
)

echo Prerequisites OK
exit /b 0

:build
call :check_prereq
if %ERRORLEVEL% neq 0 exit /b 1

echo Building %BIN_NAME%...
go build -o %BIN_NAME% ./cmd/claritask/

if %ERRORLEVEL% equ 0 (
    echo Build successful: %BIN_NAME%
) else (
    echo Build failed
    exit /b 1
)
goto :eof

:install
if not exist %BIN_NAME% (
    echo Binary not found. Building first...
    call :build
    if %ERRORLEVEL% neq 0 exit /b 1
)

echo Installing to %INSTALL_DIR%...

if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

copy /Y %BIN_NAME% "%INSTALL_DIR%\%BIN_NAME%" >nul

echo Installed: %INSTALL_DIR%\%BIN_NAME%
echo.
echo NOTE: Make sure %INSTALL_DIR% is in your PATH.
echo Run this in PowerShell to add it:
echo   [Environment]::SetEnvironmentVariable("Path", $env:PATH + ";%INSTALL_DIR%", "User")
goto :eof

:clean
echo Cleaning...
if exist %BIN_NAME% (
    del /f %BIN_NAME%
    echo Removed %BIN_NAME%
) else (
    echo Nothing to clean
)
goto :eof

:all
call :check_prereq
if %ERRORLEVEL% neq 0 exit /b 1
call :build
if %ERRORLEVEL% neq 0 exit /b 1
call :install
goto :eof
