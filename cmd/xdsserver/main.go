package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
)

var (
	port = flag.String("port", "8080", "xdsサーバのリッスンするポート")
)

const (
	maxConcurrentStreams = 10
)

func main() {
	flag.Parse()

	ctx := context.Background()

	nodeId, err := uuid.NewUUID()

	snapshot := generateSnapshot("1")
	if snapshot.Consistent(); err != nil {
		log.Fatalf("Snapshot inconsistency: %+v\n%+v", snapshot, err)
	}

	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	if err := snapshotCache.SetSnapshot(nodeId.String(), snapshot); err != nil {
		log.Fatalf("Snapshot error: %q for %+v", err, snapshot)
	}

	cb := makeCallbacks()
	srv := server.NewServer(ctx, snapshotCache, cb)

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(maxConcurrentStreams))
	// ADS
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	// EDS
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	// CDS
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	// RDS
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)
	// LDS
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)
	// SDS
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, srv)
	// RTDS
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, srv)

	lis, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %+v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err.Error())
	}
}
