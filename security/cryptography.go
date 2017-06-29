/*
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package security

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"strconv"
)

const (
	minRSABitModulus          = 2048
	errorUnsupportedSignature = "unsupported private key signature, only RSA is valid"
)

// meetCryptoRequirements validates that all the requiered cryptographic requierements
// are applied to the public/private key pairs.
func meetCryptoRequirements(cert *tls.Certificate) error {

	err := validateSignatureMethod(cert)
	if err != nil {
		return err
	}

	err = validateCertBitLength(cert)
	if err != nil {
		return err
	}

	err = validateSignatureHashAlgorithm(cert)
	if err != nil {
		return err
	}

	return nil
}

// validateSignatureMethod identify the certificate digital signature for the private key.
// Only certificates signed with RSA are acceptable.
func validateSignatureMethod(cert *tls.Certificate) error {

	switch cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		return nil
	default:
		return errors.New(errorUnsupportedSignature)
	}

}

// validateCertBitLength evaluates the size of the private key according to the signed method.
// Expected length: 2048 bits or more.
func validateCertBitLength(cert *tls.Certificate) error {

	var len int
	switch key := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		len = key.N.BitLen()
	default:
		return errors.New(errorUnsupportedSignature)
	}

	if len < minRSABitModulus {
		return errors.New("Validating certificate length, expected " + strconv.Itoa(minRSABitModulus) + " bits but the key is " + strconv.Itoa(len))
	}

	return nil
}

// validateSignatureHashAlgorithm identify the hash function used as signature algorithm.
// Supported functions: SHA256, SHA384 or SHA512.
func validateSignatureHashAlgorithm(cert *tls.Certificate) error {

	c, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}

	signature := c.SignatureAlgorithm
	if signature == x509.SHA256WithRSA || signature == x509.SHA384WithRSA || signature == x509.SHA512WithRSA {
		return nil
	}

	return errors.New("x509: " + signature.String() + " is an unsupported signature hash algorithm")

}

// SetCipherSuites will receive a tls.Config pointer and will set CipherSuites array and PreferServerCipherSuites as true.
func SetCipherSuites(config *tls.Config) {
	config.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}
	config.PreferServerCipherSuites = true
}

// GetCertificateChain generates a TLS certificate from a valid and secure pair of PEM encoded files.
func GetCertificateChain(certFile, keyFile string) (tls.Certificate, error) {

	var cert tls.Certificate

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Validate if cryptographic requirements were applied properly
	err = meetCryptoRequirements(&cert)
	if err != nil {
		return tls.Certificate{}, err
	}
	log.Println("[INFO] Security cryptographic requirements for certificates were all PASSED!")

	return cert, nil
}
