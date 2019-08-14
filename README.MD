# Welcome to kubectl-cfg 👋

![Version](https://img.shields.io/badge/version-0.0.1-blue.svg?cacheSeconds=2592000)
[![License: Apache 2](https://img.shields.io/badge/License-Apache%202-yellow.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
![Total visitor](https://visitor-count-badge.herokuapp.com/total.svg?repo_id=kubectl-cfg)
![Visitors in today](https://visitor-count-badge.herokuapp.com/today.svg?repo_id=kubectl-cfg)


> kubectl cfg is a [kubectl plugin](https://github.com/kubernetes-sigs/krew) that provides you a easier way to manage the config of your Kubernetes cluster.

### 🏠 [Homepage](https://github.com/it2911/kubectl-cfg)

## Install
Download the newest version kubectl-cfg binary file from release page. 
```sh
$ which kubectl
/usr/local/bin/kubebuilder/bin/kubectl
$ # copy the kubectl-cfg into some folder with kubectl.
$ cp kubectl-cfg /usr/local/bin/kubebuilder/bin/
$ kubectl cfg version
kubectl-cfg current version is "v0.0.1".
```
*Support the krew install way in recently.*

## Usage

```sh
kubectl cfg list context
```

![](https://raw.githubusercontent.com/77ImageHosting/repo1/master/imgrender1565764238132.gif)

More information

```sh
$ kubectl cfg -h                                
  You can use kubectl cfg command to easily manage kubeconfig file.
  The command is include list, add, delete, merge, rename and update.
  If you have some question please commit the issue to https://github.com/it2911/kubectl-cfg

Usage:
  cfg
  cfg [command]

Available Commands:
  add         Add context / cluster / authinfo into kubeconfig
  delete      Delete context / cluster / authinfo from kubeconfig
  help        Help about any command
  list        Describe the contexts / clusters / authinfos information in the kubeconfig file
  merge       Merge multi the kubeconfig files
  use         Sets the current-context in a kubeconfig file

Flags:
  -h, --help                help for cfg
      --kubeconfig string   use a particular kubeconfig file

Use "cfg [command] --help" for more information about a command.
```
## Build
If you want to build yourself version. You can try below command.
```
$ go version
go version go1.12.7 linux/amd64
$ git clone -b v0.0.1 https://github.com/it2911/kubectl-cfg.git
$ cd kubectl-cfg/cmd
$ go build -o kubectl-cfg
```
*Support makefile in recently*

## Author

👤 [@it2911](https://github.com/it2911) 👤 [@vinci7](https://github.com/vinci7)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/it2911/kubectl-cfg/issues).

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2019 [it2911](https://github.com/it2911).

This project is [Apache 2](http://www.apache.org/licenses/LICENSE-2.0.html) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_


