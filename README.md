# gcron server [In Development]
A go written tool to manage distributed cron jobs with centralized GUI. This will help you monitor gcrons if you have multiple servers.

## Features
- Supporting GCron 
- Centralized logging
- GUI interface
  - Search outputs by tags/server/text
  - Check crons status (exit code, running time, last run, ...)
  - Graphs 
  - Realtime monitoring
  - Resource usage per cron/server
- Accept syslog formats as input
- Auto detect input tags (specific syntax)
- Pipe input logs to remote syslog
- Pipe input logs to remote REST Api 

## Dev
TCP Server `go run main.go`  
TCP Server `go run main.go --host=0.0.0.0`  
TCP Server `go run main.go --host=0.0.0.0 --port=1400`  
UDP Server `go run main.go --prot=tcp`  
UNIX Server `go run main.go --prot=unix`  

## TODO
- All Features! :D



