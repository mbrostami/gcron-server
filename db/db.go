package db

import pb "github.com/mbrostami/gcron/grpc"

// DB database interface
type DB interface {
	Store(task *pb.Task) (string, error)
	Search(text string, limit int)
	Close()
}
