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
	endpoint   string
}

func NewClient(endpoint, token string) Client {
	return Client{endpoint: endpoint, httpClient: authutil.NewClientAuthzTokenSimple("", token)}
}

func (c *Client) DoJSON(data []byte, headers map[string]string) (*http.Response, error) {
	if c.httpClient == nil {
		return nil, errors.New("no auth token")
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return c.httpClient.Do(req)
}

func (c *Client) DoGraphQLString(gql string, headers map[string]string) (*http.Response, error) {
	req := QueryRequest{Query: gql}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.DoJSON(data, headers)
}

func (c *Client) DoGraphQL(gql GraphQLOperation, headers map[string]string) (*http.Response, error) {
	return c.DoGraphQLString(gql.String(), headers)
}

type QueryRequest struct {
	Query string `json:"query"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
}
