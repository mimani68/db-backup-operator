# K8S Database Operator

![](./docs/images/logo-color.png)

# Description

# Getting Started

You’ll need a Kubernetes cluster to run against. You can use KIND to get a local cluster for testing, or run against a remote cluster. Note: Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster kubectl cluster-info shows).

## Running on the cluster

```bash
make manifests
make generate
make docker-build
make docker-push
kustomize build config/default > crd.yaml
```

## Using in kubernetes 

```bash
kubectl apply -f ~/config/crd.yaml
kubectl apply -f ~/config/sample/db_v1alpha1_dbmanagment.yaml
```

## Uninstall CRDs

To delete the CRDs from the cluster:

```bash
make uninstall
```

![](./docs/images/logo-shirt.png)
## Development

```bash
vi .vscode/launch.json
```