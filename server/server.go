package server

import (
	"github.com/cerra-io/base-goutils/vipersubtree"
	"github.com/sirupsen/logrus"
	"github.com/cerra-io/base-cluster-state/clean"

	"docker.io/go-docker"
	"github.com/cerra-io/base-cluster-state/update"
	"github.com/cerra-io/base-cluster-state/vacuum"
)

var (
	logger   = logrus.WithField("module", "server")
	cleaner *clean.Clean
	updater *update.Update
	vacuumer *vacuum.Vacuum
)

func Start(conf *vipersubtree.ViperSubtree) {

	dockerClient, err := docker.NewEnvClient()

	if err != nil {
		logger.Fatalf("Unable to connect to docker, %v", err)
	}

	cleaner = &clean.Clean{
		NodeType: conf.GetString("nodeType"),
		CleanInterval: conf.GetDuration("cleanInterval"),
		Client: dockerClient,
	}

	updater = &update.Update{
		LockTable: conf.GetString("lockTable"),
		UpdateInterval: conf.GetDuration("updateInterval"),
		Region: conf.GetString("region"),
		NodeType: conf.GetString("nodeType"),
		LocalIp: conf.GetString("localIp"),
		Client: dockerClient,
	}

	vacuumer = &vacuum.Vacuum{
		Client: dockerClient,
		VacuumInterval: conf.GetDuration("vacuumInterval"),
	}

	updater.Start()
	cleaner.Start()
	vacuumer.Start()

	logger.Info("starting")
}

func Stop() {
	updater.Stop()
	cleaner.Stop()
	vacuumer.Stop()

	logger.Info("stopping")
}