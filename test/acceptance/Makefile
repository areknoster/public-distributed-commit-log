.PHONY: test
test: | compose-local-daemon
	go test -count=1 .

.PHONY: compose-local-daemon
compose-local-daemon: docker-compose.yaml
	mkdir -p data
	docker compose up -d
	sleep 10


.PHONY: cleanup
cleanup:
	docker compose down
	rm -rf data/*
