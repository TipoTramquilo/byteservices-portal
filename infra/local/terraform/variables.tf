variable "vm_name" {
  type    = string
  default = "opnsense-local"
}

variable "image_path" {
  type        = string
  description = "Ruta local al .img descomprimido de OPNsense (la deja la tarea image)."
}

variable "pool" {
  type    = string
  default = "default"
}

variable "wan_bridge" {
  type        = string
  default     = "br0"
  description = "Puente hacia la red cableada del host (vtnet0 = WAN)."
}

variable "lan_bridge" {
  type        = string
  default     = "br-lan"
  description = "Puente de visitantes, donde está el AP (vtnet1 = LAN)."
}

variable "memory_mib" {
  type    = number
  default = 4096
}

variable "vcpu" {
  type    = number
  default = 2
}
