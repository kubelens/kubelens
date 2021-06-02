package k8sv1

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

type mockWrapper struct {
	namespace string
	appname   string
	fail      bool
	innerFail bool
}

func (m *mockWrapper) GetClientSet() (clientset kubernetes.Interface, err error) {
	if m.fail {
		return nil, errors.New("GetClientSet Test Error")
	}
	name := m.appname
	if len(name) == 0 {
		name = "systemapp"
	}
	namespace := m.namespace
	if len(namespace) == 0 {
		namespace = "kube-system"
	}

	// Use a timeout to keep the test from hanging.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Create the fake client.
	client := fake.NewSimpleClientset()

	lbl := make(map[string]string)
	lbl["app"] = name

	// We will create an informer that writes added pods to a channel.
	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			pod.SetNamespace(namespace)
			pod.SetLabels(lbl)
			pod.SetName(name)
			pod.Spec.Containers = []v1.Container{
				{
					Env: []v1.EnvVar{
						{
							Name:  "testPassword",
							Value: "password",
						},
						{
							Name:  "testKey",
							Value: "key",
						},
						{
							Name:  "testSecret",
							Value: "secret",
						},
						{
							Name:  "regular_entry",
							Value: "not so sensitive",
						},
					},
				},
			}
			fmt.Printf("pod added: %s/%s\n", pod.Namespace, pod.Name)
			pods <- pod
			cancel()
		},
	})

	namepsaces := make(chan *v1.Namespace, 1)
	nsInformer := informers.Core().V1().Namespaces().Informer()
	nsInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns := obj.(*v1.Namespace)
			ns.Name = namespace
			ns.Namespace = namespace
			fmt.Printf("namespace added: %s\n", ns.Name)
			namepsaces <- ns
			cancel()
		},
	})

	services := make(chan *v1.Service, 1)
	svcInformer := informers.Core().V1().Services().Informer()
	svcInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc := obj.(*v1.Service)
			svc.SetName(name)
			svc.SetNamespace(namespace)
			svc.SetLabels(lbl)
			svc.Spec = v1.ServiceSpec{
				Selector: lbl,
			}
			svc.Status = v1.ServiceStatus{}
			fmt.Printf("service added: %s\n", svc.Name)
			services <- svc
			cancel()
		},
	})

	configMap := make(chan *v1.ConfigMap, 1)
	cmInformer := informers.Core().V1().ConfigMaps().Informer()
	cmInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			cm := obj.(*v1.ConfigMap)
			cm.SetName(name + "-cm")
			cm.SetNamespace(namespace)
			cm.SetLabels(lbl)
			fmt.Printf("config map added: %s\n", cm.Name)
			configMap <- cm
			cancel()
		},
	})

	deployment := make(chan *appsv1.Deployment, 1)
	dplInformer := informers.Apps().V1().Deployments().Informer()
	dplInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			d := obj.(*appsv1.Deployment)
			d.SetName(name)
			d.SetNamespace(namespace)
			d.SetLabels(lbl)
			d.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: lbl,
			}
			fmt.Printf("deployment added: %s\n", d.Name)
			deployment <- d
			cancel()
		},
	})

	daemonset := make(chan *appsv1.DaemonSet, 1)
	dsInformer := informers.Apps().V1().DaemonSets().Informer()
	dsInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			d := obj.(*appsv1.DaemonSet)
			d.SetName(name)
			d.SetNamespace(namespace)
			d.SetLabels(lbl)
			d.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: lbl,
			}
			fmt.Printf("daemonset added: %s\n", d.Name)
			daemonset <- d
			cancel()
		},
	})

	job := make(chan *batchv1.Job, 1)
	jobInformer := informers.Batch().V1().Jobs().Informer()
	jobInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			j := obj.(*batchv1.Job)
			j.SetName(name)
			j.SetNamespace(namespace)
			j.SetLabels(lbl)
			j.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: lbl,
			}
			j.Status = batchv1.JobStatus{
				StartTime:      &metav1.Time{Time: time.Now()},
				CompletionTime: &metav1.Time{Time: time.Now().Add(time.Second)},
			}
			fmt.Printf("job added: %s\n", j.Name)
			job <- j
			cancel()
		},
	})

	replicaset := make(chan *appsv1.ReplicaSet, 1)
	rsInformer := informers.Apps().V1().ReplicaSets().Informer()
	rsInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			s := obj.(*appsv1.ReplicaSet)
			s.SetName(name)
			s.SetNamespace(namespace)
			s.SetLabels(lbl)
			s.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: lbl,
			}
			replicas := int32(1)
			s.Spec.Replicas = &replicas
			fmt.Printf("replicaset added: %s\n", s.Name)
			replicaset <- s
			cancel()
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	if !m.innerFail {
		// This is not required in tests, but it serves as a proof-of-concept by
		// ensuring that the informer goroutine have warmed up and called List before
		// we send any events to it.
		for !podInformer.HasSynced() &&
			!nsInformer.HasSynced() &&
			!svcInformer.HasSynced() &&
			!cmInformer.HasSynced() &&
			!dplInformer.HasSynced() &&
			!dsInformer.HasSynced() &&
			!jobInformer.HasSynced() &&
			!rsInformer.HasSynced() {
			time.Sleep(100 * time.Millisecond)
		}

		a := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
		_, err = client.CoreV1().Services(namespace).Create(ctx, a, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		x := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err = client.CoreV1().Namespaces().Create(ctx, x, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		p := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl}}
		_, err = client.CoreV1().Pods(namespace).Create(ctx, p, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		cmlbl := lbl
		cmlbl["cmtest"] = "cmvalue"
		// Inject an event into the fake client.
		c1 := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: cmlbl}}
		_, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, c1, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		c2 := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: (name + "-notmatch"), Namespace: namespace, Labels: map[string]string{"asdf": "yep"}}}
		_, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, c2, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		d := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: cmlbl}}
		_, err = client.AppsV1().Deployments(namespace).Create(ctx, d, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		ds := &appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl}}
		_, err = client.AppsV1().DaemonSets(namespace).Create(ctx, ds, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		j := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl}}
		_, err = client.BatchV1().Jobs(namespace).Create(ctx, j, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		rs := &appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl}}
		specReplicas := int32(1)
		rs.Spec.Replicas = &specReplicas
		_, err = client.AppsV1().ReplicaSets(namespace).Create(ctx, rs, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	// Wait and check result.
	<-ctx.Done()
	select {
	case pod := <-pods:
		fmt.Printf("Got pod from channel: %s/%s\n", pod.Namespace, pod.Name)
		return client, nil
	case svc := <-services:
		fmt.Printf("Got service from channel: %s/%s\n", svc.Namespace, svc.Name)
		return client, nil
	case ns := <-namepsaces:
		fmt.Printf("Got namespace from channel: %s\n", ns.Name)
		return client, nil
	case cmm := <-configMap:
		fmt.Printf("Got matcher config maps from channel: %s/%s\n", cmm.Namespace, cmm.Name)
		return client, nil
	case dp := <-deployment:
		fmt.Printf("Got deployments from channel: %s/%s\n", dp.Namespace, dp.Name)
		return client, nil
	case dss := <-daemonset:
		fmt.Printf("Got daemonsets from channel: %s/%s\n", dss.Namespace, dss.Name)
		return client, nil
	case jb := <-job:
		fmt.Printf("Got jobs from channel: %s/%s\n", jb.Namespace, jb.Name)
		return client, nil
	case rss := <-replicaset:
		fmt.Printf("Got replicasets from channel: %s/%s\n", rss.Namespace, rss.Name)
		return client, nil
	default:
		fmt.Println("informer did not get the added pod")
		return client, nil
	}
}

func TestGetClient(t *testing.T) {
	w := NewWrapper()

	_, err := w.GetClientSet()

	assert.NotNil(t, err)
}
