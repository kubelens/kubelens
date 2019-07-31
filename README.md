# Kubelens

There are many tools built around [Kubernetes](https://kubernetes.io/), with many of them being one-stop tools, I found it difficult to find something more lightweight. After getting tired of running the same sequence of kubectl commands and switching contexts, I started this as a side project to make life easier while learning more about Kubernetes. As it became more useful, I decided to create Kubelens for a specific purpose; to give engineers a quick view into deployed applications.

Let's get to the quick details.

- Fully functional and has been running in multiple K8s clusters in an enterprise production environment since mid-spring 2019.
- Security focused. Authentication & Authorization flows can easily be configured at different levels if desired. More docs to come on this. 
- Intended to be generic and highly configurable to fit any organizations needs. If you find something isn't flexible enough, let's fix it for everyone. 

More docs and features to come. Any help/contributions/feedback is very much appreciated!

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
