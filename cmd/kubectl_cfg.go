package main

import (
	"flag"
	"github.com/it2911/kubectl-cfg/pkg/cmd"
	cmdutil "github.com/it2911/kubectl-for-plugin-cfg/pkg/cmd/util"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
)

func init() {
	// Initialize glog flags
	klog.InitFlags(flag.CommandLine)
	flag.CommandLine.Set("logtostderr", "true")
}

func main() {
	flags := pflag.NewFlagSet("kubectl-cfg", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdCfg(cmdutil.NewFactory(genericclioptions.NewTestConfigFlags()),
		clientcmd.NewDefaultPathOptions(),
		genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
