package fakes

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

type K8V1Wrapper struct {
	Namespace string
	AppName   string
	Fail      bool
	InnerFail bool
}

func (m *K8V1Wrapper) GetClientSet() (clientset kubernetes.Interface, err error) {
	if m.Fail {
		return nil, errors.New("GetClientSet Test Error")
	}
	name := m.AppName
	if len(name) == 0 {
		name = "systemapp"
	}
	namespace := m.Namespace
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
			svc.Spec = v1.ServiceSpec{}
			svc.Status = v1.ServiceStatus{}
			fmt.Printf("service added: %s\n", svc.Name)
			services <- svc
			cancel()
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	if !m.InnerFail {
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
	default:
		fmt.Println("informer did not get the added pod")
		return client, nil
	}
}
