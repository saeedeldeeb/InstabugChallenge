package elastic

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es     *elasticsearch.TypedClient
	config Config
}

type Config struct {
	Addresses []string
	Username  string
	Password  string
}

// NewClient creates a new typed Elasticsearch client
func NewClient(cfg Config) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// Test connection
	info, err := client.Info().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping elasticsearch: %w", err)
	}

	fmt.Printf("Connected to Elasticsearch %s\n", info.Version.Int)

	return &Client{
		es:     client,
		config: cfg,
	}, nil
}
