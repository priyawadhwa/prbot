This example contains `prbot.yaml` which the prbot binary reads for minikube.

It also includes the systemd timer and service that is running in GCE so that the prbot executes every 15 minutes.

This is the sample prbot.yaml that the prbot binary reads for the minikube PR bot.

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
