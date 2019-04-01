# Problem Domain
Prevent leakage of secret information in your containers while still maintaining an Everything as Code repository.

In a containerized workload it's very common to have config files like this example of a prometheus scrape target:
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['myServer:18245']
    metrics_path: '/prometheus/metrics'
    scheme: https
    tls_config:
      insecure_skip_verify: true
    basic_auth:
      username: 'userName'
      password: 'mYSp3ci4lPa55w0rd'
```
which contain both normal configuration settings as well as sensitive information in the form of a username and password.

One option would be to take the entirety of this config file and make it a secret.  This works - but makes it difficult to share standard configuration information that may not actually need to be kept secret.

Another option (if the application supports it) is to use environment variables:
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['myServer:18245']
    metrics_path: '/prometheus/metrics'
    scheme: https
    tls_config:
      insecure_skip_verify: true
    basic_auth:
      username: '${USER_NAME}'
      password: '${PASSWORD}'
```

Be warned though - Environment variables passed through to containers can be easily discovered:
```bash
docker run --rm -d -e USER_NAME=userName -e PASSWORD=mYSp3ci4lPa55w0rd nginx
docker inspect <containerName>

"Env": [
        "USER_NAME=userName",
        "PASSWORD=mYSp3ci4lPa55w0rd",
    ],
```

# Sinject
Sinject provides any easy way to achieve separation between config files and secrets by enabling token replacement.

Similar to using environment variables, you can mark a portion of your config file with tokens and sinject can replace those tokens with the value of a secret (default in /run/secrets).

This supports several ways of mounting secrets into a container:
* Bind Mounting (as a file)
* Docker Secrets
* Kubernetes Secrets

Sinject assumes that the secret -> token replacement is by name.  So the name of the secret must match exactly the label of the token in the config file.

Consider the following config file:
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['myServer:18245']
    metrics_path: '/prometheus/metrics'
    scheme: https
    tls_config:
      insecure_skip_verify: true
    basic_auth:
      username: '%%%USER_NAME%%%'
      password: '%%%PASSWORD%%%'
```

sinject will look for secrets named USER_NAME or USER_NAME.secret and PASSWORD or PASSWORD.secret and replace the tokens with the values of those secrets.

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
