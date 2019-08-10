package list

import (
	"fmt"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	"io"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	. "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/printers"
	"k8s.io/kubectl/pkg/util/templates"
	"sort"
	"strings"
)

var (
	listClustersLong = templates.LongDesc(`Displays one or many clusters from the kubeconfig file.`)

	listClustersExample = templates.Examples(`
		# List all the clusters in your kubeconfig file
		kubectl cfg list cluster`)
)

// ListClusterOptions contains the assignable options from the args.
type ListClusterOptions struct {
	configAccess clientcmd.ConfigAccess
	nameOnly     bool
	showHeaders  bool
	clusterNames []string

	genericclioptions.IOStreams
}

func NewCmdCfgListCluster(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &ListClusterOptions{
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
			cmdutil.CheckErr(options.RunList())
		},
	}

	_ = NewCmdConfigSetCluster(nil, nil)

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")
	return cmd
}

// Complete assigns ListClustersOptions from the args.
func (o *ListClusterOptions) Complete(cmd *cobra.Command, args []string) error {
	o.clusterNames = args
	o.nameOnly = false
	if cmdutil.GetFlagString(cmd, "output") == "name" {
		o.nameOnly = true
	}
	o.showHeaders = true
	if cmdutil.GetFlagBool(cmd, "no-headers") || o.nameOnly {
		o.showHeaders = false
	}

	return nil
}

// RunList implements all the necessary functionality for cluster retrieval.
func (o *ListClusterOptions) RunList() error {

	config, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	out, found := o.Out.(*tabwriter.Writer)
	if !found {
		out = printers.GetNewTabWriter(o.Out)
		defer out.Flush()
	}

	// Build a list of context names to print, and warn if any requested contexts are not found.
	// Do this before printing the headers so it doesn't look ugly.
	allErrs := []error{}
	toPrint := []string{}
	if len(o.clusterNames) == 0 {
		for name := range config.Clusters {
			toPrint = append(toPrint, name)
		}
	} else {
		for _, name := range o.clusterNames {
			_, ok := config.Clusters[name]
			if ok {
				toPrint = append(toPrint, name)
			} else {
				allErrs = append(allErrs, fmt.Errorf("context %v not found", name))
			}
		}
	}
	if o.showHeaders {
		err = printClusterHeaders(out, o.nameOnly)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	sort.Strings(toPrint)
	for _, name := range toPrint {
		currentContext := config.Contexts[config.CurrentContext]
		err = printCluster(name, config.Clusters[name], out, o.nameOnly, currentContext.Cluster == name)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}

func printClusterHeaders(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CURRENT", "CLUSTER_NAME", "SERVER", "STATUS_CODE", "CERTIFICATE_AUTHORITY_VALIDITY_TO"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func printCluster(name string, cluster *clientcmdapi.Cluster, w io.Writer, nameOnly, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "
	if current {
		prefix = "*"
	}

	statusCode := "UNKNOW"
	//resp, err := http.Get(cluster.Server)
	//if err != nil{
	//
	//}
	//statusCode = fmt.Sprint(resp.StatusCode)

	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", prefix, name, cluster.Server, statusCode, "")
	return err
}
