package use

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	UseContextLong = templates.LongDesc(`Choose the context in your kubeconfig file.`)

	UseContextExample = templates.Examples(`
		# Choose the context in your kubeconfig file
		kubectl cfg use example-context`)
)

var (
	useContextExample = templates.Examples(`
		# Use the context for the minikube cluster
		kubectl config use-context minikube`)
)

type useContextOptions struct {
	configAccess clientcmd.ConfigAccess
	contextName  string
}

// NewCmdConfigUseContext returns a Command instance for 'config use-context' sub command
func NewCmdCfgUseContext(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &useContextOptions{configAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   "use CONTEXT_NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets the current-context in a kubeconfig file"),
		Aliases:               []string{"use"},
		Long:                  `Sets the current-context in a kubeconfig file`,
		Example:               useContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
			fmt.Fprintf(out, "Switched to context %q.\n", options.contextName)
		},
	}

	return cmd
}

func (o useContextOptions) run() error {
	config, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	err = o.validate(config)
	if err != nil {
		return err
	}

	config.CurrentContext = o.contextName

	return clientcmd.ModifyConfig(o.configAccess, *config, true)
}

func (o *useContextOptions) complete(cmd *cobra.Command) error {
	endingArgs := cmd.Flags().Args()
	if len(endingArgs) != 1 {
		//return helpErrorf(cmd, "Unexpected args: %v", endingArgs)
		return nil
	}

	o.contextName = endingArgs[0]
	return nil
}

func (o useContextOptions) validate(config *clientcmdapi.Config) error {
	if len(o.contextName) == 0 {
		return errors.New("empty context names are not allowed")
	}

	for name := range config.Contexts {
		if name == o.contextName {
			return nil
		}
	}

	return fmt.Errorf("no context exists with the name: %q", o.contextName)
}
