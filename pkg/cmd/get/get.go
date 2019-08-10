package get

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewCmdCfgGet(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "get SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	cmd.AddCommand(NewCmdCfgGetCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgGetContext(streams, configAccess))
	cmd.AddCommand(NewCmdCfgGetUser(streams, configAccess))
	return cmd
}
