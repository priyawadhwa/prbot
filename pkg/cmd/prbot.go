package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/priyawadhwa/prbot/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "prbot",
	Short: "prbot helps you set up a prbot for your Github repo",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPRBot(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "prbot.yaml", "path to prbot config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runPRBot() error {
	c, err := config.Get(cfgFile)
	if err != nil {
		return errors.Wrap(err, "getting config file")
	}
	fmt.Println(c)
	return nil
}
