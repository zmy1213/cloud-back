# cloud-back

`cloud-back` is an empty backend scaffold modeled after `kube-nova` service layout.

## Structure

- `application/cloud-api`: HTTP API service skeleton
- `application/cloud-rpc`: RPC-like internal service skeleton
- `manifests`: Kubernetes manifests scaffold
- `dockerfile`: image build scaffold

## Quick start

```bash
cd cloud-back
go mod tidy
make run-cloud-rpc
```

Open another terminal:

```bash
cd cloud-back
make run-cloud-api
```

## Ports

- cloud-api: `:8813`
- cloud-rpc: `:30113`
