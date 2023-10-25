package asserter

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

func (a *Asserter) ThereAreAuthorizeService() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	network, err := pool.CreateNetwork("App", func(config *docker.CreateNetworkOptions) {
		config.Driver = "bridge"
	})
	time.Sleep(3 * time.Second)
	a.CloseFns = append(a.CloseFns, network.Close)

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "db",
		Repository:   "postgres",
		Tag:          "13.3",
		Networks:     []*dockertest.Network{network},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				docker.PortBinding{
					HostPort: "5432",
				},
			},
		},
		Env: []string{
			"POSTGRES_USER=user",
			"POSTGRES_DB=web",
			"POSTGRES_PASSWORD=password",
			"PGPORT=5432",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	a.CloseFns = append(a.CloseFns, pg.Close)

	web, err := pool.BuildAndRunWithOptions("./dockerfile", &dockertest.RunOptions{
		Links:        []string{"db"},
		Name:         "web_app-1",
		Networks:     []*dockertest.Network{network},
		ExposedPorts: []string{"8001"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8001": {
				docker.PortBinding{
					HostPort: "8001",
				},
			},
		},
		Env: []string{
			"DB_NAME=web",
			"HOST=db",
			"PASSWORD=password",
			"PORT=5432",
			"APP_PORT=8001",
			"USER=user",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		}
	})

	if err != nil {
		return err
	}
	a.CloseFns = append(a.CloseFns, web.Close)

	web2, err := pool.BuildAndRunWithOptions("./dockerfile", &dockertest.RunOptions{
		Name:         "web_app-2",
		Links:        []string{"db"},
		Networks:     []*dockertest.Network{network},
		ExposedPorts: []string{"8002"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8002": {
				docker.PortBinding{
					HostPort: "8002",
				},
			},
		},
		Env: []string{
			"DB_NAME=web",
			"HOST=db",
			"PASSWORD=password",
			"PORT=5432",
			"APP_PORT=8002",
			"USER=user",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		}
	})

	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	a.CloseFns = append(a.CloseFns, web2.Close)

	return nil
}
