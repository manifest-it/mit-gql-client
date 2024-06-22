package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type Client struct {
	httpClient *http.Client
	headers    map[string]string
	endpoint   string
}

func NewClient(endpoint, token string, headers map[string]string) Client {
	return Client{endpoint: endpoint, httpClient: authutil.NewClientAuthzTokenSimple("", token), headers: headers}
}

func (c *Client) DoJSON(data []byte) (*http.Response, error) {
	if c.httpClient == nil {
		return nil, errors.New("no auth token")
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)

	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	return c.httpClient.Do(req)
}

func (c *Client) DoGraphQLString(gql string) (*http.Response, error) {
	req := QueryRequest{Query: gql}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.DoJSON(data)
}

func (c *Client) DoGraphQL(gql GraphQLOperation) (*http.Response, error) {
	return c.DoGraphQLString(gql.String())
}

type QueryRequest struct {
	Query string `json:"query"`
}
