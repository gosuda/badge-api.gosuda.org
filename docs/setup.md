# Environment Setup

This guide covers direct Go launch, the shared-hosting Node.js proxy mode, and Apache CGI deployment.

## Requirements

- Go 1.26
- Node.js 24 with npm
- A Unix-like host for Unix socket proxy mode

## Install dependencies

```sh
npm ci
```

## Run the application directly

Build the SvelteKit frontend first because the Go server embeds the generated files from `dist/frontend/`.

```sh
npm run build:web
go run .
```

The default address is `http://localhost:8080`.

Check that the process is ready:

```sh
curl http://localhost:8080/healthz
```

Expected response:

```text
ok
```

## Build a production executable

```sh
npm run build
```

This command builds the frontend and creates the local executable:

```text
./badge-api.gosuda.org
```

Run it directly:

```sh
./badge-api.gosuda.org
```

## Run behind the Node.js proxy

Shared hosting providers often expose a public `PORT` to a Node.js process. The proxy runner listens on that port and connects to the Go process through a private Unix socket.

```sh
npm run build
PORT=3000 npm start
```

The process layout is:

```text
Browser → Node.js public port → private Unix socket → Go server
```

`run.js` automatically:

1. Creates a unique socket path in the operating system temporary directory.
2. Starts the Go executable with Unix socket mode enabled.
3. Proxies incoming HTTP requests to the socket.
4. Sends `SIGTERM` to the Go process during shutdown.
5. Removes the socket when the process exits.

## Use an executable outside the project directory

```sh
BADGE_API_BINARY_PATH=/home/example/bin/badge-api \
PORT=3000 \
npm start
```

The configured file must exist and be executable.

## Run through Apache CGI

Use the repository's CGI deployment package when Apache must serve every route through one executable:

```sh
npm run build:cgi
```

Upload `deployment/.htaccess` and `deployment/cgi-bin/badge-api.gosuda.org.cgi` with the same directory layout, then make the executable runnable. See [CGI Deployment](cgi.md) for the rewrite rules, cross-compilation command, and required Apache permissions.

## Validate a deployment build

```sh
npm run check
npm run test:color
go test ./...
npm run build
node --check run.js
```

## Common startup problems

### The frontend is missing

Rebuild the embedded frontend:

```sh
npm run build:web
```

### The public port is already in use

Choose another port:

```sh
PORT=3100 npm start
```

### The Go executable cannot be started

Build it again or point the proxy at the correct path:

```sh
npm run build:server
BADGE_API_BINARY_PATH="$PWD/badge-api.gosuda.org" npm start
```

### The proxy briefly returns 503

The public server can start a moment before the Go process finishes opening its Unix socket. Retry the request after the process startup log appears.
