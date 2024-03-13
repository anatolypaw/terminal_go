GOARCH=amd64  GOOS=linux go build main.go

lftp -p8701 sftp://pas:$PASS@192.168.10.89:/home/pas/  -e "put main; bye"
rm main