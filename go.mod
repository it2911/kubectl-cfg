module github.com/it2911/kubectl-cfg

require (
	github.com/fatih/color v1.7.0
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6 // indirect
	github.com/it2911/kubectl-for-plugin-cfg v0.0.0-20190808090532-d94c8ccdde3f
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	golang.org/x/arch v0.0.0-20190312162104-788fe5ffcd8c // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	k8s.io/apimachinery v0.0.0-20190806215851-162a2dabc72f
	k8s.io/cli-runtime v0.0.0-20190807063455-7df0a100ca6c
	k8s.io/client-go v0.0.0-20190807061213-4fd06e107451
	k8s.io/component-base v0.0.0-20190807101431-d6d4632c35d0
	k8s.io/klog v0.3.1
	k8s.io/kubectl v0.0.0-20190807223317-83f665480eb9 // indirect
)

replace k8s.io/kubectl => github.com/it2911/kubectl-for-plugin-cfg v0.0.0-20190807223317-83f665480eb9
