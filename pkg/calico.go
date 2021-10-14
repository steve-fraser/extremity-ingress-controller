package pkg

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func ConfigureCalico(kubeconfig string) error {

	fmt.Printf("Configuring Calcio\n")

	localip, err := GetInterfaceIpv4Addr("eth0")
	if err != nil {
		return fmt.Errorf("Error querying IP: ip %s, error: %v", localip, err)
	}

	read, err := ioutil.ReadFile("/tmp/manifests/calico-bgp.yaml")
	if err != nil {
		panic(err)
	}

	newContents := strings.Replace(string(read), "<replace>", localip, -1)

	fmt.Println(newContents)
	fmt.Printf("applying calico bgp configuration\n")
	err = ioutil.WriteFile("/tmp/manifests/calico-bgp.yaml", []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
	kubectl_command := "kubectl apply -f /tmp/manifests/calico-bgp.yaml --kubeconfig " + kubeconfig
	cmd := exec.Command("/bin/sh", "-c", kubectl_command)

	err = cmd.Run()

	if err != nil {
		fmt.Printf("ERROR: Couldn't add BGP Config\n")
		panic(err)
	}

	// Build a clientset based on the provided kubeconfig file.
	// cs, err := clientset.NewForConfig(kubeconfig)
	// if err != nil {
	// 	panic(err)
	// }
	// // List global network policies.
	// list, err := cs.ProjectcalicoV3().BGPConfigurations().List(context.Background(), v1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, gnp := range list.Items {
	// 	fmt.Printf("%#v\n", gnp)
	// }

	return nil
}
