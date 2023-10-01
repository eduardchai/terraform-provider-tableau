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
