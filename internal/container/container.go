package container

import (
	"fmt"
	"net/http"
)

type Container struct {
	Config Config

	Apps        *Apps
	Clients     *Clients
	Storages    *Storages
	Services    *Services
	AppServices *ApplicationServices
	Databus     *Broker
}

type Apps struct {
	// Application layer.
}

type Broker struct {
	// Broker.
}

type Clients struct {
	// Other client.
}

type Storages struct {
	// Persistence layer.
}

type Services struct {
	// Domain services.
}

type ApplicationServices struct {
	// Application services.
}

func NewContainer() *Container {
	cfg := NewConfig()

	conn := setupDatabase(cfg.DB)
	setupMigrations(conn)

	container := &Container{Config: cfg}

	container.Databus = makeDatabus(cfg.Amqp, cfg.Environment, cfg.Branch)
	container.Clients = makeClients(container, cfg)
	container.Storages = makeStorages(container)
	container.Services = makeServices(container)
	container.Apps = makeApps(container)
	container.AppServices = makeAppServices(container, cfg)

	return container
}

func (c *Container) Run() error {
	p := fmt.Sprintf(":%s", c.Config.HttpPort)

	return http.ListenAndServe(p, nil)
}
