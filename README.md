# prbot

prbot is run as a cron on pull requests in a specified Github repository.
prbot can run arbitrary commands, specified by the user in a config, and will comment the output of commands on valid pull requests.

## Sample Config
A sample config for the PR bot looks like this:

```yaml
github:
    owner: kubernetes
    repo: minikube
    botName: minikube-pr-bot
    accessTokenEnvVar: GITHUB_ACCESS_TOKEN
    pullRequest:
      labels:
          - ok-to-test
execute:
    setup:
        - name: update minikube repo
          command: git pull origin master
          dir: $HOME/minikube
        - name: build mkcmp/minikube
          command: make out/mkcmp out/minikube
          dir: $HOME/minikube
    track:
        - name: run mkcmp
          command: mkcmp ./minikube pr://{{.PRNumber}}
          dir: $HOME/minikube/out
    cleanup:
        - name: delete minikube
          command: out/minikube delete --all
          dir: $HOME/minikube

```
In the `github` section of the config, we specify the following:
- *required* The organization and repo that we want the PR bot to comment on
- *required* The name of the account the PR bot should comment as
- *required* An environment variable that will provide the Github access token required for the bot to comment under the specified account
- The `pullRequest` field allows for only commenting on PRs with certain labels


To setup, the prbot executes the following commands:

1. Runs `git pull origin master` to ensure we are testing @ HEAD
1. Runs `make out/minikube out/mkcmp`, to build the two binaries we will need to compare

The PR Bot will take the stdout of the `track` section in the config, and comment that on every PR.

To cleanup, we specificy that all minikube clusters should be deleted.

## Building PR Bot
To build the PR bot binary, run:

```
$ make
```

To run the binary:

```
out/prbot --config <path to prbot.yaml>
```

## Examples
To see the configuration for minikube, reference the `minikube/examples` directory.
