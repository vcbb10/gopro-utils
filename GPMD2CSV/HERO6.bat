if [%1]==[] goto :eof
:loop
START /WAIT .\ffmpeg -y -i "%~1" -codec copy -map 0:3 -f rawvideo "%~n1".bin
START /WAIT .\gpmd2csv -i "%~n1".bin
DEL "%~n1".bin
shift
if not [%1]==[] goto loop