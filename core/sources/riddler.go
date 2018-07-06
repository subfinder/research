package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "strings"
import "encoding/json"
import "bytes"
import "bufio"
import "net"
import "time"
import "errors"

type Riddler struct {
	Email    string
	Password string
	APIToken string
}

type riddlerHost struct {
	Host string `json:"host"`
}

type riddlerAuthenticationResponse struct {
	Response struct {
		User struct {
			AuthenticationToken string `json:"authentication_token"`
		} `json:"user"`
	} `json:"response"`
}

