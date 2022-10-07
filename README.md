# conversion-webhook-test-with-flux

This is the repo to reproduce the error that is posted in
- [How can the flux error related to conversion webhook be resolved? · Discussion #3105 · fluxcd/flux2](https://github.com/fluxcd/flux2/discussions/3105)
- [`dry-run` sometimes misses metadata and causes `failed to prune fields` error during CRD conversion · Issue #227 · kubernetes-sigs/structured-merge-diff](https://github.com/kubernetes-sigs/structured-merge-diff/issues/227)


## directory structure

```
.
├── README.md: README
├── flux: controllerd by flux
├── myoperator: kubebuilder sample
└── myoperator.yaml: the generated yaml
```

## How to reproduce

### Start minikube

```shell
minikube start 
```

### Deploy cert-manager

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

### Deploy myoperator

myoperator.yaml is generated with the following command:
The image is replaced with `ghcr.io/littlewat/conversion-webhook-test-with-flux/test-resource-controller:aea67ff-amd64`
For m1 mac user, please use `ghcr.io/littlewat/conversion-webhook-test-with-flux/test-resource-controller:aea67ff-arm64`

```shell
cd myoperator
kustomize build config/default > ../myoperator.yaml
```

Deploy it

```shell
kubectl apply -f myoperator.yaml
```


### Start flux

```shell
flux bootstrap github --owner=$GITHUB_USERNAME --repository=conversion-webhook-test-with-flux --branch=main --path=./flux --personal

```

```shell
kubectl apply -k ./flux/flux-system
```

### Wait until the resource is deployed

Please wait a minute.
You will see that the resource in `~/flux/flux-resources` is deployed.
 
```shell
$ kubectl  get testresource -A
NAMESPACE         NAME                           AGE
v1alpha1-flux     testresource-v1alpha1-flux     1m
v1alpha2-flux     testresource-v1alpha2-flux     1m
```

### Check the dry-run result

Please repeat the following dry-run command.
Sometimes the dry-run command gets the following error message.
```shell
$ kubectl apply --server-side --dry-run=server -f ./flux/flux-resources/cwtest_v1alpha1_testresource.yaml --field-manager kustomize-controller

namespace/v1alpha1-flux serverside-applied (server dry run)
Error from server: failed to prune fields: failed add back owned items: failed to convert pruned object at version cwtest.littlewat.github.io/v1alpha1: conversion webhook for cwtest.littlewat.github.io/v1alpha2, Kind=TestResource returned invalid metadata: invalid metadata of type <nil> in input object
```

Check the managed field.

```yaml
$ kubectl get testresource -n v1alpha1-flux testresource-v1alpha1-flux -o yaml --show-managed-fields

apiVersion: cwtest.littlewat.github.io/v1alpha2
kind: TestResource
metadata:
  creationTimestamp: "2022-10-07T04:19:26Z"
  generation: 1
  labels:
    kustomize.toolkit.fluxcd.io/name: flux-system
    kustomize.toolkit.fluxcd.io/namespace: flux-system
  managedFields:
  - apiVersion: cwtest.littlewat.github.io/v1alpha1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:labels:
          f:kustomize.toolkit.fluxcd.io/name: {}
          f:kustomize.toolkit.fluxcd.io/namespace: {}
      f:spec:
        f:foo: {}
    manager: kustomize-controller
    operation: Apply
    time: "2022-10-07T04:19:26Z"
  - apiVersion: cwtest.littlewat.github.io/v1alpha2
    fieldsType: FieldsV1
    fieldsV1:
      f:status:
        f:state: {}
    manager: manager
    operation: Update
    subresource: status
    time: "2022-10-07T04:19:26Z"
  name: testresource-v1alpha1-flux
  namespace: v1alpha1-flux
  resourceVersion: "46369"
  uid: a1ea970f-6dd9-4ae1-8995-af48d6bd4d33
spec:
  bar: v1alpha1FooBar
  foo: v1alpha1Foo
status:
  state: v1alpha1Foo-OK
```

The sample resources to deploy by hand is in `myoperator/config/samples`

Deploy the v1a1 resource.
```shell
kubectl apply --server-side -f ./myoperator/config/samples/cwtest_v1alpha1_testresource.yaml
```

Repeating the dry-run command never get failed.
```shell
kubectl apply --server-side --dry-run=server -f ./myoperator/config/samples/cwtest_v1alpha1_testresource.yaml
```
