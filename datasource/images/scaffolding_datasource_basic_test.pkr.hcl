
variable "access_key" {
  type = string
  default = env("ALICLOUD_ACCESS_KEY")
}

variable "secret_key" {
  type = string
  default = env("ALICLOUD_ACCESS_SECRET")
}

data "alicloud-images" "test" {
 
  image_name  = "ubuntu*"
  region      = "cn-shanghai"
}


source "alicloud-ecs" "packer-build" {
      access_key = var.access_key
      secret_key = var.secret_key
      
      region = "cn-shanghai"
      image_name = "test_image_${var.build_number}"
  
      source_image = "${data.alicloud-images.test}"
      ssh_username = "root"
      instance_type = "ecs.c5.large"
      io_optimized = true
      internet_charge_type = "PayByTraffic"
      image_force_delete = true

      
      tags = {
        "Name" = "Test Datasource Image"
      }
}

build {
  sources = ["sources.alicloud-ecs.packer-build"]
}
