package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigestTracker_versionIncrementsOnDigestChange(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDigestTrackerConfig("first-digest-value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "first-digest-value"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("changed-digest-value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "changed-digest-value"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "2"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("changed-digest-value"), // no change this time
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "changed-digest-value"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "2"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("yet-another-digest"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "yet-another-digest"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "3"),
				),
			},
		},
	})
}

func testAccDigestTrackerConfig(digest string) string {
	return fmt.Sprintf(`
provider "digesttracker" {}

resource "digesttracker_tracker" "test" {
  digest = "%s"
}
`, digest)
}
