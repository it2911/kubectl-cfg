package get

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdCfgGet(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "get SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdCfgGetCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgGetContext(streams, configAccess))
	cmd.AddCommand(NewCmdCfgGetUser(streams, configAccess))
	return cmd
}
