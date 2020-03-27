/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../grpc --go_out=plugins=grpc:../grpc ../groc/grpc.proto

// Package grpc implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the gcron service whose definition can be found in gcron/grpc/gcron.proto.
package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/mbrostami/gcron/grpc"
	"github.com/mbrostami/gcron/helpers"
)

type gcronServer struct {
	pb.UnimplementedGcronServer
	mux *helpers.Mutex
}

// Lock mutex lock by name
func (s *gcronServer) Lock(ctx context.Context, lockName *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("Locking ... %+v", lockName.GetValue())
	mux, _ := helpers.NewMutex(lockName.GetValue())
	s.mux = mux
	locked, err := s.mux.Lock()
	if locked {
		log.Printf("Locked! %v", lockName.GetValue())
	} else {
		log.Printf("Already locked! %v", lockName.GetValue())
	}
	boolValue := &wrappers.BoolValue{Value: locked}
	return boolValue, err
}

// Release release the lock
func (s *gcronServer) Release(ctx context.Context, lockName *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("Releasing ... %+v", lockName.GetValue())
	released, err := s.mux.Release()
	boolValue := &wrappers.BoolValue{Value: released}
	return boolValue, err
}

// Log returns the feature at the given point.
func (s *gcronServer) Log(ctx context.Context, logEntry *pb.LogEntry) (*wrappers.BoolValue, error) {
	log.Printf("Calling method Log ... %v : %v", logEntry.GUID, logEntry.Output)
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

// FinializeTask returns the feature at the given point.
func (s *gcronServer) FinializeTask(ctx context.Context, task *pb.Task) (*wrappers.BoolValue, error) {
	log.Printf("Calling method FinializeTask ... %+v", string(task.Output))
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

func newServer() *gcronServer {
	s := &gcronServer{}
	return s
}

// Run grpc server
func Run(host string, port string) {
	// flag.Parse()
	lis, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = testdata.Path("server1.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = testdata.Path("server1.key")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGcronServer(grpcServer, newServer())
	log.Printf("Started listening on : %v", host+":"+port)
	grpcServer.Serve(lis)
}
