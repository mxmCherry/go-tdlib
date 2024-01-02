# TDLib

Low-level [TDLib](https://github.com/tdlib/td) [CGo](https://pkg.go.dev/cmd/cgo) wrapper.

Roughly - Go type aliases for TDLib JSON interface, think of it more as of CGo/TDLib integration recipe.

This lib DOES NOT provide any specific TL types, thus should work with pretty much any version of TDLib with `td_create_client_id`, `td_send`, `td_receive`, `td_execute` methods. JSON payload or API methods can be changed/broken by TDLib the usual way (it's versioning is not the most convenient one), so that is left as an exercise for the lib user - good luck with that!

## `CFLAGS`/`LDFLAGS`

Quick reminder:

- `CFLAGS` mostly used to point to the TDLib `include` dir (with `*.h` files)
- `LDFLAGS` mostly used to point to the TDLib `lib` dir (with `*.{so,a}` files)

This lib hardcodes some `CFLAGS`/`LDFLAGS` right in Go code to simplfy local development/debugging. It assumes that TDLib is built into `third-party/td/tdlib`, see [Makefile](./Makefile) recipes.

This lib links against `tdjson_static` for no particular reason (just because I did it and it worked, yay).

This lib has Linux-specific `LDFLAGS` as of now. Darwin etc can be copied later from, for example, [zelenin/go-tdlib](https://github.com/zelenin/go-tdlib) if ever needed.

For external (normal, `go get`-able) lib use, you'll have to either put TDLib `include` and `lib` dirs somewhere in known places (like `/usr/src` or `/usr/lib` etc) or configure `CFLAGS`/`LDFLAGS` on your own.

### Setting `CFLAGS`/`LDFLAGS` In Go Code

You can have the CGo bits in your own Go code:

```go
package yourpkg

// Note: the followith path(s) can be either relative to current source dir or be absolute (full) ones like "/full/path/to/tdlib/{include,dir}"

/*
# cgo linux CFLAGS: -I"path/to/tdlib/include"
# cgo linux LDLAGS: -L"path/to/tdlib/lib"
*/
import "C"

import "github.com/mxmCherry/go-tdlib"

// your code
```

### Providing `CFLAGS`/`LDFLAGS` To Go Command

You can provide the CGo config when running Go commands (`run`, `build` etc):

```shell
CGO_CFLAGS='-I"/full/path/to/tdlib/include"' \
CGO_LDFLAGS='-L"/full/path/to/tdlib/lib"' \
  go run main.go
```

Keep in mind, that absolute (full) path is mandatory in this use case.
