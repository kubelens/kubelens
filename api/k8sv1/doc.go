/*Package k8sv1 is essentially a wrapper package to interact with Kubernetes client-go.
Any container can have a pod/service/daemon name that's auto generated or passed in as a static value.
Here's an example of what can be done with naming:
Let's say you have names domain-projects-1234-env or domain-projects-1234-env-passive (for blue green deployments)
The "projects-1234" could be the slug ID, which is the deployment template for the application that can
be traced back to system deploying to K8. When using Helm, names need a unique < 64 char name since it stores deployments in
the kube-system namespace so it can deploy to any other namespace. Since the deployment is basically a global name, we
have to ensure the name is always unique to an application and can be easily (as much as possible) linked back to an actual application.
*/
package k8sv1
