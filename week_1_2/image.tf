# Use existing Debian 12 image
# If you need to create it manually first, run:
# source /opt/stack/devstack/openrc admin admin
# openstack image create --disk-format qcow2 --container-format bare --public --file /opt/stack/devstack/debian-12-genericcloud-amd64.qcow2 debian12

data "openstack_images_image_v2" "debian12" {
  name        = "debian12"
  most_recent = true
}
