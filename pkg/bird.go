package pkg

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Mapping the BIRD type extracted from the peer name to the display type.
var bgpTypeMap = map[string]string{
	"Global": "global",
	"Mesh":   "node-to-node mesh",
	"Node":   "node specific",
}

func addNodeToBird(ip string) error {

	localip, err := GetInterfaceIpv4Addr("eth0")
	if err != nil {
		return fmt.Errorf("Error querying IP: ip %s, error: %v", localip, err)
	}

	dir := "/etc/bird/"

	f, err := os.Create(dir + ip + ".conf")

	check(err)

	defer f.Close()
	// # worker 1
	// protocol bgp {
	//    local 172.18.0.10 as 65000;
	//    neighbor 172.18.0.7 as 65000;
	//    direct;
	//    export none;
	// }
	n1, err := f.WriteString("protocol bgp {\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n1)
	n2, err := f.WriteString("   local " + localip + " as 65000;\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)
	n3, err := f.WriteString("   neighbor " + ip + " as 65000;\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n3)
	n4, err := f.WriteString("   export none;\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)
	n5, err := f.WriteString("}\n")
	fmt.Printf("wrote %d bytes\n", n5)
	check(err)

	err = reloadBirdConfig()
	check(err)

	return nil
}

func reloadBirdConfig() error {
	ipv := ""
	birdSuffix := ""
	if ipv == "6" {
		birdSuffix = "6"
	}
	// Try connecting to the BIRD socket in `/var/run/calico/` first to get the data
	birdSocket := fmt.Sprintf("/var/run/calico/bird%s.ctl", birdSuffix)
	if !socketFileExists(birdSocket) {
		// If that fails, try connecting to BIRD socket in `/var/run/bird` (which is the
		// default socket location for BIRD install) for non-containerized installs
		log.Debugln("Failed to connect to BIRD socket in /var/run/calico file not exists, trying /var/run/bird")
		birdSocket = fmt.Sprintf("/usr/local/var/run/bird%s.ctl", birdSuffix)
	}
	c, err := net.Dial("unix", birdSocket)
	if err != nil {
		return fmt.Errorf("Error querying BIRD: unable to connect to BIRDv%s socket: %v", ipv, err)
	}

	// To query the current state of the BGP peers, we connect to the BIRD
	// socket and send a "show protocols" message.  BIRD responds with
	// peer data in a table format.
	//
	// Send the request.
	_, err = c.Write([]byte("configure\n"))
	if err != nil {
		return fmt.Errorf("Error executing command: unable to write to BIRD socket: %s", err)
	}
	return nil
}

func verifyBirdConfig() error {
	ipv := ""
	birdSuffix := ""
	if ipv == "6" {
		birdSuffix = "6"
	}

	// Try connecting to the BIRD socket in `/var/run/calico/` first to get the data
	birdSocket := fmt.Sprintf("/var/run/calico/bird%s.ctl", birdSuffix)
	if !socketFileExists(birdSocket) {
		// If that fails, try connecting to BIRD socket in `/var/run/bird` (which is the
		// default socket location for BIRD install) for non-containerized installs
		log.Debugln("Failed to connect to BIRD socket in /var/run/calico file not exists, trying /var/run/bird")
		birdSocket = fmt.Sprintf("/var/run/bird/bird%s.ctl", birdSuffix)
	}
	c, err := net.Dial("unix", birdSocket)
	if err != nil {
		return fmt.Errorf("Error querying BIRD: unable to connect to BIRDv%s socket: %v", ipv, err)
	}

	// To query the current state of the BGP peers, we connect to the BIRD
	// socket and send a "show protocols" message.  BIRD responds with
	// peer data in a table format.
	//
	// Send the request.
	_, err = c.Write([]byte("show status\n"))
	if err != nil {
		return fmt.Errorf("Error executing command: unable to write to BIRD socket: %s", err)
	}
	return nil
}

func socketFileExists(file string) bool {
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}
