package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/xyproto/term"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	version_string = "Desktop File Generator v.0.7"
)

// Information needed to generate a .desktop file for a Window Manager
type WMStarter struct {
	Name string
	Exec string
}

// Information needed to generate a .desktop file for an application
type AppStarter struct {
	Name          string
	GenericName   string
	Comment       string
	Exec          string
	Icon          string
	UseTerminal   bool
	StartupNotify bool
	CategoryList  string
	MimeTypesList string
}

var (
	// Template for a .desktop file for starting a Window Manager
	wmTemplate, _ = template.New("WMStarter").Parse("[Desktop Entry]\nType=XSession\nExec={{.Exec}}\nTryExec={{.Exec}}\nName={{.Name}}\n")

	// Template for a .desktop file for starting an application
	appTemplate, _ = template.New("AppStarter").Parse("[Desktop Entry]\nEncoding=UTF-8\nType=Application\nName={{.Name}}\n{{if .GenericName}}GenericName={{.GenericName}}\n{{end}}Comment={{.Comment}}\nExec={{.Exec}}\nIcon={{.Icon}}\nTerminal={{if .UseTerminal}}true{{else}}false{{end}}\nStartupNotify={{if .StartupNotify}}true{{else}}false{{end}}\nCategories={{.CategoryList}};\n{{if .MimeTypesList}}MimeType={{.MimeTypesList}};\n{{end}}")
)

// Generate the contents for the .desktop file (for executing a window manager)
func createWindowManagerDesktopContents(name, exec string) *bytes.Buffer {
	var (
		buf bytes.Buffer
		wm  = WMStarter{name, exec}
	)
	// Inserting strings into the template should always work, panic if not
	if err := wmTemplate.Execute(&buf, wm); err != nil {
		panic(err)
	}
	return &buf
}

// Generate the contents for the .desktop file (for executing a desktop application)
func createDesktopContents(name, genericName, comment, exec, icon string, useTerminal, startupNotify bool, categories, mimeTypes []string) *bytes.Buffer {
	var (
		buf bytes.Buffer
		app = AppStarter{name, genericName, comment, exec, icon, useTerminal, startupNotify, strings.Join(categories, ";"), strings.Join(mimeTypes, ";")}
	)
	// Inserting strings into the template should always work, panic if not
	if err := appTemplate.Execute(&buf, app); err != nil {
		panic(err)
	}
	return &buf
}

// Write the .desktop file as generated by createWindowManagerDesktopContents
func writeWindowManagerDesktopFile(pkgname, name, exec, custom string, force bool, o *term.TextOutput) {
	buf := createWindowManagerDesktopContents(name, exec)
	if custom != "" {
		// Write the custom string to the end of the .desktop file (may contain \n)
		buf.WriteString(custom + "\n")
	}
	// Check if the file exists (and that force is not enabled)
	if _, err := os.Stat(pkgname + ".desktop"); err == nil && (!force) {
		o.Err("no")
		o.Println(pkgname + ".desktop already exists. Use -f as the first argument to gendesk to overwrite.")
		os.Exit(1)
	}

	ioutil.WriteFile(pkgname+".desktop", buf.Bytes(), 0644)
}

// Write the .desktop file as generated by createDesktopContents
func writeDesktopFile(pkgname, name, comment, exec, categories, genericName, mimeTypes, custom string, useTerminal, startupNotify, force bool, o *term.TextOutput) {
	var (
		categoryList, mimeTypeList []string
	)

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
	buf := createDesktopContents(name, genericName, comment, exec, pkgname, useTerminal, startupNotify, categoryList, mimeTypeList)
	if custom != "" {
		// Write the custom string to the end of the .desktop file (may contain \n)
		buf.WriteString(custom + "\n")
	}

	// Check if the file exists (and that force is not enabled)
	if _, err := os.Stat(pkgname + ".desktop"); err == nil && (!force) {
		o.Err("no")
		o.Println(pkgname + ".desktop already exists. Use -f as the first argument to gendesk to overwrite.")
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
func WriteDefaultIconFile(pkgname string, o *term.TextOutput) error {
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

func main() {
	var filename string
	version_help := "Show application name and version"
	nodownload_help := "Don't download anything"
	nocolor_help := "Don't use colors"
	quiet_help := "Don't output anything on stdout"
	force_help := "Overwrite .desktop files with the same name"
	windowmanager_help := "Generate a .desktop file for launching a window manager"
	pkgname_help := "The name of the package"
	pkgdesc_help := "Description of the package"
	name_help := "Name of the shortcut"
	genericname_help := "Type of application"
	comment_help := "Shortcut comment"
	exec_help := "Path to executable"
	//iconurl_help := "URL to icon"
	terminal_help := "Run the application in a terminal (default is false)"
	categories_help := "Categories, see other .desktop files for examples"
	mimetypes_help := "Mime types, see other .desktop files for examples"
	startupnotify_help := "Notifcation when the application starts (default is false)"
	custom_help := "Custom line to append at the end of the .desktop file"

	flag.Usage = func() {
		fmt.Println("\n" + version_string)
		fmt.Println("Generate .desktop files.")
		fmt.Println("\nSyntax: gendesk [flags]")
		fmt.Println("\nPossible flags:")
		fmt.Println("    --version                    " + version_help)
		fmt.Println("    -n                           " + nodownload_help)
		fmt.Println("    --nocolor                    " + nocolor_help)
		fmt.Println("    -q                           " + quiet_help)
		fmt.Println("    -f                           " + force_help)
		fmt.Println("    -wm                          " + windowmanager_help)
		fmt.Println("    --pkgname=PKGNAME            " + pkgname_help)
		fmt.Println("    --pkgdesc=PKGDESC            " + pkgdesc_help)
		fmt.Println("    --name=NAME                  " + name_help)
		fmt.Println("    --genericname=GENERICNAME    " + genericname_help)
		fmt.Println("    --comment=COMMENT            " + comment_help)
		fmt.Println("    --exec=EXEC                  " + exec_help)
		fmt.Println("    --terminal=[true|false]      " + terminal_help)
		fmt.Println("    --categories=CATEGORIES      " + categories_help)
		fmt.Println("    --mimetypes=MIMETYPES        " + mimetypes_help)
		fmt.Println("    --startupnotify=[true|false] " + startupnotify_help)
		fmt.Println("    --custom=CUSTOM              " + custom_help)
		fmt.Println("    --help                       This text")
		fmt.Println("\nNote:")
		fmt.Println("    * Just providing --pkgname is enough to generate a .desktop file.")
		fmt.Println("    * Providing a PKGBUILD filename instead of flags is a possibility.")
		fmt.Println("    * \"$startdir/PKGBUILD\" is the default PKGBUILD filename.")
		fmt.Println("    * _exec in the PKGBUILD can be used to specifiy a different executable for the")
		fmt.Println("      .desktop file. Example: _exec=('appname-gui')")
		fmt.Println("    * Split PKGBUILD packages are supported.")
		fmt.Println("    * If a .png, .svg or .xpm icon is not found as a file or in the PKGBUILD,")
		fmt.Println("      an icon will be downloaded from either the location specified in the")
		shortname := strings.Split(defaultIconSearchURL, "/")
		firstpart := strings.Join(shortname[:3], "/")
		fmt.Println("      configuration or from: " + firstpart)
		fmt.Println("      (This may or may not result in the icon you wished for).")
		fmt.Println("    * Categories are guessed based on keywords in the")
		fmt.Println("      package description, unless provided.")
		fmt.Println("    * Icons are assumed to be found in \"/usr/share/pixmaps/\" once installed.")
		fmt.Println()
	}

	version := flag.Bool("version", false, version_help)
	nodownload := flag.Bool("n", false, nodownload_help)
	nocolor := flag.Bool("nocolor", false, nocolor_help)
	quiet := flag.Bool("q", false, quiet_help)
	force := flag.Bool("f", false, force_help)
	windowmanager := flag.Bool("wm", false, windowmanager_help)
	givenPkgname := flag.String("pkgname", "", pkgname_help)
	givenPkgdesc := flag.String("pkgdesc", "", pkgdesc_help)
	name := flag.String("name", "", name_help)
	genericname := flag.String("genericname", "", genericname_help)
	comment := flag.String("comment", "", comment_help)
	exec := flag.String("exec", "", exec_help)
	terminal := flag.Bool("terminal", false, terminal_help)
	categories := flag.String("categories", "", categories_help)
	mimetypes := flag.String("mimetypes", "", mimetypes_help)
	mimetype := flag.String("mimetype", "", mimetypes_help)
	custom := flag.String("custom", "", custom_help)
	startupnotify := flag.Bool("startupnotify", false, startupnotify_help)

	flag.Parse()
	args := flag.Args()

	// New text output. Color? Enabled?
	o := term.NewTextOutput(!*nocolor, !*quiet)

	if *version {
		o.Println(version_string)
		os.Exit(0)
	}

	pkgname := *givenPkgname
	pkgdesc := *givenPkgdesc
	manualIconurl := ""

	const defaultPKGBUILD = "../PKGBUILD"

	// TODO: Write in a cleaner way, possibly by refactoring into a function. Write a test first.
	if pkgname == "" {
		if len(args) == 0 {
			if os.Getenv("pkgname") == "" {
				if os.Getenv("SRCDEST") == "" {
					filename = defaultPKGBUILD
				} else {
					// If SRCDEST is set, use that
					filename = os.Getenv("SRCDEST") + "/PKGBUILD"
				}
			} else {
				pkgname = os.Getenv("pkgname")
			}
		} else if len(args) > 0 {
			// args are non-flag arguments
			filename = args[0]
		}
	}

	// Environment variables
	dataFromEnvironment(&pkgdesc, exec, name, genericname, mimetypes, comment, categories, custom)

	var pkgnames []string
	var iconurl string

	// Several fields are stored per pkgname, hence the maps
	pkgdescMap := make(map[string]string)
	execMap := make(map[string]string)
	nameMap := make(map[string]string)
	genericNameMap := make(map[string]string)
	mimeTypesMap := make(map[string]string)
	commentMap := make(map[string]string)
	categoriesMap := make(map[string]string)
	customMap := make(map[string]string)

	if filename == "" {
		// Fill in the dictionaries using the arguments
		pkgnames = []string{pkgname}
		if pkgdesc != "" {
			pkgdescMap[pkgname] = pkgdesc
		}
		if *exec != "" {
			execMap[pkgname] = *exec
		}
		if *name != "" {
			nameMap[pkgname] = *name
		}
		if *genericname != "" {
			genericNameMap[pkgname] = *genericname
		}
		if *mimetype != "" {
			mimeTypesMap[pkgname] = *mimetype
		}
		if *mimetypes != "" {
			mimeTypesMap[pkgname] = *mimetypes
		}
		if *comment != "" {
			commentMap[pkgname] = *comment
		}
		if *categories != "" {
			categoriesMap[pkgname] = *categories
		}
		if *custom != "" {
			customMap[pkgname] = *custom
		}
	} else {
		// Check if the PKGBUILD filename is found
		if _, err := os.Stat(filename); err != nil {
			if filename != defaultPKGBUILD {
				// Not the default filename, complain that the file is missing
				o.Err("Could not find " + filename + ", provide a --pkgname or a valid PKGBUILD file")
				os.Exit(1)
			} else {
				// Could not find the default filename, complain about missing arguments
				fmt.Println(o.LightBlue("Provide a package name with --pkgname, or a valid PKGBUILD file. Use --help for more info."))
				os.Exit(1)
			}
		}
		// TODO: Use a struct per pkgname instead
		parsePKGBUILD(o, filename, &iconurl, &pkgname, &pkgnames, &pkgdescMap, &execMap, &nameMap, &genericNameMap, &mimeTypesMap, &commentMap, &categoriesMap, &customMap)
	}

	// Write .desktop and .png icon for each package
	for _, pkgname := range pkgnames {
		if strings.Contains(pkgname, "-nox") || strings.Contains(pkgname, "-cli") {
			// Don't bother if it's a -nox or -cli package
			continue
		}
		// Strip the "-git" or "-svn" suffix, if present
		if strings.HasSuffix(pkgname, "-git") || strings.HasSuffix(pkgname, "-svn") {
			pkgname = pkgname[:len(pkgname)-4]
		}
		// TODO: Find a better way for all the if checks below
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
			writeWindowManagerDesktopFile(pkgname, name, exec, custom, *force, o)
		} else {
			writeDesktopFile(pkgname, name, comment, exec, categories, genericName, mimeTypes, custom, *terminal, *startupnotify, *force, o)
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
		if (0 == (len(pngFilenames) + len(svgFilenames) + len(xpmFilenames))) && !*nodownload {
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
				err = DownloadFile(manualIconurl, iconFilename, o, *force)
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
