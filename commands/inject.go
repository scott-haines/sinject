package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	injectCmd = &cobra.Command{
		Use:   "inject [OPTIONS]",
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
	injectCmd.Flags().StringVar(&file, "file", "", "The file to inject secrets into (required)")
	injectCmd.MarkFlagRequired("file")
	injectCmd.Flags().StringVar(&prescanmode, "pre-scan-mode", "secret", "Sets the pre-scan mode.\n\nnone - sinject will not perform pre-scanning.\nsecret - sinject will output errors if there is no secret for a discovered token.\ntoken - sinject will output errors if a secret is present but there is no token.\nfull - both secret + token\n\n")
	injectCmd.Flags().StringVar(&secretspath, "secrets-path", "/run/secrets", "The path to the directory containing the secrets.")
	injectCmd.Flags().StringVar(&token, "token", "%%%", "Sets the token wrapper.")
}

func getTokensInFile(file string, token string) []string {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(fmt.Sprintf("%s(.*)%s", token, token))
	submatchall := re.FindAllString(string(fileContents), -1)

	for _, element := range submatchall {
		element = strings.Trim(element, token)
		log.WithFields(logrus.Fields{"token": element, "file": file}).Debug("Token discovered.")
	}

	return submatchall
}

func getSecrets(secretspath string) []secretData {
	files, err := ioutil.ReadDir(secretspath)
	if err != nil {
		log.Fatal(err)
	}

	secrets := []secretData{}

	for _, f := range files {
		fileName := path.Join(secretspath, f.Name())
		fileContents, _ := ioutil.ReadFile(fileName)
		d := secretData{
			SecretName:  strings.Trim(f.Name(), ".secret"),
			SecretValue: string(fileContents),
		}

		log.WithFields(logrus.Fields{"secretName": d.SecretName, "file": fileName}).Debug("Secret discovered.")
		secrets = append(secrets, d)
	}

	return secrets
}

func prescanSecret(tokens []string, secrets []secretData) bool {
	// Output errors if there is no secret for a discovered token
	log.Trace("Executing prescanSecret")
	anyMissing := false
	for _, fulltoken := range tokens {
		found := false

		for _, secret := range secrets {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				log.WithFields(logrus.Fields{"secret": secret.SecretName, "token": fulltoken}).Debug("token -> secret mapped.")
				found = true
			}
		}

		if found == false {
			log.WithField("token", fulltoken).Error("Missing secret for discovered token")
			anyMissing = true
		}
	}

	return anyMissing
}

func prescanToken(tokens []string, secrets []secretData) bool {
	// Output errors if there is no token for a discovered secret
	log.Trace("Executing prescanToken")
	anyMissing := false
	for _, secret := range secrets {
		found := false

		for _, fulltoken := range tokens {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				log.WithFields(logrus.Fields{"secret": secret.SecretName, "token": fulltoken}).Debug("secret -> token mapped.")
				found = true
			}
		}

		if found == false {
			log.WithField("secret", secret.SecretName).Error("Missing token for discovered secret")
			anyMissing = true
		}
	}

	return anyMissing
}

func replaceTokens(tokens []string, secrets []secretData) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithError(err).Fatal("An error occurred while attempting to read the file.")
	}
	buffer := input

	for _, fulltoken := range tokens {
		// Find the matching secret and replace it in the buffer
		for _, secret := range secrets {
			if strings.Trim(fulltoken, token) == secret.SecretName {
				log.WithFields(logrus.Fields{"secretName": secret.SecretName, "file": file}).Info("Injecting secret in file.")
				buffer = bytes.Replace(buffer, []byte(fulltoken), []byte(secret.SecretValue), -1)
			}
		}
	}

	if err = ioutil.WriteFile(file, buffer, 0600); err != nil {
		log.WithError(err).Fatal("An error occurred while attempting to write to the file.")
	}
}

func inject(ccmd *cobra.Command, args []string) {
	switch {
	case file == "":
		log.Error("Missing file.")
		ccmd.Help()
		os.Exit(1)
	}

	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			log.WithField("file", file).Error("The provided file does not exist.")
			os.Exit(1)
		} else {
			log.WithError(err).Fatal("An error ocurred while attempting to check if the file exists.")
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
		log.Fatal("Failed Prescan.  Please refer to warning messages.")
	} else {
		replaceTokens(tokens, secrets)
	}
}
