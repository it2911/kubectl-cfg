package list

import (
	"fmt"
	. "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/config"
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"strings"
)

var (
	listClustersLong = templates.LongDesc(`Displays one or many clusters from the kubeconfig file.`)

	listClustersExample = templates.Examples(`
		# List all the clusters in your kubeconfig file
		kubectl cfg list cluster`)
)

var printClusterHeaders = func(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CLUSTER_NAME", "SERVER", "CERTIFICATE_AUTHORITY_VALIDITY_FROM", "CERTIFICATE_AUTHORITY_VALIDITY_TO"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

var printCluster = func(name string, context *clientcmdapi.Context, w io.Writer, nameOnly, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "
	if current {
		prefix = "*"
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", prefix, name, context.Cluster, context.AuthInfo, context.Namespace)
	return err
}

func NewCmdCfgListCluster(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &ListOptions{
		configAccess: configAccess,
		IOStreams:    streams,
	}

	cmd := &cobra.Command{
		Use:                   "cluster",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many clusters",
		Long:                  listClustersLong,
		Example:               listClustersExample,
		Run: func(cmd *cobra.Command, args []string) {
			validOutputTypes := sets.NewString("", "json", "yaml", "wide", "name", "custom-columns", "custom-columns-file", "go-template", "go-template-file", "jsonpath", "jsonpath-file")
			supportedOutputTypes := sets.NewString("", "name")
			outputFormat := cmdutil.GetFlagString(cmd, "output")
			if !validOutputTypes.Has(outputFormat) {
				cmdutil.CheckErr(fmt.Errorf("output must be one of '' or 'name': %v", outputFormat))
			}
			if !supportedOutputTypes.Has(outputFormat) {
				fmt.Fprintf(options.Out, "--output %v is not available in kubectl config get-contexts; resetting to default output format\n", outputFormat)
				cmd.Flags().Set("output", "")
			}
			cmdutil.CheckErr(options.Complete(cmd, args))
			cmdutil.CheckErr(options.RunList(printClusterHeaders, printCluster))
		},
	}

	_ = NewCmdConfigSetCluster(nil, nil)

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")
	return cmd
}
