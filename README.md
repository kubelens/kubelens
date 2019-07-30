# kubelens
A lightweight lens for applications running in Kubernetes. 

It's been just me working on this so any help/contributions/feedback is very welcome! This is fully functional as is, but it's still a work in progress :)

## Minikube 

[Install Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/)

`minikube start -p minikube-local`

[Enable Ingress](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/)

`kubectl config use-context minikube`

[Install Helm](https://helm.sh/docs/using_helm/)

Some service account is needed in order for the API to self authenticate with read rights. The example uses the default service account provided by Kubernetes

`kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default`

After that, you should be set to deploy to the Minikube instance.

# Build & Deploy

[kubelens/web](https://github.com/kubelens/kubelens/tree/staging/web#build--deploy)

[kubelens/api](https://github.com/kubelens/kubelens/tree/staging/api#build--deploy)
