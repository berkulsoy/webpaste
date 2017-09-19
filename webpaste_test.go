package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"
	"mime/multipart"
	"bytes"
	"io/ioutil"
	"io"
	"strings"
	"path"
	"compress/gzip"
)

func AssertHttpStatus(expected int, received int, t *testing.T) {
	if received != expected {
		t.Errorf("Expected %d got %d\n", expected, received)
		t.FailNow()
	}
}

func UploadFile(file string, rr *httptest.ResponseRecorder) {
	infile, _ := os.Open(file)
	defer infile.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	mpart, err := writer.CreateFormFile("f", infile.Name())
	if err != nil {
		println("AAAA" + err.Error())
	}
	io.Copy(mpart, infile)
	request, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	handler := http.HandlerFunc(PasteHandler)
	writer.Close()
	handler.ServeHTTP(rr, request)
}

func DownloadFile(file string, rr *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/" + file, nil)
	handler := http.HandlerFunc(PasteHandler)
	handler.ServeHTTP(rr, request)
}

func CompareFiles(expected string, created string, t *testing.T) {
	createdfile, err := ioutil.ReadFile(created)
	expectedfile, _ := ioutil.ReadFile(expected)

	if err != nil {
		t.Errorf(t.Name() + ": " + err.Error())
		t.FailNow()
	}

	if bytes.Compare(createdfile, expectedfile) != 0 {
		t.Errorf(t.Name() + ": Uploaded file is not equal")
		t.FailNow()
	}

}

func GetFileFromResponse(rr *httptest.ResponseRecorder) string {
	return path.Base(strings.Split(rr.Body.String(),"\n")[1])
}

func TestNonGetPostRequests(t *testing.T) {
	unsupportedMethods := [7]string{http.MethodOptions, http.MethodDelete, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodPut, http.MethodTrace}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PasteHandler)
	for _, method := range unsupportedMethods {
		request, _ := http.NewRequest(method, "test", nil)
		handler.ServeHTTP(rr, request)
		AssertHttpStatus(http.StatusBadRequest, rr.Code, t)
	}
}

func TestPostBinary(t *testing.T) {
	infile, _ := ioutil.ReadFile("test/sample.txt")
	outfile, _ := ioutil.TempFile("/tmp/", "webpaste")
	gzwriter := gzip.NewWriter(outfile)
	gzwriter.Write(infile)
	gzwriter.Flush()
	gzwriter.Close()
	outfile.Close()

	rr := httptest.NewRecorder()

	UploadFile(outfile.Name(), rr)
	AssertHttpStatus(http.StatusCreated, rr.Code, t)
	CompareFiles(outfile.Name(), "/tmp/" + GetFileFromResponse(rr), t)
	os.Remove(outfile.Name())
}

func TestPostPlainText(t *testing.T) {
	rr := httptest.NewRecorder()
	UploadFile("test/sample.txt", rr)
	AssertHttpStatus(http.StatusCreated, rr.Code, t)
	CompareFiles("test/sample.txt", "/tmp/" + GetFileFromResponse(rr), t)
}

func TestGetFile(t *testing.T) {
	rrUpload := httptest.NewRecorder()
	UploadFile("test/sample.txt", rrUpload)
	AssertHttpStatus(http.StatusCreated, rrUpload.Code, t)
	rrDownload := httptest.NewRecorder()
	file := GetFileFromResponse(rrUpload)
	DownloadFile(file, rrDownload)
	AssertHttpStatus(http.StatusOK, rrDownload.Code, t)
}



