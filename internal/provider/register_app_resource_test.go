package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRegisterAppResource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id":"563419","name":"test app","website":"https://example.com/","redirect_uri":"urn:ietf:wg:oauth:2.0:oob","client_id":"TWhM-tNSuncnqN7DBJmoyeLnk6K3iJJ71KKXxgL1hPM","client_secret":"ZEaFUFmF0umgBX1qKJDjaU99Q31lDkOU8NutzTOoliw","vapid_key":"BCk-QqERU0q-CfYZjcuB6lnyyOYfJ2AifKqfeGIm7Z-HiTU5T9eTG5GxVA0_OH5mMlI4UkkDTpaZwozy0TzdZ2M="}`)
		return
	}))
	defer ts.Close()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRegisterAppResourceConfig(ts.URL, "one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_register_app.test", "id", "563419"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.client_id", "TWhM-tNSuncnqN7DBJmoyeLnk6K3iJJ71KKXxgL1hPM"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.client_secret", "ZEaFUFmF0umgBX1qKJDjaU99Q31lDkOU8NutzTOoliw"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.redirect_uri", "urn:ietf:wg:oauth:2.0:oob"),
				),
			},
			// Create and Read testing
			{
				Config: testAccRegisterAppResourceConfig(ts.URL, "two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_register_app.test", "id", "563419"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.client_id", "TWhM-tNSuncnqN7DBJmoyeLnk6K3iJJ71KKXxgL1hPM"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.client_secret", "ZEaFUFmF0umgBX1qKJDjaU99Q31lDkOU8NutzTOoliw"),
					resource.TestCheckResourceAttr("mastodon_register_app.test", "app_config.redirect_uri", "urn:ietf:wg:oauth:2.0:oob"),
				),
			},
		},
	})
}

const testAccRegisterAppResourceConfigTmplPre = `
provider "mastodon" {
	domain = %[1]q
	use_https = false
}

resource "mastodon_register_app" "test" {
    client_name = %[2]q
}
`

func testAccRegisterAppResourceConfig(tsURL string, clientName string) string {
	return fmt.Sprintf(
		testAccRegisterAppResourceConfigTmplPre,
		strings.TrimPrefix(tsURL, "http://"),
		clientName,
	)
}
