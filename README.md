# conversion-webhook-test-with-flux

This is the repo to reproduce the error that is posted in
- flux issue: [Using conversion webhooks produces the `invalid metadata of type <nil> in input object` error · Issue #3179 · fluxcd/flux2](https://github.com/fluxcd/flux2/issues/3179)
- flux question: [How can the flux error related to conversion webhook be resolved? · Discussion #3105 · fluxcd/flux2](https://github.com/fluxcd/flux2/discussions/3105)
- structured-merge-diff issue: [`dry-run` sometimes misses metadata and causes `failed to prune fields` error during CRD conversion · Issue #227 · kubernetes-sigs/structured-merge-diff](https://github.com/kubernetes-sigs/structured-merge-diff/issues/227)


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

### Deploy `myoperator`

`myoperator` is a simple operator that 
- uses the conversion webhook that converts v1alpha1 to v1alpha2
- updates the status of TestResource

There are two ways to deploy myoperator

- With Tilt
- Wighout Tilt

Tilt would be useful to debug quickly.

#### With Tilt

```shell
cd myoperator
tilt up
```

#### Without Tilt

First, deploy cert-manager:

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```

Then, deploy the operator:

```shell
kubectl apply -f myoperator.yaml
```

`./myoperator.yaml` is generated with the following command:

```shell
cd myoperator
kustomize build config/default > ../myoperator.yaml
```
(The image is replaced with `ghcr.io/littlewat/conversion-webhook-test-with-flux/test-resource-controller:aea67ff-amd64`
For m1 mac user, please use `ghcr.io/littlewat/conversion-webhook-test-with-flux/test-resource-controller:aea67ff-arm64`)


### Start flux

Export your GitHub personal access token and username:

```shell
export GITHUB_USER=<YOUR_GITHUB_USER>
export GITHUB_TOKEN=<YOUR_GITHUB_TOKEN>
```

Create the secret:

```shell
flux create secret git flux-system -u $GITHUB_USER -p $GITHUB_TOKEN --url https://github.com/LittleWat/conversion-webhook-test-with-flux.git
```

Deploy the flux-system:

```shell
kubectl apply -k ./flux/flux-system
```

### Wait until the resource is deployed by flux

Please wait a minute.
The resources in `~/flux/flux-resources` will be deployed.
 
```shell
$ kubectl  get testresource -A

NAMESPACE         NAME                           AGE
v1alpha1-flux     testresource-v1alpha1-flux     1m
v1alpha2-flux     testresource-v1alpha2-flux     1m
```

### Check the dry-run result

Please repeat the following dry-run command.
Sometimes the dry-run command gets the following error message:

```shell
$ kubectl apply --server-side --dry-run=server -f ./flux/flux-resources/cwtest_v1alpha1_testresource.yaml --field-manager kustomize-controller

namespace/v1alpha1-flux serverside-applied (server dry run)
Error from server: failed to prune fields: failed add back owned items: failed to convert pruned object at version cwtest.littlewat.github.io/v1alpha1: conversion webhook for cwtest.littlewat.github.io/v1alpha2, Kind=TestResource returned invalid metadata: invalid metadata of type <nil> in input object
```

Check the managed field:

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

Deploy the v1a1 resource:

```shell
kubectl apply --server-side -f ./myoperator/config/samples/cwtest_v1alpha1_testresource.yaml
```

Repeating the dry-run command never get failed:

```shell
kubectl apply --server-side --dry-run=server -f ./myoperator/config/samples/cwtest_v1alpha1_testresource.yaml
```

From this, it can be said that there is  a difference between the deployment in flux and `kubectl apply --server-side`.

I am glad if this repo is useful for debugging. Thank you!
