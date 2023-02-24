package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeDisk_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeDisk_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_compute_disk.foo", "google_compute_disk.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeDisk_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_disk" "foo" {
  name     = "tf-test-compute-disk-%{random_suffix}"
}

data "google_compute_disk" "foo" {
  name     = google_compute_disk.foo.name
  project  = google_compute_disk.foo.project
}
`, context)
}