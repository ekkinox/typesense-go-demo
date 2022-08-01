package client

import (
	"fmt"
	"strconv"

	"github.com/jaswdr/faker"
	"github.com/spf13/viper"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

const (
	Collection = "companies"
	BatchSize  = 1000
)

type Company struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Address      string `json:"address"`
	NumEmployees int    `json:"num_employees"`
	Country      string `json:"country"`
}

type TypesenseClient struct {
	Client *typesense.Client
}

func InitTypesenseClient() TypesenseClient {

	host := viper.GetString("TYPESENSE_HOST")
	port := viper.GetInt("TYPESENSE_PORT")
	key := viper.GetString("TYPESENSE_API_KEY")

	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("%s:%d", host, port)),
		typesense.WithAPIKey(key),
	)

	return TypesenseClient{client}
}

func (c TypesenseClient) Search(expression string, page int, size int, countries string, min int, max int) (*api.SearchResult, error) {
	searchParameters := &api.SearchCollectionParams{
		Q:                 expression,
		QueryBy:           "name,description,address",
		SortBy:            pointer.String("num_employees:desc"),
		Page:              pointer.Int(page),
		PerPage:           pointer.Int(size),
		HighlightStartTag: pointer.String("<span class=\"text-dark bg-warning\">"),
		HighlightEndTag:   pointer.String("</span>"),
		FacetBy:           pointer.String("country"),
		MaxFacetValues:    pointer.Int(30),
		FilterBy:          pointer.String(fmt.Sprintf("num_employees:>=%d&&num_employees:<=%d", min, max)),
	}

	if countries != "" {
		searchParameters.FilterBy = pointer.String(fmt.Sprintf("country:%s", countries))
	}

	return c.Client.Collection(Collection).Documents().Search(searchParameters)
}

func (c TypesenseClient) WipeAndCreateSchema() error {

	// wipe collection if exists
	_, err := c.Client.Collection(Collection).Retrieve()
	if err == nil {
		_, err = c.Client.Collection(Collection).Delete()
		if err != nil {
			return err
		}
	}

	// create collection
	schema := &api.CollectionSchema{
		Name: Collection,
		Fields: []api.Field{
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "description",
				Type: "string",
			},
			{
				Name: "address",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: pointer.True(),
			},
		},
		DefaultSortingField: pointer.String("num_employees"),
	}

	_, err = c.Client.Collections().Create(schema)
	if err != nil {
		return err
	}

	return nil
}

func (c TypesenseClient) PopulateSchema(num int) error {

	// faker
	gen := faker.New()

	// get last doc position
	collection, err := c.Client.Collection(Collection).Retrieve()
	if err != nil {
		return err
	}
	start := int(collection.NumDocuments)

	// insert new docs
	var documents []interface{}

	pos := 1
	for pos <= num {
		documents = append(documents, &Company{
			Id:           strconv.Itoa(start + pos),
			Name:         gen.Company().Name(),
			Description:  gen.Company().CatchPhrase(),
			Address:      gen.Address().Address()[1:],
			NumEmployees: gen.IntBetween(1, 10000),
			Country:      gen.Address().CountryAbbr(),
		})

		pos = pos + 1
	}

	params := &api.ImportDocumentsParams{
		Action:    pointer.String("create"),
		BatchSize: pointer.Int(BatchSize),
	}

	_, err = c.Client.Collection(Collection).Documents().Import(documents, params)
	if err != nil {
		return err
	}

	return nil
}
