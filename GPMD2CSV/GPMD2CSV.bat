if [%1]==[] goto :eof
:loop
.\ffmpeg -i "%~1" > output.txt 2>&1
for /F "delims=" %%a in ('FINDSTR "gpmd" output.txt') do set line=%%a
echo "%line%"
set stream= %line:~12,3%
echo "%stream%"
START /WAIT .\ffmpeg -y -i "%~1" -codec copy -map "%stream%" -f rawvideo "%~n1".bin
START /WAIT .\gpmd2csv -i "%~n1".bin
DEL "%~n1".bin
DEL output.txt
shift
if not [%1]==[] goto loop