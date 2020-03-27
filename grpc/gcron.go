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
	"io"
	"net"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/mbrostami/gcron/grpc"
	"github.com/mbrostami/gcron/helpers"
	log "github.com/sirupsen/logrus"
)

type gcronServer struct {
	pb.UnimplementedGcronServer
	mux *helpers.Mutex
}

// Lock mutex lock by name
func (s *gcronServer) Lock(ctx context.Context, lockName *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Debugf("Locking ... %+v", lockName.GetValue())
	mux, _ := helpers.NewMutex(lockName.GetValue())
	s.mux = mux
	locked, err := s.mux.Lock()
	if locked {
		log.Debugf("Locked! %v", lockName.GetValue())
	} else {
		log.Debugf("Already locked! %v", lockName.GetValue())
	}
	boolValue := &wrappers.BoolValue{Value: locked}
	return boolValue, err
}

// Release release the lock
func (s *gcronServer) Release(ctx context.Context, lockName *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Debugf("Releasing ... %+v", lockName.GetValue())
	released, err := s.mux.Release()
	boolValue := &wrappers.BoolValue{Value: released}
	return boolValue, err
}

// Log returns the feature at the given point.
func (s *gcronServer) StartLog(stream pb.Gcron_StartLogServer) error {
	log.Debugf("Calling method StartLog...")
	var pointCount int32
	var lastLog *pb.LogEntry
	for {
		logEntry, err := stream.Recv()
		if err == io.EOF {
			log.Tracef("Last log %v", lastLog)
			return stream.SendAndClose(&wrappers.BoolValue{Value: true})
		}
		if err != nil {
			return err
		}
		log.Tracef("Incoming stream %v", logEntry)
		pointCount++
		lastLog = logEntry
	}
}

// Done Save the latest state of task.
func (s *gcronServer) Done(ctx context.Context, task *pb.Task) (*wrappers.BoolValue, error) {
	log.Debugf("Calling method Done ... %+v", task)
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

func newServer() *gcronServer {
	return &gcronServer{}
}

// Run grpc server
func Run(host string, port string) {
	lis, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", host+":"+port, err)
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
	log.Infof("Started listening on: %s", host+":"+port)
	grpcServer.Serve(lis)
}
