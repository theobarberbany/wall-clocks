Installation:

1. Ensure [cert-manager](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html) is installed and working
2. `make docker-build` and push to registry (`make docker push`)
3. `make deploy IMG=<some-registry>/wallclocks:sha`
4. Wait for everything to stabilise. Certificates need to be provisioned. Ensure all deployments are running the correct images, Kustomize can be flaky / break in unexpected ways.

5. If there are RBAC errors in the manager logs, `kubectl create -f crb-admin.yaml`. (I don't recommend this for production, but auto generating RBAC from Kubebuilder annotations via Kustomize is a dark art.)

6. Create a Timezones object (see `test_timezones.yaml`)
7. `kubectl get wallclock -oyaml` shows wallclocks and their statuses
8. `kubectl get timezones` will show information about Timezones objects
