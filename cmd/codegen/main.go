package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime/debug"
	"slices"
	"strings"
	"text/template"
)

//go:embed templates
var templates embed.FS

type Config struct {
	WithPagination bool
	WithQueryOnly  bool
	WithOptConcur  bool
	WithSoftDelete bool
	OneToManyOnly  bool
	NewModelName   string
	NewModelPlural string
	NeedCUTime     bool
}

func main() {
	flag.Parse()

	name := flag.Arg(0)
	if name == "" {
		fmt.Println("missing domain name")
		os.Exit(1)
	}

	abbr := flag.Arg(1)
	if abbr == "" {
		fmt.Println("missing abbr name")
		os.Exit(1)
	}

	plural := flag.Arg(2)
	if plural == "" {
		fmt.Println("missing plural name")
		os.Exit(1)
	}

	cutime := flag.Arg(3)
	needCUTime := true
	if cutime == "n" || cutime == "N" {
		needCUTime = false
	}

	pagArg := flag.Arg(4)
	pag := true
	if pagArg == "n" || pagArg == "N" {
		pag = false
	}
	qryOnly := false
	if pagArg == "q" || pagArg == "Q" {
		qryOnly = true
	}

	ocArg := flag.Arg(5)
	oc := true
	if ocArg == "n" || ocArg == "N" {
		oc = false
	}

	sfArg := flag.Arg(6)
	sf := true
	if sfArg == "n" || sfArg == "N" {
		sf = false
	}

	oneToManyOnly := flag.Arg(7)
	otm := true
	if oneToManyOnly == "n" || oneToManyOnly == "N" {
		otm = false
	}
	var otmn string
	if otm {
		otmn = flag.Arg(8)
	} else {
		otmn = ""
	}
	var otmnp string
	if otm {
		otmnp = flag.Arg(9)
	} else {
		otmnp = ""
	}
	cfg := Config{
		WithPagination: pag,
		WithQueryOnly:  qryOnly,
		WithOptConcur:  oc,
		WithSoftDelete: sf,
		OneToManyOnly:  otm,
		NewModelName:   otmn,
		NewModelPlural: otmnp,
		NeedCUTime:     needCUTime,
	}

	if err := run(name, abbr, plural, cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(name string, abbr string, plural string, cfg Config) error {
	info, _ := debug.ReadBuildInfo()
	mod := info.Main.Path

	if err := addAppLayer(mod, name, abbr, plural, cfg); err != nil {
		return fmt.Errorf("adding app layer files: %w", err)
	}

	if err := addBusinessLayer(mod, name, abbr, plural, cfg); err != nil {
		return fmt.Errorf("adding bus layer files: %w", err)
	}

	if err := addStorageLayer(mod, name, abbr, plural, cfg); err != nil {
		return fmt.Errorf("adding sto layer files: %w", err)
	}

	if err := fmtCode(); err != nil {
		return fmt.Errorf("formatting code: %w", err)
	}

	fmt.Println("Done")
	return nil
}

func addAppLayer(mod, domain string, abbr string, plural string, cfg Config) error {
	const basePath = "internal/app/domain"
	var dir string
	if cfg.OneToManyOnly {
		dir = "templates/onetomany/app"
	} else {
		dir = "templates/app"
	}
	app, err := fs.Sub(templates, dir)
	if err != nil {
		return fmt.Errorf("switching to template/app folder: %w", err)
	}

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		return walkWork(mod, domain, abbr, plural, basePath, app, fileName, dirEntry, err, cfg)
	}

	fmt.Println("=======================================================")
	fmt.Println("APP LAYER CODE")

	if err := fs.WalkDir(app, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

func addBusinessLayer(mod, domain string, abbr string, plural string, cfg Config) error {
	const basePath = "internal/business/domain"
	var dir string
	if cfg.OneToManyOnly {
		dir = "templates/onetomany/business"
	} else {
		dir = "templates/business"
	}
	app, err := fs.Sub(templates, dir)
	if err != nil {
		return fmt.Errorf("switching to template/business folder: %w", err)
	}

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		return walkWork(mod, domain, abbr, plural, basePath, app, fileName, dirEntry, err, cfg)
	}

	fmt.Println("=======================================================")
	fmt.Println("BUSINESS LAYER CODE")

	if err := fs.WalkDir(app, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

func addStorageLayer(mod, domain string, abbr string, plural string, cfg Config) error {
	basePath := fmt.Sprintf("internal/business/domain/%s/stores", strings.ToLower(domain))
	var dir string
	if cfg.OneToManyOnly {
		dir = "templates/onetomany/storage"
	} else {
		dir = "templates/storage"
	}
	app, err := fs.Sub(templates, dir)
	if err != nil {
		return fmt.Errorf("switching to template/storage folder: %w", err)
	}

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		return walkWork(mod, domain, abbr, plural, basePath, app, fileName, dirEntry, err, cfg)
	}

	fmt.Println("=======================================================")
	fmt.Println("STORAGE LAYER CODE")

	if err := fs.WalkDir(app, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

func walkWork(mod string, domain string, abbr string, plural string, basePath string, app fs.FS, fileName string, dirEntry fs.DirEntry, err error, cfg Config) error {
	if err != nil {
		return fmt.Errorf("walkdir failure: %w", err)
	}

	if dirEntry.IsDir() {
		return nil
	}

	if cfg.WithQueryOnly && excludeNonQueryOnlyFiles(fileName) {
		return nil
	}

	f, err := app.Open(fileName)
	if err != nil {
		return fmt.Errorf("opening %s: %w", fileName, err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading %s: %w", fileName, err)
	}

	tmpl := template.Must(template.New("code").Parse(string(data)))

	domainVar := abbr
	domainPlural := plural
	tempDomainVar := strings.ToLower(domainVar)
	fields := projectFieldsz()
	coreFields := coreFieldsz(fields)
	jsonFields := jsonFieldsz(fields)
	dbFields := dbFieldsz(fields, cfg.OneToManyOnly)
	toAppStructFields := toAppStructFieldsz(fields, tempDomainVar)
	toCoreNewStructFields := toCoreNewStructFieldsz(fields)
	toCoreUpdateStructFields := toCoreUpdateStructFieldsz(fields)
	toDBStructFields := toDBStructFieldsz(fields, tempDomainVar)
	toCoreStructFields := toCoreStructFieldsz(fields, tempDomainVar)
	coreCreateFunction := coreCreateFunctionz(fields, tempDomainVar)
	coreUpdateFunction := coreUpdateFunctionz(fields, tempDomainVar)
	patchPointerTypeFieldsz := patchPointerTypeFields(fields)
	appUpdateStructFields := jsonFieldsz(patchPointerTypeFieldsz)
	coreUpdateStructFields := coreFieldsz(patchPointerTypeFieldsz)

	var otm string
	var otmp string
	if cfg.NewModelName != "" {
		otm = strings.ToUpper(cfg.NewModelName[0:1]) + cfg.NewModelName[1:]
	} else {
		otm = ""
	}

	if cfg.NewModelPlural != "" {
		otmp = strings.ToUpper(cfg.NewModelPlural[0:1]) + cfg.NewModelPlural[1:]
	} else {
		otmp = ""
	}

	d := struct {
		Module                   string
		DomainL                  string
		DomainU                  string
		DomainVar                string
		DomainVarU               string
		DomainVars               string
		DomainVarsU              string
		DomainNewVar             string
		DomainUpdVar             string
		DomainPlural             string
		DomainPluralU            string
		CoreFields               []string
		JSONFields               []string
		DBFields                 []string
		ToAppStructFields        []string
		ToCoreNewStructFields    []string
		ToCoreUpdateStructFields []string
		ToDBStructFields         []string
		ToCoreStructFields       []string
		CoreCreateFunction       []string
		CoreUpdateFunction       []string
		AppUpdateStructFields    []string
		CoreUpdateStructFields   []string
		DBPrefix                 string
		DBTableName              string
		OneToManyNewModelName    string
		OneToManyNewModelNameL   string
		OneToManyNewModelPlural  string
		OneToManyNewModelPluralL string
		// Options
		Config
	}{
		Module:                   mod,
		DomainL:                  strings.ToLower(domain),
		DomainU:                  strings.ToUpper(domain[0:1]) + domain[1:],
		DomainVar:                strings.ToLower(domainVar),
		DomainVarU:               strings.ToUpper(domainVar[0:1]) + strings.ToLower(domainVar[1:]),
		DomainVars:               strings.ToLower(domainVar) + "s",
		DomainVarsU:              strings.ToUpper(domainVar[0:1]) + strings.ToLower(domainVar[1:]) + "s",
		DomainNewVar:             "n" + strings.ToUpper(domainVar[0:1]) + strings.ToLower(domainVar[1:]),
		DomainUpdVar:             "u" + strings.ToUpper(domainVar[0:1]) + strings.ToLower(domainVar[1:]),
		DomainPlural:             strings.ToLower(domainPlural),
		DomainPluralU:            strings.ToUpper(domainPlural[0:1]) + domainPlural[1:],
		CoreFields:               coreFields,
		JSONFields:               jsonFields,
		DBFields:                 dbFields,
		ToAppStructFields:        toAppStructFields,
		ToCoreNewStructFields:    toCoreNewStructFields,
		ToCoreUpdateStructFields: toCoreUpdateStructFields,
		ToDBStructFields:         toDBStructFields,
		ToCoreStructFields:       toCoreStructFields,
		CoreCreateFunction:       coreCreateFunction,
		CoreUpdateFunction:       coreUpdateFunction,
		AppUpdateStructFields:    appUpdateStructFields,
		CoreUpdateStructFields:   coreUpdateStructFields,
		DBPrefix:                 dbPrefix,
		DBTableName:              dbTableName,
		OneToManyNewModelName:    otm,
		OneToManyNewModelNameL:   strings.ToLower(otm),
		OneToManyNewModelPlural:  otmp,
		OneToManyNewModelPluralL: strings.ToLower(otmp),
		// Options
		Config: cfg,
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, d); err != nil {
		return err
	}

	if err := writeFile(basePath, domain, fileName, b, cfg.NewModelName); err != nil {
		return fmt.Errorf("writing %s: %w", fileName, err)
	}

	return nil
}

func writeFile(basePath string, domain string, fileName string, b bytes.Buffer, oneToManyNewModelName string) error {
	path := basePath
	// lowerCaseFile
	parts := strings.SplitN(basePath, "/", 3)
	switch {
	case strings.HasSuffix(basePath, "stores"):
		path = fmt.Sprintf("%s/%sdb", basePath, strings.ToLower(domain))
	case parts[1] == "app":
		path = fmt.Sprintf("%s/%sapi", basePath, strings.ToLower(domain))
	case parts[1] == "business":
		path = fmt.Sprintf("%s/%s", basePath, strings.ToLower(domain))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Creating directory:", path)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("write app directory: %w", err)
		}
	}

	// Remove the suffix `.tmpl` from the file name.
	path = fmt.Sprintf("%s/%s", path, fileName[:len(fileName)-5])
	if oneToManyNewModelName != "" {
		path = strings.Replace(path, "new", oneToManyNewModelName, 1)
	} else {
		path = strings.Replace(path, "new", strings.ToLower(domain), 1)
	}
	fmt.Println("Add file:", path)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	fmt.Println("Writing code:", path)
	if _, err := f.Write(b.Bytes()); err != nil {
		return fmt.Errorf("writing bytes: %w", err)
	}

	return nil
}

func fmtCode() error {
	fmt.Println("=======================================================")
	fmt.Println("Formatting code...")
	if err := exec.Command("make", "fmt").Run(); err != nil {
		fmt.Println("command: make fmt: ", err)
	}

	return nil
}

// excludeNonQueryOnlyFiles returns a list of files that should not be generated when
// the query only flag is set.
func excludeNonQueryOnlyFiles(fileName string) bool {
	nonQueryOnlyFiles := []string{
		"filter.go.tmpl",
		"mid.go.tmpl",
		"order.go.tmpl",
		"new_test.go.tmpl",
	}

	return slices.Contains(nonQueryOnlyFiles, fileName)
}
