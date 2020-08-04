package dbhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// sendRequestJSON sends a request to DBHub.io, formatting the returned result as JSON
func sendRequestJSON(queryUrl string, data url.Values, returnStructure interface{}) (err error) {
	// Send the request
	var body io.ReadCloser
	body, err = sendRequest(queryUrl, data)
	if err != nil {
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

	// Basic error handling, based on the status code received from the server
	if resp.StatusCode != 200 {
		// The returned status code indicates something went wrong
		err = fmt.Errorf(resp.Status)
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		return
	}

	// Return the response body
	body = resp.Body
	return
}
