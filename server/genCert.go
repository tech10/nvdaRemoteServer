package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// Generate a self-signed certificate as long as the server is running
func serial_number() *big.Int {
	serial_num, serial_err := rand.Int(rand.Reader, big.NewInt(999999999999))
	if serial_err != nil {
		return big.NewInt(345098734305)
	}
	return serial_num
}

func gen_cert() (*tls.Config, error) {
	var ca = &x509.Certificate{
		SerialNumber: serial_number(),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"NVDARemote Server"},
			CommonName:   "Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, cerr := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)
	if cerr != nil {
		return nil, cerr
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}
	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	if err != nil {
		return nil, err
	}

	serverCert, serr := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if serr != nil {
		return nil, serr
	}

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return serverTLSConf, nil
}
