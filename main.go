package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const validatorUrl = "https://validator.w3.org/nu/"

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		panic("Wrong number of arguments")
	}
	filename := flag.Arg(0)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"fragment": bufio.NewReader(f),
	}
	resp, err := upload(validatorUrl, values)
	if err != nil {
		panic(err)
	}
	results, err := parse(resp)
	if err != nil {
		panic(err)
	}
	po := printOptions{descriptionsIgnored: []string{
			"Trailing slash on void elements has no effect and interacts badly with unquoted attribute values.",
		},
	}
	printResults(po, filename, results)
}

func upload(url string, values map[string]io.Reader) (io.ReadCloser, error) {
	client := http.DefaultClient
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		fw, err := w.CreateFormField(key)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return res.Body, err
}
