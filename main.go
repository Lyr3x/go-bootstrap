package main

import (
	"flag"
	"os"
	"os/exec"
	"text/template"

	"go.uber.org/zap"
)

var logger *zap.Logger
var log *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
}

var appname *string
var appdir string
var cmddir string
var example *bool

type templateInfo struct {
	Appdir  string
	Appname string
	Expose  string
}

func main() {

	InitLogger()
	defer log.Sync()
	expose := flag.String("expose", "8080", "Port to expose in docker")
	appname = flag.String("appname", "hello-world", "Name of the application")
	example = flag.Bool("example", false, "Generate rich example project")

	flag.Parse()

	goPath := os.Getenv("GOPATH")
	appdir = goPath + "/" + *appname
	cmddir = appdir + "/cmd"

	templateInfo := templateInfo{
		Appdir:  appdir,
		Appname: *appname,
		Expose:  *expose,
	}
	setupStrcuture()
	generateGoFiles(templateInfo)
	generateDockerfile(templateInfo)

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
	if *example == true {
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

	cmd := exec.Command("go", "mod", "init", *appname)
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
