package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb_jg "github.com/msqtt/sb-judger/api/pb/v1/judger"
	"github.com/msqtt/sb-judger/internal/pkg/config"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/msqtt/sb-judger/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main()  {
	conf, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalln(err)
	}	
	lcm, err := json.GetLangConfs("./configs/lang.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := service.NewJudgerServer(conf, lcm)
	mux := runtime.NewServeMux()
	ctx := context.Background()
	ctx, cf := context.WithCancel(ctx)
	defer cf()

	grpcServer := grpc.NewServer()
	pb_jg.RegisterCodeServer(grpcServer, js)
	err = pb_jg.RegisterCodeHandlerServer(ctx, mux, js)
	if err != nil {
		log.Fatalln(err)
	}

	reflection.Register(grpcServer)

	grpcL, err := net.Listen("tcp", conf.GrpcAddr)
	if err != nil {
		log.Fatalln(err)
	}
	restL, err := net.Listen("tcp", conf.HttpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("GRPC Server will start at %s", grpcL.Addr().String())
	log.Printf("REST Server will start at %s", restL.Addr().String())
	go log.Fatalln(grpcServer.Serve(grpcL))
	log.Fatalln(http.Serve(restL, mux))
}
