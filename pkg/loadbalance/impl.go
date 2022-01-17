package loadbalance

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"time"
)

// ---------- HTTPProxy ----------

type HTTPProxy struct {
	dial fasthttp.DialFunc
}

func NewHTTPProxy(username string, password string, address string, port int) *HTTPProxy {
	proxy := fmt.Sprintf("%s:%s@%s:%d", username, password, address, port)
	return &HTTPProxy{
		dial: fasthttpproxy.FasthttpHTTPDialerTimeout(proxy, time.Millisecond * 500),
	}
}

func (h HTTPProxy) Apply(client *fasthttp.Client) {
	client.Dial = h.dial
}

// ---------- NetworkAddress ----------

type NetworkAddress struct {

}
