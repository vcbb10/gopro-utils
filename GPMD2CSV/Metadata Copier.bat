@echo off
SETLOCAL EnableDelayedExpansion

:: This script simply copies the metadata from source files to edited files and exports new files.
:: Your original GoPro video files (with metadata) and the edited/color corrected files MUST be the same framerate and length.

:: To run the script, you'll want to have your input files (original GoPro files with metadata) and corrected files
:: (with the edits, color corrections, etc.) organized correctly.  Drag any source .MP4 file onto the script and it
:: will process all of the source .MP4 files in that directory that correspond to a corrected file in the "EditedVideoPath".

:: Example (using defaults):
:: Original GoPro files (source files for metadata): "C:\Source Files\GoPro Video 1.mp4"
:: Edited (color corrected files, source for video/audio): "C:\Source Files\Edited Video\GoPro Video 1.mp4"
:: The metadata corrected files will export to: "C:\Source Files\Metadata Fix\GoPro Video 1_MDF.mp4"

:: This is the name of the nested directory that contains the edited/color corrected video files.
set EditedVideoPath=Edited Video

:: Edit this to tell the script anything that has been prepended to the filename of the edited/color corrected video files.
set EditedPrepend=

:: Edit this to tell the script anything that has been postpended to the filename of the edited/color corrected video files.
set EditedPostpend=

:: Directory to export the processed files to (nested under source file directory).
set ExportVideoPath=Metadata Fix

:: Optional text to postpend to the exported video files (with the correct video/metadata)
set MetadataFixedPostPend=_MDF

:::::::::::::::::::::::::::::::::::::::::::::
:: Shouldn't need to edit below this line. ::
:::::::::::::::::::::::::::::::::::::::::::::

Set SourceScriptDirectory=%~dp0

echo %~dp1
cd %~dp1
mkdir "%~dp1\%ExportVideoPath%"

Set OutputDir=%~dp1%ExportVideoPath%

 for %%f in (*.MP4) do (
   if exist "%OutputDir%\%%~nf%MetadataFixedPostPend%.mp4" (
      ECHO.
      ECHO *************************************************
      ECHO ********** Output File already exists! **********
      ECHO *************************************************
   ) else (
         cls
         ECHO.
         ECHO *** Processing file: ***
         ECHO *** "%~dp1%%~f" ***
         ECHO.
         START "" /WAIT "%SourceScriptDirectory%bin\ffmpeg" -i "%~dp1%EditedVideoPath%\%EditedPrepend%%%~nf%EditedPostpend%.mp4" -i "%~dp1%%~f" -c copy -map_metadata 1 -map 0:0 -map 0:1 -map 1:2 -map 1:3 -map 1:4 "%OutputDir%\%%~nf%MetadataFixedPostPend%.mov"
         @echo on
         Move "%OutputDir%\%%~nf%MetadataFixedPostPend%.mov" "%OutputDir%\%%~nf%MetadataFixedPostPend%.mp4"
      )
   )
