---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tableau_group_membership Resource - terraform-provider-tableau"
subcategory: ""
description: |-
  
---

# tableau_group_membership (Resource)



## Example Usage

```terraform
resource "tableau_group" "test_group" {
  name = "Test Group"
}

resource "tableau_user" "test_user" {
  email        = "test_user@example.com"
  site_role    = "Unlicensed"
  auth_setting = "OpenID"
}

resource "tableau_group_membership" "test_group_membership" {
  group_id = tableau_group.test_group.id
  users = [
    tableau_user.test_user.email,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) Group id
- `users` (Set of String) List of user emails

## Import

Import is supported using the following syntax:

```shell
# Group membership can be imported by specifying the group identifier.
terraform import tableau_group_membership.test_group_membership de7373bd-ff18-4dab-a579-78e3dcd5ceb4
```
