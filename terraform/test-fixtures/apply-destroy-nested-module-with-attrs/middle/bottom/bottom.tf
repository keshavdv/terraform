variable bottom_param {}

resource "null_resource" "bottom" {
  triggers {
    bp = "${var.bottom_param}"
  }
}
