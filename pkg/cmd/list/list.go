package list

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	listLong = templates.LongDesc(`Displays the contexts / clusters / authinfos information in the kubeconfig file.`)

	listExample = templates.Examples(`
		# List the resources in your kubeconfig file
		kubectl cfg list SUB_COMMAND`)
)

func NewCmdCfgList(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Short:                 "Describe the contexts / clusters / authinfos information in the kubeconfig file",
		Long:                  listLong,
		Example:               listExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")

	cmd.AddCommand(NewCmdCfgListContext(streams, configAccess))
	cmd.AddCommand(NewCmdCfgListCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgListAuthInfo(streams, configAccess))
	return cmd
}
