package main

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"net/http"

	"encoding/pem"

	"go.step.sm/crypto/pemutil"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
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
	}
	mout := map[string]string{}

	for k, v := range m {
		url := fmt.Sprintf("https://%s", v)
		ca, err := getCertCaBase64(url, client)
		if err != nil {
			// log.Println(err)
			// return
		}
		mout[k] = ca
	}

	// Let's print our Kubeconfig
	fmt.Println("apiVersion: v1\nclusters:")

	// clusters
	for name, cert := range mout {
		if cert != "" {
			fmt.Println("  - cluster:")
			fmt.Printf("        certificate-authority-data: %s\n", cert)
			fmt.Printf("        server: https://%s\n", m[name])
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

func getCertCaBase64(url string, client *http.Client) (ret string, err error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	certs := resp.TLS.PeerCertificates
	p, err := pemutil.Serialize(certs[0])
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, p)
	if err != nil {
		return "", err
	}

	str := b64.StdEncoding.EncodeToString(buf.Bytes())
	return str, nil
}
