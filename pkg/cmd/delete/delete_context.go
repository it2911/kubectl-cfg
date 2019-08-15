package delete

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
	"fmt"
)

var (
	deleteContextExample = templates.Examples(`
		# Delete the context for the minikube cluster
		kubectl cfg delete context minikube`)
)

// NewCmdConfigDeleteContext returns a Command instance for 'config delete-context' sub command
func NewCmdCfgDeleteContext(out, errOut io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Delete the specified context from the kubeconfig"),
		Long:                  "Delete the specified context from the kubeconfig",
		Example:               deleteContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunDeleteContext(out, errOut, configAccess, cmd))
		},
	}

	return cmd
}

func RunDeleteContext(out, errOut io.Writer, configAccess clientcmd.ConfigAccess, cmd *cobra.Command) error {
	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	args := cmd.Flags().Args()
	if len(args) != 1 {
		cmd.Help()
		return nil
	}

	configFile := configAccess.GetDefaultFilename()
	if configAccess.IsExplicitFile() {
		configFile = configAccess.GetExplicitFile()
	}

	name := args[0]
	context, ok := config.Contexts[name]

	//backup deleted content to yaml file
	err = backup(errOut, context, "context", "context", name)
	if err != nil {
		fmt.Println("warning: backup to yaml failed.")
	}


	if !ok {
		return fmt.Errorf("cannot delete context %s, not in %s", name, configFile)
	}

	if config.CurrentContext == name {
		fmt.Fprint(errOut, "warning: this removed your active context, use \"kubectl config use-context\" to select a different one\n")
	}

	delete(config.Contexts, name)

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		return err
	}

	fmt.Fprintf(out, "deleted context %s from %s\n", name, configFile)

	return nil
}
