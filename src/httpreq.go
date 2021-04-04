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
		return err
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
		httpReqResult := buildHttpResult(0, req.URL.String(), resp.StatusCode, elapsed.Nanoseconds())

		httpResultsChannel <- httpReqResult

		return err
	} else {
		defer resp.Body.Close()

		httpReqResult := buildHttpResult(len(responseBody), req.URL.String(), resp.StatusCode, elapsed.Nanoseconds())

		processResult(httpAction, resp, variables, responseBody, elapsed.Nanoseconds())

		httpResultsChannel <- httpReqResult
	}

	return nil
}

func buildHttpResult(contentLength int, url string, status int, elapsed int64) HttpReqResult {
	httpReqResult := HttpReqResult{
		elapsed,
		contentLength,
		url,
		status,
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
			for _, pair := range httpAction.FormData {
				formKey, err2 := SubstParams(variables, pair.Key)

				if err2 != nil {
					return nil, err2
				}

				formValue, err2 := SubstParams(variables, pair.Value)

				if err2 != nil {
					return nil, err2
				}

				form.Add(formKey, formValue)
			}
		}

		reader := strings.NewReader(form.Encode())
		req, err = http.NewRequest(httpAction.Method, endpoint, reader)
	} else if httpAction.BodyType == "multipart" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		if httpAction.FormData != nil {
			// Add form data
			for _, pair := range httpAction.FormData {
				formKey, err2 := SubstParams(variables, pair.Key)
				if err2 != nil {
					return nil, err2
				}

				formValue, err2 := SubstParams(variables, pair.Value)
				if err2 != nil {
					return nil, err2
				}

				err3 := writer.WriteField(formKey, formValue)
				if err3 != nil {
					return nil, err3
				}
			}
		}

		if httpAction.Files != nil && len(httpAction.Files) > 0 {
			for _, pair := range httpAction.Files {
				fileKey, err2 := SubstParams(variables, pair.Key)
				if err2 != nil {
					return nil, err2
				}

				filePath, err2 := SubstParams(variables, pair.Value)
				if err2 != nil {
					return nil, err2
				}

				file, err2 := os.Open(filePath)
				if err2 != nil {
					return nil, err2
				}

				defer file.Close()

				part, err2 := writer.CreateFormFile(fileKey, filepath.Base(filePath))
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
		for _, pair := range httpAction.Headers {
			headerKey, err2 := SubstParams(variables, pair.Key)
			if err2 != nil {
				return nil, err2
			}

			headerValue, err2 := SubstParams(variables, pair.Value)
			if err2 != nil {
				return nil, err2
			}

			req.Header.Add(headerKey, headerValue)
		}
	}

	if httpAction.Cookies != nil {
		// Add cookies
		for _, pair := range httpAction.Cookies {
			cookieKey, err2 := SubstParams(variables, pair.Key)
			if err2 != nil {
				return nil, err2
			}

			cookieValue, err2 := SubstParams(variables, pair.Value)
			if err2 != nil {
				return nil, err2
			}

			cookie := http.Cookie{
				Name:  cookieKey,
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

	if len(httpAction.VariableName) > 0 {
		variables[httpAction.VariableName] = HttpResponseObject{
			response.StatusCode,
			body,
			headers,
			int(elapsed / 1000000),
		}
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
