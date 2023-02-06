package kubeconfigcafetch

// test
import (
	"bytes"
	b64 "encoding/base64"
	"log"
	"net/http"
	"time"

	"encoding/pem"

	"go.step.sm/crypto/pemutil"
)

func GetBase64Result(client *http.Client, name string, url string, ch chan *Base64Result) {
	result := Base64Result{name, url, "", time.Now().UnixNano()}

	result.Cert, _ = GetCertCaBase64(url, client)
	result.time_ = time.Now().UnixNano() - result.time_

	ch <- &result
}

type Base64Result struct {
	Name  string
	Url   string
	Cert  string
	time_ int64
}

func GetCertCaBase64(url string, client *http.Client) (ret string, err error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	certs := resp.TLS.PeerCertificates
	if len(certs) > 1 {
		p, err := pemutil.Serialize(certs[1])
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
	} else {
		return "", err
	}
}

func FillOutputMap(m map[string]string, out map[string]string, ch chan *Base64Result) {
	// set doLog := true to enable logging to stderr
	doLog := true
	for i := 0; i < len(m); i++ {
		c := <-ch
		name := c.Name
		cert := c.Cert
		out[name] = cert
		// Only print to stderr if logging is enabled
		if doLog {
			if c.Cert == "" {
				log.Printf("Failed to reach %s (%s) after %d ms\n", c.Name, c.Url, c.time_/1e6)
			} else {
				log.Printf("Reached %s in %d ms\n", c.Url, c.time_/1e6)
			}
		}
	}
}
