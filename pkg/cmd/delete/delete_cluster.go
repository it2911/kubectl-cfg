package delete

import (
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	kubectlconfig "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	deleteClusterExample = templates.Examples(`
		# Delete the minikube cluster
		kubectl cfg delete cluster minikube`)
)

// NewCmdConfigDeleteCluster returns a Command instance for 'config delete-cluster' sub command
func NewCmdCfgDeleteCluster(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "cluster NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Delete the specified cluster from the kubeconfig"),
		Long:                  "Delete the specified cluster from the kubeconfig",
		Example:               deleteClusterExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(kubectlconfig.RunDeleteCluster(out, configAccess, cmd))
		},
	}

	return cmd
}
