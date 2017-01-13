package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/appscode/chartify/pkg/repo"
	"github.com/ghodss/yaml"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/toolbox"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	bstore "google.golang.org/api/storage/v1"
	gcloud_gcs "google.golang.org/api/storage/v1"
	macaron "gopkg.in/macaron.v1"
	chartutil "k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	helmrepo "k8s.io/helm/pkg/repo"
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
			m.Post(pathPrefix+"/:container/", binding.MultipartForm(ChartFile{}), func(ctx *macaron.Context, chart ChartFile) {
				UploadChart(chart, ctx)
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

func UploadChart(chart ChartFile, ctx *macaron.Context) {
	url := ctx.Req.Host + ctx.Req.URL.Path
	file, err := chart.Data.Open()
	defer file.Close()
	c, err := chartutil.LoadArchive(file)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	hash, err := provenance.Digest(file)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	bucket := ctx.Params(":container")
	gceSvc, err := repo.GetGCEClient(gcloud_gcs.DevstorageReadWriteScope)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if the bucket exist or not
	_, err = gceSvc.Buckets.Get(bucket).Do()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	err = createLogFile(gceSvc, bucket)
	defer deleteLogFile(gceSvc, bucket)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := gceSvc.Objects.Get(bucket, indexfile).Download()
	allIndex := helmrepo.NewIndexFile()
	if err != nil {
		log.Println("No index file found in the repository. Creating new index...\n")
	} else {
		byteData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		allIndex, err = helmrepo.LoadIndex(byteData)
		//delete the old index
		err = gceSvc.Objects.Delete(bucket, indexfile).Do()
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	allIndex.Add(c.Metadata, c.Metadata.Name, url, "sha256:"+hash)
	updatedIndex, err := yaml.Marshal(allIndex)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	// Upload new index
	indexObject := &bstore.Object{
		Name: indexfile,
	}
	//Try again to upload if error TODO
	_, err = gceSvc.Objects.Insert(bucket, indexObject).Media(strings.NewReader(string(updatedIndex))).Do()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	container, version := getContainername(chart.Data.Filename)
	if container == "" || version == "" {
		err := errors.New("Chart version not found")
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	chartObject := &bstore.Object{
		Name:        container + "/" + version + "/" + chart.Data.Filename,
		ContentType: "application/x-compressed-tar",
	}
	f, err := chart.Data.Open()
	defer f.Close()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	//upload the tar file
	_, err = gceSvc.Objects.Insert(bucket, chartObject).Media(f).Do()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}

	//upload the untar file
	f1, err := chart.Data.Open()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	files, err := getFilesFromTar(f1)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range files {
		chartObject := &bstore.Object{
			Name: container + "/" + version + "/" + v.name,
		}
		_, err = gceSvc.Objects.Insert(bucket, chartObject).Media(strings.NewReader(string(v.data))).Do()
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getContainername(fileName string) (string, string) {
	s := strings.Trim(fileName, ".tgz")
	st := strings.SplitN(fileName, "-", 2)
	if len(st) == 2 {
		return st[0], s
	}
	return "", ""
}

func getFilesFromTar(in io.Reader) ([]*fileData, error) {
	files := []*fileData{}
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return files, err
	}
	defer unzipped.Close()
	tarReader := tar.NewReader(unzipped)
	i := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		name := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			continue
		default:
			buf := new(bytes.Buffer)
			buf.ReadFrom(tarReader)
			s := buf.String()
			file := &fileData{
				name: name,
				data: s,
			}
			files = append(files, file)
		}
		i++
	}
	return files, nil
}

func createLogFile(gceSvc *gcloud_gcs.Service, bucket string) error {
	logObject := &bstore.Object{
		Name: logfile,
	}
	i := 0
	for i = 0; i <= 5; i++ {
		_, err := gceSvc.Objects.Insert(bucket, logObject).Media(strings.NewReader(time.Now().String())).Do()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if i > 5 {
		return errors.New("Error. Try again...")
	}
	return nil
}

func deleteLogFile(gceSvc *gcloud_gcs.Service, bucket string) error {
	for i := 0; i <= 5; i++ {
		err := gceSvc.Objects.Delete(bucket, logfile).Do()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}
