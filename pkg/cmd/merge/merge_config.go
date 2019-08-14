package merge

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	kconf "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

const kubeconfigFlag string = "file"

var (
	addConfigLong = templates.LongDesc(`Merge multi the kubeconfig files.`)

	exampleString = `
    # Merge the kubeconfig into the output kubeconfig file
	kubectl cfg merge config -f import-kubeconfig01.yaml -f import-kubeconfig02.yaml > export-kubeconfig.yaml`
	addConfigExample = templates.Examples(exampleString)

	errorString = `
    Kubeconfig file path is need.
    # Merge the kubeconfig into the output kubeconfig file
	kubectl cfg merge config -f import-kubeconfig01.yaml -f import-kubeconfig02.yaml > export-kubeconfig.yaml`
	errorExample = templates.Examples(errorString)

	kubeconfigFiles []string
	kubeconfig      string
)

// NewCmdConfigView returns a Command instance for 'config view' sub command
func NewCmdCfgMergeConfig(streams genericclioptions.IOStreams, ConfigAccess clientcmd.ConfigAccess) *cobra.Command {
	o := &kconf.ViewOptions{
		PrintFlags:   genericclioptions.NewPrintFlags("").WithTypeSetter(scheme.Scheme).WithDefaultOutput("yaml"),
		ConfigAccess: ConfigAccess,
		Flatten:      true,
		IOStreams:    streams,
	}

	cmd := &cobra.Command{
		Use:     fmt.Sprintf("config [--%v=path/kubeconfg] ", kubeconfigFlag),
		Short:   i18n.T("Merge multi the kubeconfig files"),
		Long:    addConfigLong,
		Example: addConfigExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(SetFilePaths(cmd))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run())
		},
	}

	o.PrintFlags.AddFlags(cmd)

	o.Merge.Default(true)
	////cmd.Flags().Var(&o.Merge, "merge", "", "Merge the full hierarchy of kubeconfig files")
	//mergeFlag := cmd.Flags().VarPF(&o.Merge, "merge", "", "Merge the full hierarchy of kubeconfig files")
	//mergeFlag.NoOptDefVal = "true"
	//cmd.Flags().BoolVar(&o.RawByteData, "raw", o.RawByteData, "Display raw byte data")
	//cmd.Flags().BoolVar(&o.Flatten, "flatten", o.Flatten, "Flatten the resulting kubeconfig file into self-contained output (useful for creating portable kubeconfig files)")
	//cmd.Flags().BoolVar(&o.Minify, "minify", o.Minify, "Remove all information not used by current-context from the output")
	cmd.Flags().StringSliceP(kubeconfigFlag, "f", kubeconfigFiles, "Merged the kubeconfig")
	cmd.Flags().String("context", "", "The name of the kubeconfig context to use")
	return cmd
}

func SetFilePaths(cmd *cobra.Command) error {

	filePaths, err := cmd.Flags().GetStringSlice("file")

	if err != nil {
		return err
	}

	if len(filePaths) == 0 {
		return errors.New(errorExample)
	}

	importKubeconfig := strings.Replace(strings.Trim(fmt.Sprint(filePaths), "[]"), " ", ":", -1)

	err = os.Setenv("KUBECONFIG", importKubeconfig)

	if err != nil {
		return err
	}
	return nil
}
