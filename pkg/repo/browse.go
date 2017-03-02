package repo

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	"gopkg.in/macaron.v1"
)

type dataFile struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Href        string `json:"href"`
	ContentType string `json:"contentType"`
	TimeCreated string `json:"timeCreated"`
	Updated     string `json:"updated"`
	Size        uint64 `json:"size"`
	Md5Hash     string `json:"md5Hash"`
}

// BucketOptions is a struct for specifying configuration options for the macaron GCS StaticBucket middleware.
type BucketOptions struct {
	PathPrefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
}

// Static returns a middleware handler that serves static files in the given directory.
func StaticBucket(bucketOpt ...BucketOptions) macaron.Handler {
	var opt BucketOptions
	if len(bucketOpt) > 0 {
		opt = bucketOpt[0]
	}
	if opt.Expires == nil {
		opt.Expires = func() string {
			return time.Now().Add(30 * 24 * time.Hour).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		}
	}
	return func(ctx *macaron.Context, r macaron.Render, log *log.Logger) {
		if ctx.Req.Method != "GET" && ctx.Req.Method != "HEAD" {
			http.Error(ctx.Resp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		bucket := ctx.Params(":container")
		bucketPath := strings.Replace(ctx.Req.URL.Path, opt.PathPrefix+"/"+bucket+"/", "", 1)
		provider, config, err := getProviderAndConfig()
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		location, err := stow.Dial(provider, config)
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		defer location.Close()
		container, err := location.Container(bucket)
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}

		if !opt.SkipLogging {
			log.Println("[Static] Serving " + ctx.Req.URL.Path + " from " + bucketPath)
		}

		if strings.HasSuffix(bucketPath, "/") {
			items, _, err := container.Items(bucketPath, "/", "", 5000)
			if err != nil {
				http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
				return
			}
			files := make([]*dataFile, 0)
			/*			for _, folder := range objs.Prefixes { TODO Folder option is missing from stow
						files = append(files, &dataFile{
							Name: folder,
							Type: "FOLDER",
							Href: fmt.Sprintf("%s/%s/%s", opt.PathPrefix, bucket, folder),
						})
					}*/

			for _, file := range items {
				if file.Name() != bucketPath {
					f := &dataFile{
						Name: file.Name(),
						Type: "FILE",
						Href: ctx.Req.URL.Path + file.Name()[strings.LastIndex(file.Name(), "/")+1:],
						//ContentType: file.ContentType,
						//TimeCreated: file.TimeCreated,
						//Updated:     file.Updated,
						//Md5Hash:     file.Md5Hash,
					}
					size, err := file.Size()
					if err != nil {
						f.Size = uint64(size)
					}
					files = append(files, f)
				}
			}
			ctx.JSON(200, files)
			return
		} else {
			// Add an Expires header to the static content
			if opt.Expires != nil {
				ctx.Resp.Header().Set("Expires", opt.Expires())
			}

			// load file
			item, err := container.Item(bucketPath)
			if err != nil {
				http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
				return
			}
			r, err := item.Open()
			if err != nil {
				http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(ctx.Resp, r)
			r.Close()
			if err != nil {
				http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func getCredential() ([]byte, error) {
	cred, err := ioutil.ReadFile("/home/sauman/Downloads/tigerworks-kube-3803f9d609c7.json")
	if err != nil {
		return nil, err
	}
	return cred, nil
}
