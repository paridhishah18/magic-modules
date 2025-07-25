package spanner_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests

func TestAccSpannerInstance_basic(t *testing.T) {
	t.Parallel()

	idName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(idName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_basicUpdateWithProviderDefaultLabels(t *testing.T) {
	t.Parallel()

	idName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithProviderLabel(idName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithProviderLabel(idName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_noNodeCountSpecified(t *testing.T) {
	// Cannot be run in VCR because no API calls are made
	acctest.SkipIfVcr(t)
	t.Parallel()

	idName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSpannerInstance_noNodeCountSpecified(idName),
				ExpectError: regexp.MustCompile(".*one of\n`autoscaling_config,instance_type,num_nodes,processing_units` must be\nspecified.*"),
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	// Since we're testing the autogenerated name specifically here, we can't use VCR. This shouldn't be copy /
	// pasted to other configs, though.
	acctest.SkipIfVcr(t)
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s-dname", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutogenName(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "name"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_update(t *testing.T) {
	t.Parallel()

	// Update display name, but keep real name consistent, as it cannot be
	// updated after creation.
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	dName1 := fmt.Sprintf("tf-test-dname1-%s", acctest.RandString(t, 10))
	dName2 := fmt.Sprintf("tf-test-dname2-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_update(name, dName1, 1, false),
			},
			{
				ResourceName:            "google_spanner_instance.updater",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_update(name, dName2, 2, true),
			},
			{
				ResourceName:            "google_spanner_instance.updater",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_virtualUpdate(t *testing.T) {
	t.Parallel()

	dName := fmt.Sprintf("tf-test-dname1-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_virtualUpdate(dName, "true"),
			},
			{
				ResourceName: "google_spanner_instance.basic",
				ImportState:  true,
			},
			{
				Config: testAccSpannerInstance_virtualUpdate(dName, "false"),
			},
			{
				ResourceName: "google_spanner_instance.basic",
				ImportState:  true,
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingProcessingUnitConfig(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s-dname", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigs(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingProcessingUnitConfigUpdate(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigsUpdate(displayName, 1000, 2000, 65, 95),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigsUpdate(displayName, 2000, 3000, 75, 90),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingNodeConfig(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s-dname", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigs(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingNodeConfigUpdate(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigsUpdate(displayName, 1, 2, 65, 95),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigsUpdate(displayName, 2, 3, 75, 90),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingNodeConfigUpdateDisableAutoscaling(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigsUpdate(displayName, 1, 2, 65, 95),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithNodes(displayName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutoscalingUsingPrecessingUnitsConfigUpdateDisableAutoscaling(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigs(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithProcessingUnits(displayName, 1000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAsymmetricAutoscalingConfigsUpdate(t *testing.T) {
	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_main(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.main", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAsymmetricAutoscalingConfigsUpdate(displayName, 1, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.main", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_basicWithAsymmetricAutoscalingConfigsUpdate(displayName, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.main", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_spannerInstanceWithAutoscaling(t *testing.T) {

	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_spannerInstanceWithAutoscaling(context),
			},
			{
				ResourceName:            "google_spanner_instance.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccSpannerInstance_freeInstanceBasicUpdate(t *testing.T) {
	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_freeInstanceBasic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.main", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstance_freeInstanceBasicUpdate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.main", "state"),
				),
			},
			{
				ResourceName:            "google_spanner_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccSpannerInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"

  processing_units             = 100
  edition                      = "ENTERPRISE"
  default_backup_schedule_type = "NONE"
}
`, name, name)
}

func testAccSpannerInstance_basicWithNodes(name string, nodes int) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"

  num_nodes                    = %d
  edition                      = "ENTERPRISE"
  default_backup_schedule_type = "NONE"
}
`, name, name, nodes)
}

func testAccSpannerInstance_basicWithProcessingUnits(name string, processingUnits int) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"

  processing_units             = %d
  edition                      = "ENTERPRISE"
  default_backup_schedule_type = "NONE"
}
`, name, name, processingUnits)
}

func testAccSpannerInstance_basicWithProviderLabel(name string, addLabel bool) string {
	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
provider "google" {
  alias          = "with-labels"
  default_labels = {
    %s
  }
}

resource "google_spanner_instance" "basic" {
  provider     = google.with-labels
  config       = "regional-us-central1"
  name         = "%s"
  display_name = "%s"

  processing_units = 100

  labels = {
    "key1" = "value1"
  }
}
`, extraLabel, name, name)
}

func testAccSpannerInstance_noNodeCountSpecified(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"
}
`, name, name)
}

func testAccSpannerInstance_basicWithAutogenName(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  config       = "regional-us-central1"
  display_name = "%s"

  processing_units = 100
}
`, name)
}

func testAccSpannerInstance_update(name, dname string, nodes int, addLabel bool) string {
	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
resource "google_spanner_instance" "updater" {
  config       = "regional-us-central1"
  name         = "%s"
  display_name = "%s"
  num_nodes    = %d

  labels = {
    "key1" = "value1"
    %s
  }
}
`, name, dname, nodes, extraLabel)
}

func testAccSpannerInstance_virtualUpdate(name, virtual string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  display_name = "%s"
  config       = "regional-us-central1"

  processing_units = 100
  force_destroy    = "%s"
}
`, name, name, virtual)
}

func testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigs(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  autoscaling_config {
    autoscaling_limits {
      max_processing_units            = 2000
      min_processing_units            = 1000
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      storage_utilization_percent           = 95
    }
  }
  edition      = "ENTERPRISE"
}
`, name, name)
}

func testAccSpannerInstance_basicWithAutoscalerConfigUsingProcessingUnitsAsConfigsUpdate(name string, minProcessingUnits, maxProcessingUnits, cupUtilizationPercent, storageUtilizationPercent int) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  autoscaling_config {
    autoscaling_limits {
      max_processing_units            = %v
      min_processing_units            = %v
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = %v
      storage_utilization_percent           = %v
    }
  }
  edition      = "ENTERPRISE"
  default_backup_schedule_type = "AUTOMATIC"
}
`, name, name, maxProcessingUnits, minProcessingUnits, cupUtilizationPercent, storageUtilizationPercent)
}

func testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigs(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  autoscaling_config {
    autoscaling_limits {
      max_nodes            = 2
      min_nodes            = 1
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      storage_utilization_percent           = 95
    }
  }
  edition      = "ENTERPRISE"
}
`, name, name)
}

func testAccSpannerInstance_basicWithAutoscalerConfigUsingNodesAsConfigsUpdate(name string, minNodes, maxNodes, cupUtilizationPercent, storageUtilizationPercent int) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  autoscaling_config {
    autoscaling_limits {
      max_nodes           = %v
      min_nodes           = %v
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = %v
      storage_utilization_percent           = %v
    }
  }
  edition      = "ENTERPRISE"
}
`, name, name, maxNodes, minNodes, cupUtilizationPercent, storageUtilizationPercent)
}

func testAccSpannerInstance_main(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "main" {
  name         = "%s"
  config       = "nam-eur-asia3"
  display_name = "%s"
  num_nodes    = 1
  edition      = "ENTERPRISE_PLUS"
}
`, name, name)
}

func testAccSpannerInstance_basicWithAsymmetricAutoscalingConfigsUpdate(name string, minNodes, maxNodes int) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}

resource "google_spanner_instance" "main" {
  provider     = google.user-project-override
  name         = "%s"
  config       = "nam-eur-asia3"
  display_name =  "%s"
  autoscaling_config {
    autoscaling_limits {
      max_nodes = 3
      min_nodes = 1
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 75
      storage_utilization_percent           = 90
    }
    asymmetric_autoscaling_options {
      replica_selection {
        location = "europe-west1"
      }
      overrides {
        autoscaling_limits {
          min_nodes = 3
          max_nodes = 30
        }
      }
    }
    asymmetric_autoscaling_options {
      replica_selection {
        location = "asia-east1"
      }
      overrides {
        autoscaling_limits {
          min_nodes = %d
          max_nodes = %d
        }
      }
    }
  }
  edition = "ENTERPRISE_PLUS"
}`, name, name, minNodes, maxNodes)
}

func testAccSpannerInstance_spannerInstanceWithAutoscaling(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "example" {
  name         = "tf-test-spanner-instance-%{random_suffix}"
  config       = "regional-us-central1"
  display_name = "Test Spanner Instance"
  autoscaling_config {
    autoscaling_limits {
      // Define the minimum and maximum compute capacity allocated to the instance
      // Either use nodes or processing units to specify the limits,
      // but should use the same unit to set both the min_limit and max_limit.
      max_processing_units            = 3000 // OR max_nodes  = 3
      min_processing_units            = 2000 // OR min_nodes = 2
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 75
      storage_utilization_percent           = 90
    }
  }
  edition = "ENTERPRISE"

  labels = {
    "foo" = "bar"
  }
}
`, context)
}

func testAccSpannerInstance_freeInstanceBasic(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "main" {
  name          = "%s"
  config        = "regional-europe-west1"
  display_name  = "%s"
  instance_type = "FREE_INSTANCE"
}
`, name, name)
}

func testAccSpannerInstance_freeInstanceBasicUpdate(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "main" {
  name          = "%s"
  config        = "nam-eur-asia3"
  display_name  = "%s"
  edition       = "ENTERPRISE_PLUS"
  instance_type = "PROVISIONED"
  num_nodes     = 1
}
`, name, name)
}
