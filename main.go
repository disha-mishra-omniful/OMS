package main

import (
	"context"
	"fmt"
	"os"
	"time"

	appinit "awesomeProject5/OMS/init"
	"github.com/omniful/go_commons/config"
	//"github.com/omniful/go_commons/postgres/sql/migration"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	//"github.com/omniful/tenant-services/internal/workers"
	"awesomeProject5/OMS/router"
)

func main() {
	os.Setenv("CONFIG_SOURCE", "local")
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising configs, err: %v", err)
		panic(err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from configs, err: %v", err)
		panic(err)
	}
	//SetupServer(ctx)
	//appinit.Initialize(ctx)
	// ctx :=context.TODO()

	// appinit.InitializeDB(ctx)
	//fmt.Println("Initialising...")
	appinit.Initialize(ctx)
	SetupServer(context.TODO())
}
func SetupServer(ctx context.Context) {
	server := http.InitializeServer(config.GetString(ctx, "server.port"), 10*time.Second, 10*time.Second, 70*time.Second)
	err := router.Initialize(ctx, server)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	fmt.Println("hello")
	log.Infof("Starting server on port" + config.GetString(ctx, "server.port"))
	err = server.StartServer(config.GetString(ctx, "services.name"))
	if err != nil {
		log.Errorf(err.Error())
		panic(err)

	}
	fmt.Println("hello")
	<-shutdown.GetWaitChannel()

}
