```console
___ ____ _  _ ____ _  _    _  _ ____ _  _ ____ _  _ ____ ____
 |  |  | |_/  |___ |\ |     \/  |    |__| |__| |\ | | __ |___
 |  |__| | \_ |___ | \|    _/\_ |___ |  | |  | | \| |__] |___
```


## What is this?

This is a Kubernetes client credentials exec provider that enables cross Kubernetes cluster authorization using 
Kubernetes Service Account tokens and [Dex token-exchange](https://github.com/dexidp/dex/pull/2806).

Originally developed to allow ArgoCD Application Controller on one cluster to manage resources on a remote cluster in a 
multi cloud environment where using GKE / EKS IAM authentication was impractical.

## How it works

* The binary reads the local Kubernetes Service Account token.
* It sends a request to the Dex server to exchange the KSA token for a token signed by Dex that the remote cluster accepts.
* Outputs a Kubernetes `ExecCredential` object that can read by `kubectl` and other tools (e.g. ArgoCD).

### Prerequisites

* The source cluster must have published its OpenID Connect Cluster Issuer documents to a public location that Dex can read. 
  * See kOps docs: https://kops.sigs.k8s.io/cluster_spec/#service-account-issuer-discovery-and-aws-iam-roles-for-service-accounts-irsa
* Dex must be configured with an `oidc` connector for the source cluster.
  * See kOps docs: https://kops.sigs.k8s.io/cluster_spec/#oidc-flags-for-open-id-connect-tokens
* The target cluster must be configured with the `oidc` settings that allow it to trust tokens signed by Dex.

## Example 

### Dex Config

```yaml
connectors:
- id: argocd-cluster
  name: argocd-cluster
  type: oidc
  config:
    issuer: https://oidc-argocd-cluster.s3.us-east-1.amazonaws.com
    scopes:
      - openid
      - federated:id
    userNameKey: sub
    
issuer: https://dex.example.com

staticClients:
 - id: target-cluster
   name: target-cluster
   secret: not-a-secret
   public: true
```

### Target Cluster API Server OIDC configuration

```
"--oidc-client-id=target-cluster",
"--oidc-issuer-url=https://dex.example.com",
"--oidc-username-claim=sub",
"--oidc-username-prefix=oidc:",
```

### RBAC

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: remote-argocd-application-controller
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'


---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: remote-argocd-application-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: remote-argocd-application-controller
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: oidc:system:serviceaccount:argocd:argocd-application-controller
```

