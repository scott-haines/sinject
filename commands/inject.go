package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	injectCmd = &cobra.Command{
		Use:   "inject [OPTIONS]  FILE [FILE...]",
		Short: "Inject secrets into target file",
		Run:   inject,
	}
	file        string
	prescanmode string
	secretspath string
	token       string
)

func init() {
	injectCmd.Flags().StringVar(&file, "file", "", "The file to inject secrets into")
	injectCmd.Flags().StringVar(&prescanmode, "pre-scan-mode", "secret", "Sets the pre-scan mode.\n\nnone - sinject will not perform pre-scanning.\nsecret - sinject will output errors if there is no secret for a discovered token.\ntoken - sinject will output errors if a secret is present but there is no token.\nfull - both secret + token\n\n")
	injectCmd.Flags().StringVar(&secretspath, "secrets-path", "/run/secrets", "The path to the directory containing the secrets.")
	injectCmd.Flags().StringVar(&token, "token", "%%%", "Sets the token wrapper.")
}

func inject(ccmd *cobra.Command, args []string) {
	switch {
	case file == "":
		fmt.Fprintln(os.Stderr, "Missing file - please provide the file to inject secrets into")
		return
	}

	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("The file %s does not exist", file))
		} else {
			fmt.Fprintln(os.Stderr, "an error has ocurred")
		}
	}
}
