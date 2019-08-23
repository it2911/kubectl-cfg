/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	useContextExample = templates.Examples(`
		# Use the context for the minikube cluster
		kubectl config use-context minikube`)
)

type UseContextOptions struct {
	ConfigAccess clientcmd.ConfigAccess
	ContextName  string
}

// NewCmdConfigUseContext returns a Command instance for 'config use-context' sub command
func NewCmdConfigUseContext(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &UseContextOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   "use-context CONTEXT_NAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets the current-context in a kubeconfig file"),
		Aliases:               []string{"use"},
		Long:                  `Sets the current-context in a kubeconfig file`,
		Example:               useContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "Switched to context %q.\n", options.ContextName)
		},
	}

	return cmd
}

func (o UseContextOptions) Run() error {
	config, err := o.ConfigAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	err = o.validate(config)
	if err != nil {
		return err
	}

	config.CurrentContext = o.ContextName

	return clientcmd.ModifyConfig(o.ConfigAccess, *config, true)
}

func (o *UseContextOptions) Complete(cmd *cobra.Command) error {
	endingArgs := cmd.Flags().Args()
	if len(endingArgs) != 1 {
		return helpErrorf(cmd, "Unexpected args: %v", endingArgs)
	}

	o.ContextName = endingArgs[0]
	return nil
}

func (o UseContextOptions) validate(config *clientcmdapi.Config) error {
	if len(o.ContextName) == 0 {
		return errors.New("empty context names are not allowed")
	}

	for name := range config.Contexts {
		if name == o.ContextName {
			return nil
		}
	}

	return fmt.Errorf("no context exists with the name: %q", o.ContextName)
}
