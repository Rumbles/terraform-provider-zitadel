data "zitadel_machine_users" "default" {
  org_id           = "123456789012345678"
  user_name        = "example-name"
  user_name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_machine_user" "default" {
  for_each = toset(data.zitadel_machine_users.default.user_ids)
  id       = each.value
}

output "user_names" {
  value = toset([
    for user in data.zitadel_machine_user.default : user.name
  ])
}
