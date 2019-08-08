package add

import (
	"fmt"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
)

type AddConfigOptions struct {
	ConfigAccess         clientcmd.ConfigAccess
	ImportKubeconfigFile string
	ExportKubeconfigFile string
}

var (
	addConfigLong = templates.LongDesc(`Displays one or many contexts from the kubeconfig file.`)

	listContextsExample = templates.Examples(`
	# List all the contexts in your kubeconfig file
	kubectl config get-contexts
	# Describe one context in your kubeconfig file.
	kubectl config get-contexts my-context`)
)

// NewCmdConfigListContexts creates a command object for the "get-contexts" action, which
// retrieves one or more contexts from a kubeconfig.
func NewCmdCfgAddConfig(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &AddConfigOptions{
		ConfigAccess: configAccess,
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("config [INPUT_KUBECONFIG_FILE | -f] [INPUT_KUBECONFIG_FILE | -o]"),
		DisableFlagsInUseLine: true,
		Short:                 " one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(options.Complete(cmd, args))
			//cmdutil.CheckErr(options.RunListContexts())
			fmt.Fprintf(out, "Add kubeconfig file %q to %s.\n", options.ImportKubeconfigFile, options.ExportKubeconfigFile)
		},
	}

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")
	return cmd
}
