package use

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	kconf "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	UseContextLong = templates.LongDesc(`Choose the context in your kubeconfig file.`)

	UseContextExample = templates.Examples(`
		# Choose the context in your kubeconfig file
		kubectl cfg use example-context`)
)

// NewCmdCfgUseContext returns a Command instance for 'config use-context' sub command
func NewCmdCfgUseContext(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kconf.UseContextOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   "use CONTEXT_NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets the current-context in a kubeconfig file"),
		Aliases:               []string{"use"},
		Long:                  `Sets the current-context in a kubeconfig file`,
		Example:               UseContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "Switched to context %q.\n", options.ContextName)
		},
	}

	return cmd
}
