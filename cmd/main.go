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

type ParseCmd struct {
	Input string `arg:"positional"`
}

var args struct {
	Create *CreateCmd `arg:"subcommand:create"`
	Parse *ParseCmd `arg:"subcommand:parse"`
}

func main(){
	arg.MustParse(&args)
	pkg.Log.Level = logrus.DebugLevel

	switch {
	case args.Create != nil:
		entries := pkg.ParseBandFile(args.Create.Input)
		pkg.WriteComboFile(entries, args.Create.Output)
	case args.Parse != nil:
		pkg.ReadComboFile(args.Parse.Input)
	default:

	}

}