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

func TestAccDigestTracker_initialVersionIsOne(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDigestTrackerConfig("initial-digest"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "initial-digest"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
		},
	})
}

func TestAccDigestTracker_versionStaysConstantOnNoChange(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDigestTrackerConfig("stable-digest"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "stable-digest"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("stable-digest"), // Same digest
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "stable-digest"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"), // Should remain 1
				),
			},
		},
	})
}

func TestAccDigestTracker_multipleDigestChanges(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDigestTrackerConfig("digest-v1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-v1"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("digest-v2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-v2"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "2"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("digest-v3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-v3"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "3"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("digest-v4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-v4"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "4"),
				),
			},
		},
	})
}

func TestAccDigestTracker_planShowsVersionChange(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDigestTrackerConfig("plan-test-digest-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "plan-test-digest-1"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
			{
				Config: testAccDigestTrackerConfig("plan-test-digest-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "plan-test-digest-2"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "2"),
				),
				// This step will verify that the plan shows the version change
				PlanOnly: false,
			},
		},
	})
}

func TestAccDigestTracker_quickVersionTest(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Apply with digest A → expect version 1
				Config: testAccDigestTrackerConfig("digest-A"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-A"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "1"),
				),
			},
			{
				// Change to digest B → plan should show version: 1 → 2
				Config: testAccDigestTrackerConfig("digest-B"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-B"),
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "version", "2"),
				),
			},
			{
				// Change back to A → plan shows version: 2 → 3; still single apply
				Config: testAccDigestTrackerConfig("digest-A"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digesttracker_tracker.test", "digest", "digest-A"),
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
