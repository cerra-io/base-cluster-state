package update

import (
	"github.com/sirupsen/logrus"
	"time"
	"math/rand"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cerra-io/base-cluster-state/utils"
	"docker.io/go-docker"
	"context"
)

type Update struct {
	NodeType string
	LockTable string
	Region string
	UpdateInterval time.Duration
	LocalIp string
	Client *docker.Client
	dbSess *session.Session
	quit chan bool
}

var (
	logger   = logrus.WithField("module", "update")
)

func (u *Update) Start() {
	sess, err := session.NewSession(&aws.Config{
			Region: aws.String(u.Region)},
		)

	if err != nil {
		logger.Fatalf("Unable to connect to dynamodb: %v", err)
	}

	u.dbSess = sess

	go u.update()

	logger.Info("starting")
}

func (u *Update) Stop() {
	go func() {
		u.quit <- true
	}()
	logger.Info("stopping")
}

func (u *Update) update() {
	tick := time.Tick(u.UpdateInterval * time.Second)
	for {
		select {
		case <-tick:
			u.run()
		case <-u.quit:
			return
		}
	}
}


func (u *Update) run() {
	logger.Debug("starting update")
	if u.NodeType != "manager" {
		return
	}

	delay := rand.Intn(10)
	logger.Infof("starting update cycle with a %d nap", delay)
	// Sleep a random amount of time to avoid to managers trying to remove the same node.
	time.Sleep(time.Duration(delay))


	info, err := u.Client.Info(context.Background())

	if err != nil {
		logger.Errorf("unable to get docker info, %v", err)
		return
	}

	if info.Swarm.NodeID == ""{
		logger.Error("node has not joined the swarm yet")
		return
	}

	nodeInfo, _, err := u.Client.NodeInspectWithRaw(context.Background(), info.Swarm.NodeID)

	if err != nil {
		logger.Errorf("unable to get node info, %v", err)
		return
	}

	if nodeInfo.ManagerStatus.Leader {

		managerIp, err := utils.FetchManagerIp(u.dbSess, u.LockTable)

		if err != nil {
			logger.Errorf("unable to read db, %v", err)
			return
		}

		if u.LocalIp != managerIp {
			_, err := utils.SetManagerDbInfo(u.dbSess, u.LockTable, u.LocalIp)
			if err != nil {
				logger.Errorf("Unable to set manager ip in DB, %v", err)
				return
			}
		}
	}
}