echox
=====

echox is a framework for writing web applications in shell languages. It manages
a number of shell processes and executes functions in them in response to web
requests. You can use it to quickly create web services that directly execute
commands on the machine in response to user requests.

## Is this safe to do?

No, it is not.

## Should I use this for my next "microservice"?

If you dare...

## Writing services

You can write your web service logic in a bash script like so:

```bash
# A comment like the one below must appear above the function definition
# GET /cowsay/:moo
function run_cowsay() {
  cowsay $moo
}
```

Then pass it into `echox` to start the service:

```bash
echox hello.sh
```

```
> curl localhost:7171/cowsay/echox%20is%20cool
 _______________
< echox is cool >
 ---------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

## Request Handling

Your function will be invoked every time an HTTP request matches the supplied
path. Shell processes are locked when executing, so you only have to worry
about dealing with one request at a time.

### Route Parameters

Route parameters are expanded into shell variables. Using them is trivial:

```bash
# GET /greet/:name
greet_by_name() {
    echo "Hello, $name!"
}
```

### Posting data

You can also access data from the request body of POST and PUT requests in the
`body` variable.

```bash
# POST /files/:name
post_file() {
  echo $body > $name
}

# GET /files/:name
get_file() {
  cat $name
}
```

```
> curl -XPOST -d "content" localhost:7171/files/foo
> curl localhost:7171/files/foo
content
```

### Headers

Request headers are supplied in the `headers` variable. Each header is separated
by a newline, and header names and values are separated by a `:`. This allows
them to be interpreted with coreutils.

```bash
get_user_agent() {
    for h in $headers; do
        name=$(echo $h | cut -d ':' -f 1)
        if [ "$name" = "User-Agent" ]; then
            agent=$(echo $h | cut -d ':' -f 2)
            echo "Your user agent is $agent"
        fi
    done
}
```
