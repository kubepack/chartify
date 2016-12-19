# chartify
Generate Helm Charts from Kubernetes api objects

## Installation
```go
go install github.com/appscode/chartify
```

## Usage
You can provide Kubernetes objects as YAML/JSON files in a directory using --kube-dir flag. Or, you can read Kubernetes
objects from a cluster. Chartify will read objects from the current context of your local kubeconfig file.
 
You can use this as a standalone cli or a Helm plugin.

```
chartify create NAME [FLAGS]
```

### Options

```
      --chart-dir string             Specify the location where charts will be created (default "charts")
      --configmaps stringSlice       Specify the names of configmaps(configmap.namespace) to include them in chart
      --daemons stringSlice          Specify names of daemon sets(daemons.namespace)
      --jobs stringSlice             Specify names of jobs
      --kube-dir string              Specify the directory of the yaml files for Kubernetes objects
      --pods stringSlice             Specify the names of pods (podname.namespace) to include them in chart
      --pvcs stringSlice             Specify names of persistent volume claim
      --pvs stringSlice              Specify names of persistent volumes
      --rcs stringSlice              Specify the names of replication cotrollers (rcname.namespace) to include them in chart
      --replicasets stringSlice      Specify names of replica sets(replicaset_name.namespace)
      --secrets stringSlice          Specify the names of secrets(secret_name.namespace) to include them in chart
      --services stringSlice         Specify the names of services to include them in chart
      --statefulsets stringSlice     Specify names of statefulsets(statefulset_name.namespace)
      --storageclasses stringSlice   Specify names of storageclasses
```
