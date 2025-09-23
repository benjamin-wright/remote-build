# Remote Build

An operator and client combo:

Client:
- Wrap around the docker CLI
- Deploys CRs into the cluster to represent a build instance
- passes build args through to docker CLI

Operator:
- Watches for CRs and deploys buildkit single-replica statefulset and proxy pods

Proxy:
- holds incoming build TCP connections and scales up buildkit pod when a build comes in
- passes through the TCP connection when the buildkit pod is alive
- scales down the buildkit pod when there haven't been any new builds for X amount of time

```mermaid
flowchart LR
  User --> Client
  Client --> CRD:Instance
  CRD:Instance -- read by --> Operator 
  Operator -- manage --> Proxy
  Client --> Proxy
  Proxy -- pass through --> BuildKit
  Proxy -- scale up / down --> BuildKit

```
