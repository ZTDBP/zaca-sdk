module github.com/ztdbp/zaca-sdk

go 1.17

require (
	github.com/cloudflare/backoff v0.0.0-20161212185259-647f3cdfc87a
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/valyala/fasthttp v1.15.1
	github.com/ztalab/zta-tools v0.0.2
	github.com/ztdbp/cfssl v0.0.5
	golang.org/x/crypto v0.5.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
)

require (
	github.com/andybalholm/brotli v1.0.0 // indirect
	github.com/google/certificate-transparency-go v1.1.4 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/spiffe/go-spiffe/v2 v2.0.0-beta.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/weppos/publicsuffix-go v0.13.0 // indirect
	github.com/zeebo/errs v1.2.2 // indirect
	github.com/zmap/zcrypto v0.0.0-20200911161511-43ff0ea04f21 // indirect
	github.com/zmap/zlint/v2 v2.2.1 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
)

replace (
	github.com/prometheus/prometheus v2.5.0+incompatible => github.com/prometheus/prometheus/v2 v2.29.2
	github.com/zmap/rc2 v0.0.0-20131011165748-24b9757f5521 => github.com/zmap/rc2 v0.0.0-20190804163417-abaa70531248
)
