# go-bootstrap
This tool will generate a basic opinionated project setup which will serve as a starting point for lightweight go tools as bash script replacement. 

## Usage

`go run main.go -apname <app-name> -expose <port-to-be-exposed>`

This will generate the directory structure at `$GOPATH/appname` including:
* Dockerfile
* Makefile
* Example `main.go`

The following libraries are used:
* `go.uber.org/zap` for logging
* `github.com/urfave/cli/v2` to build a cli
