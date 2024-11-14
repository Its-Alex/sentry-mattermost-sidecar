# sentry-mattermost-sidecar

This tools is a sidecar to use sentry webhook on mattermost.

## Roadmap

- [x] Handle sentry Issue alerts with [webhook legacy integration](https://github.com/getsentry/sentry-webhooks)
- [ ] Handle sentry Metric alerts with [webhook legacy integration](https://github.com/getsentry/sentry-webhooks)

## How to use

First you must create a [Mattermost incoming webhook](https://docs.mattermost.com/developer/webhooks-incoming.html) integration.
![mattermost-incoming-webhook-integration-setup](docs/assets/mattermost-incoming-webhook-integration-setup.png)

Next you must deploy the [docker image](https://hub.docker.com/r/itsalex/sentry-mattermost-sidecar) (don't forget to fill `SMS_MATTERMOST_WEBHOOK_URL` environment variable with the Mattermost webhook URL) somewhere and redirect sentry webhook on it with route name defined as Mattermost channel for each projects.
![sentry-webhook-integration-setup](docs/assets/sentry-webhook-integration-setup.png)

Then you setup [sentry issue alerts](https://docs.sentry.io/product/alerts/) as you like.
![sentry-issue-alert-creation](docs/assets/sentry-issue-alert-creation.png)

## Getting started

### Requirement

- `docker`
- `bash`
- `virtualbox` (if you want to setup local mattermost and sentry instance)
- `vagrant` (if you want to setup local mattermost and sentry instance)

## Hack

To start you must launch dev environment:

```sh
$ ./scripts/up.sh
```

This will launch images in [`docker-compose.yml`](./docker-compose.yml).

An image named `workspace` with golang is used as a isolated container to
develop. You can use [`enter-workspace.sh`](./scripts/enter-workspace.sh)
to enter inside it:

```sh
$ ./scripts/enter-workspace.sh
```

From outside the container you can build with:

```sh
$ ./scripts/build.sh
```

You can test an example sentry webhook with:

```sh
$ ./scripts/test-request.sh
```

Then you can see the converted request that will be send to mattermost using:

```sh
$ ./scripts/get-last-request-result.sh
```

## Deploy

This image is automatically deployed and versionned as a docker image at [itsalex/sentry-mattermost-sidecar](https://hub.docker.com/r/itsalex/sentry-mattermost-sidecar).

To deploy a new tag use [`./scripts/create-and-push-tag.sh`](scripts/create-and-push-tag.sh):

```sh
$ ./scripts/create-and-push-tag.sh 1.0.0
```

## Setup VM with Mattermost and Sentry

You can setup a VM with Mattermost and Sentry if you want to perform real tests.
You should have at least:

- `16` CPU thread
- `20GB` RAM

If you valid those requiremnts, you can launch the VM:

```sh
$ vagrant up
```

A server will be launch with Mattermost and Sentry installed, you should now
create Sentry first user:

```sh
$ vagrant ssh -c "cd /opt/sentry && sudo docker compose run web upgrade"
...
Running hooks in /etc/ca-certificates/update.d...
done.
Running migrations for default
Operations to perform:
  Apply all migrations: auth, contenttypes, feedback, hybridcloud, nodestore, replays, sentry, sessions, sites, social_auth
Running migrations:
  No migrations to apply.
Creating missing DSNs
Correcting Group.num_comments counter
17:26:28 [INFO] sentry.outboxes: Executing outbox replication backfill
17:26:28 [INFO] sentry.outboxes: Processing sentry.ControlOutboxs...
17:26:28 [INFO] sentry.outboxes: Processing sentry.RegionOutboxs...
17:26:28 [INFO] sentry.outboxes: done

Would you like to create a user account now? [Y/n]:

```

Follow the instruction to create Sentry default user. Mattermost default user
will be asked on the first connection on Mattermost url.

You're now ready, you can new access services with the following URLs:

- Sentry: http://192.168.56.4:9000/
- Mattermost: http://192.168.56.4:8065/

The VM have a static IP so you can always access it with IP `192.168.56.4`.
You can find the IP of your computer accessible from the VM using:

```sh
$ ip a | grep 192.168.56 | awk '{print $2}'
192.168.56.1/24
```