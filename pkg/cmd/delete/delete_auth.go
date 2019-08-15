package delete

import (
	"fmt"
	"io"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"github.com/it2911/kubectl-cfg/pkg/util/yaml"

)

var (
	deleteAuthinfoExample = templates.Examples(`
		# Delete the minikube cluster
		kubectl cfg delete cluster minikube`)
)

func NewCmdCfgDeleteUser(out, errOut io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "auth AUTHINFO_NAME",
		DisableFlagsInUseLine: true,
		Short:                 "Delete the specified authinfo from the kubeconfig",
		Long:                  "Delete the specified authinfo from the kubeconfig",
		Example:               deleteAuthinfoExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunDeleteAuthInfo(out, errOut, configAccess, cmd))
		},
	}
	return cmd
}

func RunDeleteAuthInfo(out, errOut io.Writer, configAccess clientcmd.ConfigAccess, cmd *cobra.Command) error {
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

	authInfo, ok := config.AuthInfos[name]

	//backup deleted content to yaml file
	err = backup(errOut, authInfo, "auth", "user", name)
	if err != nil {
		fmt.Println("warning: backup to yaml failed.")
	} else {
		fmt.Println("info: deleted content backup to .kube/kubectl-cfg-delete-bak.yaml")
	}

	if !ok {
		return fmt.Errorf("cannot delete auth %s, not in %s", name, configFile)
	}

	context := config.Contexts[config.CurrentContext]

	if context.AuthInfo == name {
		fmt.Fprint(errOut, "warning: this removed your active context, use \"kubectl config use-context\" to select another one\n")
	}

	delete(config.AuthInfos, name)

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		return err
	}

	fmt.Fprintf(out, "deleted authinfo %s from %s\n", name, configFile)

	return nil
}


func backup(errOut io.Writer, i interface{}, op, key, name string) error {
	t := map[string] interface{}{
		key+"s" : map[string] interface{}{
			"name" : name,
			key : i,
		},
	}

	return yaml.WriteYaml(errOut, t, op, name)
}

