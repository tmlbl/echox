echox
=====

echox is a framework for writing web applications in shell languages. It manages
a number of shell processes and executes functions in them in response to web
requests.

## Basic Configuration

Configuration can be supplied from a file or to `stdin`. The configuration
language is simple and only uses a few commands:

* `include [file]` - Loads a source file into all shell processes
* `[method] [path] [func]` - Requests for the given method and path invoke the
given function

For example, this server will return the current date:

```bash
echo "get / date" | echox
```

Or, you can load a function from another file

```bash
# hello.sh
say_hello() {
    now=$(date)
    echo "Hello! It is $now"
}
```

```bash
echo "include hello.sh; get / say_hello" | echox
```

Semicolons or newlines can be used as separators. An equivalent config file
could look like this:

```
# example.txt
include hello.sh

get / say_hello
```

```bash
echox example.txt
```

## Request Handling

Your function will be invoked every time an HTTP request matches the supplied
path. Shell processes are locked when executing, so you only have to worry
about dealing with one request at a time.

### Route Parameters

Route parameters are expanded into shell variables. Using them is trivial:

```bash
# greet.sh
greet_by_name() {
    echo "Hello, $name!"
}
```

```bash
echo "include greet.sh; get /greet/:name greet_by_name" | echox
```

### Headers

Request headers are supplied in the `headers` variable. Each header is separated
by a newline, and header names and values are separated by a `:`. This allows
them to be interpreted with coreutils.

```bash
# headers.sh
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
