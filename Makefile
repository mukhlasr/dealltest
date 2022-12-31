test:
	@echo "running test"; go test ./...

minikube-deployment:
	@kubectl apply -k deployment/
	@minikube service simpleblog --url

.PHONY: test deployment
