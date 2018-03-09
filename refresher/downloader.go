package refresher

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// DownloadFile - download file by provided URL
func DownloadFile(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request object")
	}

	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "download file by URL")
	}

	return resp.Body, nil
}
