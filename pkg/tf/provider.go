package tf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"proto_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "proto",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"protogit_schemas": dataSourceSchemas(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	tagVersion := d.Get("tag_version").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	protoPath := d.Get("proto_path").(string)

	credentials := ""
	if username != "" && password != "" {
		credentials = fmt.Sprintf("%s:%s@", username, password)
	}

	fullURL := fmt.Sprintf("https://%s%s", credentials, url)

	settings := Settings{URL: fullURL, TagVersion: tagVersion, ProtoPath: protoPath}

	return settings, diag.Diagnostics{}
}
