package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
)

const (
	hostfileAreaStart = "#begin-localhost"
	hostfileAreaEnd   = "#end-localhost"
)

type config struct {
	Hosts map[string]struct {
		Service string `yaml:"Service"`
	} `yaml:"Hosts"`
	Services map[string]struct {
		Install string `yaml:"Install"`
	} `yaml:"Services"`
}

func main() {
	// load config
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// update hosts
	hosts := make([]string, 0)
	for k := range cfg.Hosts {
		hosts = append(hosts, k)
	}

	updateHosts(hosts)

	// serve http
	serve()
}

func loadConfig() (*config, error) {
	// load config
	data, err := ioutil.ReadFile("/etc/localhosts.yaml")
	if err != nil {
		return nil, err
	}

	cfg := &config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func serve() {
	mux := mux.NewRouter()
	//mux.Host(host).Handler(http.FileServer(http.Dir(path)))

	serv := http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      mux,
		Addr:         "127.0.0.1:80",
	}
	log.Println("Listening on http://" + serv.Addr)
	log.Println(serv.ListenAndServe())
}

func honk(w http.ResponseWriter, req *http.Request) {
	params := url.Values{}
	params.Add("q", strings.TrimLeft(req.URL.Path, "/"))

	http.Redirect(w, req, "https://duckduckgo.com/lite?"+params.Encode(), http.StatusTemporaryRedirect)
}

func forward(targetURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		url, _ := url.Parse(targetURL)

		proxy := httputil.NewSingleHostReverseProxy(url)
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.Host = url.Host

		proxy.ServeHTTP(w, req)
	}
}

func updateHosts(hosts []string) {
	formattedHosts := []byte(hostfileAreaStart + "\n" +
		strings.Join(addLocalhost(hosts), "\n") + "\n" + hostfileAreaEnd)

	// read hosts
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_EXCL, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// read file
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	// replace localhost area...
	re := regexp.MustCompile(`(?m)^` + hostfileAreaStart + `\n([\s\S]*?\n)?` +
		hostfileAreaEnd + `$`)
	replaced := false
	result := re.ReplaceAllFunc(bytes, func(match []byte) []byte {
		if replaced {
			return []byte("")
		}

		replaced = true
		return formattedHosts
	})

	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(result)

	// ... or create a new one
	if !replaced {
		file.Write(formattedHosts)
	}

	log.Println("Updated /etc/hosts")
}

func addLocalhost(items []string) []string {
	result := make([]string, len(items))
	for i, x := range items {
		result[i] = "127.0.0.1 " + x
	}
	return result
}
