package delete

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	deleteLong = templates.LongDesc(`Delete context / cluster / authinfo from the kubeconfig file.`)

	deleteExample = templates.Examples(`
		# Delete the resource in your kubeconfig file.
		kubectl cfg delete SUB_COMMAND`)
)

func NewCmdCfgDelete(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	// TODO add backup
	cmd := &cobra.Command{
		Use:                   "delete",
		DisableFlagsInUseLine: true,
		Short:                 "Delete context / cluster / authinfo from kubeconfig",
		Long:                  deleteLong,
		Example:               deleteExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	cmd.AddCommand(NewCmdCfgDeleteCluster(streams.Out, streams.ErrOut, configAccess))
	cmd.AddCommand(NewCmdCfgDeleteContext(streams.Out, streams.ErrOut, configAccess))
	cmd.AddCommand(NewCmdCfgDeleteUser(streams.Out, streams.ErrOut, configAccess))

	return cmd
}
