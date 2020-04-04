package k8v1

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	// watcher := watch.NewFake()
	// client.PrependWatchReactor("pods", testcore.DefaultWatchReactor(watcher, nil))

	lbl := make(map[string]string)
	lbl["app"] = name

	// We will create an informer that writes added pods to a channel.
	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			// pod.SetName(fmt.Sprintf("%s-%d", name, rand.Intn(1000)))
			pod.SetNamespace(namespace)
			pod.SetLabels(lbl)
			pod.SetName(name)
			// pod.Status = v1.PodStatus{
			// 	HostIP: "1.1.1.1",
			// }
			// watcher.Add(pod)
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
			// watcher.Add(pod)
			fmt.Printf("namespace added: %s\n", ns.Name)
			namepsaces <- ns
			cancel()
		},
	})

	svcSelector := map[string]string{}
	svcSelector["cmtest"] = "cmvalue"

	svcSelectorNotMatch := map[string]string{}
	svcSelectorNotMatch["cmtest"] = "cmNotMatchValue"

	services := make(chan *v1.Service, 1)
	svcInformer := informers.Core().V1().Services().Informer()
	svcInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc := obj.(*v1.Service)
			svc.SetName(name)
			svc.SetNamespace(namespace)
			svc.SetLabels(lbl)
			svc.Spec = v1.ServiceSpec{
				Selector: svcSelector,
			}
			svc.Status = v1.ServiceStatus{}
			// watcher.Add(pod)
			fmt.Printf("service added: %s\n", svc.Name)
			services <- svc
			cancel()
		},
	})

	configMapMatch := make(chan *v1.ConfigMap, 1)
	configMapNotMatch := make(chan *v1.ConfigMap, 1)
	cmInformer := informers.Core().V1().ConfigMaps().Informer()
	cmInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// add a config map that matches the service
			cm := obj.(*v1.ConfigMap)
			cm.SetName(name + "-cm")
			cm.SetNamespace(namespace)
			cm.SetLabels(svcSelector)
			fmt.Printf("config map added: %s\n", cm.Name)
			configMapMatch <- cm

			// add a config map the should not match the service
			cm2 := obj.(*v1.ConfigMap)
			cm2.SetName(name + "-cm-notmatch")
			cm2.SetNamespace(namespace)
			cm2.SetLabels(svcSelectorNotMatch)
			fmt.Printf("config map (not match) added: %s\n", cm2.Name)
			configMapNotMatch <- cm2

			cancel()
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	if !m.innerFail {
		// This is not required in tests, but it serves as a proof-of-concept by
		// ensuring that the informer goroutine have warmed up and called List before
		// we send any events to it.
		for !podInformer.HasSynced() || !nsInformer.HasSynced() || !svcInformer.HasSynced() {
			time.Sleep(1 * time.Millisecond)
		}

		a := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
		_, err = client.CoreV1().Services(namespace).Create(a)
		if err != nil {
			return nil, err
		}

		x := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err = client.CoreV1().Namespaces().Create(x)
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		p := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl}}
		_, err = client.CoreV1().Pods(namespace).Create(p)
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		c1 := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: (name + "-cm"), Namespace: namespace, Labels: svcSelector}}
		_, err = client.CoreV1().ConfigMaps(namespace).Create(c1)
		if err != nil {
			return nil, err
		}

		// Inject an event into the fake client.
		c2 := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: (name + "-cm-notmatch"), Namespace: namespace, Labels: svcSelectorNotMatch}}
		_, err = client.CoreV1().ConfigMaps(namespace).Create(c2)
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
	case cmm := <-configMapMatch:
		fmt.Printf("Got matcher config maps from channel: %s/%s\n", cmm.Namespace, cmm.Name)
		return client, nil
	case cmnm := <-configMapMatch:
		fmt.Printf("Got non-matcher config maps from channel: %s/%s\n", cmnm.Namespace, cmnm.Name)
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
