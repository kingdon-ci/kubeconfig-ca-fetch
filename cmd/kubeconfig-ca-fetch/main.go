package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	kcf "github.com/kingdon-ci/kubeconfig-ca-fetch"
)

var timeout = time.Duration(2 * time.Second)

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
		"cluster-01":     "cluster-01.turkey.local",
		"demo-cluster-2": "demo-cluster-2.turkey.local",
		"demo-cluster":   "demo-cluster.turkey.local",
		"hephy-stg":      "hephy-stg.turkey.local",
		"howard-space":   "howard.moomboo.space",
		"howard-stage":   "howard.moomboo.stage",
		"moo":            "moo-cluster.turkey.local",
		"vcluster":       "vcluster.turkey.local",
		"somtochi":       "somtochi.turkey.local",
		"another-test":   "another-test.turkey.local",
		"limnocentral":   "limnocentral.turkey.local",
		"management":     "loft.loft.svc.cluster.local",
	}
	// result holds a cert from certs[0], or an empty string for cert
	ch := make(chan *kcf.Base64Result)

	// Call http routine as an asynchronous function
	for k, v := range m {
		url := fmt.Sprintf("https://%s", v)
		// getBase64Result always returns a result regardless of failure
		go kcf.GetBase64Result(client, k, url, ch)
	}

	// m is the "input" map and it has the same length as the finished output map
	// but failed connections will be empty certs, get omitted from the kubeconfig
	out := map[string]string{}
	kcf.FillOutputMap(m, out, ch)
	printKubeconfig(m, out)
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
current-context: howard-space
users:
  - name: kubelogin
    user:
        auth-provider:
            config:
                client-id: weave-gitops
                client-secret: AiAImuXKhoI5ApvKWF988txjZ+6rG3S7o6X5En
                extra-scopes: groups
                idp-issuer-url: https://dex.howard.moomboo.space
            name: oidc`)
}
