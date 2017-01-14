// Package update provides dyndns2 client functions.
package update

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

func (c *Config) auth() string {
	b := []byte(fmt.Sprintf("%s:%s", c.Username, c.Password))
	return "Basic " + base64.StdEncoding.EncodeToString(b)
}

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
	return &http.Request{
		URL:   u,
		Close: true,
		Header: http.Header{
			"Authorization": []string{c.auth()},
			"User-Agent":    []string{agent},
		},
	}, nil
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
	io.Copy(os.Stdout, resp.Body)
	return nil
}
