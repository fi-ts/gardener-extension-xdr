# Gardener Cortex XDR extension

Provides a gardener extension for managing [Cortex XDR](https://www.paloaltonetworks.de/cortex/detection-and-response-10-must-haves) for a shoot cluster.

## Development

### Setup local gardener

1. Clone the Gardener Repository:
```bash
git clone git@github.com:gardener/gardener.git
```

2. Start a local Kubernetes cluster:
```bash
make kind-up
```

3. Deploy Gardener:
```bash
make gardener-up
```

4. Generate Helm Charts:
```bash
make generate
```

### Deploy the Extension

1. Apply the example configuration:
```bash
kubectl apply -k example/
```

2. Apply the shoot cluster configuration:
```bash
kubectl apply -f example/shoot.yaml
```

### Update Code Changes

When making changes to the code, build and deploy locally using:
```bash
make push-to-gardener-local
```


