package add

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	kconf "k8s.io/kubectl/pkg/cmd/config"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	createAuthInfoLong = fmt.Sprintf(templates.LongDesc(`
		Sets a user entry in kubeconfig

		Specifying a name that already exists will merge new fields on top of existing values.

		    Client-certificate flags:
		    --%v=certfile --%v=keyfile

		    Bearer token flags:
			  --%v=bearer_token

		    Basic auth flags:
			  --%v=basic_user --%v=basic_password

		Bearer token and basic auth are mutually exclusive.`), clientcmd.FlagCertFile, clientcmd.FlagKeyFile, clientcmd.FlagBearerToken, clientcmd.FlagUsername, clientcmd.FlagPassword)

	createAuthInfoExample = templates.Examples(`
		# Set only the "client-key" field on the "cluster-admin"
		# entry, without touching other values:
		kubectl config set-credentials cluster-admin --client-key=~/.kube/admin.key

		# Set basic auth for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --username=admin --password=uXFGweU9l35qcif

		# Embed client certificate data in the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --client-certificate=~/.kube/admin.crt --embed-certs=true

		# Enable the Google Compute Platform auth provider for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --auth-provider=gcp

		# Enable the OpenID Connect auth provider for the "cluster-admin" entry with additional args
		kubectl config set-credentials cluster-admin --auth-provider=oidc --auth-provider-arg=client-id=foo --auth-provider-arg=client-secret=bar

		# Remove the "client-secret" config value for the OpenID Connect auth provider for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --auth-provider=oidc --auth-provider-arg=client-secret-

		# Enable new exec auth plugin for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --exec-command=/path/to/the/executable --exec-api-version=client.authentication.k8s.io/v1beta

		# Define new exec auth plugin args for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --exec-arg=arg1 --exec-arg=arg2

		# Create or update exec auth plugin environment variables for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --exec-env=key1=val1 --exec-env=key2=val2

		# Remove exec auth plugin environment variables for the "cluster-admin" entry
		kubectl config set-credentials cluster-admin --exec-env=var-to-remove-`)
)

// NewCmdConfigSetAuthInfo returns an Command option instance for 'config set-credentials' sub command
func NewCmdCfgAddAuthInfo(out io.Writer, configAccess clientcmd.ConfigAccess) *cobra.Command {
	options := &kconf.CreateAuthInfoOptions{ConfigAccess: configAccess}
	return newCmdCfgAddAuthInfo(out, options)
}

func newCmdCfgAddAuthInfo(out io.Writer, options *kconf.CreateAuthInfoOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf(
			"auth NAME [--%v=path/to/certfile] "+
				"[--%v=path/to/keyfile] "+
				"[--%v=bearer_token] "+
				"[--%v=basic_user] "+
				"[--%v=basic_password] "+
				"[--%v=provider_name] "+
				"[--%v=key=value] "+
				"[--%v=exec_command] "+
				"[--%v=exec_api_version] "+
				"[--%v=arg] "+
				"[--%v=key=value]",
			clientcmd.FlagCertFile,
			clientcmd.FlagKeyFile,
			clientcmd.FlagBearerToken,
			clientcmd.FlagUsername,
			clientcmd.FlagPassword,
			kconf.FlagAuthProvider,
			kconf.FlagAuthProviderArg,
			kconf.FlagExecCommand,
			kconf.FlagExecAPIVersion,
			kconf.FlagExecArg,
			kconf.FlagExecEnv,
		),
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Sets a user entry in kubeconfig"),
		Long:                  createAuthInfoLong,
		Example:               createAuthInfoExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := options.Complete(cmd, out)
			if err != nil {
				cmd.Help()
				cmdutil.CheckErr(err)
			}
			cmdutil.CheckErr(options.Run())
			fmt.Fprintf(out, "User %q set.\n", options.Name)
		},
	}

	cmd.Flags().Var(&options.ClientCertificate, clientcmd.FlagCertFile, "Path to "+clientcmd.FlagCertFile+" file for the user entry in kubeconfig")
	cmd.MarkFlagFilename(clientcmd.FlagCertFile)
	cmd.Flags().Var(&options.ClientKey, clientcmd.FlagKeyFile, "Path to "+clientcmd.FlagKeyFile+" file for the user entry in kubeconfig")
	cmd.MarkFlagFilename(clientcmd.FlagKeyFile)
	cmd.Flags().Var(&options.Token, clientcmd.FlagBearerToken, clientcmd.FlagBearerToken+" for the user entry in kubeconfig")
	cmd.Flags().Var(&options.Username, clientcmd.FlagUsername, clientcmd.FlagUsername+" for the user entry in kubeconfig")
	cmd.Flags().Var(&options.Password, clientcmd.FlagPassword, clientcmd.FlagPassword+" for the user entry in kubeconfig")
	cmd.Flags().Var(&options.AuthProvider, kconf.FlagAuthProvider, "Auth provider for the user entry in kubeconfig")
	cmd.Flags().StringSlice(kconf.FlagAuthProviderArg, nil, "'key=value' arguments for the auth provider")
	cmd.Flags().Var(&options.ExecCommand, kconf.FlagExecCommand, "Command for the exec credential plugin for the user entry in kubeconfig")
	cmd.Flags().Var(&options.ExecAPIVersion, kconf.FlagExecAPIVersion, "API version of the exec credential plugin for the user entry in kubeconfig")
	cmd.Flags().StringSlice(kconf.FlagExecArg, nil, "New arguments for the exec credential plugin command for the user entry in kubeconfig")
	cmd.Flags().StringArray(kconf.FlagExecEnv, nil, "'key=value' environment values for the exec credential plugin")
	f := cmd.Flags().VarPF(&options.EmbedCertData, clientcmd.FlagEmbedCerts, "", "Embed client cert/key for the user entry in kubeconfig")
	f.NoOptDefVal = "true"

	return cmd
}
