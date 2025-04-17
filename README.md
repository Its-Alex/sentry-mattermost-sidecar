# sentry-mattermost-sidecar

This tools is a sidecar to use sentry webhook on mattermost.

## Roadmap

- [x] Handle sentry Issue alerts with [webhook legacy integration](https://github.com/getsentry/sentry-webhooks)
  (must be setup with legacy webhook integration)
- [ ] Handle sentry Metric alerts with [webhook legacy integration](https://github.com/getsentry/sentry-webhooks)
  (must be setup with legacy webhook integration)
- Handle [Custom integration - Sentry webhook](https://docs.sentry.io/organization/integrations/integration-platform/webhooks/)
  - Issues alerts
    - [x] `triggered` action
  - Issues
    - [x] `created` action
    - [x] `resolved` action
    - [x] `assigned` action
    - [x] `archived` action
    - [x] `unresolved` action
  - Comments
    - [ ] `created` action
    - [ ] `updated` action
    - [ ] `deleted` action
  - Errors 
    - [x] `created` action

## How to use

First you must create a [Mattermost incoming webhook](https://developers.mattermost.com/integrate/webhooks/incoming/)
integration:  
![mattermost-incoming-webhook-integration-setup](docs/assets/mattermost-incoming-webhook-integration-setup.png)

Next you must deploy the [docker image](https://hub.docker.com/r/itsalex/sentry-mattermost-sidecar)
(don't forget to fill `SMS_MATTERMOST_WEBHOOK_URL` environment variable with the
Mattermost webhook URL) somewhere and redirect sentry webhook on it with route
name defined as Mattermost channel for each projects, for example with
docker-compose:

```yaml
services:
  sentry-mattermost-sidecar:
    image: itsalex/sentry-mattermost-sidecar:latest
    restart: unless-stopped
    environment:
      - SMS_MATTERMOST_WEBHOOK_URL=https://mattermost.example.com/hooks/abckus71ojr9idqarirt4mr8wa
    ports:
      - 1323:1323
```

Finally you must setup sentry custom integration and alert:  
![sentry-webhook-integration-setup](docs/assets/sentry-integration-and-alert-creation.gif)

## Getting started

### Requirement

- [`mise`](https://mise.jdx.dev/) (if you want to send real errors to sentry)
- `docker`
- `bash`
- `virtualbox` (if you want to setup local mattermost and sentry instance)
- `vagrant` (if you want to setup local mattermost and sentry instance)

If you want to send real errors to sentry, and you have installed
[`mise`](https://mise.jdx.dev/), you must execute the following commands to
have everything working:

```sh
$ mise trust && mise install
$ pip install -r requirements.txt
```

This will install python with a predefined version of sentry-sdk python package
in an isolated [.venv](https://docs.python.org/3/library/venv.html) folder.

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

This is an environment aimed to reproduced real use case. If you really want to
perform tests with Mattermost and Sentry, you can do it locally following
[Setup VM with Mattermost and Sentry](#setup-vm-with-mattermost-and-sentry).

## Deploy

This image is automatically deployed and versionned as a docker image at [itsalex/sentry-mattermost-sidecar](https://hub.docker.com/r/itsalex/sentry-mattermost-sidecar).

To deploy a new tag use [`./scripts/create-and-push-tag.sh`](scripts/create-and-push-tag.sh):

```sh
$ ./scripts/create-and-push-tag.sh 1.0.0
```

## Setup VM with Mattermost and Sentry

### Setup

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

Please create a Mattermost user by going at http://192.168.56.4:8065/ and follow
instructions.

You're now ready, you can new access services with the following URLs:

- Sentry: http://192.168.56.4:9000/
- Mattermost: http://192.168.56.4:8065/

The VM have a static IP so you can always access it with IP `192.168.56.4`.
You can find the IP of your computer accessible (most likely `192.168.56.1`)
from the VM using:

```sh
$ ip a | grep 192.168.56 | awk '{print $2}'
192.168.56.1/24
```

You can now [configure webhooks](#configure-webhooks)

### Configure webhooks

This step is planned to be automatised, but for now we must do it manually. It
aim to create [Mattermost incoming webhook](https://developers.mattermost.com/integrate/webhooks/incoming/)
and [Sentry Webhook](https://docs.sentry.io/organization/integrations/integration-platform/webhooks/).

You can configure
[Mattermost incoming webhook](https://developers.mattermost.com/integrate/webhooks/incoming/)
on any channel you want. To use it in development, replace
[the content of the variable in the following line `SMS_MATTERMOST_WEBHOOK_URL=http://requests-catcher:5000`](/docker-compose.yml#L20)
by your mattermost webhook.

Restart (or start) the container by using:

```sh
$ docker compose up -d
```

Finally you must configure [Sentry Webhook](https://docs.sentry.io/organization/integrations/integration-platform/webhooks/).
The URL of the webhook should be `http://<your-ip>:1323/<channel-name>`, for
example, if your IP is `192.168.56.1` and the channel where you want to
publish is `test` your URL will be: `http://192.168.56.1:1323/test`.

You can test if everything is working by following
[Send errors to Sentry](#send-errors-to-sentry) or using your way.

## Send errors to Sentry

The repository contains a python script that can be used to push errors to sentry.
Make sure you've followed [the requirements](#requirement) before continue.

You should update the Sentry DSN in
[scripts/sentry-trigger-error.py](/scripts/sentry-trigger-error.py#L4) by
the DSN of your project.

You can now trigger an error by using:

```sh
$ python scripts/sentry-trigger-error.py
Traceback (most recent call last):
  File "/home/alex/Documents/sentry-mattermost-sidecar/scripts/sentry-trigger-error.py", line 10, in <module>
    division_by_zero = 1 / 0
                       ~~^~~
ZeroDivisionError: division by zero
```

An error should be generated on Sentry.
