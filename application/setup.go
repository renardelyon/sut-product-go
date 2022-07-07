package application

import (
	"context"
	"log"
	"net"
	"sut-product-go/config"
	"sut-product-go/lib/pkg/db"

	"google.golang.org/grpc"
)

func initGrpcServer(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		g := grpc.NewServer()
		app.GrpcServer = g
		return nil
	}
}

func grpcRun(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		lis, err := net.Listen("tcp", cfg.Port)
		if err != nil {
			return err
		}
		log.Println("Order service on Port: ", cfg.Port)
		if err := app.GrpcServer.Serve(lis); err != nil {
			return err
		}
		app.GrpcServer.GracefulStop()
		return nil
	}
}

func Setup(cfg *config.Config) (*Application, error) {
	app := new(Application)
	err := runInit(
		initDatabase(cfg),
		initGrpcClient(cfg),
		initApp(cfg))(app)

	if err != nil {
		return app, err
	}
	return app, nil
}

func runInit(appFuncs ...func(*Application) error) func(*Application) error {
	return func(app *Application) error {
		app.Context = context.Background()
		for _, fn := range appFuncs {
			if e := fn(app); e != nil {
				return e
			}
		}
		return nil
	}
}

func initGrpcClient(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		app.GrpcClients = make(map[string]*grpc.ClientConn)

		// ProductServiceCfg := cfg.ProductHost
		// conn, err := setupGrpcConnection(ProductServiceCfg)
		// if err != nil {
		// 	return err
		// }

		// app.GrpcClients["product-management-service"] = conn
		// log.Println("product-management-service" + " connected on " + ProductServiceCfg)

		log.Println("init Grpc Client done")
		return nil
	}
}

func setupGrpcConnection(cfg string) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg, grpc.WithInsecure())
}

func initDatabase(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		app.DbClients = db.Init(cfg.DBUrl)

		log.Println("init postgre database done")

		return nil
	}
}

func initApp(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		return initGrpcServer(cfg)(app)
	}
}
