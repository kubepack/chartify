package repo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	macaron "gopkg.in/macaron.v1"
	chartutil "k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	helmrepo "k8s.io/helm/pkg/repo"
)

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
	err = createLogFile(container)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer deleteLogFile(container)
	/*
		resp, err := gceSvc.Objects.Get(bucket, indexfile).Download()
	*/
	allIndex := helmrepo.NewIndexFile()
	item, err := container.Item(indexfile)
	if err != nil {
		log.Println("No index file found in the repository. Creating new index...\n")
	} else {
		r, err := item.Open()
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(r)
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		file, err := ioutil.TempFile(os.TempDir(), "index-")
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		ioutil.WriteFile(file.Name(), buf.Bytes(), 0644)
		defer os.Remove(file.Name())
		allIndex, err = helmrepo.LoadIndexFile(file.Name())
		if removeItemFromContainer(container, indexfile) != nil {
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
	_, err = uploadFileInContainer(container, indexfile, string(updatedIndex), nil)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	chartFileName, version := getChartFileNameAndVersion(chart.Data.Filename)
	if chartFileName == "" || version == "" {
		err := errors.New("Chart version not found")
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	f, err := chart.Data.Open()
	defer f.Close()
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
	}
	//upload the tar file
	name := chartFileName + "/" + version + "/" + chart.Data.Filename
	_, err = uploadFileInContainer(container, name, buf.String(), nil)
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
		name := chartFileName + "/" + version + "/" + v.name
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = uploadFileInContainer(container, name, v.data, nil)
		if err != nil {
			http.Error(ctx.Resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getChartFileNameAndVersion(fileName string) (string, string) {
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

func createLogFile(container stow.Container) error {
	i := 0
	var err error
	for i = 0; i <= 5; i++ {
		_, err = uploadFileInContainer(container, logfile, "", nil)
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	if i > 5 {
		return err
	}
	return nil
}

func deleteLogFile(container stow.Container) error {
	var err error
	for i := 0; i <= 5; i++ {
		err := container.RemoveItem(logfile)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return err
}

func uploadFileInContainer(container stow.Container, name, content string, md map[string]interface{}) (stow.Item, error) {
	buf := bytes.NewReader([]byte(content))
	return container.Put(name, buf, int64(len(content)), md)
}

func removeItemFromContainer(container stow.Container, file string) error {
	return container.RemoveItem(file)
}

func getProviderAndConfig() (string, stow.Config, error) {
	cred, err := getCredential()
	if err != nil {
		return "", stow.ConfigMap{}, err
	}
	config := stow.ConfigMap{ //TODO
		"json":       string(cred),
		"project_id": "tigerworks-kube",
	}
	return "google", config, nil
}
