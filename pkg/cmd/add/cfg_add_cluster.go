package add

import (
	"fmt"
	kubectlconfig "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/config"
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/i18n"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	addClusterLong = templates.LongDesc(`
		Sets a cluster entry in kubeconfig.
		Specifying a name that already exists will merge new fields on top of existing values for those fields.`)

	addClusterExample = templates.Examples(`
		# Set only the server field on the e2e cluster entry without touching other values.
		kubectl cfg add cluster e2e --server=https://1.2.3.4
		# Embed certificate authority data for the e2e cluster entry
		kubectl cfg add cluster e2e --certificate-authority=~/.kube/e2e/kubernetes.ca.crt
		# Disable cert checking for the dev cluster entry
		kubectl cfg add cluster e2e --insecure-skip-tls-verify=true`)
)

// NewCmdConfigSetCluster returns a Command instance for 'config set-cluster' sub command
func NewCmdCfgAddCluster(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kubectlconfig.CreateClusterOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("cluster NAME [--%v=server] [--%v=path/to/certificate/authority] [--%v=true]", clientcmd.FlagAPIServer, clientcmd.FlagCAFile, clientcmd.FlagInsecure),
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets a cluster entry in kubeconfig"),
		Long:                  addClusterLong,
		Example:               addClusterExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "Cluster %q set.\n", options.Name)
		},
	}

	options.InsecureSkipTLSVerify.Default(false)

	cmd.Flags().Var(&options.Server, clientcmd.FlagAPIServer, clientcmd.FlagAPIServer+" for the cluster entry in kubeconfig")
	f := cmd.Flags().VarPF(&options.InsecureSkipTLSVerify, clientcmd.FlagInsecure, "", clientcmd.FlagInsecure+" for the cluster entry in kubeconfig")
	f.NoOptDefVal = "true"
	cmd.Flags().Var(&options.CertificateAuthority, clientcmd.FlagCAFile, "Path to "+clientcmd.FlagCAFile+" file for the cluster entry in kubeconfig")
	cmd.MarkFlagFilename(clientcmd.FlagCAFile)
	f = cmd.Flags().VarPF(&options.EmbedCAData, clientcmd.FlagEmbedCerts, "", clientcmd.FlagEmbedCerts+" for the cluster entry in kubeconfig")
	f.NoOptDefVal = "true"

	return cmd
}
