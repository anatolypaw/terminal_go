GOARCH=arm GOARM=7 GOOS=linux go build main.go
lftp sftp://root:owen@10.0.4.10:/home/root  -e "put main; bye"