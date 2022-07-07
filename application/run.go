package application

import "sut-product-go/config"

func (app *Application) Run(cfg *config.Config) error {
	return grpcRun(cfg)(app)
}
