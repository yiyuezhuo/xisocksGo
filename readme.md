# A prototype proxy server using websocket and CDN

My bandwagon server ip running shadowsocks is blocked by GFW.
So I write the tool to connect the server again.

This is a port of [Python-version](https://github.com/yiyuezhuo/xisock0). The Python version doesn't work recently for some reason.

## Usage

* Download respective version into your computer(client) and VPS(server) from [release page](https://github.com/yiyuezhuo/xisocksGo/releases).
* In client, replace config-client.json item "your_hostname" with your hostname such as "xisocks.com", which have been "protected" by CDN such as cloudflare. 
* In client, `client.exe`(windows) or `./client`(linux)
* In server, `server.exe`(windows) or `./server`(linux)

## build

### Windows

```
$ build-linux
$ build-windows
```

### Linux

TODO
