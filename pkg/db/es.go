package db

import (
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

type EsOptions struct {
	Host string
	Port string
}

func NewEsClient(opts *EsOptions) (*elastic.Client, error) {
	esClient, err := elastic.NewClient(
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetSniff(false),
		elastic.SetURL(fmt.Sprintf("http://%s:%s/", opts.Host, opts.Port)))
	if err != nil {
		return nil, err
	}
	return esClient, nil
}
