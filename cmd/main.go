package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/denysvitali/ca-combos-editor/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel string
	mode     int

	rootCmd = &cobra.Command{
		Use:   "ca-combos-editor",
		Short: "Editor for Qualcomm NV ITEM 00028874 carrier aggregation combos",
		Long: `A combo editor for NV ITEM 00028874 (RFNV_LTE_CA_BW_CLASS_COMBO_I).

The tool can parse uncompressed 00028874 payloads and create new payloads from
human-readable band descriptions.`,
		SilenceUsage: true,
	}

	createCmd = &cobra.Command{
		Use:   "create <input> <output>",
		Short: "Create a 00028874 payload from a single bands file",
		Args:  cobra.ExactArgs(2),
		RunE:  runCreate,
	}

	createDlUlCmd = &cobra.Command{
		Use:   "create-dlul <downlink> <uplink> <output>",
		Short: "Create a 00028874 payload from separate downlink and uplink files",
		Args:  cobra.ExactArgs(3),
		RunE:  runCreateDlUl,
	}

	parseCmd = &cobra.Command{
		Use:   "parse <input>",
		Short: "Parse an extracted (uncompressed) 00028874 file",
		Args:  cobra.ExactArgs(1),
		RunE:  runParse,
	}

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	styleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	styleError = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87"))
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().IntVarP(&mode, "mode", "m", 137, "Writer mode: 137 or 201")

	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))

	viper.SetEnvPrefix("CA_COMBOS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.AddCommand(createCmd, createDlUlCmd, parseCmd)
}

func initConfig() {
	viper.SetConfigName("ca-combos-editor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/ca-combos-editor")

	// Intentionally ignore missing config files.
	_ = viper.ReadInConfig()

	configureLogger()
}

func configureLogger() {
	lvl := strings.ToLower(viper.GetString("log_level"))
	switch lvl {
	case "debug":
		pkg.Log.Level = logrus.DebugLevel
	case "info":
		pkg.Log.Level = logrus.InfoLevel
	case "warn", "warning":
		pkg.Log.Level = logrus.WarnLevel
	case "error", "err":
		pkg.Log.Level = logrus.ErrorLevel
	default:
		pkg.Log.Level = logrus.ErrorLevel
	}
}

func writerMode() pkg.ComboWriterMode {
	m := viper.GetInt("mode")
	switch m {
	case 137:
		return pkg.COMBOWRITER_137_138
	case 201:
		return pkg.COMBOWRITER_201_202
	default:
		return pkg.COMBOWRITER_137_138
	}
}

func runCreate(cmd *cobra.Command, args []string) error {
	input, output := args[0], args[1]

	entries, err := pkg.ParseBandFile(input)
	if err != nil {
		return err
	}

	if err := pkg.WriteComboFile(entries, writerMode(), output); err != nil {
		return err
	}

	fmt.Println(styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func runCreateDlUl(cmd *cobra.Command, args []string) error {
	downlink, uplink, output := args[0], args[1], args[2]

	entries, err := pkg.ParseBandDLULFile(downlink, uplink)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Name() == "DL" {
			pkg.Log.Debugf("\n")
		}
		pkg.Log.Debugf("%s: %s\n", e.Name(), e)
	}

	if err := pkg.WriteComboFile(entries, writerMode(), output); err != nil {
		return err
	}

	fmt.Println(styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func runParse(cmd *cobra.Command, args []string) error {
	input := args[0]

	fmt.Println(styleHeader.Render("Carrier Aggregation Combos"))
	fmt.Println()

	if err := pkg.ReadComboFile(input); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, styleError.Render(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}
}
