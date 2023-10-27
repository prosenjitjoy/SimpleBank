# simplebank
Backend master class [Golang + Postgres + Kubernetes + gRPC] 

## Run
Edit .env file with right configurations. Then follow one of the way to run the project.

### Kubernetes Cluster
First create docker image to use inside kubernetes.
```bash
podman/docker build -t simplebank-api
```
Then, run the following to quickly test the project.
```bash
kubectl apply -f simplebank-pod.yaml
kubectl port-forward simplebank 3000:3000
```
   
To view api specific log, run the following in a new terminal session.
```bash
kubectl logs -c simplebank-api -f simplebank
```
### Dev Environment
```bash
make dev_deploy
```

## Docs
https://dbdocs.io/prosenjitjoy/SimpleBank     
http://localhost:3000/doc/swagger