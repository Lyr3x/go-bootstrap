package main

import (
	"flag"
	"os"
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

type DockerInfo struct {
	Appdir     string
	Entrypoint string
	Expose     string
}

func main() {

	InitLogger()
	defer log.Sync()
	expose := flag.String("expose", "8080", "Port to expose in docker")
	appname = flag.String("appname", "hello-world", "Name of the application")

	flag.Parse()

	goPath := os.Getenv("GOPATH")
	appdir = goPath + "/" + *appname
	cmddir = appdir + "/cmd"

	dockerInfo := DockerInfo{
		Appdir:     appdir,
		Entrypoint: *appname,
		Expose:     *expose,
	}
	setupStrcuture()
	generateDockerfile(dockerInfo)

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
	log.Info("Generated main.go at ", cmddir)
}
func generateDockerfile(dockerInfo DockerInfo) {
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

	t.Execute(f, dockerInfo)
	var command = "docker build -t " + dockerInfo.Entrypoint + "."
	log.Infow("Dockerfile generated, you can build the image with:",
		"command", command)
}
