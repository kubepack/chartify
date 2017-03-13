package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/appscode/chartify/pkg/repo"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/toolbox"
	"github.com/spf13/cobra"
	macaron "gopkg.in/macaron.v1"
)

func NewCmdServe() *cobra.Command {
	var (
		port       = 50080
		pprofPort  = 6060
		caCertFile string
		certFile   string
		keyFile    string
		pathPrefix string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a Chart repository server",
		Run: func(cmd *cobra.Command, args []string) {
			m := macaron.New()
			m.Use(macaron.Logger())
			m.Use(macaron.Recovery())
			m.Use(toolbox.Toolboxer(m))
			m.Use(macaron.Renderer(macaron.RenderOptions{
				IndentJSON: true,
			}))
			pathPrefix = strings.Trim(pathPrefix, "/")
			if pathPrefix != "" {
				pathPrefix = "/" + pathPrefix
			}
			handler := repo.StaticBucket(repo.BucketOptions{PathPrefix: pathPrefix})
			m.Get(pathPrefix+"/:container/", handler)
			m.Get(pathPrefix+"/:container/*", handler)
			m.Post(pathPrefix+"/:container/", binding.MultipartForm(repo.ChartFile{}), func(ctx *macaron.Context, chart repo.ChartFile) {
				repo.UploadChart(chart, ctx)
			})
			log.Println("Listening on port", port)
			srv := &http.Server{
				Addr:         fmt.Sprintf(":%d", port),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				Handler:      m,
			}
			if caCertFile == "" && certFile == "" && keyFile == "" {
				log.Fatalln(srv.ListenAndServe())
			} else {
				/*
					Ref:
					 - https://blog.cloudflare.com/exposing-go-on-the-internet/
					 - http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
					 - http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
				*/
				tlsConfig := &tls.Config{
					PreferServerCipherSuites: true,
					MinVersion:               tls.VersionTLS12,
					SessionTicketsDisabled:   true,
					CipherSuites: []uint16{
						tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
						// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					},
					ClientAuth: tls.VerifyClientCertIfGiven,
				}
				if caCertFile != "" {
					caCert, err := ioutil.ReadFile(caCertFile)
					if err != nil {
						log.Fatal(err)
					}
					caCertPool := x509.NewCertPool()
					caCertPool.AppendCertsFromPEM(caCert)
					tlsConfig.ClientCAs = caCertPool
				}
				tlsConfig.BuildNameToCertificate()
				srv.TLSConfig = tlsConfig
				log.Fatalln(srv.ListenAndServeTLS(certFile, keyFile))
			}
		},
	}
	cmd.Flags().IntVar(&port, "api-port", port, "Port used to serve repository")
	cmd.Flags().IntVar(&pprofPort, "pprof-port", pprofPort, "port used to run pprof tools")
	cmd.Flags().StringVar(&caCertFile, "caCertFile", caCertFile, "File containing CA certificate")
	cmd.Flags().StringVar(&certFile, "certFile", certFile, "File container server TLS certificate")
	cmd.Flags().StringVar(&keyFile, "keyFile", keyFile, "File containing server TLS private key")
	cmd.Flags().StringVar(&pathPrefix, "path-prefix", "/charts", "Path prefix for chart repositories")
	return cmd
}
