variable "nginx_image" {
  description = "Nginx image"
  type        = string
  default     = "nginx:alpine"
}

variable "nginx_cfg" {
  description = "Nginx configure file"
  type        = string
}

variable "nginx_node_port" {
  description = "The node port"
  type        = number
  validation {
    condition     = (var.nginx_node_port) >= 30000 && (var.nginx_node_port) <= 32767
    error_message = "The port range must be in the range 30000 to 32767"
  }
}
