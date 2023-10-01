resource "tableau_user" "test_user" {
  email        = "test_user@example.com"
  site_role    = "Unlicensed"
  auth_setting = "OpenID"
}
