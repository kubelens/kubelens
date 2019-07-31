# Kubelens

There are many tools built around surfacing Kubernetes information, deploy, monitor, etc., with many of them being a one-stop tool. Kubelens was created for a specific purpose; to give developers a quick view into the applications they develop. 

This started as a side project to help me learn more about Kubernetes, while making it easier for me to run the same sequence of kubectl commands to view the components/inner-workings of deployed applications. Any help/contributions/feedback is very much appreciated!

[![CircleCI](https://circleci.com/gh/kubelens/kubelens/tree/master.svg?style=svg)](https://circleci.com/gh/kubelens/kubelens/tree/master)

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
