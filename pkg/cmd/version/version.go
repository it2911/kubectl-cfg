package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	versionLong = templates.LongDesc(`Show the version of kubectl-cfg plugin.`)

	versionExample = templates.Examples(`
		# Show the version of kubectl-cfg plugin
		kubectl cfg version`)
)

// NewCmdCfgUseContext returns a Command instance for 'config use-context' sub command
func NewCmdCfgVersion(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Show the version of kubectl-cfg plugin."),
		Long:                  versionLong,
		Example:               versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "kubectl-cfg current version is %q.\n", "v0.0.1")
		},
	}

	return cmd
}
