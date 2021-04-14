package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/tukejonny/trial-grpc-xds/pb"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type pingService struct {
	name string
}

func newPingService(name string) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer(grpc.MaxConcurrentStreams(10))
	pingSvc := &pingService{name: name}

	pb.RegisterPingServiceServer(grpcServer, pingSvc)
	healthpb.RegisterHealthServer(grpcServer, pingSvc)

	return
}

func (s *pingService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("Ping rpc has invoked: %s\n", req.Id)
	return &pb.PingResponse{
		Msg: fmt.Sprintf("Hello %s from %s", req.Id, s.name),
	}, nil
}

func (s *pingService) PingStream(req *pb.PingRequest, stream pb.PingService_PingStreamServer) error {
	log.Println("Ping bidirectional stream rpc invoked")
	stream.Send(&pb.PingResponse{Msg: fmt.Sprintf("Hello %s", req.Id)})
	stream.Send(&pb.PingResponse{Msg: fmt.Sprintf("Hello %s", req.Id)})
	return nil
}

func (s *pingService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *pingService) Watch(req *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
