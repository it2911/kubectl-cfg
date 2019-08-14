package list

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"github.com/it2911/kubectl-cfg/pkg/util/printers"
	"github.com/juju/ansiterm"
	. "github.com/logrusorgru/aurora"
)

var (
	listAuthInfoLong = templates.LongDesc(`Displays auth info from the kubeconfig file.`)

	listAuthInfoExample = templates.Examples(`
		# List all the auth info in your kubeconfig file
		kubectl cfg list auth`)
)

// ListAuthInfoOptions contains the assignable options from the args.
type ListAuthInfoOptions struct {
	configAccess clientcmd.ConfigAccess
	nameOnly     bool
	showHeaders  bool
	authInfos    []string

	genericclioptions.IOStreams
}

func NewCmdCfgListAuthInfo(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &ListAuthInfoOptions{
		configAccess: configAccess,
		IOStreams:    streams,
	}

	cmd := &cobra.Command{
		Use:                   "auth",
		DisableFlagsInUseLine: true,
		Short:                 "Describe  auth info from the kubeconfig file",
		Long:                  listAuthInfoLong,
		Example:               listAuthInfoExample,
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

// Complete assigns ListClustersOptions from the args.
func (o *ListAuthInfoOptions) Complete(cmd *cobra.Command, args []string) error {
	o.authInfos = args
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
func (o *ListAuthInfoOptions) RunList() error {

	config, err := o.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	out, found := o.Out.(*ansiterm.TabWriter)
	if !found {
		out = printers.GetNewTabWriter(o.Out)
		defer out.Flush()
	}

	// Build a list of context names to print, and warn if any requested contexts are not found.
	// Do this before printing the headers so it doesn't look ugly.
	allErrs := []error{}
	toPrint := []string{}
	if len(o.authInfos) == 0 {
		for name := range config.AuthInfos {
			toPrint = append(toPrint, name)
		}
	} else {
		for _, name := range o.authInfos {
			_, ok := config.AuthInfos[name]
			if ok {
				toPrint = append(toPrint, name)
			} else {
				allErrs = append(allErrs, fmt.Errorf("context %v not found", name))
			}
		}
	}
	if o.showHeaders {
		err = printAuthInfoHeaders(out, o.nameOnly)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	sort.Strings(toPrint)
	for _, name := range toPrint {
		currentContext := config.Contexts[config.CurrentContext]
		err = printAuthInfo(name, config.AuthInfos[name], out, o.nameOnly, currentContext.AuthInfo == name)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}

func printAuthInfoHeaders(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CURRENT", "AUTH_INFO_NAME", "USERNAME"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func printAuthInfo(name string, authInfo *clientcmdapi.AuthInfo, w io.Writer, nameOnly, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "

	var err error
	if current {
		prefix = "*"
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\n", Yellow(prefix), Yellow(name), Yellow(authInfo.Username))
	} else {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\n", prefix, name, authInfo.Username)
	}

	return err
}
