terraform {
  required_providers {
    chatbotkit = {
      source = "chatbotkit/chatbotkit"
    }
  }
}

provider "chatbotkit" {
  # token can be set via CHATBOTKIT_TOKEN environment variable
}

# Create a dataset
resource "chatbotkit_dataset" "example" {
  name        = "Example Dataset"
  description = "An example knowledge base dataset"
  type        = "text"
}

# Output the dataset ID
output "dataset_id" {
  value = chatbotkit_dataset.example.id
}
