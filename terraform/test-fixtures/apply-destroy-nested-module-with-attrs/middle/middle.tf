variable param {}

module "bottom" {
  source       = "./bottom"
  bottom_param = "${var.param}"
}

resource "null_resource" "middle" {
  val = "${var.param}"
}
