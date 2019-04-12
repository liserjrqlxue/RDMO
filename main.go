package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	port = flag.String(
		"port",
		":9091",
		"web server listen port",
	)
)

var StaticDir = make(map[string]string)

func main() {
	flag.Parse()
	http.HandleFunc("/loadMO", loadMO)
	http.HandleFunc("/updateMO", updateMO)

	StaticDir["/static"] = "static"
	StaticDir["/public"] = "public"
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			// static file server
			for prefix, staticDir := range StaticDir {
				if strings.HasPrefix(r.URL.Path, prefix) {
					file := staticDir + r.URL.Path[len(prefix):]
					fmt.Println(file)
					http.ServeFile(w, r, file)
					return
				}
			}
			index(w, r)
		},
	)

	fmt.Printf("start http://localhost%v\n", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
