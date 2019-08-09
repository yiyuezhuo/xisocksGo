# A prototype proxy server using websocket and CDN

My bandwagon server ip running shadowsocks is blocked by GFW.
So I write a tool to connect my server again.

This is a port of [Python-version](https://github.com/yiyuezhuo/xisock0). The Python version doesn't work recently for some reason.

## Usage

You can use this tool with or without hostname, CDN and TLS. The default config run the tool under `WS+CDN+TLS`.

* Download respective version into your local computer(client) and VPS(server) from [release page](https://github.com/yiyuezhuo/xisocksGo/releases).
* In client, replace config-client.json item "your_hostname" with your hostname such as "xisocks.com", which have been "protected" by CDN such as cloudflare. Then CCP internet cops can't find your real IP. 




* See some TLS 


* In client, run `client.exe`(windows) or `./client`(linux)
* In server, run `server.exe`(windows) or `./server`(linux)

### Without TLS

If you don't want be protected by TLS which may slow down your network:

* In both `config-client.json` and `config-server.json`, set `TLS` to `false`.
* Set `RemotePort` in `config-client.json` and `ListenPort` in `config-server.json` to `80`.

### And without hostname or CDN

Then you can set your `RemotePort` and `ListenPort` to any values. And use a ip in `RemoteIp` instead of hostname.
It runs faster in many situation, but can't save your ip from blocking.

### Tips

In linux, `screen` is a elegant way to run your server without connection to your server.

## build

### From Windows

```
$ go get github.com/gorilla/websocket
$ build-linux
$ build-windows
```

### From Linux

TODO

## Fake website

In fact, `client` connect to `host_name/upload` to proxy data, if a disgusting CCP internet cop try to access `host_name`,
he will get a [sadpanda](https://knowyourmeme.com/memes/sad-panda) to mock their blocking policy:
 
<img src="static/sadpanda.jpg">