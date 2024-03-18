package go_resty

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"io"
	"net/http"
	"testing"
	"time"
)

func Test_retry(t *testing.T) {
	url := "https://github.com/avast/retry-go"
	var body []byte

	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(string(body))
}
func Test_retryWithData(t *testing.T) {
	url := "https://github.com/avast/retry-go"
	var body []byte

	body, err := retry.DoWithData(
		func() ([]byte, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return body, nil
		},
	)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(body))
}
func Test_retryCount(t *testing.T) {
	err := retry.Do(
		func() error {
			return errors.New("Err")
		},
		retry.Delay(1*time.Second),
		retry.Attempts(10),
		retry.DelayType(retry.FixedDelay),
	)
	if err != nil {
		t.Error(err)
	}
}
