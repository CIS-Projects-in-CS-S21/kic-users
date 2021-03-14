.DEFAULT_GOAL := run_tests


build:
	go build -o ./bin/server ./cmd/server/server.go

push:
	docker build -t gcr.io/keeping-it-casual/kic-users:dev .
	docker push gcr.io/keeping-it-casual/kic-users:dev