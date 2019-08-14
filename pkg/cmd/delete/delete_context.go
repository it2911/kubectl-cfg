package delete

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	kubectlconfig "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	deleteContextExample = templates.Examples(`
		# Delete the context for the minikube cluster
		kubectl cfg delete context minikube`)
)

// NewCmdConfigDeleteContext returns a Command instance for 'config delete-context' sub command
func NewCmdCfgDeleteContext(out, errOut io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Delete the specified context from the kubeconfig"),
		Long:                  "Delete the specified context from the kubeconfig",
		Example:               deleteContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(kubectlconfig.RunDeleteContext(out, errOut, configAccess, cmd))
		},
	}

	return cmd
}
