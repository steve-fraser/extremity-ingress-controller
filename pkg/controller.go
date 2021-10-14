package pkg

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	informerv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset      kubernetes.Interface
	depLister      listersv1.NodeLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface, depInformer informerv1.NodeInformer) *controller {
	c := &controller{
		clientset:      clientset,
		depLister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)

	return c
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cahce %s\n", err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name %s\n", err.Error())
		return false
	}
	fmt.Printf("splitting key into namespace %s and name %s\n", ns, name)

	// check if the object has been deleted from k8s cluster
	ctx := context.Background()
	node, err := c.clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("handle delete event for dep %s\n", name)
		return true
	}

	for i := 0; i < len(node.Status.Addresses); i++ {
		if node.Status.Addresses[i].Type == "InternalIP" {
			nodeip := node.Status.Addresses[i].Address
			err := addNodeToBird(nodeip)
			if err != nil {
				fmt.Printf("Error adding node to bird %s\n", err.Error())
				return false
			}
			fmt.Println(nodeip)
		}
	}
	// 	err = c.clientset.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
	// 	if err != nil {
	// 		fmt.Printf("deleting ingrss %s, error %s\n", name, err.Error())
	// 		return false
	// 	}

	// 	return true
	// }

	// err = c.syncDeployment(ns, name)
	// if err != nil {
	// 	// re-try
	// 	fmt.Printf("syncing deployment %s\n", err.Error())
	// 	return false
	// }
	return true
}

// func (c *controller) syncDeployment(ns, name string) error {
// 	ctx := context.Background()
// 	if err != nil {
// 		fmt.Printf("getting deployment from lister %s\n", err.Error())
// 	}
// 	// create service
// 	// we have to modify this, to figure out the port
// 	// our deployment's container is listening on
// 	svc := corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      dep.Name,
// 			Namespace: ns,
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Selector: depLabels(*dep),
// 			Ports: []corev1.ServicePort{
// 				{
// 					Name: "http",
// 					Port: 80,
// 				},
// 			},
// 		},
// 	}
// 	s, err := c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
// 	if err != nil {
// 		fmt.Printf("creating service %s\n", err.Error())
// 	}
// 	// create ingress
// 	return createIngress(ctx, c.clientset, s)
// }

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("add was called")
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("del was called")
	c.queue.Add(obj)
}
