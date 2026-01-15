package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBotDataSource_basic tests reading an existing bot via data source.
// This test requires CHATBOTKIT_API_KEY to be set in the environment
// and an existing bot to be available.
func TestAccBotDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// First create a bot resource, then read it via data source
				Config: testAccBotDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the data source can read the bot
					resource.TestCheckResourceAttrPair(
						"data.chatbotkit_bot.test", "id",
						"chatbotkit_bot.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.chatbotkit_bot.test", "name",
						"chatbotkit_bot.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.chatbotkit_bot.test", "description",
						"chatbotkit_bot.test", "description",
					),
				),
			},
		},
	})
}

// TestAccBotDataSource_withBackstory tests reading a bot with backstory via data source.
func TestAccBotDataSource_withBackstory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBotDataSourceConfigWithBackstory(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.chatbotkit_bot.test", "backstory",
						"chatbotkit_bot.test", "backstory",
					),
				),
			},
		},
	})
}

func testAccBotDataSourceConfig() string {
	return `
resource "chatbotkit_bot" "test" {
  name        = "test-bot-datasource"
  description = "Test bot for data source testing"
}

data "chatbotkit_bot" "test" {
  id = chatbotkit_bot.test.id
}
`
}

func testAccBotDataSourceConfigWithBackstory() string {
	return fmt.Sprintf(`
resource "chatbotkit_bot" "test" {
  name        = "test-bot-datasource-backstory"
  description = "Test bot with backstory"
  backstory   = "You are a helpful test assistant."
}

data "chatbotkit_bot" "test" {
  id = chatbotkit_bot.test.id
}
`)
}
