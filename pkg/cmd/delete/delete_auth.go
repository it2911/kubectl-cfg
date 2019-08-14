package delete

import (
	"fmt"
	"io"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"gopkg.in/yaml.v2"
	"os"
	"time"
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

	err = backup(out, errOut, authInfo, name)
	if err != nil {
		fmt.Println("warning: backup to yaml failed.")
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


func backup(out, errOut io.Writer, i interface{}, name string) error {
	authInfos := map[string] interface{}{
		"users" : map[string] interface{}{
			"name" : name,
			"user" : i,
		},
	}

	return toYaml(out, errOut, authInfos, name)
}

func toYaml(out, errOut io.Writer, i interface{}, name string) error {
	s, err := yaml.Marshal(i)
	if err != nil {
		return err
	}


	fileName := getKubeDir(out, errOut) + "/kubectl-cfg-delete-bak.yaml"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	defer f.Close()

	if err != nil {
		fmt.Fprint(errOut, err.Error())
		return err
	}

	time := time.Now().Format("2006-01-02 15:04")
	commentary := fmt.Sprintf("# [%s] delete backup of 'kubectl cfg delete auth %s'\n", time, name)
	if _, err = f.WriteString(commentary); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if _, err = f.Write(s); err != nil {
		fmt.Fprint(errOut, err.Error())
		return err
	}

	if _, err = f.WriteString("\n"); err != nil {
		fmt.Fprint(errOut, err.Error())
		return err
	}

	return nil
}

func getKubeDir(out, errOut io.Writer,) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeDir := homeDir + "/.kube"

	if _, err := os.Stat(kubeDir); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprint(errOut, "Error: .kube directory doesn't exist.")
		} else {
			fmt.Fprint(errOut, "Error: get .kube directory error.")
		}
	}

	return kubeDir
}