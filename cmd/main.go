package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/denysvitali/ca-combos-editor/pkg"
	"github.com/denysvitali/ca-combos-editor/pkg/zlib"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := logrus.ParseLevel(viper.GetString("log_level")); err != nil {
				return fmt.Errorf("invalid log-level %q", viper.GetString("log_level"))
			}
			if _, err := comboWriterMode(viper.GetInt("mode")); err != nil {
				return err
			}
			return nil
		},
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

	compressCmd = &cobra.Command{
		Use:   "compress <input> <output>",
		Short: "Zlib-compress a raw uncompressed payload",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompress,
	}

	decompressCmd = &cobra.Command{
		Use:   "decompress <input> <output>",
		Short: "Decompress a 00028874 zlib-compressed file",
		Args:  cobra.ExactArgs(2),
		RunE:  runDecompress,
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
	rootCmd.PersistentFlags().IntVarP(&mode, "mode", "m", 137, "Writer mode: 137, 201, or 333")

	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))

	viper.SetEnvPrefix("CA_COMBOS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.AddCommand(createCmd, createDlUlCmd, parseCmd, compressCmd, decompressCmd)
}

func initConfig() {
	viper.SetConfigName("ca-combos-editor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/ca-combos-editor")

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			pkg.Log.WithError(err).Warn("failed to read config file")
		}
	}

	configureLogger()
}

func configureLogger() {
	lvl, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		pkg.Log.WithError(err).Warn("invalid log level, defaulting to error")
		lvl = logrus.ErrorLevel
	}
	pkg.Log.SetLevel(lvl)
}

func comboWriterMode(mode int) (pkg.ComboWriterMode, error) {
	switch mode {
	case 137:
		return pkg.COMBOWRITER_137_138, nil
	case 201:
		return pkg.COMBOWRITER_201_202, nil
	case 333:
		return pkg.COMBOWRITER_333_334, nil
	default:
		return 0, fmt.Errorf("invalid mode %d: must be 137, 201, or 333", mode)
	}
}

func runCreate(cmd *cobra.Command, args []string) error {
	input, output := args[0], args[1]

	entries, err := pkg.ParseBandFile(input)
	if err != nil {
		return fmt.Errorf("parse band file %q: %w", input, err)
	}

	mode, err := comboWriterMode(viper.GetInt("mode"))
	if err != nil {
		return err
	}

	if err := pkg.WriteComboFile(entries, mode, output); err != nil {
		return fmt.Errorf("write combo file %q: %w", output, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func runCreateDlUl(cmd *cobra.Command, args []string) error {
	downlink, uplink, output := args[0], args[1], args[2]

	entries, err := pkg.ParseBandDLULFile(downlink, uplink)
	if err != nil {
		return fmt.Errorf("parse DL/UL band files: %w", err)
	}

	for _, e := range entries {
		pkg.Log.Debugf("%s: %s", e.Name(), e)
	}

	mode, err := comboWriterMode(viper.GetInt("mode"))
	if err != nil {
		return err
	}

	if err := pkg.WriteComboFile(entries, mode, output); err != nil {
		return fmt.Errorf("write combo file %q: %w", output, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func runParse(cmd *cobra.Command, args []string) error {
	input := args[0]

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), styleHeader.Render("Carrier Aggregation Combos"))
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	if err := pkg.ReadComboFile(input, cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("read combo file %q: %w", input, err)
	}
	return nil
}

func runCompress(cmd *cobra.Command, args []string) error {
	input, output := args[0], args[1]

	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open input file %q: %w", input, err)
	}
	defer func() { _ = f.Close() }()

	compressed, err := zlib.Compress(f)
	if err != nil {
		return fmt.Errorf("compress %q: %w", input, err)
	}

	if err := os.WriteFile(output, compressed, 0o644); err != nil {
		return fmt.Errorf("write output file %q: %w", output, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func runDecompress(cmd *cobra.Command, args []string) error {
	input, output := args[0], args[1]

	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open input file %q: %w", input, err)
	}
	defer func() { _ = f.Close() }()

	decompressed, err := zlib.Decompress(f)
	if err != nil {
		return fmt.Errorf("decompress %q: %w", input, err)
	}

	if err := os.WriteFile(output, decompressed, 0o644); err != nil {
		return fmt.Errorf("write output file %q: %w", output, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), styleSuccess.Render(fmt.Sprintf("✓ wrote %s", output)))
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, styleError.Render(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}
}
