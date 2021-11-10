package main

import (
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var logger *zap.Logger
var log *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
}

var appdir string
var cmddir string
var example bool
var username string
var appname string

type templateInfo struct {
	Appdir   string
	Appname  string
	Expose   string
	Username string
}

func main() {

	InitLogger()
	defer log.Sync()
	app := cli.NewApp()
	app.Name = "go-bootstrap"
	app.Usage = "CLI tool to boot strap a fresh go project including build files"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "appname",
			Value:       "default",
			Usage:       "app name as string",
			Required:    true,
			Destination: &appname,
		},
		&cli.StringFlag{
			Name:        "username",
			Value:       "user",
			Usage:       "username as string",
			Required:    true,
			Destination: &username,
		},
		&cli.BoolFlag{
			Name:        "example",
			Usage:       "true: Full example, false: minimal working example ",
			Required:    false,
			Destination: &example,
		},
	}

	app.Action = func(c *cli.Context) error {
		goPath := os.Getenv("GOPATH")
		appdir = goPath + "/" + appname
		cmddir = appdir + "/cmd"

		templateInfo := templateInfo{
			Appdir:   appdir,
			Appname:  appname,
			Username: strings.ToLower(username),
		}
		setupStrcuture()
		generateGoFiles(templateInfo)
		generateDockerfile(templateInfo)
		generateMakefile(templateInfo)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func setupStrcuture() {

	log.Info("Setting up structure.")

	err := os.MkdirAll(appdir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created project directory ", appdir)
	err = os.MkdirAll(cmddir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Creatd cmd directory ", cmddir)

	os.Create(appdir + "/cmd/" + "main.go")
	if err != nil {
		log.Fatal(err)
	}
}
func generateGoFiles(templateInfo templateInfo) {
	var t *template.Template
	var err error
	if example == true {
		log.Info("Creating full example main.go")
		t, err = template.ParseFiles("template/main.go.example")
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Info("Creating simple example main.go")
		t, err = template.ParseFiles("template/main.go.clean")
		if err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(cmddir + "/main.go")
	if err != nil {
		log.Infof("Error writing main.go %v", err.Error())
		return
	}
	t.Execute(f, templateInfo)

	cmd := exec.Command("go", "mod", "init", appname)
	cmd.Dir = *&appdir

	err = cmd.Run()
	if err != nil {
		log.Infof("Error initializing go module %v", err.Error())
		return
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = *&appdir

	err = cmd.Run()
	if err != nil {
		log.Infof("Error tidying go module %v", err.Error())
		return
	}

}
func generateDockerfile(templateInfo templateInfo) {
	t, err := template.ParseFiles("template/Dockerfile")
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(appdir + "/Dockerfile")
	if err != nil {
		log.Infof("Error wrinting Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, templateInfo)
	var command = "docker build -t " + templateInfo.Appname + "."
	log.Infow("Dockerfile generated, you can build the image with:",
		"command", command)
}

func generateMakefile(templateInfo templateInfo) {
	t, err := template.ParseFiles("template/Makefile")
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(appdir + "/Makefile")
	if err != nil {
		log.Infof("Error wrinting Makefile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, templateInfo)
	log.Info("Makefile generated")
}
