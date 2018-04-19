package main

import "github.com/cerra-io/base-cluster-state/cmd"

var (
	version string
)

func main() {
	cmd.Version = version
	cmd.Execute()
}