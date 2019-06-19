package vacuum

import (
	"docker.io/go-docker"
	"github.com/sirupsen/logrus"
	"time"
	"math/rand"
	"context"
	"docker.io/go-docker/api/types/filters"
)

type Vacuum struct {
	Client *docker.Client
	VacuumInterval time.Duration
	quit chan bool
}

var (
	logger   = logrus.WithField("module", "vacuum")
)

func (v *Vacuum) Start() {
	go v.vacuum()

	logger.Info("starting")
}

func (v *Vacuum) Stop() {
	go func() {
		v.quit <- true
	}()
	logger.Info("stopping")
}

func (v *Vacuum) vacuum() {
	tick := time.Tick(v.VacuumInterval * time.Hour)
	for {
		select {
		case <-tick:
			v.run()
		case <-v.quit:
			return
		}
	}
}

func (v *Vacuum) run() {
	logger.Debug("starting vacuum")

	delay := rand.Intn(3600)
	logger.Infof("starting update cycle with a %d nap", delay)
	// Sleep a random amount of time to avoid to managers trying to remove the same node.
	time.Sleep(time.Duration(delay))

	ctx := context.Background()

	v.Client.ContainersPrune(ctx, filters.NewArgs())
	v.Client.NetworksPrune(ctx, filters.NewArgs())
	v.Client.ImagesPrune(ctx, filters.NewArgs())

	/*
		this -> v.Client.VolumesPrune(ctx, filters.NewArgs()) does not work because of https://github.com/docker/for-linux/issues/389
	 */

	 f := filters.NewArgs(filters.Arg("dangling", "true"))
	 volumes, err := v.Client.VolumeList(ctx, f)
	 if err != nil {
	 	logger.Errorf("Unable to obtain dangling list of volumes")
	 	return
	 }

	 for _, vol := range volumes.Volumes {
	 	err = v.Client.VolumeRemove(ctx, vol.Name, true)
	 	logger.Errorf("unable to vacuum volume %s", vol.Name)
	 }
}