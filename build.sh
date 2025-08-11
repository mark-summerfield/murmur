cd cmd
go build -o murmur .
strip murmur
upx -q --best --lzma murmur
mv murmur ..
cd ..
