package resty

import (
	"github.com/go-resty/resty/v2"
	"github.com/lexcao/genapi"
)

var DefaultClient = New(resty.New())

func New(client *resty.Client) *HttpClient {
	return &HttpClient{client: client}
}

type HttpClient struct {
	client *resty.Client
}

func (c *HttpClient) SetConfig(config genapi.Config) {
	c.client.SetBaseURL(config.BaseURL)
	c.client.Header = config.Header
}

func (c *HttpClient) Do(req *genapi.Request) (*genapi.Response, error) {
	restyReq := c.client.NewRequest().
		SetDoNotParseResponse(true).
		SetPathParams(req.PathParams).
		SetContext(req.Context).
		SetBody(req.Body)

	for key, value := range req.Header {
		for _, v := range value {
			restyReq.Header.Add(key, v)
		}
	}

	restyReq.QueryParam = req.Queries
	restyReq.URL = req.Path
	restyReq.Method = req.Method

	resp, err := restyReq.Send()
	return resp.RawResponse, err
}
