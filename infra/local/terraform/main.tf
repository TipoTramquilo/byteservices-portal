resource "libvirt_volume" "disk" {
  name   = "${var.vm_name}.raw"
  pool   = var.pool
  source = var.image_path
  format = "raw"
}

resource "libvirt_domain" "opnsense" {
  name   = var.vm_name
  type   = "kvm"
  memory = var.memory_mib
  vcpu   = var.vcpu

  cpu {
    mode = "host-passthrough"
  }

  disk {
    volume_id = libvirt_volume.disk.id
  }

  network_interface {
    bridge = var.wan_bridge
  }

  network_interface {
    bridge = var.lan_bridge
  }

  video {
    type = "qxl"
  }

  graphics {
    type        = "spice"
    listen_type = "address"
  }
}
