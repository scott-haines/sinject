# Install
## Linux
Install from source: `make install` - Requires Golang and will build and install into `$GOPATH/bin`

or

Download precompiled binary from https://github.com/scott-haines/sinject/releases

## MacOS (Darwin)
Download precompiled binary from https://github.com/scott-haines/sinject/releases

# Sinject
Sinject is a portable, easy to use secret & token injector.  It's designed to be easy to use in a Docker environment, but can easily be adapted to be used for anything.

Sinject supports several ways of mounting secrets into a container:
* Bind Mounting (as a file)
* Docker Secrets
* Kubernetes Secrets

Sinject assumes that the secret -> token replacement is by name.  So the name of the secret must match exactly the label of the token in the config file.

Consider the following config file:
```yaml
myconfig:
    auth:
      username: '%%%USER_NAME%%%'
      password: '%%%PASSWORD%%%'
```

sinject will look for files named USER_NAME or USER_NAME.secret and PASSWORD or PASSWORD.secret and replace the tokens with the values of those secrets.

By default, it will look in `/run/secrets` but this can be overridden with the flag `--secrets-path`
By default, the token is assumed to be `%%%` but this be overridden with the flag `--token`

# Pre Scan Mode
sinject supports 4 modes of file pre scanning with the flag `--pre-scan-mode`:
* none - sinject will not perform pre-scanning.
* secret - sinject will output errors if there is no secret for a discovered token.
* token - sinject will output errors if a secret is present but there is no token.
* full - both secret + token.

The default pre scan mode is 'secret'

# Help
```bash
$ sinject --help

$ sinject [COMMAND] --help
```

# Examples
```bash
$ sinject inject --file /opt/app1/myConfig.yml

$ sinject inject --file exampleConfig.txt --secrets-path $(pwd)/secrets --pre-scan-mode none
```
