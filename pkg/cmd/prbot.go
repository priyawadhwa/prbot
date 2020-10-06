package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/priyawadhwa/prbot/pkg/config"
	"github.com/priyawadhwa/prbot/pkg/execute"
	"github.com/priyawadhwa/prbot/pkg/github"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "prbot",
	Short: "prbot helps you set up a prbot for your Github repo",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPRBot(context.Background()); err != nil {
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

func runPRBot(ctx context.Context) error {
	cfg, err := config.Get(cfgFile)
	if err != nil {
		return errors.Wrap(err, "getting config file")
	}
	prs, err := github.ListPRs(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "listing prs")
	}
	log.Printf("Got PRs: %v", prs)
	for _, pr := range prs {
		contents, err := execute.Execute(cfg, execute.NewConfig(pr))
		if err != nil {
			return errors.Wrap(err, "executing")
		}
		if err := github.Comment(ctx, cfg, contents); err != nil {
			return errors.Wrap(err, "commenting on github")
		}
	}
	return nil
}
