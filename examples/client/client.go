/*
Copyright 2022-present The ZTDBP Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ztalab/zta-tools/logger"
	"github.com/ztalab/zta-tools/pkg/spiffe"
	"github.com/ztdbp/zaca-sdk/caclient"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	caAddr     = flag.String("ca", "https://127.0.0.1:8081", "CA Server")
	ocspAddr   = flag.String("ocsp", "http://127.0.0.1:8082", "Ocsp Server")
	serverAddr = flag.String("server", "https://127.0.0.1:6066", "")
	authKey    = "0739a645a7d6601d9d45f6b237c4edeadad904f2fce53625dfdd541ec4fc8134"
)

// go run server.go -ca https://127.0.0.1:8081 -ocsp http://127.0.0.1:8082 -server https://127.0.0.1:6066

func main() {
	flag.Parse()
	client, err := NewMTLSClient()
	if err != nil {
		logger.Fatalf("Client init error: %v", err)
	}
	ticker := time.Tick(time.Second)
	for i := 0; i < 1000; i++ {
		<-ticker

		resp, err := client.Get(*serverAddr)
		if err != nil {
			logger.Errorf("%v", err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Infof("Request result: %v, %s", resp.StatusCode, body)
	}
}

// mTLS Client Use example
func NewMTLSClient() (*http.Client, error) {
	c := caclient.NewCAI(
		caclient.WithCAServer(caclient.RoleDefault, *caAddr),
		caclient.WithAuthKey(authKey),
		caclient.WithOcspAddr(*ocspAddr),
		caclient.WithLogger(logger.StandardLogger()),
	)
	ex, err := c.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "client1",
	})
	if err != nil {
		return nil, errors.Wrap(err, "Exchanger initialization failed")
	}
	cfger, err := ex.ClientTLSConfig("supreme")
	if err != nil {
		panic(err)
	}
	cfger.BindExtraValidator(func(identity *spiffe.IDGIdentity) error {
		fmt.Println("id: ", identity.String())
		return nil
	})
	tlsCfg := cfger.TLSConfig()
	//tlsCfg.VerifyConnection = func(state tls.ConnectionState) error {
	//	cert := state.PeerCertificates[0]
	//	fmt.Println("Server certificate generation time: ", cert.NotBefore.String())
	//	return nil
	//}
	client := httpClient(tlsCfg)
	go ex.RotateController().Run()
	// util.ExtractCertFromExchanger(ex)

	resp, err := client.Get("http://www.google.com")
	if err != nil {
		panic(err)
	}

	fmt.Println("baidu test: ", resp.StatusCode)

	return client, nil
}

func httpClient(cfg *tls.Config) *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   cfg,
			DisableKeepAlives: true,
		},
	}
	return &client
}
