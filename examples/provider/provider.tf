terraform {
  required_providers {
    stratoid = {
      source = "stratoid.dev/azure/stratoid"
    }
  }
}

provider "stratoid" {

  use_cli = false
}

data "stratoid_user_flow" "name" {
  id = "97c40800-12ba-4520-b46a-637ca9e484a5"
}


output "user_flow" {
  value = data.stratoid_user_flow.name.display_name
  
}