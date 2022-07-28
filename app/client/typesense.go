package client

import (
	"fmt"
	"strconv"

	"github.com/jaswdr/faker"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

const (
	Collection = "companies"
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

	//host := viper.GetString("TYPESENSE_HOST")
	//port := viper.GetInt("TYPESENSE_PORT")
	//key := viper.GetString("TYPESENSE_API_KEY")

	host := "http://localhost"
	port := 8108
	key := "secret"

	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("%s:%d", host, port)),
		typesense.WithAPIKey(key),
	)

	return TypesenseClient{client}
}

func (c TypesenseClient) Search(expression string) (*api.SearchResult, error) {
	searchParameters := &api.SearchCollectionParams{
		Q:       expression,
		QueryBy: "name,description",
		//FilterBy: pointer.String("num_employees:>100"),
		SortBy:            pointer.String("num_employees:desc"),
		PerPage:           pointer.Int(250),
		HighlightStartTag: pointer.String("<code>"),
		HighlightEndTag:   pointer.String("</code>"),
		FacetBy:           pointer.String("country"),
	}

	return c.Client.Collection(Collection).Documents().Search(searchParameters)
}

func (c TypesenseClient) CreateSchema() error {

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

	gen := faker.New()
	pos := 1

	var documents []interface{}

	for pos <= num {

		documents = append(documents, &Company{
			Id:           strconv.Itoa(pos),
			Name:         gen.Company().Name(),
			Description:  gen.Company().CatchPhrase(),
			Address:      gen.Address().Address(),
			NumEmployees: gen.IntBetween(1, 10000),
			Country:      gen.Address().CountryAbbr(),
		})

		pos = pos + 1
	}

	params := &api.ImportDocumentsParams{
		Action:    pointer.String("create"),
		BatchSize: pointer.Int(num),
	}

	_, err := c.Client.Collection(Collection).Documents().Import(documents, params)
	if err != nil {
		return err
	}

	return nil
}
