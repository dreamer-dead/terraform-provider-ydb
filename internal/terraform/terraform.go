package terraform

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	Endpoint string
	Token    string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ydb_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ydb_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ydb_topic": ydbTopicDataSource(),
			"ydb_table": ydbTableDataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ydb_topic": ydbTopicResource(),
			"ydb_table": ydbTableResource(),
		},
	}
}

func defaultTimeouts() *schema.ResourceTimeout {
	return &schema.ResourceTimeout{
		Create:  schema.DefaultTimeout(time.Minute * 20),
		Read:    schema.DefaultTimeout(time.Minute * 20),
		Update:  schema.DefaultTimeout(time.Minute * 20),
		Delete:  schema.DefaultTimeout(time.Minute * 20),
		Default: schema.DefaultTimeout(time.Minute * 20),
	}
}
