package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAccountDataSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"username":"user","acct":"user","display_name":"Cool Dude","created_at":"2019-06-09T00:00:00.000Z","url":"https://example.com/user/tyr","discoverable":true}"`)
		return
	}))
	defer ts.Close()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountDataSourceConfig(ts.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mastodon_account.test", "username", "user"),
					resource.TestCheckResourceAttr("data.mastodon_account.test", "account", "user"),
					resource.TestCheckResourceAttr("data.mastodon_account.test", "display_name", "Cool Dude"),
					resource.TestCheckResourceAttr("data.mastodon_account.test", "created_at", "2019-06-09 00:00:00 +0000 UTC"),
					resource.TestCheckResourceAttr("data.mastodon_account.test", "url", "https://example.com/user/tyr"),
					resource.TestCheckResourceAttr("data.mastodon_account.test", "discoverable", "true"),
				),
			},
		},
	})
}

const testAccAccountDataSourceConfigTmplPre = `
provider "mastodon" {
	domain = %[1]q
	use_https = false
}

data "mastodon_account" "test" {
	id = "1"
}
`

func testAccAccountDataSourceConfig(tsURL string) string {
	return fmt.Sprintf(testAccAccountDataSourceConfigTmplPre, strings.TrimPrefix(tsURL, "http://"))
}
