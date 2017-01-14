// Package update provides dyndns2 client functions.
package update

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Remote   string
	Host     string
	Username string
	Password string
	IP       string
	Offline  bool
}

const agent = "kurin/dnsupdate 1.0"

func (c *Config) query() string {
	var params []string
	params = append(params, fmt.Sprintf("hostname=%s", c.Host))
	if c.IP != "" {
		params = append(params, fmt.Sprintf("myip=%s", c.IP))
	}
	if c.Offline {
		params = append(params, "offline=YES")
	}
	return strings.Join(params, "&")
}

func (c *Config) url() (*url.URL, error) {
	u, err := url.Parse(c.Remote)
	if err != nil {
		return nil, err
	}
	u.Scheme = "https"
	u.RawQuery = c.query()
	return u, nil
}

func (c *Config) request() (*http.Request, error) {
	u, err := c.url()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("User-Agent", agent)
	return req, nil
}

func (c *Config) Update() error {
	req, err := c.request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("update failed: %s", resp.Status)
	}
	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	fmt.Println(buf.String())
	return nil
}
