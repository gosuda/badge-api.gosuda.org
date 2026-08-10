# Environment Variables

## Address selection order

The Go server chooses its listener in this order:

1. A non-empty `BADGE_API_HTTP_ADDR` value.
2. The `PORT`, `HOST`, and `IP` PaaS variables.
3. The default address `:8080`.

## Variable reference

| Variable | Used by | Purpose |
| --- | --- | --- |
| `GATEWAY_INTERFACE` | Go CGI mode | Standard CGI variable supplied by Apache. Any non-empty value selects CGI mode before listener setup. |
| `PORT` | Go and `run.js` | Direct Go mode uses it as the listener port. Proxy mode uses it as the public Node.js port. |
| `HOST` | Go | Hostname used with `PORT`. Ignored when `PORT` is not set. |
| `IP` | Go | Valid IPv4 or IPv6 address that overrides `HOST`. Ignored when `PORT` is not set. |
| `BADGE_API_HTTP_ADDR` | Go | Explicit listener target. Supports Unix sockets and TCP addresses. |
| `BADGE_API_EXIT_ON_STDIN_EOF` | Go | Shuts down when the parent process closes stdin. |
| `BADGE_API_BINARY_PATH` | `run.js` | Overrides the Go executable path used by the Node.js proxy. |

## `GATEWAY_INTERFACE`

Apache sets this variable automatically when it starts a CGI executable, normally to `CGI/1.1`:

```text
GATEWAY_INTERFACE=CGI/1.1
```

When it is present, Tiny Badge serves one request with `net/http/cgi` and exits. It does not inspect listener variables, monitor stdin for parent shutdown, create a socket, or start a long-running HTTP server. Do not set this variable for direct server or Node.js proxy mode.

## `PORT`

Use a platform-provided port:

```sh
PORT=8080 ./badge-api.gosuda.org
```

The resulting address is:

```text
:8080
```

In proxy mode, `PORT` belongs to the public Node.js server:

```sh
PORT=3000 npm start
```

`run.js` supplies an explicit private Unix socket address to the Go process, so the inherited public `PORT` does not cause a second TCP listener.

## `HOST`

Combine a hostname with `PORT`:

```sh
HOST=127.0.0.1 PORT=8080 ./badge-api.gosuda.org
```

Result:

```text
127.0.0.1:8080
```

`HOST` alone does not replace the default address. `PORT` must also be set.

## `IP`

A valid `IP` value overrides `HOST`.

IPv4 example:

```sh
HOST=localhost IP=127.0.0.1 PORT=8080 ./badge-api.gosuda.org
```

Result:

```text
127.0.0.1:8080
```

IPv6 example:

```sh
IP=::1 PORT=8080 ./badge-api.gosuda.org
```

Result:

```text
[::1]:8080
```

Invalid IP values are ignored, leaving `HOST` unchanged.

## `BADGE_API_HTTP_ADDR`

This variable overrides `PORT`, `HOST`, and `IP` when it is non-empty.

Unix socket:

```sh
BADGE_API_HTTP_ADDR=unix:/tmp/tiny-badge.sock ./badge-api.gosuda.org
```

Explicit TCP address:

```sh
BADGE_API_HTTP_ADDR=tcp:127.0.0.1:8080 ./badge-api.gosuda.org
```

A plain TCP address is also accepted:

```sh
BADGE_API_HTTP_ADDR=127.0.0.1:8080 ./badge-api.gosuda.org
```

An empty `unix:` or `tcp:` address stops startup with an error.

Unix sockets are created with mode `0600`. Existing files at the configured socket path are removed before startup and the socket is removed again during a normal shutdown.

## `BADGE_API_EXIT_ON_STDIN_EOF`

Enable parent-pipe shutdown:

```sh
BADGE_API_EXIT_ON_STDIN_EOF=true ./badge-api.gosuda.org
```

Accepted true values follow Go boolean parsing:

```text
1, t, T, TRUE, true, True
```

Accepted false values:

```text
0, f, F, FALSE, false, False
```

Any other value stops startup with an error.

The Node.js proxy sets this variable to `true`. If the parent process disappears and closes the pipe, the Go server begins graceful shutdown.

## `BADGE_API_BINARY_PATH`

Choose the executable launched by `run.js`:

```sh
BADGE_API_BINARY_PATH=/home/example/bin/tiny-badge npm start
```

Without this variable, the default path is:

```text
./badge-api.gosuda.org
```

## Shared-hosting example

```sh
npm ci
npm run build
PORT=3000 \
BADGE_API_BINARY_PATH="$PWD/badge-api.gosuda.org" \
npm start
```

`run.js` adds these values for the child process automatically:

```text
BADGE_API_HTTP_ADDR=unix:<temporary socket path>
BADGE_API_EXIT_ON_STDIN_EOF=true
```
