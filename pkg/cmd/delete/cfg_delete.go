package delete

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdCfgDelete(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "delete SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdCfgDeleteCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgDeleteContext(streams, configAccess))
	cmd.AddCommand(NewCmdCfgDeleteUser(streams, configAccess))
	return cmd
}
