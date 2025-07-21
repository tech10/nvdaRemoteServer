package server

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// Generate a self-signed certificate as long as the server is running.
func serial_number() *big.Int {
	serialNumLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial_num, serial_err := rand.Int(rand.Reader, serialNumLimit)
	if serial_err != nil {
		return big.NewInt(time.Now().UnixNano())
	}
	return serial_num
}

func gen_cert() (*tls.Config, error) {
	ca := &x509.Certificate{
		SerialNumber: serial_number(),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"NVDARemote Server"},
			CommonName:   "Root CA",
		},
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	keyID := sha1.Sum(pubKeyBytes)
	ca.SubjectKeyId = keyID[:]
	ca.AuthorityKeyId = keyID[:]
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

	mpk, merr := x509.MarshalPKCS8PrivateKey(priv)
	if merr != nil {
		return nil, merr
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: mpk,
	})
	if err != nil {
		return nil, err
	}

	serverCert, serr := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if serr != nil {
		return nil, serr
	}

	gen_cert_file(gencertfile, certPEM.Bytes(), certPrivKeyPEM.Bytes())

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return serverTLSConf, nil
}

func gen_cert_file(file string, cert, key []byte) {
	if default_gen_cert_file(file) {
		return
	}
	Log(LOG_DEBUG, "Attempting to write certificate to file "+file)
	err := file_rewrite(file, append(key, cert...))
	if err != nil {
		Log_error("Failed to write certificate.\n" + err.Error())
		Launch_fail()
		return
	}
	Log(LOG_DEBUG, "Certificate and key successfully written to "+file)
}
