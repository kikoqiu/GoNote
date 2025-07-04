@echo off
setlocal enableDelayedExpansion

REM ============================================================================
REM == Multi-Platform Build Script for Go Application (with Frontend Assets)  ==
REM ============================================================================
REM == This script:                                                           ==
REM == 1. Deletes and recreates the 'embed' directory for fresh assets.       ==
REM == 2. Deletes only previously generated binaries from the 'build' directory,
REM ==    preserving any other files (like configs).                          ==
REM ============================================================================

REM --- Determine absolute paths ---
SET "BACKEND_DIR=%~dp0"
for %%I in ("%BACKEND_DIR%..") do set "PROJECT_ROOT=%%~fI"

REM --- Configuration (using absolute paths for clarity) ---
SET "APP_NAME=gonote"
SET "FRONTEND_SRC_DIR=%PROJECT_ROOT%\frontend\dist"
SET "EMBED_DIR=%BACKEND_DIR%embed"
SET "OUTPUT_DIR=%BACKEND_DIR%build"

REM --- Define Target Platforms (Format: GOOS/GOARCH) ---
SET "PLATFORMS=windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

REM ============================================================================

TITLE Building %APP_NAME% with Frontend Assets

echo.
echo [INFO] Starting build for %APP_NAME%...
echo [INFO] Detected Project Root: %PROJECT_ROOT%
echo.

REM --- Step 1: Recreate 'embed' directory with frontend assets ---
echo [STEP 1] Preparing frontend assets for embedding...

REM Check if frontend build output exists
IF NOT EXIST "%FRONTEND_SRC_DIR%" (
    echo [ERROR] Frontend build output not found at:
    echo [ERROR] %FRONTEND_SRC_DIR%
    echo [ERROR] Please run 'npm run build' in the 'frontend' directory first.
    echo.
    pause
    exit /b 1
)

REM --- MODIFIED: Delete and recreate the entire embed directory ---
IF EXIST "%EMBED_DIR%" (
    echo [INFO] Deleting old assets directory: %EMBED_DIR%
    rmdir /s /q "%EMBED_DIR%"
)
echo [INFO] Creating new assets directory: %EMBED_DIR%
mkdir "%EMBED_DIR%"

REM Copy files using absolute paths
echo [INFO] Copying assets from 'frontend\dist' to 'backend\embed'...
xcopy "%FRONTEND_SRC_DIR%\*" "%EMBED_DIR%\" /E /I /Q /Y > nul

IF ERRORLEVEL 1 (
    echo [ERROR] Failed to copy frontend assets. Aborting.
    echo.
    pause
    exit /b 1
)

echo [OK] Frontend assets prepared successfully.
echo.


REM --- Step 2: Clean up old binaries and run Go build ---
echo [STEP 2] Building Go application...

REM --- Ensure the build directory exists, create if not ---
IF NOT EXIST "%OUTPUT_DIR%" (
    echo [INFO] Creating build output directory: %OUTPUT_DIR%
    mkdir "%OUTPUT_DIR%"
)

REM --- MODIFIED: Precisely delete only the target binaries from the build directory ---
echo [INFO] Cleaning previously generated binaries from: %OUTPUT_DIR%
FOR %%P IN (%PLATFORMS%) DO (
    call :delete_binary %%P
)
echo.

REM --- Main Build Loop ---
FOR %%P IN (%PLATFORMS%) DO (
    call :build_platform %%P
)

echo.
echo [SUCCESS] All builds completed!
echo Binaries are located in the 'build' directory.
echo.
pause
goto :eof


REM --- Subroutine for deleting a single binary ---
:delete_binary
    setlocal
    FOR /F "tokens=1,2 delims=/" %%G IN ("%1") DO (
        SET "GOOS=%%G"
        SET "GOARCH=%%H"
    )
    SET "FILENAME_NO_EXT=%APP_NAME%-!GOOS!-!GOARCH!"
    
    REM Delete both with and without .exe to cover all platforms
    IF EXIST "%OUTPUT_DIR%\!FILENAME_NO_EXT!" del "%OUTPUT_DIR%\!FILENAME_NO_EXT!"
    IF EXIST "%OUTPUT_DIR%\!FILENAME_NO_EXT!.exe" del "%OUTPUT_DIR%\!FILENAME_NO_EXT!.exe"
    endlocal
    goto :eof


REM --- Subroutine for building a single platform ---
:build_platform
    setlocal
    FOR /F "tokens=1,2 delims=/" %%G IN ("%1") DO (
        SET "GOOS=%%G"
        SET "GOARCH=%%H"
    )

    SET "FILENAME=%APP_NAME%-!GOOS!-!GOARCH!"
    IF "!GOOS!"=="windows" SET "FILENAME=!FILENAME!.exe"

    SET "FINAL_OUTPUT_PATH=build\%FILENAME%"

    echo [BUILD] Building for !GOOS!/!GOARCH!...
    
    set "CGO_ENABLED=0"
    set "GOOS=!GOOS!"
    set "GOARCH=!GOARCH!"

    go build -v -o "%FINAL_OUTPUT_PATH%" -ldflags="-s -w" .

    IF ERRORLEVEL 1 (
        echo.
        echo [ERROR] Build failed for !GOOS!/!GOARCH%!
        echo.
    ) else (
        echo [OK] Successfully created: %FINAL_OUTPUT_PATH%
        echo.
    )
    
    endlocal
    goto :eof