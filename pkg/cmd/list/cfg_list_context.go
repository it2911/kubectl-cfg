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
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/printers"
	"k8s.io/kubectl/pkg/util/templates"
	"sort"
	"strings"
	"github.com/fatih/color"
)

var (
	listContextsLong = templates.LongDesc(`Displays one or many contexts from the kubeconfig file.`)

	listContextsExample = templates.Examples(`
		# List all the contexts in your kubeconfig file
		kubectl cfg list context`)
)

// ListContextsOptions contains the assignable options from the args.
type ListContextOptions struct {
	configAccess clientcmd.ConfigAccess
	nameOnly     bool
	showHeaders  bool
	contextNames []string

	genericclioptions.IOStreams
}

// NewCmdConfigListContexts creates a command object for the "get-contexts" action, which
// retrieves one or more contexts from a kubeconfig.
func NewCmdCfgListContext(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &ListContextOptions{
		configAccess: configAccess,
		IOStreams:    streams,
	}

	cmd := &cobra.Command{
		Use:                   "context",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  listContextsLong,
		Example:               listContextsExample,
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

	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: name")
	return cmd
}

// Complete assigns ListContextsOptions from the args.
func (o *ListContextOptions) Complete(cmd *cobra.Command, args []string) error {
	o.contextNames = args
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

// RunListContexts implements all the necessary functionality for context retrieval.
func (o *ListContextOptions) RunList() error {

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
	if len(o.contextNames) == 0 {
		for name := range config.Contexts {
			toPrint = append(toPrint, name)
		}
	} else {
		for _, name := range o.contextNames {
			_, ok := config.Contexts[name]
			if ok {
				toPrint = append(toPrint, name)
			} else {
				allErrs = append(allErrs, fmt.Errorf("context %v not found", name))
			}
		}
	}
	if o.showHeaders {
		err = printContextHeaders(out, o.nameOnly)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	sort.Strings(toPrint)
	for _, name := range toPrint {
		err = printContext(name, config.Contexts[name], out, o.nameOnly, config.CurrentContext == name)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}

func printContextHeaders(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CURRENT", "CONTEXT_NAME", "CLUSTER_NAME", "AUTH_INFO", "DEFAULT_NAMESPACE"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func printContext(name string, context *clientcmdapi.Context, w io.Writer, nameOnly, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "
	var err error
	if current {
		prefix = "*"
		//	_, err = color.New(color.FgHiYellow).Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", prefix, name, context.Cluster, context.AuthInfo, context.Namespace)
		//}else{
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", prefix, name, context.Cluster, context.AuthInfo, context.Namespace)
	}
	return err
}
