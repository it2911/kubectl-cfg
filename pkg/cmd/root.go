package cmd

import (
	"fmt"
	"github.com/it2911/kubectl-cfg/pkg/cmd/add"
	"github.com/it2911/kubectl-cfg/pkg/cmd/delete"
	"github.com/it2911/kubectl-cfg/pkg/cmd/list"
	"github.com/it2911/kubectl-cfg/pkg/cmd/merge"
	"github.com/it2911/kubectl-cfg/pkg/cmd/use"
	"github.com/it2911/kubectl-cfg/pkg/cmd/version"
	"k8s.io/kubectl/pkg/util/templates"
	"strconv"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
)

// NewCmdConfig creates a command object for the "config" action, and adds all child commands to it.
func NewCmdCfg(f cmdutil.Factory, pathOptions *clientcmd.PathOptions, streams genericclioptions.IOStreams) *cobra.Command {
	if len(pathOptions.ExplicitFileFlag) == 0 {
		pathOptions.ExplicitFileFlag = clientcmd.RecommendedConfigPathFlag
	}

	cmd := &cobra.Command{
		Use:                   "cfg",
		DisableFlagsInUseLine: true,

		Short: i18n.T("Easily manage kubeconfig files"),
		Long: templates.Examples(`You can use kubectl cfg command to easily manage kubeconfig file.
			The command is include list, add, delete, merge, rename and update.
			If you have some question please commit the issue to https://github.com/it2911/kubectl-cfg `),

		Run: cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	// file paths are common to all sub commands
	cmd.PersistentFlags().StringVar(&pathOptions.LoadingRules.ExplicitPath, pathOptions.ExplicitFileFlag, pathOptions.LoadingRules.ExplicitPath, "use a particular kubeconfig file")

	// TODO(juanvallejo): update all subcommands to work with genericclioptions.IOStreams
	cmd.AddCommand(add.NewCmdCfgAdd(streams, pathOptions))
	cmd.AddCommand(delete.NewCmdCfgDelete(streams, pathOptions))
	//cmd.AddCommand(get.NewCmdCfgGet(streams, pathOptions))
	cmd.AddCommand(list.NewCmdCfgList(streams, pathOptions))
	cmd.AddCommand(use.NewCmdCfgUseContext(streams.Out, pathOptions))
	cmd.AddCommand(merge.NewCmdCfgMerge(streams, pathOptions))
	cmd.AddCommand(version.NewCmdCfgVersion(streams.Out, pathOptions))

	return cmd
}

func toBool(propertyValue string) (bool, error) {
	boolValue := false
	if len(propertyValue) != 0 {
		var err error
		boolValue, err = strconv.ParseBool(propertyValue)
		if err != nil {
			return false, err
		}
	}

	return boolValue, nil
}

func helpErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	cmd.Help()
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s", msg)
}
