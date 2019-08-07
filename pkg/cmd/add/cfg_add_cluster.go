package add

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdCfgAddCluster(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "cluster [(-o|--output=)name)]",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "listContextsLong",
		Example:               "listContextsExample",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
