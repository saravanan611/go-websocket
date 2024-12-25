package base

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func GenerateSelfSignedCert2() (lCert, lKey string, lErr error) {
	priv, lErr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if lErr != nil {
		return lCert, lKey, lErr
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumber, lErr := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if lErr != nil {
		return lCert, lKey, lErr
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Example Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, lErr := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if lErr != nil {
		return lCert, lKey, lErr
	}

	privBytes, lErr := x509.MarshalECPrivateKey(priv)
	if lErr != nil {
		return lCert, lKey, lErr
	}

	lCert = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}))
	lKey = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))

	return lCert, lKey, nil
}

func GenerateSelfSignedCert() (lCertPath, lKeyPath string, lErr error) {
	// Generate private key
	priv, lErr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if lErr != nil {
		return "", "", lErr
	}

	// Certificate validity period
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year validity

	// Serial number for the certificate
	serialNumber, lErr := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if lErr != nil {
		return "", "", lErr
	}

	// Certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Example Org"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Generate the certificate
	derBytes, lErr := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if lErr != nil {
		return "", "", lErr
	}

	// Marshal the private key
	privBytes, lErr := x509.MarshalECPrivateKey(priv)
	if lErr != nil {
		return "", "", lErr
	}

	// Define file paths
	lCertPath = "cert.pem"
	lKeyPath = "key.pem"

	// Save the certificate to a file
	certFile, lErr := os.Create(lCertPath)
	if lErr != nil {
		return "", "", lErr
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return "", "", err
	}

	// Save the private key to a file
	keyFile, lErr := os.Create(lKeyPath)
	if lErr != nil {
		return "", "", lErr
	}
	defer keyFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return "", "", err
	}

	return lCertPath, lKeyPath, nil
}
