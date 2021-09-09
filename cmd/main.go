package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	// "github.com/improbable-eng/grpc-web/go/grpcweb"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mintak21/grpc-sample/server/service"
	pb "github.com/mintak21/proto/sample/golang"
)

var (
	port          uint = 8080
	useReflection bool = true
	useGRPCWeb    bool = true
)

const (
	serviceName = "pancake baker"
)

func main() {
	log.Println("starting service.")

	app := cli.NewApp()
	app.Name = serviceName
	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:        "port, p",
			Usage:       "listening port number",
			Value:       port,
			Destination: &port,
		},
		&cli.BoolFlag{
			Name:        "reflection, r",
			Usage:       "which to use gRPC reflection",
			Value:       useReflection,
			Destination: &useReflection,
		},
		&cli.BoolFlag{
			Name:        "web, w",
			Usage:       "which to use gRPC web",
			Value:       useGRPCWeb,
			Destination: &useGRPCWeb,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				grpc_validator.UnaryServerInterceptor(),
			),
		),
	)
	// TODO wire gen
	pancakeBaker := service.NewBakePancakeService()
	pb.RegisterBakePancakeServiceServer(server, pancakeBaker)
	if useReflection {
		reflection.Register(server)
	}
	if useGRPCWeb {
		// TODO
	}

	go func() {
		log.Printf("start %s server. port: %d", serviceName, port)
		if err := server.Serve(lis); err != nil {
			log.Fatalln("failed to serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Printf("stop %s server..", serviceName)
	defer server.GracefulStop()

	return nil
}
