package add

import (
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	addLong = templates.LongDesc(`Displays one or many clusters from the kubeconfig file.`)

	addExample = templates.Examples(`
		# List all the clusters in your kubeconfig file
		kubectl cfg add cluster
		# List all the contexts in your kubeconfig file
		kubectl cfg add context`)
)

func NewCmdCfgAdd(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "add",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  addLong,
		Example:               addExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	cmd.AddCommand(NewCmdCfgAddConfig(streams.Out, configAccess))
	cmd.AddCommand(NewCmdCfgAddContext(streams.Out, configAccess))
	cmd.AddCommand(NewCmdCfgAddCluster(streams.Out, configAccess))
	cmd.AddCommand(NewCmdCfgAddAuthInfo(streams.Out, configAccess))
	return cmd
}
