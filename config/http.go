package config

import "time"

var (
	GetHTTPServerAddress = getHTTPServerAddress
	GetHTTPReadTimeout   = getHTTPReadTimeout
	GetHTTPWriteTimeout  = getHTTPWriteTimeout
)

func getHTTPServerAddress() string {
	return getConfigString("http.server_address")
}

func getHTTPReadTimeout() time.Duration {
	return getConfigDuration("http.read_timeout")
}

func getHTTPWriteTimeout() time.Duration {
	return getConfigDuration("http.write_timeout")
}
