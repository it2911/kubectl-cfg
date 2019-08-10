package add

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	addLong = templates.LongDesc(`Create context / cluster / authinfo into the kubeconfig file.`)

	addExample = templates.Examples(`
		# Add resources into your kubeconfig file
		kubectl cfg add SUB_COMMAND`)
)

func NewCmdCfgAdd(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "add",
		DisableFlagsInUseLine: true,
		Short:                 "Add context / cluster / authinfo into kubeconfig",
		Long:                  addLong,
		Example:               addExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	//cmd.AddCommand(merge.NewCmdCfgAddConfig(cmdutil.NewFactory(genericclioptions.NewTestConfigFlags()), streams, configAccess))
	cmd.AddCommand(NewCmdCfgAddContext(streams.Out, configAccess))
	cmd.AddCommand(NewCmdCfgAddCluster(streams.Out, configAccess))
	cmd.AddCommand(NewCmdCfgAddAuthInfo(streams.Out, configAccess))

	return cmd
}
