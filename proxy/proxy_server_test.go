package proxy_test

import (
	"crypto/tls"
	"go-proxy/proxy"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	certPEM = `-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIUQARYlwkXpzJFaWvnuC6gLsH1uhcwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDYyODIwMTEwNFoXDTI1MDYy
OTIwMTEwNFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAm41lGadkOY9Xci86OqiVo//axgGAXzUFF/KPAj2PosKJ
vcCB4aHCfWhcqsF9Ak18kYJ79ywHGYrN2B152I+kl0Ri4WMDFZJJ756U7g1FtP33
M7AU4hX+v29J0gXGbDwhWGtlJAeoF3Y3NtRSiw5DemYsU/W8EGtHLWkUS/paEMeJ
r1Y2aF5UAyHHgeBEohXO9d9OhPPrh3B0xJQjR/194/Z6f+qtZB+3H9nEU/+PFj9r
+Z3jsCO+okjsJkzTbCuTrIeNVLIKgb7/ThtHHPSlSOQcUIKIYqLaSLVVnl4XPxXP
B+L1ou08iNVUW8DKj64CuhhF0gCG8J3q9Uzqcjv3yQIDAQABo1MwUTAdBgNVHQ4E
FgQUd9xJuM49LxqKqAJAY1FXe7na9QcwHwYDVR0jBBgwFoAUd9xJuM49LxqKqAJA
Y1FXe7na9QcwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAJu4s
bGoByg9l6+hNcq5S272xPCr83a1w0gbM/RqGlxX92MMjcNzJrPahX26qSLpXe8r7
4H84+tBTvXegnliqo5dS5fVQ8Eu/j/ryTWLo4qxfOtVazt0ph93nPiBjPio0mBar
YAcesLEzMY8sMgsqPY/2inrAbBjT1Ejp1pFmFo6EdLiXtceorLhj/iJ0OwtAuEm4
HvqF7XsBkwiP18O4eU+m1HOyUON6AVE4yaqWq9Hbt9dAj4kBLWIcPfO9JG9pkqPG
zobxO080Ung4YPP57uHNLXbjZNQWqX/gAqWkheQugJNCwzwqKfgLaogozpUtvu2E
rsm5pPBL8b9N0ojaMg==
-----END CERTIFICATE-----`

	keyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCbjWUZp2Q5j1dy
Lzo6qJWj/9rGAYBfNQUX8o8CPY+iwom9wIHhocJ9aFyqwX0CTXyRgnv3LAcZis3Y
HXnYj6SXRGLhYwMVkknvnpTuDUW0/fczsBTiFf6/b0nSBcZsPCFYa2UkB6gXdjc2
1FKLDkN6ZixT9bwQa0ctaRRL+loQx4mvVjZoXlQDIceB4ESiFc71306E8+uHcHTE
lCNH/X3j9np/6q1kH7cf2cRT/48WP2v5neOwI76iSOwmTNNsK5Osh41UsgqBvv9O
G0cc9KVI5BxQgohiotpItVWeXhc/Fc8H4vWi7TyI1VRbwMqPrgK6GEXSAIbwner1
TOpyO/fJAgMBAAECggEAJ2iRmge+DRUjLuMHKhBG5QcvimWb+8LnePMzRvqIts/5
bJDmdt7v8qyreXlOrfQqoITB2lOVsuNnFh8VnQd7R/WD5Z0bjW4D/Eida5gCNoH3
DGnKSKMQ6LgNqD9dT1OGkSuYMqIb9GG3SFVPWjpxXwOerGC/1hsPq9II2dzhz4kQ
1KNjqlIME50iQrfdlD039OIwrz9DlqJSHUtvg/y4njHS1q+7aOD9p42jgGbjqxZD
/fQAw598ERU/AhUPYTm0Y5oOqgtsO+kYCmJ5d0BFZ6UXHFrvxU4rrfn34Iz5lDc8
JBioRejB5XNtYfRWbv6TG8O/jsr1R/0XcstdIt4ijwKBgQDMet8B65yCcmm7wZ9Y
R8+di6NykUg9uAHZ4QNnkk1QYC+jFddkVLCwdst0xyIEneQaQETwC8fa9pxkjHMv
AMknt1FZeWNemsV8ZO1KZVyYPFfBBB+5JcnnHOsGHC4fIE1PPnGPSrNp4nqW5nto
3Ns+mPTmocaRCkxoKrZtm+EeSwKBgQDCvqZlbBYkxLbgQ/ydcIDDv/FRVVbpMipY
AnqCr5I133Kr9RRDUqrVZ87qwGkxOoWoM+C7F8z3GMcWJ58jOVe2O0A0fLyrhu2b
7Xc6IBS+w2FR4Ok5uqW65jJnFztic1c8z9JUbWkLeKQ6Fk4z5JgJJRRtKCeUpBki
ZaVv1/4luwKBgQCHL7TAUETo+TtuJlRyyQc54VfuJp7cLwsKQPk/QDpdKTpVV9tP
Oa6W+/MHAaA77SchM9xf12oKGYDL+Q8txBc5arkdrmND6I8n7pHy3ZCaFUrvEQro
HVOeuD+pinfznCeAfIgXdAuptVHW8golCd7pQ7alw87DlUtuks6JKMVsgQKBgQCf
xdquiyb/s4R2KlEuugZqkydhGyra15V171KjtXe3S0PBYKjnMwOFYk2Yu5OSF/lg
Lm/KD5TRhTqRKqCdPYaAs8vRRCVmdKSssP6IaZmbiKBnlKbD/iXKWOIxQhYuh4Kj
Gb3uFnWAO9JA9dvjJ0C0//7qL2+Ju7gDSHGaeRLB0QKBgQCRms/16H+FG4WJt5PZ
lq7ACXt8ffM5SLlcSwMq8J7eJo8piTpaBfvpshWcp69Oybva7RUbemsC5DxsyU8z
teJaEEp/Z0eAc/Ea7No1S5SrBvjFispew33Rs7hWLYt42PYTn74NhSshCa8yVuyL
4ICulLHllpAmmz3le0iKIu5vQQ==
-----END PRIVATE KEY-----
`
)

func TestProxyTunnel(t *testing.T) {
	t.Run("Connects to the right destination", func(t *testing.T) {
		const (
			successfulStatusCode = 200
			proxyServerPort      = 3000
		)
		mockedSourceServer := setupMockedServer(t)
		defer mockedSourceServer.Close()

		proxyServer := proxy.NewServer(mockedSourceServer.URL)

		proxyRequest, err := http.NewRequest("GET", "/test", nil)

		if err != nil {
			t.Errorf("An unexpected error occurred: %+v", err)
		}

		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)
		// _, err := http.Get(mockedSourceServer.URL + "/test")

		if err != nil {
			t.Errorf("An unexpected error occurred: %+v", err)
		}

		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

	})
}

func TestProxyHandler(t *testing.T) {
	mockedServer := setupMockedServer(t)
	defer mockedServer.Close()

	t.Run("Strips the proxy headers", func(t *testing.T) {
		proxyServer := proxy.NewServer(mockedServer.URL)
		proxyRequest, err := http.NewRequest("GET", "/test", nil)

		if err != nil {
			t.Fatalf("An unexpected error occurred %+v", err)
		}

		// Attach proxy headers to the request
		proxyHeaderList := map[string]string{
			"Proxy-Connection":    "Keep-Alive",
			"Proxy-Authorization": "Basic",
		}
		for key, value := range proxyHeaderList {
			proxyRequest.Header.Add(key, value)
		}
		// Make Request
		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)

		// Check for headers starting with Proxy
		assertHeaderPrefixNotExist(t, rr.Result().Header, "proxy-")

	})

	t.Run("Attaches the correct proxy headers", func(t *testing.T) {
		proxyServer := proxy.NewServer(mockedServer.URL)
		proxyRequest, err := http.NewRequest("GET", "http://localhost:3001/test", nil)

		if err != nil {
			t.Errorf("An unepxected error occurred : %+v", err)
		}

		proxyRequest.Header.Add("X-Forwarded-Host", "dummyjson.com")

		t.Logf("HeaderList:\n")
		for key, value := range proxyRequest.Header {
			t.Logf(key + ":" + strings.Join(value, ","))
		}

		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)
	})

	t.Run("Returns the correct content headers", func(t *testing.T) {
		jsonMockServer := setupMockedContentServer(t)
		proxyServer := proxy.NewServer(jsonMockServer.URL)
		proxyRequest, err := http.NewRequest("GET", "http://localhost:3001/test", nil)

		if err != nil {
			t.Errorf("An error occurred while creating the test request: %+v", err)
		}

		proxyRequest.Header.Add("Accept", "application/json")
		rr := httptest.NewRecorder()

		proxyServer.ServeHTTP(rr, proxyRequest)
		contentTypeHeader := rr.Header().Get("Content-Type")

		if contentTypeHeader != "application/json" {
			t.Errorf("Expecting 'application/json' for 'Content-Type' Header, got '%+v'", contentTypeHeader)
		}

		t.Logf("Content-Type Header Value: %v", contentTypeHeader)

	})
}

func assertHeaderPrefixNotExist(t *testing.T, headerList http.Header, prefix string) {
	for key := range headerList {
		lowerCaseKey := strings.ToLower(key)
		if strings.HasPrefix(lowerCaseKey, prefix) {
			t.Errorf("Key: %v shouldn't be present in the response", key)
		}
	}
}

func setupMockedServer(t *testing.T) *httptest.Server {
	mockHttpHandler := func(w http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)

		if err != nil {
			w.WriteHeader(500)
			t.Errorf("An error occurred in mocked http handler : %+v", err)
		}

		for key, value := range request.Header {
			t.Logf("%+v -> %+v\n", key, value)
			w.Header()[key] = value
		}
		w.WriteHeader(200)
		w.Write(body)
	}

	return httptest.NewServer(http.HandlerFunc(mockHttpHandler))
}

func setupMockedContentServer(t *testing.T) *httptest.Server {
	mockHttpHandler := func(w http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)

		if err != nil {
			w.WriteHeader(500)
			t.Errorf("An error occurred in mocked http handler : %+v", err)
		}

		for key, value := range request.Header {
			t.Logf("%+v -> %+v\n", key, value)
			w.Header()[key] = value
		}

		w.Header()["Content-Type"] = []string{"application/json"}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	}

	return httptest.NewServer(http.HandlerFunc(mockHttpHandler))
}

func setupMockedTLSServer(t *testing.T) string {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		t.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "localhost:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				c.Write([]byte("Hello from TLS server\n"))
				c.Close()
			}(conn)
		}
	}()
	return ln.Addr().String()
}
