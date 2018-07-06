package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "encoding/json"
import "bufio"
import "bytes"

type CrtSh struct{}

type crtshObject struct {
	NameValue string `json:"name_value"`
}

