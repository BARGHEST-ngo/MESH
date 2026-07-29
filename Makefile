.PHONY: analyst

all: dev-init dev-ui apikey

dev-init: dev-config dev-certs

dev-config:
	cp control-plane/frps.toml.dev control-plane/frps.toml
	cp control-plane/frpc.toml.dev control-plane/frpc.toml
	cp control-plane/headscale/config.yaml.dev control-plane/headscale/config.yaml
	cp .env.dev .env

dev-certs:
	mkdir certs
	mkcert -install
	mkcert -cert-file certs/server.crt -key-file certs/server.key  localhost frps

dev: dev-ui analyst

dev-ui:
	docker compose up -d frps frpc ui

apikey:
	docker compose exec headscale headscale apikeys create --expiration 3h

analyst:
	docker compose up -d analyst
	docker compose exec --user mesh analyst /bin/bash

down:
	docker compose down

clean:
	docker compose down -v
	rm -f control-plane/frps.toml
	rm -f control-plane/frpc.toml
	rm -f control-plane/headscale/config.yaml
	rm -f .env
	rm -rf certs
