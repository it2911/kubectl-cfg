package add

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
	kubectlconfig "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	addContextLong = templates.LongDesc(`
		Sets a context entry in kubeconfig
		Specifying a name that already exists will merge new fields on top of existing values for those fields.`)

	addContextExample = templates.Examples(`
		# Set the user field on the gce context entry without touching other values
		kubectl cfg add context gce --user=cluster-admin`)
)

// NewCmdConfigSetContext returns a Command instance for 'config set-context' sub command
func NewCmdCfgAddContext(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kubectlconfig.CreateContextOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("context [NAME | --current] [--%v=cluster_nickname] [--%v=user_nickname] [--%v=namespace]", clientcmd.FlagClusterName, clientcmd.FlagAuthInfoName, clientcmd.FlagNamespace),
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets a context entry in kubeconfig"),
		Long:                  addContextLong,
		Example:               addContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			name, exists, err := options.Run()
			cmdutil.CheckErr(err)
			if exists {
				fmt.Fprintf(out, "Context %q modified.\n", name)
			} else {
				fmt.Fprintf(out, "Context %q created.\n", name)
			}
		},
	}

	cmd.Flags().BoolVar(&options.CurrContext, "current", options.CurrContext, "Modify the current context")
	cmd.Flags().Var(&options.Cluster, clientcmd.FlagClusterName, clientcmd.FlagClusterName+" for the context entry in kubeconfig")
	cmd.Flags().Var(&options.AuthInfo, clientcmd.FlagAuthInfoName, clientcmd.FlagAuthInfoName+" for the context entry in kubeconfig")
	cmd.Flags().Var(&options.Namespace, clientcmd.FlagNamespace, clientcmd.FlagNamespace+" for the context entry in kubeconfig")

	return cmd
}
