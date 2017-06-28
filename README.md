GoPro Metadata Format Parser
============================

I spent some time trying to reverse-engineer the GoPro Metadata Format (GPMD or GPMDF) that is stored in GoPro Hero 5 cameras if GPS is enabled. This is what I found.

Part of this code is in production on [Earthscape](https://public.earthscape.com/); for an example of what you can do with the extracted data, see [this video](https://public.earthscape.com/videos/10231).

If you enjoy working on this sort of thing, please see our [careers page](https://churchillnavigation.com/careers/).

### Dependencies:

- Golang
- FFmpeg
- [Go.geo](https://github.com/paulmach/go.geo)
- [Go.geojson](https://github.com/paulmach/go.geojson)

The go dependencies can be installed with: 

```
go get github.com/paulmach/go.geo
go get github.com/paulmach/go.geojson
go get github.com/stilldavid/gopro-utils/telemetry
```

Extracting the Metadata File
----------------------------

The metadata stream is stored in the `.mp4` video file itself alongside the video and audio streams. We can use `ffprobe` to find it:

```
[computar][100GOPRO] ➔ ffprobe GOPR0008.MP4
ffprobe version 3.2.4 Copyright (c) 2007-2017 the FFmpeg developers
[SNIP]
    Stream #0:3(eng): Data: none (gpmd / 0x646D7067), 33 kb/s (default)
    Metadata:
      creation_time   : 2016-11-22T23:42:41.000000Z
      handler_name    : 	GoPro MET
[SNIP]
```

We can identify it by the `gpmd` in the tag string - in this case it's id 3. We can then use `ffmpeg` to extract the metadata stream into a binary file for processing:

`ffmpeg -y -i GOPR0001.MP4 -codec copy -map 0:3 -f rawvideo out-0001.bin`

This leaves us with a binary file with the data.

## Build and run:

There are several utilities in this repo that generate different type of data:

- ```gopro2json``` exports a JSON file with all the labels
- ```gpmd2csv``` exports a CSV database for acceleration, gps lat/long/alt, gyroscope and camera temperature
- ```gpmdinfo``` prints the metadata
- ```gps2kml``` exports a KML file of the GPS track in the video
- ```gpmd2srt``` exports a .srt file with acceleration readings (X,Y,Z)
- ```stats``` prints stats from your video

Clone the repo and build each of the utilities by doing ```go build``` inside the corresponding directory.

### Which cameras are supported:

- HERO5 Black: Gyro, Accelerometer, GPS (Lat/Long/Altitude/Speed/Speed3d), GPS Fix & accuracy, temperature
- HERO5 Session: Gyro, Accelerometer, Temperature

Data We Get
-----------

* ~400 Hz 3-axis gyro readings
* ~200 Hz 3-axis accelerometer readings
* ~18 Hz GPS position (lat/lon/alt/spd)
* 1 Hz GPS timestamps
* 1 Hz GPS accuracy (cm) and fix (2d/3d)
* 1 Hz temperature of camera

---


The Protocol
------------

Data starts with a label that describes the data following it. Values are all big endian, and floats are IEEE 754. Everything is packed to 4 bytes where applicable, padded with zeroes so it's 32-bit aligned.

 * **Labels** - human readable types of proceeding data
 * **Type** - single ascii character describing data
 * **Size** - how big is the data type
 * **Count** - how many values are we going to get
 * **Length** = size * count

Labels include:

 * `ACCL` - accelerometer reading x/y/z
 * `DEVC` - device 
 * `DVID` - device ID, possibly hard-coded to 0x1
 * `DVNM` - device name, string "Camera"
 * `EMPT` - empty packet
 * `GPS5` - GPS data (lat, lon, alt, speed, 3d speed)
 * `GPSF` - GPS fix (none, 2d, 3d)
 * `GPSP` - GPS positional accuracy in cm
 * `GPSU` - GPS acquired timestamp; potentially different than "camera time"
 * `GYRO` - gryroscope reading x/y/z
 * `SCAL` - scale factor, a multiplier for subsequent data
 * `SIUN` - SI units; strings (m/s², rad/s)
 * `STRM` - ¯\\\_(ツ)\_/¯
 * `TMPC` - temperature
 * `TSMP` - total number of samples
 * `UNIT` - alternative units; strings (deg, m, m/s)

Types include:

 * `c` - single char
 * `L` - unsigned long
 * `s` - signed short
 * `S` - unsigned short
 * `f` - 32 float

For implementation details, see `reader.go` and other corresponding files in `telemetry/`.

### Credits:

- @StillDavid
- @JuanIrache
- @dnewman-gpsw
