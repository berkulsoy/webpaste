package main

import (
	"flag"
	"fmt"
	"net/http"
	"io"
	"os"
	"bufio"
	"time"
	"strconv"
	"io/ioutil"
	"github.com/dustinkirkland/golang-petname"
	"path"
	"math/rand"
)

var (
	port = flag.String("port", "8290", "Port number to listen")
	tmpdir = flag.String("dir", "/tmp/", "Tmp directory to save files")
)

type RequestLog struct {
	code int
	method string
	message string
}


func Log(severity string, log RequestLog) {
	fmt.Printf("%s %s %d %s %s\n", time.Now(), severity, log.code, log.method, log.message)
}


func Get(writer http.ResponseWriter, request *http.Request) {
	outfile, err := ioutil.ReadFile(*tmpdir + path.Base(request.URL.Path))

	if err != nil {
		http.Error(writer, "Use following in order to paste:\ncurl -F \"f=@path_to_my_file\" http://" + request.Host, http.StatusBadRequest)
	} else {
		Log("INFO", RequestLog{http.StatusOK, http.MethodGet,request.URL.String()})
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Disposition", "attachment; filename=" + path.Base(request.URL.Path))
		writer.Header().Set("Content-Type", http.DetectContentType(outfile))
		writer.Write(outfile)
	}
}

func Post(writer http.ResponseWriter, request *http.Request) {
	infile, header, err := request.FormFile("f")
	var newfile string
	if err != nil {
		Log("ERROR", RequestLog{http.StatusInternalServerError, http.MethodPost, err.Error()})
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	} else {
		defer infile.Close()
		bufReader := bufio.NewReader(infile)
		extension := path.Ext(header.Filename)

		if extension != "" {
			newfile = *tmpdir + petname.Generate(2, "-") + extension
		} else {
			newfile = *tmpdir + petname.Generate(2, "-")
		}

		outfile, _ := os.Create(newfile)
		defer outfile.Close()
		bufWriter := bufio.NewWriter(outfile)
		outfilesize, err := io.Copy(bufWriter, bufReader)
		bufWriter.Flush()

		if err != nil {
			Log("ERROR", RequestLog{http.StatusInternalServerError, http.MethodPost, err.Error()})
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			Log("INFO", RequestLog{http.StatusCreated, http.MethodPost,header.Filename + " " + outfile.Name() + " " + strconv.FormatInt(outfilesize, 10)})
			writer.WriteHeader(http.StatusCreated)
			writer.Write([]byte("You can download your file at:\nhttp://" + request.Host + "/" + path.Base(newfile) + "\n"))
		}
	}
}

func Others(writer http.ResponseWriter, request *http.Request) {
	Log("INFO", RequestLog{http.StatusBadRequest, request.Method,"Unsupported " + request.Method + " " + request.URL.Path})
	http.Error(writer, "Unsupported", http.StatusBadRequest)
}

func PasteHandler(writer http.ResponseWriter, request *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	switch method := request.Method; method {
		case http.MethodGet:
			Get(writer, request)

		case http.MethodPost:
			Post(writer, request)

		default:
			Others(writer, request)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", PasteHandler)
	fmt.Printf("%s %s\n", "Webpaste listening on port", *port)
	err := http.ListenAndServe(":" + *port, nil)

	if err != nil {
		println("ERROR: ", err.Error())
	}
}
