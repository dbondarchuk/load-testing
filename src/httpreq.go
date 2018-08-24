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
func DoHttpRequest(httpAction HttpAction, httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	req, err := buildHttpRequest(httpAction, variables)
	if err != nil {
		return nil
	}

	start := time.Now()
	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	resp, err := DefaultTransport.RoundTrip(req)

	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal(err)
		log.Printf("Reading HTTP response failed: %s\n", err)
		httpReqResult := buildHttpResult(0, req.URL.String(), resp.StatusCode, elapsed.Nanoseconds(), httpAction.Name)

		httpResultsChannel <- httpReqResult

		return err
	} else {
		defer resp.Body.Close()

		if httpAction.StoreCookie != "" {
			for _, cookie := range resp.Cookies() {

				if cookie.Name == httpAction.StoreCookie {
					variables["____"+cookie.Name] = cookie.Value
				}
			}
		}

		httpReqResult := buildHttpResult(len(responseBody), req.URL.String(), resp.StatusCode, elapsed.Nanoseconds(), httpAction.Name)

		processResult(httpAction, resp, variables, responseBody, elapsed.Nanoseconds())

		httpResultsChannel <- httpReqResult
	}

	return nil
}

func buildHttpResult(contentLength int, url string, status int, elapsed int64, name string) HttpReqResult {
	httpReqResult := HttpReqResult{
		elapsed,
		contentLength,
		url,
		status,
		name,
		time.Since(SimulationStart).Nanoseconds(),
	}
	return httpReqResult
}

func buildHttpRequest(httpAction HttpAction, variables map[string]interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	endpoint, err := SubstParams(variables, httpAction.Endpoint)
	if err != nil {
		return nil, err
	}

	if httpAction.BodyType == "raw" && httpAction.RawData != "" {
		rawData, err2 := SubstParams(variables, httpAction.RawData)
		if err2 != nil {
			return nil, err
		}

		reader := strings.NewReader(rawData)
		req, err = http.NewRequest(httpAction.Method, endpoint, reader)
	} else if httpAction.BodyType == "dataform" {
		form := url.Values{}

		if httpAction.FormData != nil {
			// Add form data
			for key, value := range httpAction.FormData {
				formValue, err2 := SubstParams(variables, value.(string))

				if err2 != nil {
					return nil, err
				}

				form.Add(key, formValue)
			}
		}

		reader := strings.NewReader(form.Encode())
		req, err = http.NewRequest(httpAction.Method, endpoint, reader)
	} else if httpAction.BodyType == "dataform" || httpAction.BodyType == "multipart" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		if httpAction.FormData != nil {
			// Add form data
			for key, value := range httpAction.FormData {
				formValue, err2 := SubstParams(variables, value.(string))
				if err2 != nil {
					return nil, err2
				}

				err3 := writer.WriteField(key, formValue)
				if err3 != nil {
					return nil, err3
				}
			}
		}

		if httpAction.Files != nil && len(httpAction.Files) > 0 {
			for key, path := range httpAction.Files {
				file, err2 := os.Open(path.(string))
				if err2 != nil {
					return nil, err2
				}

				defer file.Close()

				part, err2 := writer.CreateFormFile(key, filepath.Base(path.(string)))
				if err2 != nil {
					return nil, err2
				}

				_, err2 = io.Copy(part, file)
				err2 = writer.Close()

				if err2 != nil {
					return nil, err2
				}
			}
		}

		req, err = http.NewRequest(httpAction.Method, endpoint, body)
	} else {
		req, err = http.NewRequest(httpAction.Method, endpoint, nil)
	}

	if err != nil {
		return nil, err
	}

	if httpAction.Headers != nil {
		// Add headers
		for key, value := range httpAction.Headers {
			headerValue, err2 := SubstParams(variables, value.(string))
			if err2 != nil {
				return nil, err2
			}

			req.Header.Add(key, headerValue)
		}
	}

	if httpAction.Cookies != nil {
		// Add cookies
		for key, value := range httpAction.Cookies {
			k := key
			if strings.HasPrefix(key, "____") {
				k = key[4:len(key)]
			}

			cookieValue, err2 := SubstParams(variables, value.(string))
			if err2 != nil {
				return nil, err2
			}

			cookie := http.Cookie{
				Name:  k,
				Value: cookieValue,
			}

			req.AddCookie(&cookie)
		}
	}

	return req, err
}

func processResult(httpAction HttpAction, response *http.Response, variables map[string]interface{}, responseBody []byte, elapsed int64) {
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

	if err != nil {
		err = xml.Unmarshal(responseBody, body)

		if err != nil {
			body = string(responseBody[:])
		}
	}

	headers := make(map[string]string)
	for name, value := range response.Header {
		headers[name] = value[0]
	}

	variables[httpAction.VariableName] = HttpResponseObject{
		response.StatusCode,
		body,
		headers,
		int(elapsed / 1000000),
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
