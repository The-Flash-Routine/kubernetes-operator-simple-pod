# SimplePod Kubernetes Custom Resource and Operator

### Brief
The repo contains code for: 
1. Custom Resource Definition (CRD) for a new type of object in Kubernetes called SimplePod.
2. Operator associated with the SimplePod resource for managing it.

### Implementation Brief

Kubernetes being highly extensible allows creation of custom API's and associated representational Objects.
These custom objects can then make use of native Kubernetes resources to obtain a desried functionality.

An Operator is needed to manage the lifecycle and translation of these objects to native Kubernetes resources

I've used the Kubebuilder Framework for the purpose.

### SimplePod

Meta information around SimplePod resource:

- **API Group**: pod.routine.kat
- **Version**: v1
- **Kind**: SimplePod  
- **Spec**: Specification resembles that of a Pod specification
- **Status**: Status consists of a field **PodIp** holding information around managed pod's IP.

### Behaviour of SimplePod

The object creates a simple Kubernetes Pod with specification mentioned in the spec of SimplePod resource.

It then monitors the created pod and displays it's IP in it's own status.

In case the pod is deleted, it recreates the pod with same specification.

### Basic commands

1. Creation of basic folder structure 
    
```
    kubebuilder init --domain routine.kat --repo routine.kat/simple-pod-operator
```

2. SimplePod API Creation

```
kubebuilder create api --group pod --version v1 --kind SimplePod
```
3. Creation of Custom Resource Defnition API Specs

```
make manifests
```

4. Installing CRD's into cluster

```
make install
```

5. Running Operator locally

```
make run
```

6. Building Operator Image and Deploying on Cluster

```
make docker-build docker-push IMG=routinekat/kubernetes-operator-simple-pod:v1.0.0
make deploy IMG=routinekat/kubernetes-operator-simple-pod:v1.0.0
```

### Side Information

1. I'm using a local cluster created using kind

2. Sample yaml for SimplePod can be found at config/samples

3. Logic for controller can be found at internal/controller/ simplepod_controller.go

4. CRD API-Specs & RoleBindings  for SimplePod can be found config/crd and config/rbac

5. Go Representation of SimePlod object can be found at api/v1/simplepod_types.go






