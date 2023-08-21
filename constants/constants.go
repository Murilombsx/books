package constants

import "time"

const (
	MAX_RETRIES     = 3
	SERVER_ADDR     = ":8080"
	DEFAULT_TIMEOUT = 15 * time.Second
)

var URL = "https://www.googleapis.com/books/v1/volumes"
