package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	//"fmt"

	"bytes"
	"crypto/tls"
)

// DoHttpRequest - Accepts a Httpaction and a one-way channel to write the results to.
func DoHttpRequest(httpAction HttpAction, resultsChannel chan HttpReqResult, variables map[string]interface{}) {
	req := buildHttpRequest(httpAction, variables)

	start := time.Now()
	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	resp, err := DefaultTransport.RoundTrip(req)

	if err != nil {
		log.Printf("HTTP request failed: %s", err)
	} else {
		elapsed := time.Since(start)
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Fatal(err)
			log.Printf("Reading HTTP response failed: %s\n", err)
			httpReqResult := buildHttpResult(0, resp.StatusCode, elapsed.Nanoseconds(), httpAction.Name)

			resultsChannel <- httpReqResult
		} else {
			defer resp.Body.Close()

			if httpAction.StoreCookie != "" {
				for _, cookie := range resp.Cookies() {

					if cookie.Name == httpAction.StoreCookie {
						variables["____"+cookie.Name] = cookie.Value
					}
				}
			}

			httpReqResult := buildHttpResult(len(responseBody), resp.StatusCode, elapsed.Nanoseconds(), httpAction.Name)

			processResult(httpAction, resp, variables, responseBody)

			resultsChannel <- httpReqResult
		}
	}
}

func buildHttpResult(contentLength int, status int, elapsed int64, name string) HttpReqResult {
	httpReqResult := HttpReqResult{
		"HTTP",
		elapsed,
		contentLength,
		status,
		name,
		time.Since(SimulationStart).Nanoseconds(),
	}
	return httpReqResult
}

func buildHttpRequest(httpAction HttpAction, variables map[string]interface{}) *http.Request {
	var req *http.Request
	var err error

	if httpAction.BodyType == "raw" && httpAction.RawData != "" {
		reader := strings.NewReader(SubstParams(variables, httpAction.RawData))
		req, err = http.NewRequest(httpAction.Method, SubstParams(variables, httpAction.Endpoint), reader)
	} else if httpAction.BodyType == "dataform" {
		form := url.Values{}

		// Add form data
		for key, value := range httpAction.FormData {
			form.Add(key, SubstParams(variables, value))
		}

		reader := strings.NewReader(form.Encode())
		req, err = http.NewRequest(httpAction.Method, SubstParams(variables, httpAction.Endpoint), reader)
	} else if httpAction.BodyType == "dataform" || httpAction.BodyType == "multipart" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add form data
		for key, value := range httpAction.FormData {
			err2 := writer.WriteField(key, SubstParams(variables, value))
			if err2 != nil {
				log.Fatal(err2)
			}
		}

		if len(httpAction.Files) > 0 {
			for key, path := range httpAction.Files {
				file, err2 := os.Open(path)
				if err2 != nil {
					log.Fatal(err2)
				}

				defer file.Close()

				part, err2 := writer.CreateFormFile(key, filepath.Base(path))
				if err2 != nil {
					log.Fatal(err2)
				}

				_, err2 = io.Copy(part, file)
				err2 = writer.Close()

				if err2 != nil {
					log.Fatal(err2)
				}
			}
		}

		req, err = http.NewRequest(httpAction.Method, SubstParams(variables, httpAction.Endpoint), body)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Add headers
	for key, value := range httpAction.Headers {
		req.Header.Add(key, SubstParams(variables, value))
	}

	// Add cookies
	for key, value := range httpAction.Cookies {
		k := key
		if strings.HasPrefix(key, "____") {
			k = key[4:len(key)]
		}

		cookie := http.Cookie{
			Name:  k,
			Value: SubstParams(variables, value),
		}

		req.AddCookie(&cookie)
	}

	return req
}

func processResult(httpAction HttpAction, response *http.Response, variables map[string]interface{}, responseBody []byte) {
	var body interface{} = nil
	var err error

	if len(responseBody) == 0 {
		body = ""
	}

	// if json
	obj := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &obj)
	if err != nil {
		arr := make([]interface{}, 0)
		err = json.Unmarshal(responseBody, &arr)
		if err != nil {
			var b bool
			err = json.Unmarshal(responseBody, &b)
			if err != nil {
				var f float64
				err = json.Unmarshal(responseBody, &f)
				if err != nil {
					var s string
					err = json.Unmarshal(responseBody, &s)
					if err != nil {
						if len(responseBody) > 0 {
							err = errors.New("Not a json")
						} else {
							body = ""
						}
					} else {
						body = s
					}
				} else {
					body = f
				}
			} else {
				body = b
			}
		} else {
			body = arr
		}
	} else {
		body = obj
	}

	err = xml.Unmarshal(responseBody, body)

	if err != nil {
		body = string(responseBody[:])
	}

	headers := make(map[string]string)
	for name, value := range response.Header {
		headers[name] = value[0]
	}

	variables[httpAction.VariableName] = HttpResponseObject{
		response.StatusCode,
		body,
		headers,
	}
}

/**
 * Trims leading and trailing byte r from string s
 */
func trimChar(s string, r byte) string {
	sz := len(s)

	if sz > 0 && s[sz-1] == r {
		s = s[:sz-1]
	}
	sz = len(s)
	if sz > 0 && s[0] == r {
		s = s[1:sz]
	}
	return s
}
