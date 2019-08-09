package cmd

import (
	"fmt"
	"github.com/it2911/kubectl-cfg/pkg/cmd/add"
	"github.com/it2911/kubectl-cfg/pkg/cmd/delete"
	"github.com/it2911/kubectl-cfg/pkg/cmd/get"
	"github.com/it2911/kubectl-cfg/pkg/cmd/list"
	"github.com/it2911/kubectl-cfg/pkg/cmd/use"
	"k8s.io/kubectl/pkg/util/templates"
	"path"
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
		Short:                 i18n.T("Modify kubeconfig files"),
		Long: templates.LongDesc(`
			Modify kubeconfig files using subcommands like "kubectl config set current-context my-context"
			The loading order follows these rules:
			1. If the --` + pathOptions.ExplicitFileFlag + ` flag is set, then only that file is loaded. The flag may only be set once and no merging takes place.
			2. If $` + pathOptions.EnvVar + ` environment variable is set, then it is used as a list of paths (normal path delimiting rules for your system). These paths are merged. When a value is modified, it is modified in the file that defines the stanza. When a value is created, it is created in the first file that exists. If no files in the chain exist, then it creates the last file in the list.
			3. Otherwise, ` + path.Join("${HOME}", pathOptions.GlobalFileSubpath) + ` is used and no merging takes place.`),

		Run: cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	// file paths are common to all sub commands
	cmd.PersistentFlags().StringVar(&pathOptions.LoadingRules.ExplicitPath, pathOptions.ExplicitFileFlag, pathOptions.LoadingRules.ExplicitPath, "use a particular kubeconfig file")

	// TODO(juanvallejo): update all subcommands to work with genericclioptions.IOStreams
	cmd.AddCommand(add.NewCmdCfgAdd(streams, pathOptions))
	cmd.AddCommand(delete.NewCmdCfgDelete(streams, pathOptions))
	cmd.AddCommand(get.NewCmdCfgGet(streams, pathOptions))
	cmd.AddCommand(list.NewCmdCfgList(streams, pathOptions))
	cmd.AddCommand(use.NewCmdCfgUseContext(streams.Out, pathOptions))

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
