package service

import (
	"context"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Fetcher interface {
	Fetch(url string) (*goquery.Document, error)
	IsAccessible(ctx context.Context, pr *model.WorkerWrapper) (*model.WorkerWrapper, error)
}

type FetcherService struct {
	client *http.Client

	sync.WaitGroup
}

func NewFetcherService(client *http.Client) *FetcherService {
	return &FetcherService{
		client: client,
	}
}

func (p *FetcherService) Fetch(url string) (*goquery.Document, error) {
	// Load the HTML document
	return goquery.NewDocument(url)
}

func (p *FetcherService) IsAccessible(ctx context.Context, pr *model.WorkerWrapper) (*model.WorkerWrapper, error) {
	req, err := http.NewRequest(http.MethodGet, pr.Url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, errors.Wrap(err, "creating preprocess request failed")
	}
	const requestTimeout = 10 * time.Second
	ctx1, _ := context.WithTimeout(ctx, requestTimeout)

	response, err := p.client.Do(req.WithContext(ctx1))
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "creating preprocess api failed")
	}
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading preprocess body failed")
	}
	if response.StatusCode >= 200 || response.StatusCode < 400 {
		pr.Result = true
	}
	return pr, nil
}
