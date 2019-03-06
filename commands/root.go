package commands

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scott-haines/sinject/version"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                   "sinject [OPTIONS] COMMAND [ARG...]",
	Short:                 "secrets injection utility for containers",
	Long:                  ``,
	Version:               fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
	DisableFlagsInUseLine: false,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	lumber.Trace("Initializing the root command.")

	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [OPTIONS] COMMAND [ARG...]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	// set config defaults
	cobra.OnInitialize(initConfig)
	verbosity := "INFO"

	// cli flags
	rootCmd.PersistentFlags().String("verbosity", verbosity, "Verbosity level of messages (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")

	// bind config to cli flags
	viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity"))

	// cli-only flags
	rootCmd.Flags().Bool("version", false, "Print version information and quit")

	// commands
	rootCmd.AddCommand(injectCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sinject" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sinject")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
