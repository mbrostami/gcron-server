# gcron server [In Development]
A go written tool to manage distributed cron jobs with centralized GUI. This will help you monitor gcrons if you have multiple servers.


## Requirements 
 - [gcron client](https://github.com/mbrostami/gcron)


## TODO
- [ ] Clean code!
- [ ] Write tests
- [ ] Pick distributed high performance database to store all logs (search optimized, hash O(1) read support)
- [ ] GUI
  - [ ] Authentication
  - [ ] Search logs (tag, hostname, uid, command, guid, output)
- [ ] Log stream proxy... (remote third party log server, REST Api, tcp/udp)
- [x] Remote mutex lock for clients
- [ ] TLS enabled over RPC
- [ ] Client authentication + (caching system)
- [ ] Async write (Get stream logs and write in database async)
- [ ] Handle timeouts
- [ ] Customized taging clientside
- [ ] Support different agents
  


## Dev
`go run main.go`
```
      --log.enable           Enable log in file
      --log.level string     Log level (default "warning")
      --log.path string      Log file path (default "/var/log/gcron/gcron-server.log")
      --server.host string   Server host (default "localhost")
      --server.port string   Server port (default "1400")
```

