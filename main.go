package main

import (
	"flag"
	"io/ioutil"
	"os"
	"github.com/vinkdong/gox/log"
	"gopkg.in/yaml.v2"
	"net/url"
	"net/http"
	"fmt"
	"net/http/httputil"
	"strconv"
)

var(
	config = flag.String("conf","","Config file")
)


type Config struct {
	Server []Server `yaml:"server"`
}

type Server struct {
	Name        string
	BindPort    int `yaml:"bind_port"`
	LocalPort   int `yaml:"local_port"`
	LocalDomain string `yaml:"local_domain"`
	Tls         bool
}

func main() {
	flag.Parse()

	if *config == "" {
		log.Error("must special a config file")
		os.Exit(127)
	}

	content, err := ioutil.ReadFile(*config)
	if err != nil{
		log.Error(err)
	}

	c := &Config{}

	err = yaml.Unmarshal(content,c)

	fmt.Println(c)

	if err != nil{
		log.Error(err)
		os.Exit(127)
	}

	for _, vv := range c.Server {
		v := vv
		go func() {
			scheme := "http"
			if v.Tls {
				scheme = "https"
			}
			host := fmt.Sprintf("%s://%s:%d", scheme, v.LocalDomain, v.LocalPort)
			fmt.Println(host)
			url, err := url.Parse(host)
			if err != nil {
				log.Error(err)
				os.Exit(127)
			}

			server := &http.Server{
				Addr:    ":" + strconv.Itoa(v.BindPort),
				Handler: httputil.NewSingleHostReverseProxy(url),
			}

			err = server.ListenAndServe()
			if err != nil {
				log.Error(err)
			}
		}()
	}

	err = http.ListenAndServe(":10000",nil)
	if err != nil{
		log.Error(err)
	}

}