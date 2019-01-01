@ECHO OFF

:: Name of the folder that is created to output to (Optional)
Set BatchOutputFolder=GoPro Metadata Extract

:: You shouldn't need to edit below this line

Set SourceScriptDirectory=%~dp0
Set SourceFile=%1

cd "%~dp1"

if exist "%BatchOutputFolder%\%~n1" goto :eof

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
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gopro2gpx" -i "%~n1".bin -o "%~n1"/"%~n1".gpx
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gopro2json" -i "%~n1".bin -o "%~n1"/"%~n1".json
START "" /WAIT /MIN "%SourceScriptDirectory%bin\gps2kml" -i "%~n1".bin -o "%~n1"/"%~n1".kml
DEL "%~n1".bin
DEL output.txt
Mkdir "%BatchOutputFolder%"
Move /Y "%~n1" "%BatchOutputFolder%"
shift
if not [%1]==[] goto loop