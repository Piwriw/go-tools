package go_ip

import "testing"

func TestGetInternetIP(t *testing.T) {
	t.Log("PublicIP", GetInternetIP())
}

func TestGetLocalIP(t *testing.T) {
	t.Log("LocalIP", GetLocalIP())
}
