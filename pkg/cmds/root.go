package cmds

import (
	"flag"

	"github.com/appscode/go/flags"
	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
	"pharmer.dev/cloud/pkg/apis"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "pharmer-tools",
		Short:             `Pharmer tools by Appscode`,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			flags.DumpAll(c.Flags())
			cli.SendAnalytics(c, v.Version.Version)
			return apis.AddToScheme(scheme.Scheme)
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	logs.ParseFlags()
	rootCmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "Send analytical events to Google Analytics")

	rootCmd.AddCommand(NewCmdGenData())
	rootCmd.AddCommand(NewCmdKubeSupport())
	rootCmd.AddCommand(v.NewCmdVersion())

	return rootCmd
}
