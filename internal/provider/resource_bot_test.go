package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBotResource_basic tests the basic lifecycle of a bot resource.
// This test requires CHATBOTKIT_API_KEY to be set in the environment.
// See the README.md for instructions on obtaining an API key.
func TestAccBotResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBotResourceConfig("test-bot", "Test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "name", "test-bot"),
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "description", "Test description"),
					resource.TestCheckResourceAttrSet("chatbotkit_bot.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "chatbotkit_bot.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccBotResourceConfig("test-bot-updated", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "name", "test-bot-updated"),
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "description", "Updated description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccBotResource_withBackstory tests creating a bot with a backstory.
func TestAccBotResource_withBackstory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBotResourceConfigWithBackstory("bot-with-backstory", "A bot with backstory", "You are a helpful assistant."),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "name", "bot-with-backstory"),
					resource.TestCheckResourceAttr("chatbotkit_bot.test", "backstory", "You are a helpful assistant."),
				),
			},
		},
	})
}

func testAccBotResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "chatbotkit_bot" "test" {
  name        = %q
  description = %q
}
`, name, description)
}

func testAccBotResourceConfigWithBackstory(name, description, backstory string) string {
	return fmt.Sprintf(`
resource "chatbotkit_bot" "test" {
  name        = %q
  description = %q
  backstory   = %q
}
`, name, description, backstory)
}
