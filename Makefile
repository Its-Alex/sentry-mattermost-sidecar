build-bin:
	go build -o ./build/sentry-mattermost-sidecar ./cmd/sms/main.go

build-docker:
	docker build -t sentry-mattermost-sidecar:local .

build-clean:
	rm -rf ./build

attach:
	docker run --rm -it -v $(shell pwd):/app/ --name sentry-mattermost-sidecar --entrypoint 'sh' sentry-mattermost-sidecar:local

test:
	go test -v ./...

lint:
	golangci-lint run --out-format=github-actions

mod-update:
	go get -u && go mod tidy
