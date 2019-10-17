package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/denysvitali/ca-combos-editor/pkg"
	"io/ioutil"
	"log"
	"os"
)

func main(){
	pkg.Log.Level = logrus.DebugLevel
	args := os.Args[1:]
	result, err := ioutil.ReadFile(args[0])

	if err != nil {
		log.Fatal(err)
	}

	ce := pkg.NewComboEdit(result)
	cf := ce.Parse()

	for _, e := range cf.Entries {
		logrus.Info("Entry " + fmt.Sprintf("%s: %v", e.Name(),
			e.Bands()))
	}

}