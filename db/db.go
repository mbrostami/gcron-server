package db

import pb "github.com/mbrostami/gcron/grpc"

// DB database interface
type DB interface {
	Store(task *pb.Task) (string, error)
	Get(uid uint32, start int, stop int)
	Close()
}
