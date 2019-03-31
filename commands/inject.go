package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/jcelliott/lumber"

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

type secretData struct {
	SecretName  string
	SecretValue string
}

func init() {
	injectCmd.Flags().StringVar(&file, "file", "", "The file to inject secrets into")
	injectCmd.Flags().StringVar(&prescanmode, "pre-scan-mode", "secret", "Sets the pre-scan mode.\n\nnone - sinject will not perform pre-scanning.\nsecret - sinject will output errors if there is no secret for a discovered token.\ntoken - sinject will output errors if a secret is present but there is no token.\nfull - both secret + token\n\n")
	injectCmd.Flags().StringVar(&secretspath, "secrets-path", "/run/secrets", "The path to the directory containing the secrets.")
	injectCmd.Flags().StringVar(&token, "token", "%%%", "Sets the token wrapper.")
}

func getTokensInFile(file string, token string) []string {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	re := regexp.MustCompile(fmt.Sprintf("%s(.*)%s", token, token))
	submatchall := re.FindAllString(string(fileContents), -1)

	for _, element := range submatchall {
		element = strings.Trim(element, token)
		lumber.Trace(element)
	}

	return submatchall
}

func getSecrets(secretspath string) []secretData {
	files, err := ioutil.ReadDir(secretspath)
	if err != nil {

	}

	secrets := []secretData{}

	for _, f := range files {
		fileName := path.Join(secretspath, f.Name())
		fileContents, _ := ioutil.ReadFile(fileName)
		d := secretData{
			SecretName:  strings.Trim(f.Name(), ".secret"),
			SecretValue: string(fileContents),
		}

		secrets = append(secrets, d)
	}

	return secrets
}

func prescanSecret(tokens []string, secrets []secretData) bool {
	// Output errors if there is no secret for a discovered token
	anyMissing := false
	for _, fulltoken := range tokens {
		found := false

		for _, secret := range secrets {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				found = true
			}
		}

		if found == false {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Missing secret for token: %s", fulltoken))
			anyMissing = true
		}
	}

	return anyMissing
}

func prescanToken(tokens []string, secrets []secretData) bool {
	// Output errors if there is no token for a discovered secret
	anyMissing := false
	for _, secret := range secrets {
		found := false

		for _, fulltoken := range tokens {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				found = true
			}
		}

		if found == false {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Missing token for secret: %s", secret.SecretName))
			anyMissing = true
		}
	}

	return anyMissing
}

func replaceTokens(tokens []string, secrets []secretData) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	buffer := input

	for _, fulltoken := range tokens {
		// Find the matching secret and replace it in the buffer
		for _, secret := range secrets {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				buffer = bytes.Replace(buffer, []byte(fulltoken), []byte(secret.SecretValue), -1)
			}
		}
	}

	if err = ioutil.WriteFile(file, buffer, 0600); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
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

	tokens := getTokensInFile(file, token)
	secrets := getSecrets(secretspath)

	failedPrescan := false
	switch {
	case prescanmode == "secret":
		failedPrescan = prescanSecret(tokens, secrets)
	case prescanmode == "token":
		failedPrescan = prescanToken(tokens, secrets)
	case prescanmode == "full":
		// Avoid short circuiting.
		failedSecretPrescan := prescanSecret(tokens, secrets)
		failedTokenPrescan := prescanToken(tokens, secrets)
		failedPrescan = failedSecretPrescan || failedTokenPrescan
	}

	if failedPrescan {
		fmt.Fprintln(os.Stderr, "Failed Prescan")
		os.Exit(1)
	} else {
		replaceTokens(tokens, secrets)
	}
}
