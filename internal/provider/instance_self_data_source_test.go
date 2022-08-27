package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccInstanceSelfDataSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"email":"user@example.com","thumbnail":"https://example.com/image.jpg","title":"Example Title","uri":"example.com","version":"1.2.4"}`)
		return
	}))
	defer ts.Close()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceSelfDataSourceConfig(ts.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "id", "example.com"),
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "email", "user@example.com"),
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "thumbnail", "https://example.com/image.jpg"),
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "title", "Example Title"),
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "uri", "example.com"),
					resource.TestCheckResourceAttr("data.mastodon_instance_self.test", "version", "1.2.4"),
				),
			},
		},
	})
}

const testAccInstanceSelfDataSourceConfigTmplPre = `
provider "mastodon" {
	domain = %[1]q
	use_https = false
}

data "mastodon_instance_self" "test" {
}
`

func testAccInstanceSelfDataSourceConfig(tsURL string) string {
	return fmt.Sprintf(testAccInstanceSelfDataSourceConfigTmplPre, strings.TrimPrefix(tsURL, "http://"))
}
