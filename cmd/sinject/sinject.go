package main

import "flag"
import "fmt"

func main() {

	//var path = *flag.String("path", "/run/secrets", "The path to the directory containing the secrets.  Default is /run/secrets")
	//var secrets = *flag.String("secrets", "", "Explicit list of filtered secrets to use for injection.  Array should be separated by commas.  Default is none; resulting in an unfiltered secret list.")
	//var token = *flag.String("token", "%%%", "Sets the token wrapper.  Default is %%%")
	//var errorMode = *flag.String("errorMode", "secret", "Sets the error mode for token pre-scanning.")

	version := flag.Bool("version", false, "Prints version information and quits.")

	flag.Parse()

	if *version {
		fmt.Println("Version Information")
	}
}
