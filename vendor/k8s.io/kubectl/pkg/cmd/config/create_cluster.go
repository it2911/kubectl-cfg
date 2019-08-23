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
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

type CreateClusterOptions struct {
	ConfigAccess          clientcmd.ConfigAccess
	Name                  string
	Server                cliflag.StringFlag
	InsecureSkipTLSVerify cliflag.Tristate
	CertificateAuthority  cliflag.StringFlag
	EmbedCAData           cliflag.Tristate
}

var (
	createClusterLong = templates.LongDesc(`
		Sets a cluster entry in kubeconfig.

		Specifying a name that already exists will merge new fields on top of existing values for those fields.`)

	createClusterExample = templates.Examples(`
		# Set only the server field on the e2e cluster entry without touching other values.
		kubectl config set-cluster e2e --server=https://1.2.3.4

		# Embed certificate authority data for the e2e cluster entry
		kubectl config set-cluster e2e --certificate-authority=~/.kube/e2e/kubernetes.ca.crt

		# Disable cert checking for the dev cluster entry
		kubectl config set-cluster e2e --insecure-skip-tls-verify=true`)
)

// NewCmdConfigSetCluster returns a Command instance for 'config set-cluster' sub command
func NewCmdConfigSetCluster(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &CreateClusterOptions{ConfigAccess: configAccess}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("set-cluster NAME [--%v=server] [--%v=path/to/certificate/authority] [--%v=true]", clientcmd.FlagAPIServer, clientcmd.FlagCAFile, clientcmd.FlagInsecure),
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets a cluster entry in kubeconfig"),
		Long:                  createClusterLong,
		Example:               createClusterExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "Cluster %q set.\n", options.Name)
		},
	}

	options.InsecureSkipTLSVerify.Default(false)

	cmd.Flags().Var(&options.Server, clientcmd.FlagAPIServer, clientcmd.FlagAPIServer+" for the cluster entry in kubeconfig")
	f := cmd.Flags().VarPF(&options.InsecureSkipTLSVerify, clientcmd.FlagInsecure, "", clientcmd.FlagInsecure+" for the cluster entry in kubeconfig")
	f.NoOptDefVal = "true"
	cmd.Flags().Var(&options.CertificateAuthority, clientcmd.FlagCAFile, "Path to "+clientcmd.FlagCAFile+" file for the cluster entry in kubeconfig")
	cmd.MarkFlagFilename(clientcmd.FlagCAFile)
	f = cmd.Flags().VarPF(&options.EmbedCAData, clientcmd.FlagEmbedCerts, "", clientcmd.FlagEmbedCerts+" for the cluster entry in kubeconfig")
	f.NoOptDefVal = "true"

	return cmd
}

func (o CreateClusterOptions) Run() error {
	err := o.validate()
	if err != nil {
		return err
	}

	config, err := o.ConfigAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	startingStanza, exists := config.Clusters[o.Name]
	if !exists {
		startingStanza = clientcmdapi.NewCluster()
	}
	cluster := o.modifyCluster(*startingStanza)
	config.Clusters[o.Name] = &cluster

	if err := clientcmd.ModifyConfig(o.ConfigAccess, *config, true); err != nil {
		return err
	}

	return nil
}

// cluster builds a Cluster object from the options
func (o *CreateClusterOptions) modifyCluster(existingCluster clientcmdapi.Cluster) clientcmdapi.Cluster {
	modifiedCluster := existingCluster

	if o.Server.Provided() {
		modifiedCluster.Server = o.Server.Value()
	}
	if o.InsecureSkipTLSVerify.Provided() {
		modifiedCluster.InsecureSkipTLSVerify = o.InsecureSkipTLSVerify.Value()
		// Specifying insecure mode clears any certificate authority
		if modifiedCluster.InsecureSkipTLSVerify {
			modifiedCluster.CertificateAuthority = ""
			modifiedCluster.CertificateAuthorityData = nil
		}
	}
	if o.CertificateAuthority.Provided() {
		caPath := o.CertificateAuthority.Value()
		if o.EmbedCAData.Value() {
			modifiedCluster.CertificateAuthorityData, _ = ioutil.ReadFile(caPath)
			modifiedCluster.InsecureSkipTLSVerify = false
			modifiedCluster.CertificateAuthority = ""
		} else {
			caPath, _ = filepath.Abs(caPath)
			modifiedCluster.CertificateAuthority = caPath
			// Specifying a certificate authority file clears certificate authority data and insecure mode
			if caPath != "" {
				modifiedCluster.InsecureSkipTLSVerify = false
				modifiedCluster.CertificateAuthorityData = nil
			}
		}
	}

	return modifiedCluster
}

func (o *CreateClusterOptions) Complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) != 1 {
		return helpErrorf(cmd, "Unexpected args: %v", args)
	}

	o.Name = args[0]
	return nil
}

func (o CreateClusterOptions) validate() error {
	if len(o.Name) == 0 {
		return errors.New("you must specify a non-empty cluster name")
	}
	if o.InsecureSkipTLSVerify.Value() && o.CertificateAuthority.Value() != "" {
		return errors.New("you cannot specify a certificate authority and insecure mode at the same time")
	}
	if o.EmbedCAData.Value() {
		caPath := o.CertificateAuthority.Value()
		if caPath == "" {
			return fmt.Errorf("you must specify a --%s to embed", clientcmd.FlagCAFile)
		}
		if _, err := ioutil.ReadFile(caPath); err != nil {
			return fmt.Errorf("could not read %s data from %s: %v", clientcmd.FlagCAFile, caPath, err)
		}
	}

	return nil
}
