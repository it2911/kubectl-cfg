package yaml

import (
	"io"
	"os"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

func WriteYaml(errOut io.Writer, i interface{}, op, name string) error {
	s, err := yaml.Marshal(i)
	if err != nil {
		return err
	}


	fileName := getKubeDir(errOut) + "/kubectl-cfg-delete-bak.yaml"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	defer f.Close()

	if err != nil {
		fmt.Fprint(errOut, err.Error())
		return err
	}

	time := time.Now().Format("2006-01-02 15:04")
	commentary := fmt.Sprintf("# [%s] delete backup of 'kubectl cfg delete %s %s'\n", time, op, name)
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

func getKubeDir(errOut io.Writer) string {
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