## gurl - simple curl written in Go
```bash
kmille@linbox:gurl gurl -w outputfile.txt -H "X-Go: <3" -A "gurl 0.1" -X POST -d my_post_data https://httpbin.org/anything
kmille@linbox:gurl cat outputfile.txt
{
  "args": {},
  "data": "my_post_data",
  "files": {},
  "form": {},
  "headers": {
    "Accept-Encoding": "gzip",
    "Content-Length": "12",
    "Host": "httpbin.org",
    "User-Agent": "gurl 0.1",
    "X-Amzn-Trace-Id": "Root=1-603e2fd3-63d8eada4bd4ae4c4908f6ac",
    "X-Go": "<3"
  },
  "json": null,
  "method": "POST",
  "origin": "185.104.140.190",
  "url": "https://httpbin.org/anything"
}
```

## FUTURE WORK
- proxy support
- resolve DNS
- choose between ipv4 and ipv6
- show ssl stuff
- verbose output
- tests
