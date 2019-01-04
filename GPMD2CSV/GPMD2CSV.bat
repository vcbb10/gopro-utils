@ECHO OFF

:: Name of the folder that is created to output to (Optional)
Set BatchOutputFolder=GoPro Metadata Extract

:: Make Individual Subdirectories for each file? Type in Yes or No
Set IndividualSubDirs=Yes

:: Choose GPS accuracy filter. For example 500 (high accuracy) or 10000 (very low accuracy)
Set AccuracyFilter=1000

:: Choose GPS fix filter. 3 (3D Fix, best case scenario), 2 (2D fix) or 0 (no fix, but there might still be some useful data)
Set FixFilter=3

:: ==========================================
:: You shouldn't need to edit below this line
::===========================================

:IndividualSubDirCheck
if '%IndividualSubDirs%'=='Yes' goto IndSubDirY
if '%IndividualSubDirs%'=='yes' goto IndSubDirY
if '%IndividualSubDirs%'=='Y' goto IndSubDirY
if '%IndividualSubDirs%'=='y' goto IndSubDirY
if '%IndividualSubDirs%'=='No' goto IndSubDirN
if '%IndividualSubDirs%'=='no' goto IndSubDirN
if '%IndividualSubDirs%'=='N' goto IndSubDirN
if '%IndividualSubDirs%'=='n' goto IndSubDirN
Goto IndividualSubDirCheckError

:IndSubDirY
Set IndividualSubDirs=1
GoTo RunIt

:IndSubDirN
Set IndividualSubDirs=0
GoTo RunIt

:RunIt

Set SourceScriptDirectory=%~dp0
Set SourceFile=%1

cd "%~dp1"

If '%IndividualSubDirs%'=='1' (
	If exist ".\%BatchOutputFolder%\%~n1\%~n1".gpx goto :eof
) Else (
	If exist ".\%BatchOutputFolder%\%~n1".gpx goto :eof )

mkdir "%~n1"
if [%1]==[] goto :eof
:loop
"%SourceScriptDirectory%\bin\ffmpeg" -i "%~1" > output.txt 2>&1
for /F "delims=" %%a in ('FINDSTR "gpmd" output.txt') do set line=%%a
echo "%line%"
set stream= %line:~12,3%
echo "%stream%"
CLS
ECHO.
ECHO **************************************************
ECHO **************** Processing file: ****************
ECHO "%~nx1"
ECHO ***************** In Directory: ******************
ECHO "%~dp1"
ECHO **************************************************
@ECHO ON
START "" /WAIT /MIN "%SourceScriptDirectory%bin\ffmpeg" -y -i "%~1" -codec copy -map "%stream%" -f rawvideo "%~n1".bin
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gpmd2csv" -i "%~n1".bin -o "%~n1"/"%~n1".csv
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gopro2json" -i "%~n1".bin -o "%~n1"/"%~n1".json
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gps2kml" -i "%~n1".bin -a %AccuracyFilter% -f %FixFilter% -o "%~n1"/"%~n1".kml
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gopro2gpx" -i "%~n1".bin -a %AccuracyFilter% -f %FixFilter% -o "%~n1"/"%~n1".gpx
DEL "%~n1".bin
DEL output.txt
Mkdir "%BatchOutputFolder%"

If '%IndividualSubDirs%'=='1' (
	Move /Y "%~n1" ".\%BatchOutputFolder%\"
) Else (
	Move /Y "%~n1\*" ".\%BatchOutputFolder%\"
	RMDIR /S /Q "%~n1\" )

shift
if not [%1]==[] goto loop
GoTo eof

:IndividualSubDirCheckError
cls
ECHO **********************************************************
ECHO ************************* ERROR **************************
ECHO ***** "IndividualSubDirs" is not formatted correctly *****
ECHO *** Check "GPMD2CSV.bat" and verify it is = Yes or No ****
ECHO **********************************************************
pause
exit
