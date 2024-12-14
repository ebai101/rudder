live/templ:
	templ generate --watch --proxy="http://localhost:4040" --open-browser=false

live/server:
	air

live/tailwind:
	cd tailwind && npm run dev

# live/sync_assets:
# 	arelo -t './assets' -p '**/*.css' -p '**/*.js' -- templ generate --notify-proxy

live/sync_sqlc:
	arelo -p '**/*.sql' -p 'sqlc.yml' -- sqlc generate

live:
	make -j4 live/templ live/server live/tailwind live/sync_sqlc
