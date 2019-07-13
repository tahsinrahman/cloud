package cmds

import (
	"github.com/appscode/go/term"
	"github.com/spf13/cobra"
	"pharmer.dev/cloud/pkg/providers"
)

func NewCmdGenData() *cobra.Command {
	opts := providers.NewOptions()
	cmd := &cobra.Command{
		Use:               "gendata",
		Short:             "Load Kubernetes cluster data for a given cloud provider",
		Example:           "",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.ValidateFlags(cmd, args); err != nil {
				term.Fatalln(err)
			}
			cloudProvider, err := providers.NewCloudProvider(*opts)
			if err != nil {
				term.Fatalln(err)
			}
			err = providers.MergeAndWriteCloudProvider(cloudProvider)
			if err != nil {
				term.Fatalln(err)
			} else {
				term.Successln("Data successfully written for ", opts.Provider)
			}
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}
