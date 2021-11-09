/// 2>/dev/null; exec gorun "$0" "$@"
package main

import (
	"io/ioutil"

	"go.uber.org/zap"
)

var logger *zap.Logger
var log *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
}
func main() {
	InitLogger()
	defer log.Sync()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%v", files)

}
