package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/xyproto/env"
	"github.com/xyproto/textoutput"
)

const (
	versionString = "Desktop File Generator 1.0.9"

	versionHelp       = "Show application name and version"
	nodownloadHelp    = "Don't download anything"
	nocolorHelp       = "Don't use colors"
	quietHelp         = "Don't output anything on stdout"
	forceHelp         = "Overwrite .desktop files with the same name"
	windowmanagerHelp = "Generate a .desktop file for launching a window manager"
	pkgnameHelp       = "The name of the package"
	pkgdescHelp       = "Description of the package"
	pathHelp          = "Starting directory"
	nameHelp          = "Name of the shortcut"
	genericnameHelp   = "Type of application"
	commentHelp       = "Shortcut comment"
	execHelp          = "Path to executable"
	terminalHelp      = "Run the application in a terminal (default is false)"
	categoriesHelp    = "Categories, see other .desktop files for examples"
	mimetypesHelp     = "Mime types, see other .desktop files for examples"
	startupnotifyHelp = "Notification when the application starts (default is false)"
	customHelp        = "Custom line to append at the end of the .desktop file"
	iconHelp          = "Specify a filename that will be used for the icon"

	defaultPKGBUILD = "../PKGBUILD"
)

// WMStarter contains the information needed to generate
// a .desktop file for a Window Manager
type WMStarter struct {
	Name, Exec string
}

// AppStarter contains the information needed to generate
// a .desktop file for an application
type AppStarter struct {
	Name, GenericName, Comment, Exec, Icon, Path string
	UseTerminal, StartupNotify                   bool
	CategoryList, MimeTypesList                  string
}

var (
	// Template for a .desktop file for starting a Window Manager
	wmTemplate, _ = template.New("WMStarter").Parse("[Desktop Entry]\nType=XSession\nExec={{.Exec}}\nTryExec={{.Exec}}\nName={{.Name}}\n")

	// Template for a .desktop file for starting an application
	appTemplate, _ = template.New("AppStarter").Parse("[Desktop Entry]\nVersion=1.0\nType=Application\nName={{.Name}}\n{{if .GenericName}}GenericName={{.GenericName}}\n{{end}}Comment={{.Comment}}\nExec={{.Exec}}\nIcon={{.Icon}}{{if .Path}}\nPath={{.Path}}{{end}}\nTerminal={{if .UseTerminal}}true{{else}}false{{end}}\nStartupNotify={{if .StartupNotify}}true{{else}}false{{end}}\nCategories={{.CategoryList}};\n{{if .MimeTypesList}}MimeType={{.MimeTypesList}};\n{{end}}")
)

// Generate the contents for the .desktop file (for executing a window manager)
func createWindowManagerDesktopContents(name, execCommand string) *bytes.Buffer {
	var (
		buf bytes.Buffer
		wm  = WMStarter{name, execCommand}
	)
	// Inserting strings into the template should always work, panic if not
	if err := wmTemplate.Execute(&buf, wm); err != nil {
		panic(err)
	}
	return &buf
}

// Generate the contents for the .desktop file (for executing a desktop application)
func createDesktopContents(name, genericName, comment, execCommand, icon, path string, useTerminal, startupNotify bool, categories, mimeTypes []string) *bytes.Buffer {
	var (
		buf bytes.Buffer
		app = AppStarter{name, genericName, comment, execCommand, icon, path, useTerminal, startupNotify, strings.Join(categories, ";"), strings.Join(mimeTypes, ";")}
	)
	// Inserting strings into the template should always work, panic if not
	if err := appTemplate.Execute(&buf, app); err != nil {
		panic(err)
	}
	return &buf
}

// Write the .desktop file as generated by createWindowManagerDesktopContents
func writeWindowManagerDesktopFile(pkgname, name, execCommand, custom string, force bool, o *textoutput.TextOutput) {
	buf := createWindowManagerDesktopContents(name, execCommand)
	if custom != "" {
		// Write the custom string to the end of the .desktop file (may contain \n)
		buf.WriteString(custom + "\n")
	}
	// Check if the file exists (and that force is not enabled)
	if _, err := os.Stat(pkgname + ".desktop"); err == nil && (!force) {
		o.Err("no")
		o.Fprintln(os.Stderr, pkgname+".desktop already exists. Use -f as the first argument to overwrite it.")
		os.Exit(1)
	}
	ioutil.WriteFile(pkgname+".desktop", buf.Bytes(), 0644)
}

// Write the .desktop file as generated by createDesktopContents
func writeDesktopFile(pkgname, name, comment, execCommand, icon, path, categories, genericName, mimeTypes, custom string, useTerminal, startupNotify, force bool, o *textoutput.TextOutput) {
	var categoryList, mimeTypeList []string

	if len(categories) == 0 {
		categoryList = []string{"Application"}
	} else {
		categoryList = strings.Split(categories, ";")
	}

	// mimeTypeList is an empty []string, or...
	if len(mimeTypes) != 0 {
		mimeTypeList = strings.Split(mimeTypes, ";")
	}

	// Use the pkgname is the icon name if no icon is specified
	if icon == "" {
		icon = pkgname
	}

	// mimeTypes may be empty. Disabled terminal
	// and startupnotify for now.
	buf := createDesktopContents(name, genericName, comment, execCommand, icon, path, useTerminal, startupNotify, categoryList, mimeTypeList)
	if custom != "" {
		// Write the custom string to the end of the .desktop file (may contain \n)
		buf.WriteString(custom + "\n")
	}

	// Check if the file exists (and that force is not enabled)
	if _, err := os.Stat(pkgname + ".desktop"); err == nil && (!force) {
		o.Err("no")
		o.Fprintln(os.Stderr, pkgname+".desktop already exists. Use -f as the first argument to overwrite it.")
		os.Exit(1)
	}

	ioutil.WriteFile(pkgname+".desktop", buf.Bytes(), 0644)
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

// WriteDefaultIconFile copies /usr/share/pixmaps/default.png to pkgname + ".png"
func WriteDefaultIconFile(pkgname string, o *textoutput.TextOutput) error {
	defaultIconFilename := "/usr/share/pixmaps/default.png"
	b, err := ioutil.ReadFile(defaultIconFilename)
	if err != nil {
		o.Err("could not read " + defaultIconFilename + "!")
	}
	filename := pkgname + ".png"
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		o.Err("could not write icon to " + filename + "!")
	}
	return nil
}

func usage() {
	shortname := strings.Split(defaultIconSearchURL, "/")
	firstpart := "INVALID ICON SEARCH URL"
	if len(shortname) >= 3 {
		firstpart = strings.Join(shortname[:3], "/")
	}
	fmt.Println(`
` + versionString + `
Generate .desktop files.

Syntax: gendesk [flags]

Possible flags:
    --version                    ` + versionHelp + `
    -n                           ` + nodownloadHelp + `
    --nocolor                    ` + nocolorHelp + `
    -q                           ` + quietHelp + `
    -f                           ` + forceHelp + `
    -wm                          ` + windowmanagerHelp + `
    --pkgname=PKGNAME            ` + pkgnameHelp + `
    --pkgdesc=PKGDESC            ` + pkgdescHelp + `
    --path=PATH                  ` + pathHelp + `
    --name=NAME                  ` + nameHelp + `
    --genericname=GENERICNAME    ` + genericnameHelp + `
    --comment=COMMENT            ` + commentHelp + `
    --exec=EXEC                  ` + execHelp + `
    --icon=FILENAME              ` + iconHelp + `
    --terminal=[true|false]      ` + terminalHelp + `
    --categories=CATEGORIES      ` + categoriesHelp + `
    --mimetypes=MIMETYPES        ` + mimetypesHelp + `
    --startupnotify=[true|false] ` + startupnotifyHelp + `
    --custom=CUSTOM              ` + customHelp + `
    --help                       This text

Note:
    * Just providing a package name is enough to generate a .desktop file.
    * Providing a PKGBUILD filename instead of flags is a possibility.
    * "$startdir/PKGBUILD" is the default PKGBUILD filename.
    * _exec in the PKGBUILD can be used to specify a different executable for
      the .desktop file. Example: _exec=('appname-gui')
    * Split PKGBUILD packages are supported.
    * If a .png, .svg or .xpm icon is not found as a file or in the PKGBUILD,
      an icon will be downloaded from either the location specified in the
      configuration or from: ` + firstpart + `
      (This may or may not result in the icon you wished for).
    * Categories are guessed based on keywords in the
      package description, unless provided.
    * Icons are assumed to be found in "/usr/share/pixmaps/" once installed.
`)
}

func main() {
	flag.Usage = usage

	var (
		version       = flag.Bool("version", false, versionHelp)
		nodownload    = flag.Bool("n", false, nodownloadHelp)
		nocolor       = flag.Bool("nocolor", false, nocolorHelp)
		quiet         = flag.Bool("q", false, quietHelp)
		force         = flag.Bool("f", false, forceHelp)
		windowmanager = flag.Bool("wm", false, windowmanagerHelp)
		givenPkgname  = flag.String("pkgname", "", pkgnameHelp)
		givenPkgdesc  = flag.String("pkgdesc", "", pkgdescHelp)
		name          = flag.String("name", "", nameHelp)
		path          = flag.String("path", "", pathHelp)
		genericname   = flag.String("genericname", "", genericnameHelp)
		comment       = flag.String("comment", "", commentHelp)
		execCommand   = flag.String("exec", "", execHelp)
		icon          = flag.String("icon", "", iconHelp)
		terminal      = flag.Bool("terminal", false, terminalHelp)
		categories    = flag.String("categories", "", categoriesHelp)
		mimetypes     = flag.String("mimetypes", "", mimetypesHelp)
		mimetype      = flag.String("mimetype", "", mimetypesHelp)
		custom        = flag.String("custom", "", customHelp)
		startupnotify = flag.Bool("startupnotify", false, startupnotifyHelp)

		manualIconurl string
		filename      string
		pkgnames      []string
		iconurl       string
	)

	flag.Parse()

	var (
		args    = flag.Args()
		pkgname = *givenPkgname
		pkgdesc = *givenPkgdesc

		// New text output struct.
		// The first bool is if color should be enabled or disabled.
		// The second bool is if any output should be enabled at all.
		o = textoutput.NewTextOutput(!*nocolor, !*quiet)
	)

	// Output the version number and quit if --version is given
	if *version {
		o.Println(versionString)
		return
	}

	// TODO: Write in a cleaner way, possibly by refactoring into a function. Write a test first.
	if pkgname == "" {
		if len(args) == 0 {
			if envPkgname := env.Str("pkgname"); envPkgname == "" {
				if envSrcdest := env.Str("SRCDEST"); envSrcdest != "" {
					// If SRCDEST is set, use that
					filename = filepath.Join(envSrcdest, "PKGBUILD")
				} else {
					filename = defaultPKGBUILD
				}
			} else {
				pkgname = envPkgname
			}
		} else {
			// args are non-flag arguments
			filename = args[0]
		}
	}

	// Environment variables
	dataFromEnvironment(&pkgdesc, execCommand, name, genericname, mimetypes, comment, categories, custom)

	// Several fields are stored per pkgname
	pkgdescMap := make(map[string]string)
	execMap := make(map[string]string)
	nameMap := make(map[string]string)
	genericNameMap := make(map[string]string)
	mimeTypesMap := make(map[string]string)
	commentMap := make(map[string]string)
	categoriesMap := make(map[string]string)
	customMap := make(map[string]string)

	// Strip the "-bin", "-git", "-hg" or "-svn" suffix from the name, if present
	for _, suf := range []string{"bin", "git", "hg", "svn"} {
		pkgname = strings.TrimSuffix(pkgname, "-"+suf)
	}

	if filename != "" {
		// Check if the given filename is found
		if !exists(filename) {
			// If --pkgname is not given and the file does not exist, use the base name as the pkgname
			pkgname = filepath.Base(filename)
			// Clear the filename variable, since the file was not found
			filename = ""
		} else {
			// TODO: Use a struct per pkgname instead
			parsePKGBUILD(o, filename, &iconurl, &pkgname, &pkgnames, &pkgdescMap, &execMap, &nameMap, &genericNameMap, &mimeTypesMap, &commentMap, &categoriesMap, &customMap)
		}
	}

	// Fill in the dictionaries using the given arguments. This overrides values from the PKGBUILD.
	pkgnames = []string{pkgname}

	// Set a value if the value is not an empty string
	setv := func(m *map[string]string, value string) {
		if value != "" {
			(*m)[pkgname] = value
		}
	}

	noExecSpecified := *execCommand == ""

	setv(&pkgdescMap, pkgdesc)
	setv(&execMap, *execCommand)
	setv(&nameMap, *name)
	setv(&genericNameMap, *genericname)
	setv(&mimeTypesMap, *mimetype)
	setv(&mimeTypesMap, *mimetypes)
	setv(&commentMap, *comment)
	setv(&categoriesMap, *categories)
	setv(&customMap, *custom)

	// Write .desktop and .png icon for each package
	for _, pkgname := range pkgnames {
		if strings.Contains(pkgname, "-nox") || strings.Contains(pkgname, "-cli") {
			// Don't bother if it's a -nox or -cli package
			continue
		}

		// TODO: Find a better way for all the if checks below
		pkgdesc, found := pkgdescMap[pkgname]
		if !found {
			// Fall back on the package name
			pkgdesc = pkgname
		}
		execCommand, found := execMap[pkgname]
		if !found {
			// Fall back on the package name
			execCommand = pkgname
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
		mimeTypes, found := mimeTypesMap[pkgname]
		if !found {
			// Fall back on no mime type
			mimeTypes = ""
		}
		custom, found := customMap[pkgname]
		if !found {
			// Fall back on no custom additional lines
			custom = ""
		}
		categories, found := categoriesMap[pkgname]
		if !found {
			categories = GuessCategory(pkgdesc)
		}

		// For the "Email" category: add "%u" to exec, if no exec command has been specified
		if strings.Contains(categories, "Email") && noExecSpecified && !strings.HasSuffix(execCommand, "%u") {
			// %u is added to be able to open mailto: links with e-mail applications
			execCommand += " %u"
		}

		// TODO: Refactor into a function
		const nSpaces = 32
		spaces := strings.Repeat(" ", nSpaces)[:nSpaces-min(nSpaces, len(pkgname))]
		if o.IsEnabled() {
			fmt.Printf("%s%s%s%s%s ",
				o.DarkGray("["), o.LightBlue(pkgname),
				o.DarkGray("]"), spaces,
				o.DarkGray("Generating desktop file..."))
		}

		if *windowmanager {
			writeWindowManagerDesktopFile(pkgname, name, execCommand, custom, *force, o)
		} else {
			writeDesktopFile(pkgname, name, comment, execCommand, *icon, *path, categories, genericName, mimeTypes, custom, *terminal, *startupnotify, *force, o)
		}

		if o.IsEnabled() {
			fmt.Printf("%s\n", o.DarkGreen("ok"))
		}

		// TODO: Refactor into a function
		// Download an icon if it's not downloaded by
		// the PKGBUILD and not there already (.png, .svg or .xpm)
		pngFilenames, _ := filepath.Glob("*.png")
		svgFilenames, _ := filepath.Glob("*.svg")
		xpmFilenames, _ := filepath.Glob("*.xpm")
		if (len(pngFilenames)+len(svgFilenames)+len(xpmFilenames) == 0) && !*nodownload {
			if len(pkgname) < 1 {
				o.Err("No pkgname, can't download icon")
			}
			fmt.Printf("%s%s%s%s%s ",
				o.DarkGray("["), o.LightBlue(pkgname),
				o.DarkGray("]"), spaces,
				o.DarkGray("Downloading icon..."))
			var err error
			if manualIconurl == "" {
				err = WriteIconFile(pkgname, o, *force)
			} else {
				// Default filename
				iconFilename := pkgname + ".png"
				// Get the last part of the URL, after the "/" to use as the filename
				if strings.Contains(manualIconurl, "/") {
					pos := strings.LastIndex(manualIconurl, "/")
					iconFilename = manualIconurl[pos+1:]
				}
				MustDownloadFile(manualIconurl, iconFilename, o, *force)
			}
			if err == nil {
				if o.IsEnabled() {
					fmt.Printf("%s\n", o.LightCyan("ok"))
				}
			} else {
				if o.IsEnabled() {
					fmt.Printf("%s\n", o.DarkYellow("no"))
					fmt.Printf("%s%s%s%s%s ",
						o.DarkGray("["),
						o.LightBlue(pkgname),
						o.DarkGray("]"),
						spaces,
						o.DarkGray("Using default icon instead..."))
				}
				if err := WriteDefaultIconFile(pkgname, o); (err == nil) && o.IsEnabled() {
					fmt.Printf("%s\n", o.LightPurple("yes"))
				}
			}
		}
	}
}
