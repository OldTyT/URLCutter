# Link shortener.

---
## How to use

### Generate new URLCut

You need to make a request to the endpoint `/url/create` with the url key, as a result, the id will be returned. Example:

```shell script
root@gusev:~# curl localhost:8080/url/create?url=github.com
http://localhost:8080/url?id=zqQvHHHx1b
```

### Use redirect

You need to make a request to the endpoint `/url` with the id key. Example:

```shell script
root@gusev:~# curl http://localhost:8080/url?id=zqQvHHHx1b
<a href="http://github.com">See Other</a>.

root@gusev:~# curl -ILs http://localhost:8080/url?id=zqQvHHHx1b
HTTP/1.1 303 See Other
Content-Type: text/html; charset=utf-8
Location: http://github.com
Date: Wed, 20 Jul 2022 14:31:41 GMT

HTTP/1.1 301 Moved Permanently
Content-Length: 0
Location: https://github.com/

HTTP/2 200 
server: GitHub.com
```
**Clone and run URL Cutter**

```shell script
git clone https://github.com/OldTyT/URLCutter.git
cd URLCutter
go run .
```

**Clone and build URL Cutter**

```shell script
git clone https://github.com/OldTyT/URLCutter.git
cd URLCutter
go build .
./URLCutter
```