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
	deleteClusterExample = templates.Examples(`
		# Delete the minikube cluster
		kubectl cfg delete cluster minikube`)
)

// NewCmdConfigDeleteCluster returns a Command instance for 'config delete-cluster' sub command
func NewCmdCfgDeleteCluster(out, errOut io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "cluster NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Delete the specified cluster from the kubeconfig"),
		Long:                  "Delete the specified cluster from the kubeconfig",
		Example:               deleteClusterExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunDeleteCluster(out, errOut, configAccess, cmd))
		},
	}

	return cmd
}

func RunDeleteCluster(out, errOut io.Writer, configAccess clientcmd.ConfigAccess, cmd *cobra.Command) error {
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
	cluster, ok := config.Clusters[name]

	//backup deleted content to yaml file
	err = backup(errOut, cluster, "cluster", "cluster", name)
	if err != nil {
		fmt.Println("warning: backup to yaml failed.")
	}

	if !ok {
		return fmt.Errorf("cannot delete cluster %s, not in %s", name, configFile)
	}

	delete(config.Clusters, name)

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		return err
	}

	fmt.Fprintf(out, "deleted cluster %s from %s\n", name, configFile)

	return nil
}
