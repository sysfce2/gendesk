# VT100

[![Build Status](https://travis-ci.com/xyproto/vt100.svg?branch=master)](https://travis-ci.com/xyproto/vt100) [![GoDoc](https://godoc.org/github.com/xyproto/vt100?status.svg)](https://godoc.org/github.com/xyproto/vt100) [![License](https://img.shields.io/badge/license-BSD-blue.svg?style=flat)](https://raw.githubusercontent.com/xyproto/vt100/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/vt100)](https://goreportcard.com/report/github.com/xyproto/vt100)

* Supports colors and attributes.
* Supports platforms with VT100 support and a `/dev/tty` device.
* Can detect the terminal size.
* Can get key-presses, including arrow keys (252, 253, 254, 255) and pgup/pgdn (251, 250).
* Has a Canvas struct, for drawing only the updated lines to the terminal.
* Uses the a reference document directly, but memoizes the commands sent to the terminal, for performance.
* Could be used for making an alternative to the `dialog` or `whiptail` utilities.

### Editor

For an editor that uses this module, take a look at [orbiton](https://github.com/xyproto/orbiton).

### Images

![shooter example](img/shooter.gif)

Screen recording of the [`shooter`](cmd/shooter) example, where you can control a small character with the arrow keys and shoot with `space`.

---

![menu example](img/menu.gif)

Screen recording of the [`menu`](cmd/menu) example, which uses VT100 terminal codes and demonstrates a working menu.

---

![VT100 terminal](https://upload.wikimedia.org/wikipedia/commons/thumb/9/99/DEC_VT100_terminal.jpg/300px-DEC_VT100_terminal.jpg)

A physical VT100 terminal. Photo by [Jason Scott](https://www.flickr.com/photos/54568729@N00/9636183501), [CC BY 2.0](https://creativecommons.org/licenses/by/2.0)

### Requirements

* Go 1.17 or later.

### Features and limitations

* Can detect letters, arrow keys and space. F12 and similar keys are not supported (they are supported by VT220 but not VT100).
* Resizing the terminal when using the Canvas struct may cause artifacts, for a brief moment.
* Holding down a key may trigger key repetition which may speed up the main loop.
* As an exception, pgup and pgdown are supported.

### Simple use

Output "hi" in blue:

```go
vt100.Blue.Output("hi")
```

Erase the current line:

```go
vt100.Do("Erase Line")
```

Move the cursor 3 steps up (it's a bit verbose, but it is memoized for better performance):

```go
vt100.Set("Cursor Up", map[string]string{"{COUNT}": "3"})
```

The full overview of possible commands are at the top of `vt100.go`.

### Another example

See `cmd/move` for a more advanced example, where a character can be moved around with the arrow keys.

### A small editor using `vt100`

The [Orbiton editor](https://github.com/xyproto/orbiton) uses `vt100`:

Quick installation:

    go install github.com/xyproto/orbiton/v2@latest

### General info

* Version: 1.20.0
* Licence: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
