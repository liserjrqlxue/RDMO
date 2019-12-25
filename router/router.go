package router

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/liserjrqlxue/simple-util"
)

// os
var (
	ex, _        = os.Executable()
	exPath       = filepath.Dir(ex)
	pSep         = string(os.PathSeparator)
	templatePath = exPath + pSep + "template" + pSep
)

type Infos struct {
	Img     string
	Src     string
	Token   string
	Title   string
	Message string
}

func md5sum(str string) string {
	byteStr := []byte(str)
	sum := md5.Sum(byteStr)
	sumStr := fmt.Sprintf("%x", sum)
	return sumStr
}

func createToken() string {
	// token
	return md5sum(strconv.FormatInt(time.Now().Unix(), 10))
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, " method:", r.Method)

	t, err := template.ParseFiles(templatePath+"header.html", templatePath+"footer.html", templatePath+"index.html")
	simple_util.CheckErr(err)

	var Info Infos
	Info.Title = "Home Page"
	Info.Token = createToken()
	err = t.ExecuteTemplate(w, "index", Info)
	if err != nil {
		log.Print(err)
	}
}

func LoadMO(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, " method:", r.Method)

	t, err := template.ParseFiles(templatePath+"header.html", templatePath+"footer.html", templatePath+"loadMO.html")
	simple_util.CheckErr(err)

	var Info Infos
	Info.Title = "研发MO"
	Info.Token = createToken()
	err = t.ExecuteTemplate(w, "loadMO", Info)
	if err != nil {
		log.Print(err)
	}
}

func UpdateMO(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, " method:", r.Method)

	t, err := template.ParseFiles(templatePath+"header.html", templatePath+"footer.html", templatePath+"updateMO.html")
	simple_util.CheckErr(err)

	var Info Infos
	Info.Title = "更新MO"
	Info.Token = createToken()

	if r.Method == "POST" {
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Print(err)
			return
		}
		file, handler, err := r.FormFile("uploadfile")
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(file)
		//Info.Message=fmt.Sprint(handler.Header)
		f, err := os.Create("data" + pSep + handler.Filename)
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(f)
		_, err = io.Copy(f, file)
		simple_util.CheckErr(err)
		_, mapArray := simple_util.Sheet2MapArray("data"+pSep+handler.Filename, "研发领料")
		var db = map[string][]map[string]string{
			"data": mapArray,
		}
		jsonByte, err := json.MarshalIndent(db, "", "\t")
		simple_util.CheckErr(err)
		simple_util.CheckErr(simple_util.Json2file(jsonByte, "public"+pSep+"MO.json"))
		Info.Message = "update done!"
	}
	simple_util.CheckErr(t.ExecuteTemplate(w, "updateMO", Info))
}
