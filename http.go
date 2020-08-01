package dbhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// sendRequest sends the query to the server.  It exists because http.PostForm() doesn't seem to have a way
// to change header values.
func sendRequest(queryUrl string, data url.Values, returnStructure interface{}) (err error) {
	var req *http.Request
	var resp *http.Response
	var client http.Client
	req, err = http.NewRequest(http.MethodPost, queryUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("go-dbhub v%s", LibraryVersion))
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

	// Unmarshall the JSON response into the structure provided by the caller
	err = json.NewDecoder(resp.Body).Decode(returnStructure)
	if err != nil {
		return
	}
	return
}
