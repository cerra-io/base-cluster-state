package cmd

import (
	"runtime"

	"github.com/cerra-io/base-goutils/vipersubtree"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"


	"os"
	"os/signal"
	"fmt"
	"github.com/cerra-io/base-cluster-state/server"
)

func configRuntime() {
	nuCPU := runtime.NumCPU()
	old := runtime.GOMAXPROCS(nuCPU)
	logger.WithFields(logrus.Fields{"old": old, "new": nuCPU}).Info("set-go-max-procs")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the application",
	Long:  `Starts all components required by the manager`,
	Run: func(cmd *cobra.Command, arg []string) {

		conf := vipersubtree.Subtree(viper.GetViper(), "server")

		configRuntime()

		conf = vipersubtree.Subtree(viper.GetViper(), "config")

		// Arrange for clean shutdown on keyboard interrupt
		interrupt := make(chan os.Signal, 1)
		done := make(chan bool)
		signal.Notify(interrupt, os.Interrupt)
		go func() {
			<-interrupt // wait for signal
			server.Stop()
			done <- true
		}()

		server.Start(conf)

		// Wait till all cleaned up
		<-done

		// Return only when all server components are cleanly shut down
		fmt.Printf("\nbye\n")

		// able to detect leaking goroutines
		if conf.GetBool("showGoRoutineLeaksOnExit") {
			panic("not really a panic; just to show any leak (you should see only one goroutine printed)")
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
