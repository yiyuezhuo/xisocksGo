go env
go build -ldflags "-w" -o bin_windows/server.exe github.com/yiyuezhuo/xisocksGo/server
go build -ldflags "-w" -o bin_windows/client.exe github.com/yiyuezhuo/xisocksGo/client
copy config-client.json bin_windows
copy config-server.json bin_windows
xcopy static bin_windows\static\ /y
7z a -tzip release/windows.zip bin_windows