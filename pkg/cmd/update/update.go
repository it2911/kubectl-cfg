package update

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
	setLong = templates.LongDesc(`
	Sets an individual value in a kubeconfig file

	PROPERTY_NAME is a dot delimited name where each token represents either an attribute name or a map key.  Map keys may not contain dots.

	PROPERTY_VALUE is the new value you wish to set. Binary fields such as 'certificate-authority-data' expect a base64 encoded string unless the --set-raw-bytes flag is used.

	Specifying a attribute name that already exists will merge new fields on top of existing values.`)

	setExample = templates.Examples(`
	# Set server field on the my-cluster cluster to https://1.2.3.4
	kubectl config set clusters.my-cluster.server https://1.2.3.4

	# Set certificate-authority-data field on the my-cluster cluster.
	kubectl config set clusters.my-cluster.certificate-authority-data $(echo "cert_data_here" | base64 -i -)

	# Set cluster field in the my-context context to my-cluster.
	kubectl config set contexts.my-context.cluster my-cluster

	# Set client-key-data field in the cluster-admin user using --set-raw-bytes option.
	kubectl config set users.cluster-admin.client-key-data cert_data_here --set-raw-bytes=true`)
)

// NewCmdConfigSet returns a Command instance for 'config set' sub command
func NewCmdConfigSet(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kubectlconfig.SetOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   "set PROPERTY_NAME PROPERTY_VALUE",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets an individual value in a kubeconfig file"),
		Long:                  setLong,
		Example:               setExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "Property %q set.\n", options.PropertyName)
		},
	}

	f := cmd.Flags().VarPF(&options.SetRawBytes, "set-raw-bytes", "", "When writing a []byte PROPERTY_VALUE, write the given string directly without base64 decoding.")
	f.NoOptDefVal = "true"
	return cmd
}
