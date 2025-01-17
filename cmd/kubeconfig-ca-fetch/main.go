package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	kcf "github.com/kingdon-ci/kubeconfig-ca-fetch"
)

var timeout = time.Duration(10 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func main() {
	client := &http.Client{
		// NB: This is not the timeout we needed!
		// Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Dial: dialTimeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	m := map[string]string{
		"admin@cozy":     "10.17.13.253:6443",
		"harvey":         "harvey.test.moomboo.space",
		"test":           "kubernetes-cluster.test.moomboo.space",
		"moo":            "moo-cluster.turkey.local",
		"mop":            "mop-cluster.turkey.local",
		"vcluster":       "vcluster-cluster.turkey.local",
	}
	out := map[string]string{}

	doHeavyLifting(client, m, out)
	printKubeconfig(m, out)
}

func doHeavyLifting(client *http.Client, m map[string]string, out map[string]string) {
	// result holds a cert from certs[0], or an empty string for cert
	ch := make(chan *kcf.Base64Result)
	wg := sync.WaitGroup{}

	// Call http routine as an asynchronous function
	for k, v := range m {
		wg.Add(1)
		url := fmt.Sprintf("https://%s", v)
		// getBase64Result always returns a result regardless of failure
		go kcf.GetBase64Result(client, k, url, ch, &wg)
	}

	wg.Wait()

	// m is the "input" map and it has the same length as the finished output map
	// but failed connections will be empty certs, get omitted from the kubeconfig
	kcf.FillOutputMap(m, out, ch)
}

func printKubeconfig(min map[string]string, mout map[string]string) {
	fmt.Println("apiVersion: v1\nclusters:")

	// clusters
	for name, cert := range mout {
		if cert != "" {
			fmt.Println("  - cluster:")
			fmt.Printf("        certificate-authority-data: %s\n", cert)
			fmt.Printf("        server: https://%s\n", min[name])
			fmt.Printf("    name: %s\n", name)
		}
	}

	fmt.Println("contexts:")
	for name, cert := range mout {
		if cert != "" {
			fmt.Println("  - context:")
			fmt.Printf("        cluster: %s\n", name)
			fmt.Printf("        user: kubelogin\n")
			fmt.Printf("    name: %s\n", name)
		}
	}

	fmt.Println(`kind: Config
preferences: {}
current-context: moo
users:
- name: kubelogin
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - oidc-login
      - get-token
      - --oidc-issuer-url=https://dex.harvey.moomboo.space
      - --oidc-client-id=wg-kubelogin
      - --oidc-client-secret=AiAImuXKhoI5ApvKWF988txjZ+6rG3S7o6X5En
      - --oidc-extra-scope=groups
      - --oidc-extra-scope=offline_access
      command: kubectl
      env: null
      provideClusterInfo: false`)
}
