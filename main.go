package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func init() {

}

var (
	ServiceName = "golang-docker-demo"
	port        = ":8100"
	portssl     = ":443"
	certsPath   = "certs.cert"
	keyPath     = "certs.key"
)

func Run(addr string, sslAddr string, ssl map[string]string) chan error {

	errs := make(chan error)

	// Starting HTTP server
	go func() {
		log.Printf("Staring HTTP service on %s ...", addr)

		if err := http.ListenAndServe(addr, nil); err != nil {
			errs <- err
		}

	}()

	// Starting HTTPS server
	go func() {
		log.Printf("Staring HTTPS service on %s ...", sslAddr)
		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], nil); err != nil {
			errs <- err
		}
	}()

	return errs
}

func main() {
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(1000)
		fmt.Println("ping ok!")
		log.Println("logs ping ok!")
		w.Write([]byte(fmt.Sprintf("ping test ok %s", ServiceName)))
	})
	errs := Run(port, portssl, map[string]string{
		"cert": viper.GetString(certsPath),
		"key":  viper.GetString(keyPath),
	})

	select {
	case err := <-errs:
		log.Printf("Could not start serving service due to (error: %s)", err)
	}
}
