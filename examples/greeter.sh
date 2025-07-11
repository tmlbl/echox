# GET /greet/:name
say_hello() {
  echo "Hello, $name!"
}

# GET /containers
list_containers() {
  containers=$(docker ps -q)
  echo $containers
}

# GET /cowsay/:moo
function run_cowsay() {
  cowsay $moo
}
