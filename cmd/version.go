package cmd

import (
"fmt"

"github.com/sirupsen/logrus"
"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Long:  `Print the version number of cluster-state`,
	Run: func(cmd *cobra.Command, arg []string) {
		logrus.Debug("printing version")
		fmt.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
