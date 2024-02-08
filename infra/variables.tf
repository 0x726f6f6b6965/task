variable "nginx_image" {
  description = "Nginx image"
  type        = string
  default     = "nginx:1.23.2"
}

variable "nginx_cfg" {
  description = "Nginx configure file"
  type        = string
  default     = "nginx.conf"
}

variable "redis_image" {
  description = "Redis image"
  type        = string
  default     = "redis:7.2.3"
}

variable "task_image" {
  description = "Task image"
  type        = string
  default     = "task-svc:v0.0.1"
}

variable "env_config" {
  description = "The CONFIG of envirable value"
  type        = string
  default     = "/app/application.yaml"
}
