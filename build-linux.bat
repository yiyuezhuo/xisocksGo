set GOOS=linux
go env
go build -ldflags "-w" -o bin/server github.com/yiyuezhuo/xisocksGo/server
go build -ldflags "-w" -o bin/client github.com/yiyuezhuo/xisocksGo/client
copy config-client.json bin
copy config-server.json bin
xcopy static bin\static\ /y
7z a -tzip release/linux.zip bin