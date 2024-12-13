air -c ./.air.toml &
air_pid=$!
npx tailwindcss \
  -i 'tailwind/base.css' \
  -o 'assets/css/main.css' \
  --watch &
tailwind_pid=$!
read -r -d '' _ </dev/tty
kill "$air_pid"
kill "$tailwind_pid"