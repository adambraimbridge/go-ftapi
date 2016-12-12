package ftapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	Key  string
	Auth string
}

func (c *Client) FromURL(url string, obj interface{}) (*[]byte, error) {
	return c.do(url, nil, nil, obj)
}

func (c *Client) FromURLWithBody(url string, body []byte, obj interface{}) (*[]byte, error) {
	return c.do(url, body, nil, obj)
}

func (c *Client) FromPath(path string, obj interface{}) (*[]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) FromURLWithCookie(url string, obj interface{}, cookie *http.Cookie) (*[]byte, error) {
	return c.do(url, nil, cookie, obj)
}

func (c *Client) do(url string, body []byte, cookie *http.Cookie, obj interface{}) (*[]byte, error) {
	client := &http.Client{}

	var req *http.Request
	var err error

	log.Println(url)
	if body == nil {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Failed to build a GET request for ", url)
			return nil, err
		}
	} else {
		req, err = http.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		log.Printf("POST %s\nbody: %s\n", url, string(body))
		if err != nil {
			log.Println("Failed to build a POST request for ", url)
			log.Println(body)
			return nil, err
		}
	}

	req.Header.Add("X-API-Key", c.Key)

	if c.Auth != "" {
		req.Header.Add("Authorization", c.Auth)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute request for %s:%s\n", url, err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if err != nil {
		log.Printf("Failed to get %s:%s\n", url, err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", resp.Status, http.StatusText(resp.StatusCode))
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body from ", url)
	}

	if err := json.Unmarshal(rbody, obj); err != nil {
		log.Println("Failed to decode JSON from ", url)
		return nil, err
	}

	return &rbody, nil
}
