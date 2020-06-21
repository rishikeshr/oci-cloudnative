# Copyright (c) 2019 Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at http://oss.oracle.com/licenses/upl.
# 

# creates an ATP database
resource "oci_database_autonomous_database" "mushop_autonomous_database" {
  admin_password                                 = random_string.autonomous_database_wallet_password.result
  compartment_id                                 = var.compartment_ocid
  cpu_core_count                                 = 1
  data_storage_size_in_tbs                       = 1
  db_name                                        = var.database_name
  display_name                                   = "${var.database_name}${random_id.mushop_id.dec}"
  freeform_tags                                  = local.common_tags
  is_free_tier                                   = true
  license_model                                  = "LICENSE_INCLUDED"
  is_preview_version_with_service_terms_accepted = false
}
