package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	targetUrl := "https://10.10.102.96:6443/api/v1/namespaces/default/services"

	req, _ := http.NewRequest("GET", targetUrl, nil)

	req.Header.Add("Authorization", "Bearer xxx")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	sprintf := fmt.Sprintf("%s", s)
	fmt.Println(sprintf)
}
