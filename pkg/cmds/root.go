package cmds

import (
	"flag"

	"github.com/appscode/go/flags"
	v "github.com/appscode/go/version"
	"github.com/pharmer/pharmer/config"
	"github.com/spf13/cobra"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "pharmer-tools",
		Short:             `Pharmer by Appscode - Manages farms`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			flags.DumpAll(c.Flags())
			cli.SendAnalytics(c, v.Version.Version)
		},
	}
	config.AddFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	logs.ParseFlags()
	rootCmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "Send analytical events to Google Analytics")

	rootCmd.AddCommand(NewCmdGenData())
	rootCmd.AddCommand(NewCmdKubeSupport())
	rootCmd.AddCommand(NewCmdGenNPM())
	rootCmd.AddCommand(v.NewCmdVersion())

	return rootCmd
}
