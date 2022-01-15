# Facade

A cached content proxy server -- piping requests through proxies & different IPs, to alleviate rate limiting.

Plans:
- Implement Unix Domain Sockets to reduce http & tcp overhead, for if the server is running on the same machine as its user
- Make a configuration system to run Facade for any basic service with no extra Go code
- Use fasthttp for the htto clients
- Change cache
