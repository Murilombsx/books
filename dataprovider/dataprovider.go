package dataprovider

import (
	"books/constants"
	"books/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ HttpClient = &http.Client{}

type DataProvider struct {
	URL    string
	Client HttpClient
}

type BooksParams struct {
	Title string `schema:"title,required"`
	Limit uint8  `schema:"limit,required"`
}

func NewDataProvider(url string, client HttpClient) DataProvider {
	return DataProvider{
		URL:    url,
		Client: client,
	}
}

type fetchResult struct {
	data model.Items
	err  error
}

func (f DataProvider) Fetch(params *BooksParams) ([]model.BookInformation, error) {
	ch := make(chan fetchResult)
	go f.fetchURL(f.URL, params, ch)

	data, errors := []model.BookInformation{}, []string{}
	r := <-ch

	if r.err != nil {
		errors = append(errors, r.err.Error())
	} else {
		data = append(data, r.data.Items...)
	}

	if len(errors) > 0 {
		return []model.BookInformation{}, fmt.Errorf("%s", strings.Join(errors, ", "))
	}
	return data, nil
}

func (f DataProvider) fetchURL(url string, params *BooksParams, c chan<- fetchResult) {
	var resp *http.Response
	var err error
	attempts := 1

	url = applyParams(url, params)

	// Creating GET HTTP request based on the url provided
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Error("failed to create HTTP request")
		c <- fetchResult{err: err}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", constants.Key)

	// Retry requests if they fail, I'll use a simple linear backoff
	// Backoff will wait for the current number of attempts in seconds
	for {
		log.Infof("attempt %d/%d to fecht url %s", attempts, constants.MAX_RETRIES, url)
		resp, err = f.Client.Do(req)
		if err != nil {
			log.WithError(err).Error("http client call failed")
			time.Sleep(time.Duration(attempts) * time.Second)
		} else if err == nil && resp.StatusCode != http.StatusOK {
			log.Warnf("http client call succeeded, but with unexpected status code %d", resp.StatusCode)
			time.Sleep(time.Duration(attempts) * time.Second)
		} else {
			break
		}
		if attempts >= constants.MAX_RETRIES {
			err := fmt.Errorf("reached maximum number of retries for %s", path.Base(url))
			log.WithError(err).Errorf("exhausted retries for url %s", url)
			c <- fetchResult{err: err}
			return
		}
		attempts++
	}

	data := model.Items{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.WithError(err).Error("returned json doesn't comply with expected struct")
		c <- fetchResult{err: err}
		return
	}

	log.Infof("returning fetched data for %s", url)
	c <- fetchResult{data: data}
}

func applyParams(url string, params *BooksParams) string {
	search := "?q="
	title := params.Title + "+title"
	return url + search + title
}
