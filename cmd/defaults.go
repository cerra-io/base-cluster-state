package cmd

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"
)

// Defaults is the system-wide default config, also showing what is configurable
var Defaults = []byte(`
logging:
    # Application-wide default log level
    #
    # Permitted values are:
    #   debug
    #   info
    #   warn or warning
    #   error
    #   critical
    #   panic
    defaultLevel: info

    # Time-stamp format
    # Examples:
    #   '2006-01-02 15:04:05' - for full time-stamp with second resolution
    #   '2006-01-02T15:04:05.000000000Z07:00' - same, but nanoseconds and time zone offset
    timeFormat: '2006-01-02 15:04:05'

    # Per-subsystem log level overrides
    manager: info
    client: info

    format: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'

config:
    cleanInterval: 20
    updateInterval: 30
    region: ""
    nodeType: ""
    localIp: ""
    lockTable: ""

server:
    # Setting for runtime.GOMAXPROCS(n). If n < 1, then it has no effect (so we
    # use system default)
    gomaxprocs: 0

    # Frequency of self-diagnostics
    diagInterval: "5s"

    # Show go-routine leaks at exit
    showGoRoutineLeaksOnExit: false

    # Shard number for this instance (must be between 0...4095)
    shard: 0
`)

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Print defaults as a yaml file",
	Long:  `Print the defaults to initialize a new config file`,
	Run: func(cmd *cobra.Command, arg []string) {
		for _, line := range strings.Split(string(Defaults), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				// print empty line
				fmt.Println()
			} else if strings.HasPrefix(trimmed, "#") {
				// already a comment, print it as is
				fmt.Println(line)
			} else if strings.HasSuffix(trimmed, ":") {
				// a container line (dictionary), print it as is
				fmt.Println(line)
			} else {
				// an actual value line, print with a "# " prefixed in front
				leftTrimmed := strings.TrimLeft(line, " \t")
				spaces := len(line) - len(leftTrimmed)
				fmt.Printf("%s#%s\n", line[0:spaces], leftTrimmed)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultsCmd)
}

