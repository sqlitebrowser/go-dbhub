package dbhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// sendRequestJSON sends a request to DBHub.io, formatting the returned result as JSON
func sendRequestJSON(queryUrl string, data url.Values, returnStructure interface{}) (err error) {
	// Send the request
	var body io.ReadCloser
	body, err = sendRequest(queryUrl, data)
	if body != nil {
		defer body.Close()
	}
	if err != nil {
		if body != nil {
			// If there's useful error info in the returned JSON, return that as the error message
			var z JSONError
			err = json.NewDecoder(body).Decode(&z)
			if err != nil {
				return
			}
			err = fmt.Errorf("%s", z.Msg)
		}
		return
	}

	// Unmarshall the JSON response into the structure provided by the caller
	err = json.NewDecoder(body).Decode(returnStructure)
	if err != nil {
		return
	}
	return
}

// sendRequest sends a request to DBHub.io.  It exists because http.PostForm() doesn't seem to have a way of changing
// header values.
func sendRequest(queryUrl string, data url.Values) (body io.ReadCloser, err error) {
	var req *http.Request
	var resp *http.Response
	var client http.Client
	req, err = http.NewRequest(http.MethodPost, queryUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("go-dbhub v%s", version))
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	// Return the response body, even if an error occurred.  This lets us return useful error information provided as
	// JSON in the body of the message
	body = resp.Body

	// Basic error handling, based on the status code received from the server
	if resp.StatusCode != 200 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		return
	}
	return
}

// sendUpload uploads a database to DBHub.io.  It exists because the DBHub.io upload end point requires multi-part data
func sendUpload(queryUrl string, data *url.Values, dbBytes *[]byte) (body io.ReadCloser, err error) {
	// Prepare the database file byte stream
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	dbName := data.Get("dbname")
	var wri io.Writer
	if dbName != "" {
		wri, err = w.CreateFormFile("file", dbName)
	} else {
		wri, err = w.CreateFormFile("file", "database.db")
	}
	if err != nil {
		return
	}
	_, err = wri.Write(*dbBytes)
	if err != nil {
		return
	}

	// Add the headers
	for i, j := range *data {
		wri, err = w.CreateFormField(i)
		if err != nil {
			return
		}
		_, err = wri.Write([]byte(j[0]))
		if err != nil {
			return
		}
	}
	err = w.Close()
	if err != nil {
		return
	}

	// Prepare the request
	var req *http.Request
	var resp *http.Response
	var client http.Client
	req, err = http.NewRequest(http.MethodPost, queryUrl, &buf)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("go-dbhub v%s", version))
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Upload the database
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	// Return the response body, even if an error occurred.  This lets us return useful error information provided as
	// JSON in the body of the message
	body = resp.Body

	// Basic error handling, based on the status code received from the server
	if resp.StatusCode != 201 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		return
	}
	return
}
