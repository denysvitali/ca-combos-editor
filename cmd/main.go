package main

import (
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/ca-combos-editor/pkg"
	"github.com/sirupsen/logrus"
	"strings"
)

type CreateCmd struct {
	Input  string              `arg:"positional" help:"Input file"`
	Output string              `arg:"positional" help:"Output file"`
	Mode   pkg.ComboWriterMode `arg:"-m" help:"Mode to use: 137 or 201" default:"137"`
}

type Create2Cmd struct {
	Downlink string              `arg:"positional" help:"Downlink file"`
	Uplink   string              `arg:"positional" help:"Uplink file"`
	Output   string              `arg:"positional" help:"Output file"`
	Mode     pkg.ComboWriterMode `arg:"-m" help:"Mode to use: 137 or 201" default:"137"`
}

type ParseCmd struct {
	Input string `arg:"positional"`
}

var args struct {
	Create     *CreateCmd  `arg:"subcommand:create" help:"Create a file using a single file which contains the combos"`
	CreateDlUl *Create2Cmd `arg:"subcommand:create_dlul" help:"Create a file using a downlink.txt file and an uplink.txt"`
	Parse      *ParseCmd   `arg:"subcommand:parse" help:"Parse an extracted 00028874 file"`
	LogLevel   string      `arg:"-l" help:"log level (Debug, Info, Warn, Err)"`
}

func main() {
	arg.MustParse(&args)

	switch strings.ToLower(args.LogLevel) {
	case "debug":
		pkg.Log.Level = logrus.DebugLevel
	case "info":
		pkg.Log.Level = logrus.InfoLevel
	case "warn":
		pkg.Log.Level = logrus.WarnLevel
	case "err":
		pkg.Log.Level = logrus.ErrorLevel
	default:
		pkg.Log.Level = logrus.ErrorLevel
	}

	switch {
	case args.Create != nil:
		entries := pkg.ParseBandFile(args.Create.Input)
		pkg.WriteComboFile(entries, args.Create.Mode, args.Create.Output)
	case args.CreateDlUl != nil:
		entries := pkg.ParseBandDLULFile(args.CreateDlUl.Downlink, args.CreateDlUl.Uplink)
		for _, e := range entries {
			if e.Name() == "DL" {
				pkg.Log.Debugf("\n")
			}
			pkg.Log.Debugf("%s: %s\n", e.Name(), e)
		}
		pkg.WriteComboFile(entries, args.CreateDlUl.Mode, args.CreateDlUl.Output)
	case args.Parse != nil:
		pkg.ReadComboFile(args.Parse.Input)
	default:

	}

}
