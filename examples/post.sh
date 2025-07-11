# POST /files/:name
post_file() {
  echo $body >$name
}

# GET /files/:name
get_file() {
  cat $name
}
