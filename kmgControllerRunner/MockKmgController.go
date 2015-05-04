package kmgControllerRunner

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func CallApiByHttp(uri string, c *kmgHttp.Context) string {
	return CallApiByHttpWithUploadFile(uri, c, map[string]string{})
}

func CallApiByHttpWithUploadFile(uri string, c *kmgHttp.Context, uploadFileList map[string]string) string {
	server := httptest.NewServer(HttpHandler)
	defer server.Close()
	var response *http.Response
	var err error
	uri = server.URL + uri
	if c.Method == "GET" {
		response, err = http.Get(uri)
		handleErr(err)
	} else {
		buf := &bytes.Buffer{}
		formDataWriter := multipart.NewWriter(buf)
		defer formDataWriter.Close()
		for key, fullFilePath := range uploadFileList {
			formFilePart, err := formDataWriter.CreateFormFile(key, filepath.Base(fullFilePath))
			handleErr(err)
			file, err := os.Open(fullFilePath)
			defer file.Close()
			handleErr(err)
			_, err = io.Copy(formFilePart, file)
			handleErr(err)
		}
		for key, value := range c.Request {
			formDataWriter.WriteField(key, value)
		}
		contentType := formDataWriter.FormDataContentType()
		formDataWriter.Close()
		response, err = http.Post(uri, contentType, buf)
		handleErr(err)
	}
	respContent, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	handleErr(err)
	return string(respContent)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
