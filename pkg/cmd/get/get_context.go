package get

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

var (
//listContextsLong = templates.LongDesc(`Displays one or many contexts from the kubeconfig file.`)
//
//listContextsExample = templates.Examples(`
//	# List all the contexts in your kubeconfig file
//	kubectl config get-contexts
//	# Describe one context in your kubeconfig file.
//	kubectl config get-contexts my-context`)
)

// NewCmdConfigListContexts creates a command object for the "get-contexts" action, which
// retrieves one or more contexts from a kubeconfig.
func NewCmdCfgGetContext(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "context [(-o|--output=)name)]",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}
