package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "Inject secrets into target file",

		Run: inject,
	}
	file string
)

func init() {
	injectCmd.Flags().StringVarP(&file, "file", "f", "", "The file to inject secrets into")
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
