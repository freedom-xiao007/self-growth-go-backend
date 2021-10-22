package main

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func initMongodb() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "phone_record", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initMongodb()
	var eg errgroup.Group

	// 一进程多端口
	insecureServer := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 9 * time.Second,
	}

	//secureServer := &http.Server{
	//	Addr:         "192.168.1.3:8443",
	//	Handler:      router(),
	//	ReadTimeout:  4 * time.Second,
	//	WriteTimeout: 9 * time.Second,
	//}

	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	//eg.Go(func() error {
	//	err := secureServer.ListenAndServeTLS("D:\\temp\\https\\server.crt", "D:\\temp\\https\\server_no_passwd.key")
	//	if err != nil && err != http.ErrServerClosed {
	//		log.Fatal(err)
	//	}
	//	return err
	//})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}