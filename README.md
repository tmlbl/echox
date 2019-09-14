echox
=====

echox is a framework for writing web applications in shell languages. It manages
a number of shell processes and executes functions in them in response to web
requests.

## Usage

Configuration can be supplied from a file or to `stdin`. The configuration
language is simple and only uses a few commands:

* `include [file]` - Loads a source file into all shell processes
* `handle [path] [func]` - Responds to the given path with the given function

For example, this server will return the current date:

```bash
echo "handle / date" | echox
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
echo "include hello.sh; handle / say_hello" | echox
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

More features to come!
