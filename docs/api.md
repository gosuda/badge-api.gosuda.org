# Badge API

The service returns SVG badges from either a query-based URL or a path-based URL.

Examples below use this placeholder origin:

```text
https://badge.example.com
```

## Query endpoint

```text
GET /badge.svg
HEAD /badge.svg
```

A message is required:

```text
https://badge.example.com/badge.svg?message=ready
```

Add a label:

```text
https://badge.example.com/badge.svg?label=build&message=passing
```

## Path endpoint

```text
GET /badge/:label/:message
HEAD /badge/:label/:message
```

Example:

```text
https://badge.example.com/badge/release/stable.svg
```

Use `_` for an empty path label:

```text
https://badge.example.com/badge/_/available.svg
```

The path endpoint removes a final `.svg` suffix from the message. The query endpoint preserves it:

```text
/badge/release/version.svg        → message is "version"
/badge.svg?message=version.svg    → message is "version.svg"
```

## Query parameters

| Parameter | Required | Default | Description |
| --- | --- | --- | --- |
| `message` | Yes | — | Main badge text. Maximum 128 characters. |
| `label` | No | Empty | Left-side badge text. Maximum 64 characters. |
| `style` | No | `flat` | Badge rendering style. |
| `labelColor` | No | `555555` | Label background color. |
| `color` | No | `44cc11` | Message background color. |
| `labelTextColor` | No | `ffffff` | Label text color. |
| `textColor` | No | `ffffff` | Message text color. |

Spaces and punctuation must be URL encoded. For example:

```text
https://badge.example.com/badge.svg?label=release&message=ready%20to%20ship
```

## Styles

| Style | Notes |
| --- | --- |
| `flat` | Rounded 20 px badge and the default style. |
| `flat-square` | Flat badge with square corners. |
| `plastic` | Compact badge with a highlight layer. |
| `round` | Taller capsule shape. |
| `outline` | Framed badge surface. |
| `neon` | Badge with a focused glow. |
| `glass` | Layered highlight and border. |
| `flatbar` | 28 px square badge with uppercase text. |

Style example:

```text
https://badge.example.com/badge.svg?label=mood&message=loud&style=flatbar
```

## Colors

Colors accept 3-digit or 6-digit hexadecimal values. Omit the leading `#`:

```text
color=d6ef53
labelColor=292724
```

If a leading `#` is included, encode it as `%23` because an unescaped `#` starts a URL fragment.

Named colors are also supported:

```text
brightgreen, green, yellowgreen, yellow, orange, red, blue,
grey, gray, lightgrey, lightgray, success, important, critical,
informational, inactive
```

Color example:

```text
https://badge.example.com/badge.svg?label=status&message=ready&labelColor=292724&color=d6ef53&labelTextColor=ffffff&textColor=292724
```

## Markdown example

```md
![Status](https://badge.example.com/badge.svg?label=status&message=ready&style=round)
```

## HTML example

```html
<img
  src="https://badge.example.com/badge.svg?label=status&message=ready&style=round"
  alt="Status: ready"
/>
```

## cURL examples

Download a badge:

```sh
curl -o badge.svg \
  'https://badge.example.com/badge.svg?label=build&message=passing&style=flatbar'
```

Inspect response headers:

```sh
curl -I \
  'https://badge.example.com/badge.svg?label=build&message=passing'
```

## Response headers

Successful badge responses include:

```text
Content-Type: image/svg+xml; charset=utf-8
Cache-Control: public, max-age=315360000, immutable
CDN-Cache-Control: public, max-age=315360000, immutable
Surrogate-Control: public, max-age=315360000, immutable
Access-Control-Allow-Origin: *
ETag: "<content hash>"
```

Use the ETag for a conditional request:

```sh
curl -I \
  -H 'If-None-Match: "<etag>"' \
  'https://badge.example.com/badge.svg?message=ready'
```

An unchanged badge returns:

```text
304 Not Modified
```

## Errors

Invalid input returns `400 Bad Request` with `Cache-Control: no-store`.

Examples of invalid input:

```text
/badge.svg?message=
/badge.svg?message=ready&style=unknown
/badge.svg?message=ready&color=not-a-color
```

The health endpoint is available at:

```text
GET /healthz
```

Successful response:

```text
ok
```
