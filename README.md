# conversion-webhook-test-with-flux

## Start minikube

```shell
minikube start 
```

## Deploy myoperator

```shell
cd myoperator
make install
make deploy
```


## Start flux
```shell
kubectl apply -k ./flux/flux-system
```

## Check the dry-run result

```shell
kubectl apply --server-side --dry-run=server -f ./flux/flux-resources/cwtest_v1alpha1_testresource.yaml --field-manager kustomize-controller
```
