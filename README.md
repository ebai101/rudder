# rudder

personal finance tracker built with Go, [Templ](https://templ.guide/), [HTMX](https://htmx.org/), and PostgreSQL. pulls transaction data automatically from bank accounts via [SimpleFIN Bridge](https://beta-bridge.simplefin.org/) and presents it in a self-hosted web dashboard.

mostly abandoned but fairly functional.

![screenshot](https://github.com/user-attachments/assets/dee82f01-ed5b-4c79-8a44-7fbf4612a28a)

## features

- automatic transaction sync from bank accounts via SimpleFIN Bridge
- scheduled pulls on configurable intervals (hourly, daily, weekly, setup)
- transaction history stored in PostgreSQL with decimal-accurate accounting
- interactive charts via [go-echarts](https://github.com/go-echarts/go-echarts)
- server-side rendered UI with Templ + HTMX (no frontend build step)
- optional local caching of SimpleFIN responses for offline/debug use

## tech stack

| layer | technology |
|---|---|
| backend | Go 1.23, Echo v4 |
| templating | Templ + HTMX |
| database | PostgreSQL (pgx v5, sqlc) |
| charts | go-echarts |
| scheduling | gocron v2 |
| config | YAML |

## setup

**prerequisites:** Go 1.23+, PostgreSQL, a [SimpleFIN Bridge token](https://beta-bridge.simplefin.org/my-account/tokens/create)

```bash
git clone https://github.com/ebai101/rudder
cd rudder
cp config.example.yml config.yml
# edit config.yml with your db url, simplefin token, and timezone
make run
```

### config

```yaml
timezone: America/New_York
sfin_bridge_token: your_token_here
db_url: postgres://user:pass@localhost:5432/rudder
setup_pull_days: 180   # initial historical pull
weekly_pull_days: 30
daily_pull_days: 7
hourly_pull_days: 3
send_requests: true    # set false to use cached data
save_cache: true
```

## development

uses [air](https://github.com/air-verse/air) for live reloading and [sqlc](https://sqlc.dev/) for type-safe database queries.

```bash
make        # start with live reload via air
make sqlc   # regenerate db query code
```
