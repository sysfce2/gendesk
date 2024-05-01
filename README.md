## Desktop File Generator

![Build](https://github.com/xyproto/gendesk/workflows/Build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/gendesk)](https://goreportcard.com/report/github.com/xyproto/gendesk)

Generates `.desktop` files and downloads `.png` icons based on command line arguments.

See `gendesk --help` or the man page for more info.

Pull requests are welcome.

[![Packaging status](https://repology.org/badge/vertical-allrepos/gendesk.svg)](https://repology.org/project/gendesk/versions)

## Changes from 1.0.9 to 1.0.10

## Changes from 1.0.8 to 1.0.9

* Update documentation.
* Expand variables, ref #16.
* Add a `--path` flag for setting the starting directory, ref #17.
* Add keywords for detecting the `AudioVideo` category, and for detecting e-mail related applications.
* Only set `noExecSpecififed` if the flag was not given.
* If `--exec` is not specified for e-mail related applications, add ` %u` to the `Exec` field.
* Consider the `Email;Network;Office;` categories, ref #19.
* Update dependencies.

## Changes from 1.0.7 to 1.0.8

* Update dependencies.

## Changes from 1.0.6 to 1.0.7

* Update documentation.
* Also strip the `-bin` suffix.
* Update CI configuration.
* Minor improvement to the `--help` output.
* Add additional categories.
* Update dependencies.

## Changes from 1.0.5 to 1.0.6

* Fix an issue with the `-hg` suffix (thanks Michael Straube).
* Fix an issue with the package description (thanks Simon Dierl).
* Minor changes.
* Update dependencies.

## Changes from 1.0.4 to 1.0.5

* If the first argument is not a file that exists, use it as the package name.
* Update dependencies.

## Changes from 1.0.3 to 1.0.4

* Include go.mod and go.sum in the release package.
* Minor changes to the output message when no arguments are given.
* Update dependencies.

## Changes from 1.0.2 to 1.0.3

* Let flags overrides the values from a given PKGBUILD file.
* Switch from [goconf](https://github.com/akrennmair/goconf) to [goconfig](https://github.com/unknwon/goconfig).
* Switch from [term](https://github.com/xyproto/term) to [textoutput](https://github.com/xyproto/textoutput).
* Requires Go 1.10 or later.

## Changes from 1.0.1 to 1.0.2

* Set version to 1.0 instead of 1.2 when generating `.desktop` files, to support a wider range of distributions.

## Changes from 1.0.0 to 1.0.1

* Fix a typo in the `--help` output.
* Update the release script.

## Changes from 0.7.0 to 1.0.0

* Add `--icon` flag, ref #7.
* Update to the desktop-entry-spec 1.2 format (remove `Encoding` and specify `Version`), ref #8.
* Several minor changes, as suggested by the `golint` utility.
* Tested with Go 1.11.

## Changes from 0.6.5 to 0.7.0

* Updated vendored dependencies.
* Added support for [goreleaser](https://github.com/goreleaser/goreleaser).
* Improved handling of icons, if an icon is missing.
* Minor changes and refactoring.

## Changes from 0.6.4 to 0.6.5

* Ignore the `-svn` suffix in package names (same as for `-git`, thanks @mstraube).
* Use `text/template` for generating the `.desktop` file contents.
* Minor changes to the command line output/documentation.
* Some refactoring.
* Tested with Go 1.9.

## Changes from 0.6.3 to 0.6.4

* Fix bug where some flags could not be overridden.

## Changes from 0.6.2 to 0.6.3

* Will now ignore the `-git` suffix if it is part of a package name.

## Changes from 0.6.1 to 0.6.2

* Added the possibility of having a configuration file for specifying a different URL for searching for missing icons.
* Remove the `--iconurl` flag.
* Refactored out some code to an external package.

## Changes from 0.6 to 0.6.1

* Support for `StartupNotify=true`/`false`
* Both `--mimetype` and `--mimetypes` are allowed
* Guesses more categories than before

## Changes from 0.5.5 to 0.6

* Added an option for generating .desktop files for launching window managers

## Changes from 0.5.4 to 0.5.5

* Bug fix when generating .desktop files from PKGBUILD files.

## Changes from 0.5.3 to 0.5.4

* Added a `-f` flag for overwriting files (will not overwrite without it).
* Some refactoring

## Changes from 0.5.2 to 0.5.3

* Added a `--terminal` flag for specifying if the application should be run in a terminal.
* Some refactoring.

## Changes from 0.5.1 to 0.5.2

* Support for additional environment variables.

## Changes from 0.5.0 to 0.5.1

* Support for `$pkgname` and `$pkgdesc`.
* Updated the man page.
* Will try to download icons specified with `--iconurl`.

## Changes from 0.4.4 to 0.5.0

* Command line options, no need to specify a PKGBUILD.

## Changes from 0.4.3 to 0.4.4

* Changed the URL for searching for icons from Fedora to Open Icon Library

## Changes from 0.4.2 to 0.4.3

* Fixed minor bug where puzzle games were not placed in the right category
* Added \_categories=()

## Changes from 0.4.1 to 0.4.2

* Added category "Graphics;3DGraphics;" for 3D modellers
* Added category "System;" for sensor monitors
* Added category "Game;BoardGame;" for kw "board", "chess", "goban" or "chessboard"
* Added category "Office" for kw "e-book" and "ebook"
* Doesn't use ".png" by default when specifying an icon

## Changes from 0.4 to 0.4.1

* Fixed a bug where \_name=() and \_comment=() didn't work as they should

## Changes from 0.3 to 0.4

* Added \_genericname=()
* Added \_comment=()
* Added \_mimetype=()
* Added Type=Application
* Added category "Game;LogicGame" for keyword "puzzle"
* Added category "Game;ArcadeGame" for keyword "fighting"
* Fixed weird formatting in --help output
* Added \_custom=() for adding custom fields at the end of the .desktop file
* Glob for existing .svg icons too
* Shorter lines
* Moved functions and settings related to terminal output to a separate file

## Changes from 0.2 to 0.3

* New flag: -q for quiet
* New flag: --nocolor for no color
* New flag: -n for not downloading anything (only generate a .desktop file)
* New flag: -q for quiet (no stdout output)
* Added \_name=('Name') to be able to specify a name that isn't only lowercase (like "ZynAddSubFX" or "jEdit")
* kw "synthesizer" is now category AudioVideo
* kw "editor" is now category TextEditor and/or Development;TextEditor
* kw "emulator" is now category "Game"
* kw "game" is now category "Game"
* kw "combat" is now be category "Game;ArcadeGame"
* kw "GPS" or "inspecting" is now category "Application;Science"
* kw "player" is now category "Application;Game;"
* kw "shooter" is now "Application;Game;ActionGame;"
* kw "roguelike" is now "Application;Game;AdventureGame;"
* kw "git" is now category Development;RevisionControl

## Requirements

* Go 1.17 or later.

## Troubleshooting

* If you get something like `GLIBC_3.32 not found` on Linux, try the `gendesk-1.x.x-linux_static` release.

## General information

* Version: 1.0.10
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: BSD-3
