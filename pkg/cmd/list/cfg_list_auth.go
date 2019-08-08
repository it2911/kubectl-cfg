package list

import (
	"fmt"
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/printers"
	"github.com/it2911/kubectl-for-plugin-cfg/pkg/util/templates"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	"io"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sort"

	"strings"
)

var (
	listUsersLong = templates.LongDesc(`Displays one or many users from the kubeconfig file.`)

	listUsersExample = templates.Examples(`
		# List all the users in your kubeconfig file
		kubectl cfg list user`)
)

// ListUsersOptions contains the assignable options from the args.
type ListUserOptions struct {
	configAccess clientcmd.ConfigAccess
	nameOnly     bool
	showHeaders  bool
	authInfos    []string

	genericclioptions.IOStreams
}

func NewCmdCfgListUser(streams genericclioptions.IOStreams, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &ListUserOptions{
		configAccess: configAccess,
		IOStreams:    streams,
	}

	cmd := &cobra.Command{
		Use:                   "auth",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many users",
		Long:                  listUsersLong,
		Example:               listUsersExample,
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
func (o *ListUserOptions) Complete(cmd *cobra.Command, args []string) error {
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
func (o *ListUserOptions) RunList() error {

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
		err = printUserHeaders(out, o.nameOnly)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	sort.Strings(toPrint)
	for _, name := range toPrint {
		currentContext := config.Contexts[config.CurrentContext]
		err = printUser(name, config.AuthInfos[name], out, o.nameOnly, currentContext.AuthInfo == name)
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return utilerrors.NewAggregate(allErrs)
}

func printUserHeaders(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CURRENT", "AUTH_INFO_NAME", "USERNAME"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func printUser(name string, authInfo *clientcmdapi.AuthInfo, w io.Writer, nameOnly, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "
	if current {
		prefix = "*"
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\n", prefix, name, authInfo.Username)
	return err
}
