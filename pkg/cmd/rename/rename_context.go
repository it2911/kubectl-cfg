package rename

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	kubectlconfig "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	renameContextUse = "rename context CONTEXT_NAME NEW_CONTEXT_NAME"

	renameContextShort = "Renames a context from the kubeconfig file."
)

var (
	renameContextLong = templates.LongDesc(`
		Renames a context from the kubeconfig file.

		CONTEXT_NAME is the context name that you wish to change.

		NEW_CONTEXT_NAME is the new name you wish to set.

		Note: In case the context being renamed is the 'current-context', this field will also be updated.`)

	renameContextExample = templates.Examples(`
		# Rename the context 'old-name' to 'new-name' in your kubeconfig file
		kubectl cfg rename context old-name new-name`)
)

// NewCmdConfigRenameContext creates a command object for the "rename-context" action
func NewCmdCfgRenameContext(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kubectlconfig.RenameContextOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   renameContextUse,
		DisableFlagsInUseLine: true,
		Short:                 renameContextShort,
		Long:                  renameContextLong,
		Example:               renameContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			// remove redundant args "context"
			if args[0] == "context" {
				args = args[1:]
			}
			cmdutil.CheckErr(options.Complete(cmd, args, streams.Out))
			cmdutil.CheckErr(options.Validate())
			cmdutil.CheckErr(options.RunRenameContext(streams.Out))
		},
	}
	return cmd
}
