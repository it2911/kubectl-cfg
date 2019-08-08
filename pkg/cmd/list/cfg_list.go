package list

import (
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	listLong = templates.LongDesc(`Displays one or many clusters from the kubeconfig file.`)

	listExample = templates.Examples(`
		# List all the clusters in your kubeconfig file
		kubectl cfg list cluster
		# List all the contexts in your kubeconfig file
		kubectl cfg list context`)
)

func NewCmdCfgList(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  listLong,
		Example:               listExample,
		Run:                   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")

	cmd.AddCommand(NewCmdCfgListContext(streams, configAccess))
	cmd.AddCommand(NewCmdCfgListCluster(streams, configAccess))
	cmd.AddCommand(NewCmdCfgListUser(streams, configAccess))
	return cmd
}
