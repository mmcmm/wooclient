package woocommerce

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func oauth_generator(_KEY, _SECRET, HttpMethod, RequestUrl string, params url.Values) (string, *url.Values) {
	HASH := "HMAC-SHA256"

	SignatureParameters := url.Values{}
	SignatureParameters.Add("oauth_consumer_key", _KEY)
	SignatureParameters.Add("oauth_nonce", getSha1Nonce())
	SignatureParameters.Add("oauth_signature_method", HASH)
	SignatureParameters.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	SignatureParameters = addParams(params, SignatureParameters)

	SignatureBaseString := strings.Join([]string{strings.ToUpper(HttpMethod), url.QueryEscape(RequestUrl),
		url.QueryEscape(SignatureParameters.Encode())}, "&")

	return get_HMAC_SHA256(_SECRET, SignatureBaseString), &SignatureParameters
}

func getSha1Nonce() string {
	nonce := make([]byte, 16)
	rand.Read(nonce)
	sha1Nonce := fmt.Sprintf("%x", sha1.Sum(nonce))
	return sha1Nonce
}

func get_HMAC_SHA256(consumer_secret, SignatureBaseString string) string {
	mac := hmac.New(sha256.New, []byte(consumer_secret+"&"))
	mac.Write([]byte(SignatureBaseString))
	SignatureBytes := mac.Sum(nil)
	Signature := base64.StdEncoding.EncodeToString(SignatureBytes)
	return Signature
}

func Oauth_generator(_KEY, _SECRET, HttpMethod, RequestUrl string, params url.Values) (string, *url.Values) {
	return oauth_generator(_KEY, _SECRET, HttpMethod, RequestUrl, params)
}
