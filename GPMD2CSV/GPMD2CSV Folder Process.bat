@echo off

Set SourceScriptDirectory=%~dp0
Set SourceFile=%1

For %%f in (*.MP4) do (
	CLS
	ECHO.
	ECHO **************************************************
	ECHO **************** Processing file: ****************
	ECHO "%~dp1%%~f"
	ECHO **************************************************
	Call "%~dp0\GPMD2CSV.bat" "%~dp1%%~f"
)