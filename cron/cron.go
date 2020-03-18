package cron

import "time"

// Task keeps cronjob information
type Task struct {
	Pid        int
	Username   string
	Command    string
	StartTime  time.Time
	EndTime    time.Time
	ExitCode   int
	Output     []byte
	SystemTime time.Duration
	UserTime   time.Duration
	Success    bool
}
