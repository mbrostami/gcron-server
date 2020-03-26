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

//go:generate protoc -I ../routeguide --go_out=plugins=grpc:../routeguide ../routeguide/route_guide.proto

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/mbrostami/gcron/grpc"
)

type gcronServer struct {
	pb.UnimplementedGcronServer
	mu sync.Mutex // protects routeNotes
}

// InitializeTask returns the feature at the given point.
func (s *gcronServer) InitializeTask(ctx context.Context, guid *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("Calling method InitializeTask ... %+v", guid)
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

// InitializeTask returns the feature at the given point.
func (s *gcronServer) Log(ctx context.Context, output *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("Calling method Log ... %v", output)
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

// InitializeTask returns the feature at the given point.
func (s *gcronServer) FinializeTask(ctx context.Context, task *pb.Task) (*wrappers.BoolValue, error) {
	log.Printf("Calling method Log ... %+v", task)
	boolValue := &wrappers.BoolValue{Value: true}
	return boolValue, nil
}

func newServer() *gcronServer {
	s := &gcronServer{}
	return s
}

// Run grpc server
func Run() {
	// flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 1402))
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
	grpcServer.Serve(lis)
}
