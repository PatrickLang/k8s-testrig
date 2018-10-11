# azure-k8s-testrig

`testrig` makes spinning up and managing test Kubernetes clusters on Azure simple.

### Usage:
```
testrig create --location=eastus myCluster
```

This will spin up a full kubernetes cluster on Azure.
This can be customized to some extent with the available CLI flags:

```
$ testrig create --help
Create a new kubernetes cluster on Azure

Usage:
  testrig create [flags]

      --acs-engine-path string                    Set the path to use for acs-engine (default "acs-engine")
  -h, --help                                      help for create
      --kubernetes-version string                 set the kubernetes version to use (default "1.10")
      --linux-agent-availability-profile string   set the availabiltiy profile for agent nodes (default "VirtualMachineScaleSets")
      --linux-agent-count int                     sets the number of nodes for agent pools (default 3)
      --linux-agent-node-sku string               sets sku to use for agent nodes (default "Standard_DS2_v2")
      --linux-leader-count int                    sets the number of nodes for the leader pool (default 3)
      --linux-leader-node-sku string              sets sku to use for agent nodes (default "Standard_DS2_v2")
  -l, --location string                           Set the location to deploy to
      --network-plugin string                     set the network plugin to use for the cluster (default "azure")
      --network-policy string                     set the network policy to use for the cluster (default "azure")
      --runtime string                            sets the containe runtime to use
      --ssh-key sshKey                            set public SSH key to install as authorized keys in cluster nodes
      -s, --subscription string                       Set the subscription to use to deploy with
      -u, --user string                               set the username to use for nodes (default "azureuser")
```

When creating a cluster you can provide your own (public) ssh key or a key pair will be generated for you.

#### Authentication

`testrig` attempts to setup authentcation in the following order:
1. Read service principal from the `AZURE_AUTH_LOCATION`
2. Receive a bearer token through azure-cli

If method `1` fails (e.g. if the value is unset), then method 2 is attempted.

To make a service principal suitable for `1`, run:

```
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP_NAME}"
```

The output of this should be saved to a file and the file path passed into `testrig` as an environment variable as `AZURE_AUTH_LOCATION`.

For method `2`, you must make sure the azure-cli is logged in: `az login`.

Commands will use whatever subscription you are logged in with, or you can pass in a custom subscription.

#### Management commands

```
$ testrig --help
Quickly create and manage test Kubernetes clusters on Azure for testing purposes

Usage:
  testrig [command]

Available Commands:
  create      Create a new kubernetes cluster on Azure
  help        Help about any command
  inspect     Get details about an existing cluster
  kubeconfig  Get the path to the kubeconfig file for the specified cluster
  ls          List available clusters
  rm          Remove a cluster
  ssh         ssh into a running cluster

Flags:
  -h, --help               help for testrig
      --state-dir string   directory to store state information to
```

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

