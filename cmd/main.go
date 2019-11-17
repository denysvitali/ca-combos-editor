package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/ca-combos-editor/pkg"
)

type CreateCmd struct {
	Input string `arg:"positional"`
	Output string `arg:"positional"`
}

type Create2Cmd struct {
	Downlink string `arg:"positional"`
	Uplink string `arg:"positional"`
	Output string `arg:"positional"`
}

type ParseCmd struct {
	Input string `arg:"positional"`
}

var args struct {
	Create *CreateCmd `arg:"subcommand:create"`
	Create2 *Create2Cmd `arg:"subcommand:create2"`
	Parse *ParseCmd `arg:"subcommand:parse"`
}

func main(){
	arg.MustParse(&args)
	pkg.Log.Level = logrus.DebugLevel

	switch {
	case args.Create != nil:
		entries := pkg.ParseBandFile(args.Create.Input)
		pkg.WriteComboFile(entries, args.Create.Output)
	case args.Create2 != nil:
		entries := pkg.ParseBandDLULFile(args.Create2.Downlink, args.Create2.Uplink)
		for _, e := range entries {
			if e.Name() == "DL" {
				pkg.Log.Debugf("\n")
			}
			pkg.Log.Debugf("%s: %s\n", e.Name(), e)
		}
		pkg.WriteComboFile(entries, args.Create2.Output)
	case args.Parse != nil:
		pkg.ReadComboFile(args.Parse.Input)
	default:

	}

}