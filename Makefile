.PHONY: cluster
cluster:
	mkdir -p ./.scratch
	k3d cluster create remote-build \
		--servers 1 \
		--agents 0 \
		--k3s-arg "--disable=traefik@server:*" \
		--registry-create remote-build-registry:5001 \
		--port "8080:8080@server:0" \
		--kubeconfig-update-default=false \
		--wait
	k3d kubeconfig get remote-build > ./.scratch/kubeconfig

.PHONY: clean
clean:
	k3d cluster delete remote-build || true
	rm -rf ./.scratch
