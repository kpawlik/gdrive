set GOOS=windows
set GOARCH=amd64
set GO111MODULE=on
cd uploadGD
go build -o ..\\bin\\uploadgd.exe
cd ..\downloadGD
go build -o ..\\bin\\downloadgd.exe
cd ..
