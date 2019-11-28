package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/tyk-sync/clients/objects"
	"github.com/ongoingio/urljoin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

func (c *Client) CreateCertificate(cert []byte) (string, error) {
	fullPath := urljoin.Join(c.url, endpointCerts)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("cert", "cert.pem")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, ioutil.NopCloser(bytes.NewReader(cert)))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fullPath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", c.secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API Returned error: %v", string(rBody))
	}

	dbResp := objects.CertResponse{}
	if err := json.Unmarshal(rBody, &dbResp); err != nil {
		return "", err
	}

	if strings.ToLower(dbResp.Status) != "ok" {
		return "", fmt.Errorf("API request completed, but with error: %v", dbResp.Message)
	}

	return dbResp.Id, nil
}
