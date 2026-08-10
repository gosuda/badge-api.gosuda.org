# CGI Deployment

Tiny Badge detects CGI mode from Apache's standard `GATEWAY_INTERFACE` variable. CGI mode is checked before signal handling, stdin EOF monitoring, TCP configuration, Unix socket cleanup, or listener creation.

## Build

Build the embedded frontend and CGI executable:

```sh
npm run build:cgi
```

The generated executable is:

```text
deployment/cgi-bin/badge-api.gosuda.org.cgi
```

Production builds use `-trimpath -ldflags='-s -w'` to remove local source paths, symbol tables, and DWARF debug data from the executable.

For a Linux AMD64 shared host, use the reproducible cross-compilation script:

```sh
npm run build:cgi:linux-amd64
```

Confirm the operating system and architecture required by the hosting provider before uploading.

## Upload layout

Copy the deployment files into the hosting document root:

```text
public_html/
├── .htaccess
└── cgi-bin/
    └── badge-api.gosuda.org.cgi
```

Make the executable runnable:

```sh
chmod 755 public_html/cgi-bin/badge-api.gosuda.org.cgi
```

The supplied `deployment/.htaccess` uses `badge-api.gosuda.org` as the canonical host. It rewrites every public path to the CGI executable. The executable itself is the only rewrite exception, which prevents recursion.

Apache supplies `GATEWAY_INTERFACE` when it invokes the executable as CGI. Do not set `BADGE_API_EXIT_ON_STDIN_EOF` in CGI deployment; CGI uses standard input for each request.

The frontend is embedded in the Go executable. Do not upload `dist/frontend` separately for this deployment mode.

## Required Apache permissions

The hosting configuration must permit these `.htaccess` directives:

- `Options +ExecCGI`
- `AddHandler cgi-script .cgi`
- `AcceptPathInfo On`
- `RewriteEngine On`

The server also needs `mod_rewrite` and either `mod_cgi` or `mod_cgid`. Some shared hosts preconfigure `cgi-bin` and prohibit overriding `Options`; remove only the prohibited activation directives when the provider already enables CGI.
