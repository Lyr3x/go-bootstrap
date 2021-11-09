package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

const (
	buildTemplate = `# syntax=docker/dockerfile:1

	FROM golang:1.16-alpine
	
	WORKDIR {{.Appdir}}
	
	COPY go.mod ./
	COPY go.sum ./
	RUN go mod download
	
	COPY *.go ./
	
	RUN go build -o {{.Entrypoint}}
	
	EXPOSE {{.Expose}}
	
	ENTRYPOINT ["{{.Entrypoint}}"]
`
)

var appPath *string
var appName *string

type DockerInfo struct {
	Appdir     string
	Entrypoint string
	Expose     string
}

func main() {

	expose := flag.String("expose", "8080", "Port to expose in docker")
	appPath = flag.String("apppath", "", "Path to the application")
	appName = flag.String("name", "", "Name of the application")

	flag.Parse()

	goPath := os.Getenv("GOPATH")
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	appdir := strings.Replace(dir, goPath, "", 1)

	_, entrypoint := path.Split(appdir)

	dockerInfo := DockerInfo{
		Appdir:     appdir,
		Entrypoint: entrypoint,
		Expose:     *expose,
	}

	generateDockerfile(dockerInfo)

}
func setupStrcuture() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Setting up structure")

	os.Mkdir(*appPath+"./cmd", 0755)
	os.Create(*appPath + "./cmd/main.go")
}
func generateDockerfile(dockerInfo DockerInfo) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	t := template.Must(template.New("buildTemplate").Parse(buildTemplate))

	f, err := os.Create("Dockerfile")
	if err != nil {
		sugar.Infof("Error wrinting Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, dockerInfo)
	var command = "docker build -t " + dockerInfo.Entrypoint + "."
	sugar.Infow("Dockerfile generated, you can build the image with:",
		"command", command)
}
