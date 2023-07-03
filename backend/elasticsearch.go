package backend

import (
	"context"
    "fmt"
    "appstore/constants"
    "github.com/olivere/elastic/v7"
)

var (
	ESBackend *ElasticsearchBackend 
) 

type ElasticsearchBackend struct {
	client *elastic.Client
}

func InitElasticsearchBackend() {
	
	// New client
	client, err := elastic.NewClient(
        elastic.SetURL(constants.ES_URL),
        elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD),
        elastic.SetSniff(false))
    if err != nil {
		fmt.Println("Fail to create new ES client.")
        panic(err)
    }


	// Check if app index exists
	exists, err := client.IndexExists(constants.APP_INDEX).Do(context.Background())
    if err != nil {
		fmt.Println("Fail to check if app index exists.")
        panic(err)
    }

	// Create app index if it does not exist
	if !exists {
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "title":      { "type": "text"},
                    "description":  { "type": "text" },
                    "price":      { "type": "keyword", "index": false },
                    "url":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := client.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
			fmt.Println("Fail to create app index.")
            panic(err)
        }
    }

	// Check if user index exists
	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
		fmt.Println("Fail to check if user index exists.")
        panic(err)
    }

	// Create user index if it does not exist
    if !exists {
        mapping := `{
                     "mappings": {
                         "properties": {
                            "username": {"type": "keyword"},
                            "password": {"type": "keyword"},
                            "age": {"type": "long", "index": false},
                            "gender": {"type": "keyword", "index": false}
                         }
                    }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
			fmt.Println("Fail to create user index.")
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

	// ES backend pointer 
    ESBackend = &ElasticsearchBackend{client: client}
}

// Method for type ElasticsearchBackend
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }
    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}

func (backend *ElasticsearchBackend) DeleteFromES(id, index string) error {
    _, err := backend.client.Delete().
    Index(index).
    Pretty(true).
    Id(id).
    Do(context.Background())
    return err
}


