.PHONY: default
default: live ;

live/db:
	cd db && docker compose up --remove-orphans

live/templ:
	templ generate --watch --proxy="http://localhost:4040" --open-browser=false

live/server:
	air

live/sync_assets:
	arelo -t './assets' -p '**/*.css' -p '**/*.js' -- templ generate --notify-proxy

live/sync_sqlc:
	arelo -p '**/*.sql' -p 'sqlc.yml' -- sqlc generate

live:
	make -j4 live/templ live/server live/sync_assets live/sync_sqlc

sync:
	go run cmd/server/main.go -update -days 60
