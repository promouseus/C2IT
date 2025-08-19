/*
Copyright © 2025 00010110 B.V. input@00010110.nl
*/
package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// performerCmd represents the performer command
var performerCmd = &cobra.Command{
	Use:   "performer",
	Short: "Runs the C2IT performer HTTP/2 service",
	Long: `The performer provides an embedded HTTPS server that uses HTTP/2 with 
a self-signed certificate, and streams real-time telemetry data using SSE to Chrome/Firefox clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("performer: HTTP/2 server starting on https://localhost:8443")

		// Ensure TLS certificate is ready
		cert := loadOrCreateTLSCert()

		// Serve static files from /web (index.html)
		http.Handle("/", http.FileServer(http.Dir("web")))

		// Server-Sent Events endpoint
		http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Streaming not supported", http.StatusInternalServerError)
				return
			}

			// Set SSE headers
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			log.Println("📡 SSE client connected")

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-r.Context().Done():
					log.Println("SSE client disconnected")
					return
				case t := <-ticker.C:
					fmt.Fprintf(w, "data: Hello at %s\n\n", t.Format(time.RFC3339))
					flusher.Flush()
				}
			}
		})

		// Action endpoint for client fetch
		http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Action received from client")
			w.Write([]byte("Server received your action"))
		})

		// Start the HTTPS HTTP/2 server
		server := &http.Server{
			Addr:      ":8443",
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		}

		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(performerCmd)
}

// loadOrCreateTLSCert loads TLS cert from disk or generates a new self-signed EC cert
func loadOrCreateTLSCert() tls.Certificate {
	const certFile = "cert.pem"
	const keyFile = "key.pem"

	// Return existing cert if found
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err == nil {
				log.Println("Using existing TLS certificate")
				return cert
			}
		}
	}

	log.Println("Generating new self-signed TLS certificate...")

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	writeFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to marshal EC key: %v", err)
	}
	writeFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b}))

	log.Println("TLS certificate created successfully")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load generated cert: %v", err)
	}
	return cert
}

// writeFile writes content to a file with proper permissions
func writeFile(path string, data []byte) {
	err := os.WriteFile(filepath.Clean(path), data, 0600)
	if err != nil {
		log.Fatalf("Could not write file %s: %v", path, err)
	}
}
