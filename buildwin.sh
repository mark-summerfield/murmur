cd bin
env GOOS=windows GOARCH=386 go build -o murmur.exe .
mv murmur.exe ..
cd ..
