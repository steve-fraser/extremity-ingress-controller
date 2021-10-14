package main

import (
	"flag"
	"fmt"

	// "time"
	pkg "github.com/steve-fraser/extremity-ingress-controller/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/root/.kube/config", "location to your kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// handle error
		fmt.Printf("error %s, creating clientset\n", err.Error())
	}
	pkg.ConfigureCalico("/root/.kube/config")
	ch := make(chan struct{})
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Nodes()
	c := pkg.NewController(clientset, informer)
	factory.Start(ch)
	c.Run(ch)
	fmt.Println(informer)
}
