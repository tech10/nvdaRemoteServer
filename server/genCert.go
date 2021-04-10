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
	"os"
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

	gen_cert_file(gencertfile, certPEM.Bytes(), certPrivKeyPEM.Bytes())

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return serverTLSConf, nil
}

func gen_cert_file(file string, cert, key []byte) {
	if file == "" {
		return
	}
	Log(LOG_DEBUG, "Attempting to write certificate to file "+file)
	w, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		Log_error("Unable to create or open file for writing certificate information.\r\n", err.Error())
		Launch_fail()
		return
	}
	_, err = w.Write(append(key, cert...))
	if err != nil {
		Log_error("Unable to write to the file " + file + "\r\n" + err.Error())
		Launch_fail()
	}
	err = w.Close()
	if err != nil {
		Log_error("The file at " + file + " was unable to close. Information may not have been written to it correctly.\r\n" + err.Error())
		Launch_fail()
	}
	Log(LOG_DEBUG, "Certificate and key successfully written to "+file)
}
