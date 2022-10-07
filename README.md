# conversion-webhook-test-with-flux

## Start minikube

```shell
minikube start 
```

## Deploy certmanager

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

## Deploy myoperator

```shell
cd myoperator

make install
make deploy
```


## Start flux

```shell
flux bootstrap github --owner=$GITHUB_USERNAME --repository=conversion-webhook-test-with-flux --branch=main --path=./flux --personal

```

```shell
kubectl apply -k ./flux/flux-system
```

## Check the dry-run result

```shell
kubectl apply --server-side --dry-run=server -f ./flux/flux-resources/cwtest_v1alpha1_testresource.yaml --field-manager kustomize-controller
```
