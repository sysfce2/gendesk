package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	version_string  = "Desktop File Generator v.0.4.4"
	//icon_search_url = "https://admin.fedoraproject.org/pkgdb/appicon/show/%s"
    icon_search_url = "http://openiconlibrary.sourceforge.net/gallery2/open_icon_library-full/icons/png/48x48/apps/%s.png"
)

var (
	model3d_kw    = []string{"rendering", "modeling", "modeler", "render", "raytracing"}
	multimedia_kw = []string{"non-linear", "audio", "sound", "graphics", "draw", "demo"}
	network_kw    = []string{"network", "p2p", "browser"}
	audiovideo_kw = []string{"synth", "synthesizer"}
	office_kw     = []string{"ebook", "e-book"}
	editor_kw     = []string{"editor"}
	science_kw    = []string{"gps", "inspecting"}
	vcs_kw        = []string{"git"}
	// Emulator and player aren't always for games, but those cases should be
	// picked up by one of the other categories first
	game_kw          = []string{"game", "rts", "mmorpg", "emulator", "player"}
	arcadegame_kw    = []string{"combat", "arcade", "racing", "fighting", "fight"}
	actiongame_kw    = []string{"shooter", "fps"}
	adventuregame_kw = []string{"roguelike", "rpg"}
	logicgame_kw     = []string{"puzzle"}
	boardgame_kw     = []string{"board", "chess", "goban", "chessboard"}
	programming_kw   = []string{"code", "ide", "programming", "develop", "compile"}
	system_kw        = []string{"sensor"}

	// Global flags
	use_color = true
	verbose   = true
)

// Generate the contents for the .desktop file
func createDesktopContents(name string, genericName string, comment string,
	exec string, icon string, useTerminal bool,
	categories []string, mimeTypes []string,
	startupNotify bool) *bytes.Buffer {
	var buf []byte
	b := bytes.NewBuffer(buf)
	b.WriteString("[Desktop Entry]\n")
	b.WriteString("Encoding=UTF-8\n")
	b.WriteString("Type=Application\n")
	b.WriteString("Name=" + name + "\n")
	if genericName != "" {
		b.WriteString("GenericName=" + genericName + "\n")
	}
	b.WriteString("Comment=" + comment + "\n")
	b.WriteString("Exec=" + exec + "\n")
	b.WriteString("Icon=" + icon + "\n")
	if useTerminal {
		b.WriteString("Terminal=true\n")
	} else {
		b.WriteString("Terminal=false\n")
	}
	if startupNotify {
		b.WriteString("StartupNotify=true\n")
	} else {
		b.WriteString("StartupNotify=false\n")
	}
	b.WriteString("Categories=" + strings.Join(categories, ";") + ";\n")
	if len(mimeTypes) > 0 {
		b.WriteString("MimeType=" + strings.Join(mimeTypes, ";") + ";\n")
	}
	return b
}

// Capitalize a string or return the same if it is too short
func capitalize(s string) string {
	if len(s) >= 2 {
		return strings.ToTitle(s[0:1]) + s[1:]
	}
	return s
}

// Write the .desktop file as generated by createDesktopContents
func writeDesktopFile(pkgname string, name string, comment string, exec string,
	categories string, genericName string, mimeTypes string, custom string) {
	var categoryList []string
	var mimeTypeList []string

	if len(categories) == 0 {
		categoryList = []string{"Application"}
	} else {
		categoryList = strings.Split(categories, ";")
	}
	// mimeTypeList is an empty []string, or...
	if len(mimeTypes) != 0 {
		mimeTypeList = strings.Split(mimeTypes, ";")
	}

	// mimeTypes may be empty. Disabled terminal
	// and startupnotify for now.
	buf := createDesktopContents(name, genericName, comment, exec, pkgname,
		false, categoryList, mimeTypeList, false)
	if custom != "" {
		// Write the custom string to the end of the .desktop file (may contain \n)
		buf.WriteString(custom + "\n")
	}
	ioutil.WriteFile(pkgname+".desktop", buf.Bytes(), 0666)
}

// Checks if a trimmed line starts with a specific word
func startsWith(line string, word string) bool {
	return 0 == strings.Index(strings.TrimSpace(line), word)
}

// Return what's between two strings, "a" and "b", in another string
func between(orig string, a string, b string) string {
	if strings.Contains(orig, a) && strings.Contains(orig, b) {
		posa := strings.Index(orig, a) + len(a)
		posb := strings.LastIndex(orig, b)
		return orig[posa:posb]
	}
	return ""
}

// Return the contents between "" or '' (or an empty string)
func betweenQuotes(orig string) string {
	var s string
	for _, quote := range []string{"\"", "'"} {
		s = between(orig, quote, quote)
		if s != "" {
			return s
		}
	}
	return ""
}

// Return the string between the quotes or after the "=", if possible
// or return the original string
func betweenQuotesOrAfterEquals(orig string) string {
	s := betweenQuotes(orig)
	// Check for exactly one "="
	if (s == "") && (strings.Count(orig, "=") == 1) {
		s = strings.TrimSpace(strings.Split(orig, "=")[1])
	}
	return s
}

// Does a keyword exist in the string?
func has(s string, kw string) bool {
	lowercase := strings.ToLower(s)
	// Remove the most common special characters
	massaged := strings.Trim(lowercase, "-_.,!?()[]{}\\/:;+@")
	words := strings.Split(massaged, " ")
	for _, word := range words {
		if word == kw {
			return true
		}
	}
	return false
}

// Check if a keyword appears in a package description
func keywordsInDescription(pkgdesc string, keywords []string) bool {
	for _, keyword := range keywords {
		if has(pkgdesc, keyword) {
			return true
		}
	}
	return false
}

// Download icon from the search url in icon_search_url
func writeIconFile(pkgname string, o *Output) error {
	// Only supports png icons
	filename := pkgname + ".png"
	var client http.Client
	resp, err := client.Get(fmt.Sprintf(icon_search_url, pkgname))
	if err != nil {
		o.errText("Could not download icon")
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		o.errText("Could not dump body")
		os.Exit(1)
	}

	var h hash.Hash = md5.New()
	h.Write(b)
	//fmt.Printf("Icon MD5: %x\n", h.Sum())

	// If the icon is the "No icon found" icon (known hash), return with an error
	if fmt.Sprintf("%x", h.Sum(nil)) == "12928aa3233965175ea30f5acae593bf" {
		return errors.New("No icon found")
	}

	if b[0] == 60 && b[1] == 104 && b[2] == 116 {
		// if it starts with "<ht", it's not a png
		return errors.New("No icon found")
	}

	err = ioutil.WriteFile(filename, b, 0666)
	if err != nil {
		o.errText("Could not write icon to " + filename + "!")
		os.Exit(1)
	}
	return nil
}

func writeDefaultIconFile(pkgname string, o *Output) error {
	defaultIconFilename := "/usr/share/pixmaps/default.png"
	b, err := ioutil.ReadFile(defaultIconFilename)
	if err != nil {
		o.errText("could not read " + defaultIconFilename + "!")
		os.Exit(1)
	}
	filename := pkgname + ".png"
	err = ioutil.WriteFile(filename, b, 0666)
	if err != nil {
		o.errText("could not write icon to " + filename + "!")
		os.Exit(1)
	}
	return nil
}

// Download a file
func downloadFile(url string, filename string, o *Output) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		o.errText("Could not download file")
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		o.errText("Could not dump body")
		os.Exit(1)
	}
	err = ioutil.WriteFile(filename, b, 0666)
	if err != nil {
		o.errText("Could not write data to " + filename + "!")
		os.Exit(1)
	}
}

// Use a function for each element in a string list and
// return the modified list
func stringMap(f func(string) string, stringlist []string) []string {
	newlist := make([]string, len(stringlist))
	for i, elem := range stringlist {
		newlist[i] = f(elem)
	}
	return newlist
}

// Return a list of pkgnames for split packages
// or just a list with the pkgname for regular packages
func pkgList(splitpkgname string) []string {
	center := between(splitpkgname, "(", ")")
	if center == "" {
		center = splitpkgname
	}
	if strings.Contains(center, " ") {
		unquoted := strings.Replace(center, "\"", "", -1)
		unquoted = strings.Replace(unquoted, "'", "", -1)
		return strings.Split(unquoted, " ")
	}
	return []string{splitpkgname}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {

	var filename string
	version_help := "Show application name and version"
	nodownload_help := "Don't download anything"
	nocolor_help := "Don't use colors"
	quiet_help := "Don't output anything on stdout"
	flag.Usage = func() {
		fmt.Println()
		fmt.Println(version_string)
		fmt.Println("generates .desktop files from a PKGBUILD")
		fmt.Println()
		fmt.Println("Syntax: gendesk [flags] filename")
		fmt.Println()
		fmt.Println("Possible flags:")
		fmt.Println("    * --version        " + version_help)
		fmt.Println("    * -n               " + nodownload_help)
		fmt.Println("    * --nocolor        " + nocolor_help)
		fmt.Println("    * -q               " + quiet_help)
		fmt.Println("    * --help           This text")
		fmt.Println()
		fmt.Println("Note:")
		fmt.Println("    * \"../PKGBUILD\" is the default filename")
		fmt.Println("    * _exec in the PKGBUILD can be used to specific a")
		fmt.Println("      different executable for the .desktop file")
		fmt.Println("      Example: _exec=('appname-gui')")
		fmt.Println("    * Split packages are supported")
		fmt.Println("    * If a .png or .svg icon is not found as a file or in")
		fmt.Println("      the PKGBUILD, an icon will be downloaded from:")
		fmt.Println("      " + icon_search_url)
		fmt.Println("      This may or may not result in the icon you wished for.")
		fmt.Println("    * Categories are guessed based on keywords in the")
		fmt.Println("      package description, but there's also _categories=().")
		fmt.Println("    * Icons are assumed to be installed to")
		fmt.Println("      \"/usr/share/pixmaps/\" by the PKGBUILD")
		fmt.Println()
	}
	version := flag.Bool("version", false, version_help)
	nodownload := flag.Bool("n", false, nodownload_help)
	nocolor := flag.Bool("nocolor", false, nocolor_help)
	quiet := flag.Bool("q", false, quiet_help)
	flag.Parse()
	args := flag.Args()

	// New output. Color? Enabled?
	o := NewOutput(!*nocolor, !*quiet)

	if *version {
		o.Println(version_string)
		os.Exit(0)
	} else if len(args) == 0 {
		filename = "../PKGBUILD"
	} else if len(args) == 1 {
		filename = args[0]
	} else {
		o.errText("Too many arguments")
		os.Exit(1)
	}

	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		o.errText("Could not read " + filename)
		os.Exit(1)
	}
	filetext := string(filedata)

	var pkgname string
	var pkgnames []string
	var iconurl string

	// Several fields are stored per pkgname, hence the maps
	pkgdescMap := make(map[string]string)
	execMap := make(map[string]string)
	nameMap := make(map[string]string)
	genericNameMap := make(map[string]string)
	mimeTypeMap := make(map[string]string)
	commentMap := make(map[string]string)
	customMap := make(map[string]string)
	categoriesMap := make(map[string]string)

	for _, line := range strings.Split(filetext, "\n") {
		if startsWith(line, "pkgname") {
			pkgname = betweenQuotesOrAfterEquals(line)
			pkgnames = pkgList(pkgname)
			// Select the first pkgname in the array as the "current" pkgname
			if len(pkgnames) > 0 {
				pkgname = pkgnames[0]
			}
		} else if startsWith(line, "package_") {
			pkgname = between(line, "_", "(")
		} else if startsWith(line, "pkgdesc") {
			// Description for the package
			pkgdesc := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				pkgdescMap[pkgname] = pkgdesc
			}
		} else if startsWith(line, "_exec") {
			// Custom executable for the .desktop file per (split) package
			exec := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				execMap[pkgname] = exec
			}
		} else if startsWith(line, "_name") {
			// Custom Name for the .desktop file per (split) package
			name := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				nameMap[pkgname] = name
			}
		} else if startsWith(line, "_genericname") {
			// Custom GenericName for the .desktop file per (split) package
			genericName := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if (pkgname != "") && (genericName != "") {
				genericNameMap[pkgname] = genericName
			}
		} else if startsWith(line, "_mimetype") {
			// Custom MimeType for the .desktop file per (split) package
			mimeType := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				mimeTypeMap[pkgname] = mimeType
			}
		} else if startsWith(line, "_comment") {
			// Custom Comment for the .desktop file per (split) package
			comment := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				commentMap[pkgname] = comment
			}
		} else if startsWith(line, "_custom") {
			// Custom string to be added to the end
			// of the .desktop file in question
			custom := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				customMap[pkgname] = custom
			}
		} else if startsWith(line, "_categories") {
			categories := betweenQuotesOrAfterEquals(line)
			// Use the last found pkgname as the key
			if pkgname != "" {
				categoriesMap[pkgname] = categories
			}
		} else if strings.Contains(line, "http://") &&
			strings.Contains(line, ".png") {
			// Only supports png icons downloaded over http,
			// picks the first fitting url
			if iconurl == "" {
				iconurl = "h" + between(line, "h", "g") + "g"
				if strings.Contains(iconurl, "$pkgname") {
					iconurl = strings.Replace(iconurl,
						"$pkgname", pkgname, -1)
				}
				if strings.Contains(iconurl, "${pkgname}") {
					iconurl = strings.Replace(iconurl,
						"${pkgname}", pkgname, -1)
				}
				if strings.Contains(iconurl, "$") {
					// If there are more $variables, don't bother (for now)
					// TODO: replace all defined $variables...
					iconurl = ""
				}
			}
		}
	}

	//o.Println("pkgnames:", pkgnames)

	// Write .desktop and .png icon for each package
	for _, pkgname := range pkgnames {
		if strings.Contains(pkgname, "-nox") ||
			strings.Contains(pkgname, "-cli") {
			// Don't bother if it's a -nox or -cli package
			continue
		}
		pkgdesc, found := pkgdescMap[pkgname]
		if !found {
			// Fall back on the package name
			pkgdesc = pkgname
		}
		exec, found := execMap[pkgname]
		if !found {
			// Fall back on the package name
			exec = pkgname
		}
		name, found := nameMap[pkgname]
		if !found {
			// Fall back on the capitalized package name
			name = capitalize(pkgname)
		}
		genericName, found := genericNameMap[pkgname]
		if !found {
			// Fall back on no generic name
			genericName = ""
		}
		comment, found := commentMap[pkgname]
		if !found {
			// Fall back on pkgdesc
			comment = pkgdesc
		}
		mimeType, found := mimeTypeMap[pkgname]
		if !found {
			// Fall back on no mime type
			mimeType = ""
		}
		custom, found := customMap[pkgname]
		if !found {
			// Fall back on no custom additional lines
			custom = ""
		}
		categories, found := categoriesMap[pkgname]
		if !found {
			// Approximately identify various categories
			categories = ""
			if keywordsInDescription(pkgdesc, model3d_kw) {
				categories = "Application;Graphics;3DGraphics"
			} else if keywordsInDescription(pkgdesc, multimedia_kw) {
				categories = "Application;Multimedia"
			} else if keywordsInDescription(pkgdesc, network_kw) {
				categories = "Application;Network"
			} else if keywordsInDescription(pkgdesc, audiovideo_kw) {
				categories = "Application;AudioVideo"
			} else if keywordsInDescription(pkgdesc, office_kw) {
				categories = "Application;Office"
			} else if keywordsInDescription(pkgdesc, editor_kw) {
				categories = "Application;Development;TextEditor"
			} else if keywordsInDescription(pkgdesc, science_kw) {
				categories = "Application;Science"
			} else if keywordsInDescription(pkgdesc, vcs_kw) {
				categories = "Application;Development;RevisionControl"
			} else if keywordsInDescription(pkgdesc, arcadegame_kw) {
				categories = "Application;Game;ArcadeGame"
			} else if keywordsInDescription(pkgdesc, actiongame_kw) {
				categories = "Application;Game;ActionGame"
			} else if keywordsInDescription(pkgdesc, adventuregame_kw) {
				categories = "Application;Game;AdventureGame"
			} else if keywordsInDescription(pkgdesc, logicgame_kw) {
				categories = "Application;Game;"
			} else if keywordsInDescription(pkgdesc, boardgame_kw) {
				categories = "Application;Game;BoardGame"
			} else if keywordsInDescription(pkgdesc, game_kw) {
				categories = "Application;Game"
			} else if keywordsInDescription(pkgdesc, programming_kw) {
				categories = "Application;Development"
			} else if keywordsInDescription(pkgdesc, system_kw) {
				categories = "Application;System"
			}
		}
		const nSpaces = 32
		spaces := strings.Repeat(" ", nSpaces)[:nSpaces-min(nSpaces, len(pkgname))]
		if o.isEnabled() {
			fmt.Printf("%s%s%s%s%s ",
				o.darkGrayText("["), o.lightBlueText(pkgname),
				o.darkGrayText("]"), spaces,
				o.darkGrayText("Generating desktop file..."))
		}
		writeDesktopFile(pkgname, name, comment, exec,
			categories, genericName, mimeType, custom)
		if o.isEnabled() {
			fmt.Printf("%s\n", o.darkGreenText("ok"))
		}

		// Download an icon if it's not downloaded by
		// the PKGBUILD and not there already (.png or .svg)
		pngFilenames, _ := filepath.Glob("*.png")
		svgFilenames, _ := filepath.Glob("*.svg")
		if ((0 == (len(pngFilenames) + len(svgFilenames))) && (iconurl == "")) && (*nodownload == false) {
			if len(pkgname) < 1 {
				o.errText("No pkgname, can't download icon")
				os.Exit(1)
			}
			fmt.Printf("%s%s%s%s%s ",
				o.darkGrayText("["), o.lightBlueText(pkgname),
				o.darkGrayText("]"), spaces,
				o.darkGrayText("Downloading icon..."))
			err := writeIconFile(pkgname, o)
			if err == nil {
				if o.isEnabled() {
					fmt.Printf("%s\n", o.lightCyanText("ok"))
				}
			} else {
				if o.isEnabled() {
					fmt.Printf("%s\n", o.darkYellowText("no"))
					fmt.Printf("%s%s%s%s%s ",
						o.darkGrayText("["),
						o.lightBlueText(pkgname),
						o.darkGrayText("]"),
						spaces,
						o.darkGrayText("Using default icon instead..."))
				}
				err := writeDefaultIconFile(pkgname, o)
				if err == nil {
					if o.isEnabled() {
						fmt.Printf("%s\n", o.lightPurpleText("yes"))
					}
				}
			}
		}
	}
}
