package tf

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/curve-technology/terraform-provider-protogit/pkg/schemas"
	"github.com/curve-technology/terraform-provider-protogit/pkg/store"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchemas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchemasRead,
		Schema: map[string]*schema.Schema{
			"entries": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:     schema.TypeString,
							Required: true,
						},
						"section": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "value",
							ValidateDiagFunc: ValidateSection,
						},
						"filepath": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"schemas_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schema": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"references": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func ValidateSection(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Diagnostics{diag.Diagnostic{Summary: "invalid"}}
	}

	if v != "key" && v != "value" {
		return diag.Diagnostics{diag.Diagnostic{Summary: fmt.Sprintf("section argument can only be 'key' or 'value', got %q instead", v)}}
	}
	return nil
}

func dataSourceSchemasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	settings := m.(Settings)

	repo, err := store.NewGitRepo(settings.URL, settings.TagVersion, "/tmp/schemas", settings.ProtoPath)
	if err != nil {
		return diag.FromErr(err)
	}

	// Build entries
	entries := schemas.Entries{}

	entriesInput := d.Get("entries")
	entriesInputSet, _ := entriesInput.(*schema.Set)
	for _, entryInput := range entriesInputSet.List() {
		entryInputMap := entryInput.(map[string]interface{})

		entry := schemas.Entry{
			Topic:    entryInputMap["topic"].(string),
			Section:  schemas.Section(entryInputMap["section"].(string)),
			Filepath: entryInputMap["filepath"].(string),
		}
		entries = append(entries, entry)
	}

	records, err := schemas.BuildRecords(repo, entries)
	if err != nil {
		return diag.FromErr(err)
	}

	schemasIntfs := make([]map[string]interface{}, 0)
	for _, record := range records {
		item := make(map[string]interface{})
		item["subject"] = record.Subject
		item["schema"] = record.Schema

		refIntfs := make([]interface{}, len(record.References))
		for i, ref := range record.References {
			refIntfs[i] = ref
		}

		item["references"] = refIntfs

		schemasIntfs = append(schemasIntfs, item)
	}

	if err := d.Set("schemas_mapping", schemasIntfs); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
