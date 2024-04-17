package kubeconfigcafetch

// test
import (
	"bytes"
	b64 "encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"

	"encoding/pem"

	"go.step.sm/crypto/pemutil"
)

func GetBase64Result(client *http.Client, name string, url string, ch chan *Base64Result, wg *sync.WaitGroup) {
	result := Base64Result{name, url, "", time.Now().UnixNano()}

	result.Cert, _ = GetCertCaBase64(url, client, wg)
	result.time_ = time.Now().UnixNano() - result.time_

	ch <- &result
}

type Base64Result struct {
	Name  string
	Url   string
	Cert  string
	time_ int64
}

func GetCertCaBase64(url string, client *http.Client, wg *sync.WaitGroup) (ret string, err error) {
	defer wg.Done()
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	// The server will give you its entire cert chain.
	certs := resp.TLS.PeerCertificates

	// Some valid APIs won't have an ephemeral cert, they might not be rotating (?)
	// Then, go ahead and try certs[0] since it will probably also work. Probably.
	p, err := pemutil.Serialize(certs[0])

	log.Printf("%v has %v certificates", url, len(certs))

	// The above values are discarded in this branch:
	if len(certs) > 1 {
		// If there is a parent cert, prefer certs[1], (it is the self-signed cert.)
		// Use it instead, of certs[0] which is rotated and will be expired quickly.
		p, err = pemutil.Serialize(certs[1])
	}
	if err != nil {
		return "", err
	}
	// Encode the certificate as a PEM data
	var buf bytes.Buffer
	err = pem.Encode(&buf, p)
	if err != nil {
		return "", err
	}

	// Return the encoded PEM as base64 standard encoding
	str := b64.StdEncoding.EncodeToString(buf.Bytes())
	return str, nil
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
