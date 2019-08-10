package merge

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	listLong = templates.LongDesc(`Merge multi the kubeconfig files.`)

	listExample = templates.Examples(`
		# Merge the kubeconfig into the output kubeconfig file
		kubectl cfg merge config -f import-kubeconfig01.yaml -f import-kubeconfig02.yaml > export-kubeconfig.yaml`)
)

func NewCmdCfgMerge(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "merge",
		DisableFlagsInUseLine: true,
		Short:                 "Merge multi the kubeconfig files",
		Long:                  listLong,
		Example:               listExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")

	cmd.AddCommand(NewCmdCfgMergeConfig(streams, configAccess))
	return cmd
}
