---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "protogit_schemas Data Source - terraform-provider-protogit"
subcategory: ""
description: |-
  
---

# protogit_schemas (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **entries** (Block Set, Min: 1) (see [below for nested schema](#nestedblock--entries))

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **schemas_mapping** (List of Object) (see [below for nested schema](#nestedatt--schemas_mapping))

<a id="nestedblock--entries"></a>
### Nested Schema for `entries`

Required:

- **filepath** (String)
- **topic** (String)

Optional:

- **section** (String)


<a id="nestedatt--schemas_mapping"></a>
### Nested Schema for `schemas_mapping`

Read-Only:

- **references** (List of String)
- **schema** (String)
- **subject** (String)


