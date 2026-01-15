terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # api_key = "..." # Or set CHATBOTKIT_API_KEY env var
}

resource "chatbotkit_bot" "example" {
  name        = "My Test Bot"
  description = "Created via Terraform"
  backstory   = "You are a helpful assistant."
}
