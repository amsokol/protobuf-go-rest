package main

import (
	"context"
	"log"
	"net/http"

	"github.com/amsokol/protobuf-rest/examples/gw-hello-world/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type greeterServer struct {
	proto.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func main() {
	var ctx context.Context

	mux := runtime.NewServeMux()

	var srv greeterServer

	if err := proto.RegisterGreeterHandlerServer(ctx, mux, &srv); err != nil {
		log.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8080")
	log.Fatalln(gwServer.ListenAndServe())
}
