package main

import (
	"flag"
	"fmt"
	"github.com/liserjrqlxue/RDMO/router"
	"log"
	"net/http"
	"os"
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
	_ = os.Mkdir("data", 0755)
	http.HandleFunc("/loadMO", router.LoadMO)
	http.HandleFunc("/updateMO", router.UpdateMO)

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
			router.Index(w, r)
		},
	)

	fmt.Printf("start http://localhost%v\n", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
