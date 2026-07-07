package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetRootCmd(t *testing.T) {
	t.Helper()
	rootCmd.SetArgs(nil)
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	require.NoError(t, rootCmd.PersistentFlags().Set("log-level", "error"))
	require.NoError(t, rootCmd.PersistentFlags().Set("mode", "137"))
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() { os.Stdout = old }()
	fn()
	require.NoError(t, w.Close())
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	return buf.String()
}

func execute(t *testing.T, args ...string) string {
	t.Helper()
	resetRootCmd(t)
	rootCmd.SetArgs(args)
	out := captureStdout(t, func() {
		require.NoError(t, rootCmd.Execute())
	})
	return out
}

func executeExpectErr(t *testing.T, args ...string) string {
	t.Helper()
	resetRootCmd(t)
	rootCmd.SetArgs(args)
	var out strings.Builder
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	err := rootCmd.Execute()
	require.Error(t, err)
	return out.String()
}

func TestRootCommand(t *testing.T) {
	outDir := t.TempDir()

	tests := []struct {
		name   string
		args   []string
		assert func(t *testing.T, out string)
	}{
		{
			name: "parse subcommand on fixture",
			args: []string{"parse", "../test/resources/2019-10-17/extracted"},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "Carrier Aggregation Combos")
				assert.Contains(t, out, "DL")
				assert.Contains(t, out, "UL")
			},
		},
		{
			name: "create subcommand from bands.txt",
			args: []string{"create", "../test/resources/2019-10-17/bands.txt", filepath.Join(outDir, "created")},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "wrote")
				assert.FileExists(t, filepath.Join(outDir, "created"))
			},
		},
		{
			name: "create-dlul subcommand from fixtures",
			args: []string{
				"create-dlul",
				"../test/resources/2019-11-15/iphone-xr-bands.txt",
				"../test/resources/2019-11-15/iphone-xr-bands-ul.txt",
				filepath.Join(outDir, "created-dlul"),
			},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "wrote")
				assert.FileExists(t, filepath.Join(outDir, "created-dlul"))
			},
		},
		{
			name: "--log-level flag",
			args: []string{"--log-level", "debug", "parse", "../test/resources/2019-10-17/extracted"},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "Carrier Aggregation Combos")
			},
		},
		{
			name: "--mode flag 201",
			args: []string{"--mode", "201", "create", "../test/resources/2019-10-17/bands.txt", filepath.Join(outDir, "created-mode-201")},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "wrote")
				assert.FileExists(t, filepath.Join(outDir, "created-mode-201"))
			},
		},
		{
			name: "--mode flag 137",
			args: []string{"--mode", "137", "create", "../test/resources/2019-10-17/bands.txt", filepath.Join(outDir, "created-mode-137")},
			assert: func(t *testing.T, out string) {
				assert.Contains(t, out, "wrote")
				assert.FileExists(t, filepath.Join(outDir, "created-mode-137"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := execute(t, tt.args...)
			tt.assert(t, out)
		})
	}
}

func TestRootCommandInvalidArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "create with too few args",
			args: []string{"create", "only-one-arg"},
		},
		{
			name: "create with too many args",
			args: []string{"create", "a", "b", "c"},
		},
		{
			name: "create-dlul with too few args",
			args: []string{"create-dlul", "a", "b"},
		},
		{
			name: "parse with too few args",
			args: []string{"parse"},
		},
		{
			name: "parse with too many args",
			args: []string{"parse", "a", "b"},
		},
		{
			name: "unknown subcommand",
			args: []string{"unknown"},
		},
		{
			name: "missing required input file for create",
			args: []string{"create", "does-not-exist.txt", filepath.Join(t.TempDir(), "out")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executeExpectErr(t, tt.args...)
		})
	}
}
