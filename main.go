package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

//go:embed index.html
var htmlFile string

const validatorUrl = "https://validator.w3.org/nu/"

func main() {

    remoteURL := validatorUrl
    client := http.DefaultClient

    //prepare the reader instances to encode
    values := map[string]io.Reader{
        "fragment": strings.NewReader(htmlFile),
    }
    err := Upload(client, remoteURL, values)
    if err != nil {
        panic(err)
    }
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    for key, r := range values {
        var fw io.Writer
        if x, ok := r.(io.Closer); ok {
            defer x.Close()
        }
		if fw, err = w.CreateFormField(key); err != nil {
			return
		}
        if _, err = io.Copy(fw, r); err != nil {
            return err
        }

    }
    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
        return
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Submit the request
    res, err := client.Do(req)
    if err != nil {
        return
    }

    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    }
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
    return
}
