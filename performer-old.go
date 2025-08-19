/*
Copyright © 2025 00010110 B.V.
input@00010110.nl
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
	"net"
	"net/http"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/spf13/cobra"
)

// performerCmd represents the performer command
var performerCmd = &cobra.Command{
	Use:   "performer",
	Short: "Start a telemetry server using HTTP/2 and HTTP/3 (QUIC)",
	Long: `The performer command launches an embedded server that supports
both HTTP/2 and HTTP/3. It streams a simple message to clients and shows protocol details.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("performer: Starting H2 and H3 server on :4433")

		// TLS config with EC certificate
		tlsConf := generateSelfSignedECCert()

		// Shared handler mux
		mux := http.NewServeMux()

		// Static web UI
		mux.Handle("/", http.FileServer(http.Dir("web")))

		// Streaming endpoint with protocol info
		mux.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Client connected over %s", r.Proto)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming unsupported", http.StatusInternalServerError)
				return
			}

			_, _ = fmt.Fprintf(w, "Connected over protocol: %s\n", r.Proto)
			flusher.Flush()

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-r.Context().Done():
					log.Println("Client disconnected")
					return
				case <-ticker.C:
					_, _ = fmt.Fprintf(w, "hello (%s)\n", time.Now().Format("15:04:05"))
					flusher.Flush()
				}
			}
		})

		// Start H2 server (fallback)
		go func() {
			server := &http.Server{
				Addr:      ":4433",
				Handler:   mux,
				TLSConfig: tlsConf,
			}
			log.Println("Listening (HTTP/2) on https://localhost:4433 ...")
			if err := server.ListenAndServeTLS("", ""); err != nil {
				log.Fatalf("H2 server failed: %v", err)
			}
		}()

		// Start H3 server (QUIC)
		h3 := http3.Server{
			Addr:      ":4433",
			Handler:   mux,
			TLSConfig: tlsConf,
		}
		log.Println("🚀 Listening (HTTP/3 QUIC) on https://localhost:4433 ...")
		if err := h3.ListenAndServe(); err != nil {
			log.Fatalf("H3 server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(performerCmd)
}

func generateSelfSignedECCert() *tls.Config {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate EC key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("failed to create cert: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("failed to marshal EC key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("failed to load cert pair: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "h3"},
	}
}
