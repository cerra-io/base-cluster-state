package clean

import (
	"github.com/sirupsen/logrus"
	"time"
	"math/rand"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"context"
	"docker.io/go-docker/api/types/swarm"
)

type Clean struct {
	NodeType string
	CleanInterval time.Duration
	quit chan bool
	Client *docker.Client
}

var (
	logger   = logrus.WithField("module", "clean")
)

func (c *Clean) Start() {

	go c.clean()

	logger.Info("starting")
}

func (c *Clean) Stop() {
	go func() {
		c.quit <- true
	}()
	logger.Info("stopping")
}

func (c *Clean) clean() {
	tick := time.Tick(c.CleanInterval * time.Second)
	for {
		select {
			case <-tick:
				c.run()
			case <-c.quit:
				return
		}
	}
}

func (c *Clean) run() {
	logger.Debug("starting clean")
	if c.NodeType != "manager" {
		return
	}
	delay := rand.Intn(10)
	logger.Infof("Starting cleaning cycle with a %d nap", delay)
	// Sleep a random amount of time to avoid to managers trying to remove the same node.
	time.Sleep(time.Duration(delay))

	nodeList, err := c.nodeList()

	if err != nil {
		logger.Errorf("Unable to get manager list, %v", err)
		return
	}

	c.pruneNodes(nodeList)
}

func (c *Clean) nodeList() ([]swarm.Node, error) {
	opts := types.NodeListOptions{}

	list, err := c.Client.NodeList(context.Background(), opts)

	if err != nil {
		logger.Errorf("Unable to get node list, %v", err)
		return nil, err
	}

	return list, nil
}

func (c *Clean) pruneNodes(nodes []swarm.Node) {
	logger.Info("attempting to prune nodes")
	for _, node := range nodes {
		if node.Status.State == swarm.NodeStateDown {
			logger.Infof("node %s is down; pruning", node.ID)

			if node.Spec.Role != swarm.NodeRoleWorker {
				node.Spec.Role = swarm.NodeRoleWorker
				err := c.Client.NodeUpdate(context.Background(), node.ID, node.Version, node.Spec)
				if err != nil {
					logger.Errorf("unable to demote node, %v", err)
					return
				}
			}

			err := c.Client.NodeRemove(context.Background(), node.ID, types.NodeRemoveOptions{
				Force: true,

			})

			if err != nil {
				logger.Errorf("unable to remove node, %v", err)
				return
			}

			logger.Infof("successfully removed node %s", node.ID)
		}
	}
}