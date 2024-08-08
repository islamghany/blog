# Blog

This project is a monolithic Go application with PostgreSQL, deployed on Kubernetes using Minikube for local development. The application has two main components: the admin entry point for administrative tasks and the blog-api for serving the blog content.

## Prerequisites

- [Make](https://www.gnu.org/software/make/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)

## Deployment Steps

1. Install Minikube and start the cluster.

```bash
minikube start
```

2. build the Docker image and push it to the Minikube Docker daemon.

```bash
make dev-up
```

3. Deploy the application to the Minikube cluster.

```bash
make dev-apply
```

4. Show the Cluster status.

```bash
make dev-status-all
```

5. Show the application logs.

```bash
make dev-logs
```

6. Open the application in the browser.

```bash
minikube service blog-admin -n blog-system
```

A table like the following will be displayed:

| NAMESPACE   | NAME       | TARGET PORT | URL                                |
| ----------- | ---------- | ----------- | ---------------------------------- |
| blog-system | blog-admin | 8000        | http://http://192.168.59.104:30131 |
| blog-system | blog-api   | 8000        | http://192.168.59.104:32413        |

then open the URL in the browser.

7. down the application from the Minikube cluster.

```bash
make dev-down
```

8. delete the Minikube cluster.

```bash
minikube delete
```

## TODO List

- [x] Run InitContainer to initialize the database migrations/seed.
- [x] Liveness and Readiness probes.
- [ ] Resources Requests and Limits.
- [ ] Horizontal Pod Autoscaler.
- [ ] Ingress Controller.
- [x] Monitoring with Prometheus and Grafana.
- [ ] CI/CD Integration.
- [ ] Helm Charts.
- [ ] Add more tests.
- [ ] Network Policies.
- [x] Secrets Management.
- [ ] Security Best Practices.
