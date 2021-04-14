package main

import (
	"context"
	"log"

	api "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

func makeCallbacks() server.CallbackFuncs {
	return server.CallbackFuncs{
		// sotw (State of The World)
		StreamOpenFunc: func(ctx context.Context, id int64, url string) error {
			log.Printf("[+] Open stream: id=%d, url=%s\n", id, url)
			return nil
		},
		StreamClosedFunc: func(int64) {
			log.Println("[-] Close stream")
		},
		StreamRequestFunc: func(id int64, req *api.DiscoveryRequest) error {
			log.Printf("[<==] id=%d, req=%s", id, req.String())
			return nil
		},
		StreamResponseFunc: func(id int64, req *api.DiscoveryRequest, resp *api.DiscoveryResponse) {
			log.Printf("[==>] id=%d, req=%s, resp=%s\n", id, req.String(), resp.String())
		},
		// rest
		FetchRequestFunc: func(ctx context.Context, req *api.DiscoveryRequest) error {
			log.Printf("[<==] Fetch: req=%s\n", req.String())
			return nil
		},
		FetchResponseFunc: func(req *api.DiscoveryRequest, resp *api.DiscoveryResponse) {
			log.Printf("[==>] Fetch: req=%s, resp=%s\n", req.String(), resp.String())
		},
	}
}
