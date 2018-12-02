set GOOS=linux
set GOARCH=amd64
set GO111MODULE=on
cd uploadGD
go build -o ..\\bin\\uploadgd
cd ..\downloadGD
go build -o ..\\bin\\downloadgd
cd ..
