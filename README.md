v1p
-----------

Help:
```shell
[irocha@irrlab v1p (master)]$ ./v1p 
v1p version 0.1 (ivan.ribeiro@gmail.com)
v1p [-s][-h][-t] -l <addr:port> -r <addr:port>
  -c="": config file
  -h=false: help
  -l="": saddr:port (local)
  -r="": raddr:port (remote)
  -s=false: syslog (enabled/disabled)
  -t=0: timeout (seconds)
```

Using command line options:
```shell
[irocha@irrlab v1p (master)]$ ./v1p -l localhost:7777 -r google.com:80
[v1p] 2014/07/30 09:20:41 proxying localhost:7777 to [google.com:80] (t:0)...
...
[v1p] 2014/07/30 09:21:44 127.0.0.1:7777 > 172.22.33.140:33283 173.194.118.41:80 168 [OK]
[v1p] 2014/07/30 09:21:44 127.0.0.1:7777 < 173.194.118.41:80 172.22.33.140:33283 287 [OK]
...
```

Log format:
```shell
[v1p] [DATE] [LOCAL ADDR/PORT] [IN="<"|OUT=">"] [FROM] [TO] [bytes] [OK|ERR]
```
```shell
[irocha@irrlab v1p (master)]$ curl -I localhost:7777
HTTP/1.1 302 Found
Location: http://www.google.com/
Cache-Control: private
Content-Type: text/html; charset=UTF-8
X-Content-Type-Options: nosniff
Date: Wed, 30 Jul 2014 12:21:44 GMT
Server: sffe
Content-Length: 219
X-XSS-Protection: 1; mode=block
Alternate-Protocol: 80:quic
```

Using configuration file:
```shell
[irocha@irrlab v1p (master)]$ cat vcfg/vcfg.json 
[ {"Local":"127.0.0.1:7777", "Remote":["www.uol.com.br:80", "www.bol.com.br:80"], "Timeout":10}, 
  {"Local":"127.0.0.1:8888", "Remote":["www.google.com:80"]} ]
```
```shell
[irocha@irrlab v1p (master)]$ ./v1p -c vcfg/vcfg.json 
[v1p] 2014/07/30 09:23:31 proxying 127.0.0.1:7777 to [www.uol.com.br:80 www.bol.com.br:80] (t:10)...
[v1p] 2014/07/30 09:23:31 proxying 127.0.0.1:8888 to [www.google.com:80] (t:0)...
...
[v1p] 2014/07/30 09:24:30 127.0.0.1:7777 < 200.147.67.142:80 172.22.33.140:42256 400 [OK]
[v1p] 2014/07/30 09:24:30 127.0.0.1:7777 > 172.22.33.140:42256 200.147.67.142:80 167 [OK]
[v1p] 2014/07/30 09:24:37 127.0.0.1:8888 > 172.22.33.140:43645 173.194.115.82:80 167 [OK]
[v1p] 2014/07/30 09:24:37 127.0.0.1:8888 < 173.194.115.82:80 172.22.33.140:43645 506 [OK]
...
```
```shell
[irocha@irrlab v1p (master)]$ curl localhost:7777
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>302 Found</title>
</head><body>
<h1>Found</h1>
<p>The document has moved <a href="http://www.uol.com.br/">here</a>.</p>
</body></html>
[irocha@irrlab v1p (master)]$ curl localhost:8888
<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>302 Moved</TITLE></HEAD><BODY>
<H1>302 Moved</H1>
The document has moved
<A HREF="http://www.google.com/">here</A>.
</BODY></HTML>
```

Retrieving **JSON** statistics (last 15 minutes) using 1 slot per minute:
```shell
[irocha@irrlab v1p (master)]$ curl localhost:1972 |python -mjson.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
101  4160    0  4160    0     0  3548k      0 --:--:-- --:--:-- --:--:-- 4062k
[
    {
        "Counters": {
            "127.0.0.1:7777": {
                "BytesIn": 1346, 
                "BytesOut": 501, 
                "Errors": 0, 
                "Remote": [
                    "www.uol.com.br:80", 
                    "www.bol.com.br:80"
                ], 
                "Success": 6
            }, 
            "127.0.0.1:8888": {
                "BytesIn": 1518, 
                "BytesOut": 501, 
                "Errors": 0, 
                "Remote": [
                    "www.google.com:80"
                ], 
                "Success": 6
            }
        }, 
        "Date": "2014-07-30T09:27:29.198259579-03:00"
    }, 
    {
        "Counters": {
            "127.0.0.1:7777": {
                "BytesIn": 0, 
                "BytesOut": 0, 
                "Errors": 0, 
                "Remote": [
                    "www.uol.com.br:80", 
                    "www.bol.com.br:80"
                ], 
                "Success": 0
            }, 
            "127.0.0.1:8888": {
                "BytesIn": 0, 
                "BytesOut": 0, 
                "Errors": 0, 
                "Remote": [
                    "www.google.com:80"
                ], 
                "Success": 0
            }
        }, 
        "Date": "2014-07-30T09:26:29.198259579-03:00"
    }, 
    ...
]
```

Copyright and License
---------------------
Copyright 2014 Ivan Ribeiro Rocha

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

