package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cerra-io/base-goutils/vipersubtree"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is a global version flag, filled in compile time
var Version string

var (
	cfgFile          string // alternate config file requested at command line
	printVersionFlag bool
	logger           = logrus.WithField("module", "root")
)

// Setup logging based in the viper config "logging" sub-section
// in the config file (possibly overridden from command line or
// environment variables)
func setupLogging(conf *vipersubtree.ViperSubtree) {

	// Setup time format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = conf.GetString("timeFormat")
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// Setup default log level
	levelName := conf.GetString("defaultLevel")
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		if strings.Contains(err.Error(), "not a valid logrus Level") {
			logger.Errorf("%s: not a valid log level name (valid names are:"+
				"debug, info, warn, warning, error, fatal, panic)", levelName)
			os.Exit(1)
		}
		panic(err)
	}
	logrus.SetLevel(level)
}

func printBanner() {
	logger.Debug("started:", os.Args[0])
	logger.WithField("file", viper.ConfigFileUsed()).Debug("using-config-file")
	if logrus.GetLevel() == logrus.DebugLevel {
		bytes, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
		if err != nil {
			panic(err)
		}
		logger.Debug("viper-config-starts")
		for _, line := range strings.Split(string(bytes), "\n") {
			logger.Debug(line)
		}
		logger.Debug("viper-config-ends")
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cluster-state",
	Short: "Automation Service for a cerra cluster",
	Long:  `Run help to see all subcommands`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Setup logging using the "logging" section from viper config
		// (which by now are optionally overridden by bound values from command
		// line args and environment variables
		setupLogging(vipersubtree.Subtree(viper.GetViper(), "logging"))

		// Print banner and debug info as needed
		printBanner()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if printVersionFlag {
			fmt.Println(Version)
		} else {
			cmd.HelpFunc()(cmd, args)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Load all defaults, thus creating a full viper tree
	viper.SetConfigType("yml")
	err := viper.ReadConfig(bytes.NewBuffer(Defaults))
	if err != nil {
		panic(err)
	}

	// Call init after main() but before any cobra commands are executed
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	flags := RootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default is .dcm.yml)")
	flags.BoolVarP(&printVersionFlag, "version", "v", false, "Display current version")
	flags.StringP("logging-level", "l", "",
		"Default logging level (one of panic, error, warn, info, or debug)")
	viper.BindPFlag("logging.defaultLevel", flags.Lookup("logging-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// set up viper to handle config
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(".bcs") // name of config file (without extension)
	viper.AddConfigPath(".")        // adding local directory as second
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.SetEnvPrefix("bcs")   // all environment variables start with dcm prefix
	viper.AutomaticEnv()            // read in environment variables that match
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer) // Use underscore (and not dot) in environment variables

	// If a config file is found, read it in
	if err := viper.MergeInConfig(); err != nil {
		logger.Warnf("no-or-bad-config-file: %s", err)
	}
}

