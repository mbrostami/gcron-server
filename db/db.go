package db

import pb "github.com/mbrostami/gcron/grpc"

// TaskCollection list of tasks
type TaskCollection struct {
	Tasks map[string]*pb.Task
}

// DB database interface
type DB interface {
	Store(task *pb.Task) (string, error)
	Get(uid uint32, start int, stop int) *TaskCollection
	Close()

	Lock(key string, timeout int32) (bool, error)
	Release(key string) (bool, error)
}
