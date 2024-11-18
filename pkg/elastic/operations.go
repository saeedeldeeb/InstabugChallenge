package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"

	_ "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Index indexes a document
func (c *Client) Index(ctx context.Context, index string, id string, document interface{}) error {
	res, err := c.es.Index(index).Id(id).Document(document).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	if res.Result != result.Created && res.Result != result.Updated {
		return fmt.Errorf("unexpected index result: %s", res.Result)
	}

	return nil
}

// Get retrieves a document by ID (internal method)
func (c *Client) get(ctx context.Context, index string, id string) ([]byte, error) {
	res, err := c.es.Get(index, id).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return res.Source_, nil
}

// Get retrieves a document by ID
func Get[T any](c *Client, ctx context.Context, index string, id string) (*T, error) {
	source, err := c.get(ctx, index, id)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(source, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &result, nil
}

// SearchOptions represents search parameters
type SearchOptions struct {
	From  *int
	Size  *int
	Query *types.Query
}

// Search performs a search query with type safety
func Search[T any](c *Client, ctx context.Context, index string, opts SearchOptions) ([]T, error) {
	// Build search request
	req := c.es.Search().Index(index)

	if opts.From != nil {
		req.From(*opts.From)
	}
	if opts.Size != nil {
		req.Size(*opts.Size)
	}
	if opts.Query != nil {
		req.Query(opts.Query)
	}

	// Execute search
	res, err := req.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	// Process results
	var results []T
	for _, hit := range res.Hits.Hits {
		var item T
		if err := json.Unmarshal(hit.Source_, &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal hit: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// BuildMatchQuery creates a match query
func BuildMatchQuery(field string, value interface{}) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			field: {Query: fmt.Sprintf("%v", value)},
		},
	}
}

// BuildTermQuery creates a term query
func BuildTermQuery(field string, value interface{}) *types.Query {
	return &types.Query{
		Term: map[string]types.TermQuery{
			field: {Value: value},
		},
	}
}
