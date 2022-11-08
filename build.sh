if [[ -f "_assets/source.mp4" ]]; then
    echo "Assets already exist, not recreating."
else
    echo "Assets do not exist, creating..."
    mkdir -p _assets/frames/
    curl -o _assets/source.mp4 https://raw.githubusercontent.com/vcokltfre/badapplefs/master/assets/badapple.mp4
    ffmpeg -i _assets/source.mp4 -vf fps=30 _assets/frames/%04d.png
fi

rm -rf _output/
mkdir -p _output/

rm output.mp4 final.mp4

go run main.go

ffmpeg -framerate 30 -pattern_type glob -i '_output/*.png' \
  -c:v libx264 -pix_fmt yuv420p output.mp4

ffmpeg -i _assets/source.mp4 -i output.mp4 -c copy -map 0:a:0 -map 1:v:0 final.mp4

rm output.mp4
