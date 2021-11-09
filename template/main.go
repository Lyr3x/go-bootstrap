/// 2>/dev/null; exec gorun "$0" "$@"
package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var logger *zap.Logger
var log *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
}

var dirname string
var namepattern string

func main() {
	InitLogger()
	defer log.Sync()

	app := cli.NewApp()
	app.Name = "{{.Appname}}"
	app.Usage = "Define what your app does"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "findfile",
			Value:       ".",
			Usage:       "Find file in given dir",
			Required:    false,
			Destination: &dirname,
		},
		&cli.StringFlag{
			Name:        "pattern",
			Usage:       "file name pattern",
			Required:    false,
			Destination: &namepattern,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "pwd",
			Aliases: []string{"p"},
			Usage:   "Lists the current directory",
			Action: func(c *cli.Context) error {
				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				log.Info("Currnt director: ", dir)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		list, err := findFile(dirname, namepattern)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Found file: ", list)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func findFile(root, pattern string) ([]string, error) {
	log.Info("Find file in dir", zap.String("dir", root), zap.String("pattern", pattern))
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
