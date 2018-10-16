package main

import (
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"net/http"
	"flag"
	"net"
	"log"
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
)

func main() {
	serverName := flag.String(`ServerName`, ``, `server name`)
	certType := flag.String(`CertType`, `ecdsa`, `cert type[rsa/ecdsa]`)
	outFileName := flag.String(`OutFile`, ``, `output file, default ServerName.cert/ServerName.key`)

	flag.Parse()
	if *serverName == `` {
		flag.Usage()
		return
	}
	if *outFileName == `` {
		*outFileName = *serverName
	}
	mgr := autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}
	listenOk := make(chan struct{})
	go func() {
		ln, err := net.Listen(`tcp`, ":80")
		panicIfError(err)
		close(listenOk)

		err = http.Serve(ln, mgr.HTTPHandler(nil))
		panicIfError(err)
	}()

	<-listenOk
	log.Println(`update cert of`, *serverName, *certType)
	var err error
	var cert *tls.Certificate
	if *certType == `rsa` {
		cert, err = mgr.GetCertificate(&tls.ClientHelloInfo{
			ServerName: *serverName,
		})
	} else {
		cert, err = mgr.GetCertificate(&tls.ClientHelloInfo{
			ServerName: *serverName,
			SignatureSchemes: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			},
		})
	}
	panicIfError(err)
	saveCertFile(cert, *outFileName+".cert", *outFileName+".key")
	log.Println(`done`)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func saveCertFile(cert *tls.Certificate, certFileName string, keyFileName string) {
	var err error
	// contains PEM-encoded data
	var buf bytes.Buffer

	// private
	switch key := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		err = encodeECDSAKey(&buf, key)
		panicIfError(err)
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(key)
		pb := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
		err = pem.Encode(&buf, pb)
		panicIfError(err)
	default:
		panic("acme/autocert: unknown private key type")
	}

	err = ioutil.WriteFile(keyFileName, buf.Bytes(), 0600)
	panicIfError(err)

	buf.Reset()
	// public
	for _, b := range cert.Certificate {
		pb := &pem.Block{Type: "CERTIFICATE", Bytes: b}
		err = pem.Encode(&buf, pb)
		panicIfError(err)
	}
	err = ioutil.WriteFile(certFileName, buf.Bytes(), 0600)
	panicIfError(err)
}

func encodeECDSAKey(w io.Writer, key *ecdsa.PrivateKey) error {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	return pem.Encode(w, pb)
}

