mkdir %~dp0exports
if [%1]==[] goto :eof
:loop
%~dp0bin\ffmpeg -i "%~1" > output.txt 2>&1
for /F "delims=" %%a in ('FINDSTR "gpmd" output.txt') do set line=%%a
echo "%line%"
set stream= %line:~12,3%
echo "%stream%"
START /WAIT %~dp0bin\ffmpeg -y -i "%~1" -codec copy -map "%stream%" -f rawvideo "%~n1".bin
START /WAIT %~dp0bin\gpmd2csv -i "%~n1".bin -o exports/"%~n1".csv
START /WAIT %~dp0bin\gopro2gpx -i "%~n1".bin -o exports/"%~n1".gpx
START /WAIT %~dp0bin\gopro2json -i "%~n1".bin -o exports/"%~n1".json
START /WAIT %~dp0bin\gps2kml -i "%~n1".bin -o exports/"%~n1".kml
DEL "%~n1".bin
DEL output.txt
shift
if not [%1]==[] goto loop