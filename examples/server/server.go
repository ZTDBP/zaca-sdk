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
	"github.com/ztdbp/zaca-sdk/examples/util"
	"net"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/ztalab/zta-tools/logger"
	"github.com/ztalab/zta-tools/pkg/keygen"
	"github.com/ztalab/zta-tools/pkg/spiffe"
	"github.com/ztdbp/zaca-sdk/caclient"
)

var (
	caAddr   = flag.String("ca", "https://127.0.0.1:8081", "CA Server")
	ocspAddr = flag.String("ocsp", "http://127.0.0.1:8082", "Ocsp Server")
	addr     = flag.String("addr", ":6066", "")
	authKey  = "0739a645a7d6601d9d45f6b237c4edeadad904f2fce53625dfdd541ec4fc8134"
)

// go run server.go -ca https://127.0.0.1:8081 -ocsp http://127.0.0.1:8082

func main() {
	flag.Parse()
	err := NewMTLSServer()
	if err != nil {
		logger.Errorf("%v", err)
		return
	}
	select {}
}

// NewMTLSServer mTLS Server Use example
func NewMTLSServer() error {
	c := caclient.NewCAI(
		caclient.WithCAServer(caclient.RoleDefault, *caAddr),
		caclient.WithOcspAddr(*ocspAddr),
		caclient.WithAuthKey(authKey),
		caclient.WithLogger(logger.StandardLogger()),
		caclient.WithCSRConf(keygen.CSRConf{
			SNIHostnames: []string{"supreme"},
			IPAddresses:  []string{"10.10.10.10"},
		}),
	)
	ex, err := c.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "server1",
	})
	if err != nil {
		return errors.Wrap(err, "Exchanger initialization failed")
	}

	// Start certificate rotation
	go ex.RotateController().Run()

	cfger, err := ex.ServerTLSConfig()
	if err != nil {
		panic(err)
	}
	cfger.BindExtraValidator(func(identity *spiffe.IDGIdentity) error {
		fmt.Println("id: ", identity)
		return nil
	})
	tlsCfg := cfger.TLSConfig()
	tlsCfg.VerifyConnection = func(state tls.ConnectionState) error {
		fmt.Println("test state connection")
		return nil
	}
	go func() {
		httpsServer(tlsCfg)
	}()
	util.ExtractCertFromExchanger(ex)
	return nil
}

func httpsServer(cfg *tls.Config) {
	ln, err := net.Listen("tcp4", *addr)
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	lnTLS := tls.NewListener(ln, cfg)

	if err := fasthttp.Serve(lnTLS, func(ctx *fasthttp.RequestCtx) {
		str := ctx.Request.String()
		logger.Infof("Recv: %s", str)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte("Hello " + str))
	}); err != nil {
		panic(err)
	}
}
