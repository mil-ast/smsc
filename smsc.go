package smsc

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	TIMEOUT_DURATION time.Duration = time.Second * 20
	CHARSET          string        = "utf-8"
)

type SMSC struct {
	url    string
	client *http.Client
	params *url.Values
}

func New(addr, login, password string) (*SMSC, error) {
	if addr == "" {
		return nil, errors.New("error:url")
	}

	if login == "" {
		return nil, errors.New("error:login")
	}

	url_rapse, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	sms := new(SMSC)
	sms.url = url_rapse.String()
	sms.params = &url.Values{}

	sms.AddParam("login", login)
	sms.AddParam("psw", password)
	sms.AddParam("charset", "utf-8")

	sms.client = &http.Client{Timeout: TIMEOUT_DURATION}

	if url_rapse.Scheme == "https" {
		sms.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
	}

	return sms, nil
}

/*
	управление параметрами
*/
func (this *SMSC) AddParam(key, value string) {
	this.params.Add(key, value)
}
func (this *SMSC) SetParam(key, value string) {
	this.params.Set(key, value)
}
func (this *SMSC) DelParam(key string) {
	this.params.Del(key)
}

/*
	POST
*/
func (this *SMSC) post(addr string, params *url.Values) (string, error) {
	var body string = params.Encode()

	req, err := http.NewRequest("POST", addr, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	resp, err := this.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

/*
	Отправка сообщения
*/
func (this *SMSC) Send(method string, phones []string, message string) (string, error) {
	if method == "GET" || method == "get" {
		return this.send_sms_get(phones, message)
	}

	return this.send_sms_post(phones, message)
}

/*
	отправка методом POST
*/
func (this *SMSC) send_sms_post(phones []string, message string) (string, error) {
	var params *url.Values = this.params

	params.Add("phones", strings.Join(phones, ","))
	params.Add("mes", message)

	return this.post(this.url, params)
}

/*
	отправка методом GET
*/
func (this *SMSC) send_sms_get(phones []string, message string) (string, error) {
	var params *url.Values = this.params
	params.Add("phones", strings.Join(phones, ","))
	params.Add("mes", message)

	u, err := url.Parse(this.url)
	if err != nil {
		return "", err
	}

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

/*
	Выполнить запрос
*/
func (this *SMSC) Request() (string, error) {
	return this.post(this.url, this.params)
}
