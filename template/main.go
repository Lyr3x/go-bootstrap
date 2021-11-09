/// 2>/dev/null; exec gorun "$0" "$@"
package main

import (
	"io/ioutil"
	"os"
)

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
}
