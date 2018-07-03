package woocommerce

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	UserAgent     = "WooCommerce API Client-PHP"
	HashAlgorithm = "HMAC-SHA256"
)

type Client struct {
	storeURL  *url.URL
	ck        string
	cs        string
	option    *Options
	rawClient *http.Client
}

func NewClient(store, ck, cs string, option *Options) (*Client, error) {
	storeURL, err := url.Parse(store)
	if err != nil {
		return nil, err
	}

	if option == nil {
		option = &Options{}
	}
	if option.OauthTimestamp.IsZero() {
		option.OauthTimestamp = time.Now()
	}

	if option.Version == "" {
		option.Version = "v2"
	}
	path := "/wp-json/wc/"
	if option.API {
		path = option.APIPrefix
	}
	path = path + option.Version + "/"
	storeURL.Path += path

	rawClient := &http.Client{}
	rawClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: option.VerifySSL},
	}
	return &Client{
		storeURL:  storeURL,
		ck:        ck,
		cs:        cs,
		option:    option,
		rawClient: rawClient,
	}, nil
}

func (c *Client) basicAuth(params url.Values) string {
	params.Add("consumer_key", c.ck)
	params.Add("consumer_secret", c.cs)
	return params.Encode()
}
func (c *Client) oauth(method, urlStr string, params url.Values) string {
	signature, signatureParams := Oauth_generator(c.ck, c.cs, strings.ToUpper(method), urlStr, params)
	signatureParams.Add("oauth_signature", signature)

	//Arrangement of url parameter(specUrl) and Oauth(oauth_generator) parameter will affect result of Signature
	specURL := "oauth_consumer_key=" + signatureParams.Get("oauth_consumer_key") + 
	"&oauth_nonce=" + signatureParams.Get("oauth_nonce") + 
	"&oauth_signature=" + url.QueryEscape(signature) + 
	"&oauth_signature_method=" + signatureParams.Get("oauth_signature_method") + 
	"&oauth_timestamp=" + signatureParams.Get("oauth_timestamp");
	specURL = addURLParams(signatureParams, specURL)
	
	return specURL
}

func (c *Client) request(method, endpoint string, params url.Values, data interface{}) (io.ReadCloser, error) {
	urlstr := c.storeURL.String() + endpoint
	if params == nil {
		params = make(url.Values)
	}
	if c.storeURL.Scheme == "https" {
		urlstr += "?" + c.basicAuth(params)
	} else {
		urlstr += "?" + c.oauth(method, urlstr, params)
	}
	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, fmt.Errorf("Method is not recognised: %s", method)
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("Request failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) Post(endpoint string, data interface{}) (io.ReadCloser, error) {
	return c.request("POST", endpoint, nil, data)
}

func (c *Client) Put(endpoint string, data interface{}) (io.ReadCloser, error) {
	return c.request("PUT", endpoint, nil, data)
}

func (c *Client) Get(endpoint string, params url.Values) (io.ReadCloser, error) {
	return c.request("GET", endpoint, params, nil)
}

func (c *Client) Delete(endpoint string, params url.Values) (io.ReadCloser, error) {
	return c.request("POST", endpoint, params, nil)
}

func (c *Client) Options(endpoint string) (io.ReadCloser, error) {
	return c.request("OPTIONS", endpoint, nil, nil)
}