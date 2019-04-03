package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/scott-haines/sinject/version"
)

var cfgFile string
var log = logrus.New()
var verbosity string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                   "sinject [OPTIONS] COMMAND [ARG...]",
	Short:                 "secrets injection utility for containers",
	Long:                  ``,
	Version:               fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
	DisableFlagsInUseLine: false,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetString("verbosity")
		switch {
		case verbosity == "TRACE":
			log.SetLevel(logrus.TraceLevel)
		case verbosity == "DEBUG":
			log.SetLevel(logrus.DebugLevel)
		case verbosity == "INFO":
			log.SetLevel(logrus.InfoLevel)
		case verbosity == "WARN":
			log.SetLevel(logrus.WarnLevel)
		case verbosity == "ERROR":
			log.SetLevel(logrus.ErrorLevel)
		case verbosity == "FATAL":
			log.SetLevel(logrus.FatalLevel)
		case verbosity == "PANIC":
			log.SetLevel(logrus.PanicLevel)
		default:
			log.SetLevel(logrus.InfoLevel)
		}

		log.WithField("verbosity", verbosity).Trace("Log Level Initialized.")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Trace("Execution of default command.")
		log.Trace("Printing help as default command is NOOP.")
		cmd.Help()
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("An error occurred when executing the root command.")
	}
}

func init() {
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

	// global flags
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "", "INFO", "Verbosity level of messages (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)")

	// cli-only flags
	rootCmd.Flags().Bool("version", false, "Print version information and quit")

	// commands
	rootCmd.AddCommand(injectCmd)

	// logging
	log.Out = os.Stderr
}
