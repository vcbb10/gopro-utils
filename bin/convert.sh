for i in *.MP4;
        do name=`echo $i | cut -d'.' -f1`+S;
        echo $name;
		ffmpeg -y -i $i -codec copy -map 0:3 -f rawvideo $name.bin
done