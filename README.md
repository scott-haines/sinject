# sinject
Secrets Injector and utility for containerized solutions.

NAME
    inject - secrets injection utility

SYNOPSIS
    inject [OPTIONS] [FILE(s)...]

    inject --help

DESCRIPTION
    is a utility for injecting secrets into configuration files.

    inject will accept a list of files in which to inject secrets.  Within the files tokens can be defined which will be replaced by their corresponding secrets.

    an example of a token within a file: %%%database_username%%%
    will be replaced with the secret with the same name (database_username or database_username.secret)

OPTIONS
    --path="/run/secrets"
      The path to the directory containing the secrets.  Default is /run/secrets

    --secrets=""
      Explicit list of filtered secrets to use for injection.  Array should be separated by commas.  Default is none; resulting in an unfiltered secret list.

    --token="%%%"
      Sets the token wrapper.  Default is %%%

    --errormode=none|secret|token|all
      Sets the error mode for token pre-scanning.

      none - inject will never output pre-scanning errors
      secret - output errors if there is no secret for a discovered token
      token - output errors if a secret is present but there is no token
      all - both secret + token

      Default is secret.

    --help
      Print this help page and quit.

    --version
      Display version information and quit.

EXAMPLES
    inject /opt/my_app/config/config.yml

    inject --path="/mnt/secrets" /opt/my_app/config/config.yml

    inject --secrets="database_user,database_password,api_token" /opt/my_app/config/config.yml

    inject /opt/my_app/config/config.yml /opt/my_app/config/options.xml

    ls /opt/my_app/config_files/ | inject
