package add

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdCfgAdd(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "add SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	//cmd.AddCommand(NewCmdCfgAddCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgAddConfig(streams, configAccess))
	cmd.AddCommand(NewCmdCfgAddContext(streams.Out, configAccess))
	//cmd.AddCommand(NewCmdCfgAddUser(streams, configAccess))
	return cmd
}
