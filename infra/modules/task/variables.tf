variable "task_image" {
  description = "Task image"
  type        = string
  default     = "task-svc:v0.0.1"
}

variable "task_cfg" {
  description = "Task configure file"
  type        = string
  default     = "application.yaml"
}

variable "env_config" {
  description = "The CONFIG of envirable value"
  type        = string
}

variable "label_app" {
  description = "The deployment label"
  type        = string
}

variable "service_name" {
  description = "The service name"
  type        = string
}

variable "deployment_name" {
  description = "The deployment name"
  type        = string
}
