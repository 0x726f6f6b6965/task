provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "docker-desktop"
}

module "redis" {
  source      = "./modules/redis"
  redis_image = var.redis_image
}


module "nginx" {
  source          = "./modules/nginx"
  nginx_image     = var.nginx_image
  nginx_cfg       = "${path.cwd}/${var.nginx_cfg}"
  nginx_node_port = 30001
  depends_on      = [module.redis, module.task_1, module.task_2]
}

module "task_1" {
  source          = "./modules/task"
  task_image      = var.task_image
  task_cfg        = "${path.cwd}/application-1.yaml"
  env_config      = var.env_config
  label_app       = "task-app-1"
  service_name    = "taks-svc1"
  deployment_name = "task-deployment-1"
  depends_on      = [module.redis]
}

module "task_2" {
  source          = "./modules/task"
  task_image      = var.task_image
  task_cfg        = "${path.cwd}/application-2.yaml"
  env_config      = var.env_config
  label_app       = "task-app-2"
  service_name    = "taks-svc2"
  deployment_name = "task-deployment-2"
  depends_on      = [module.redis]
}
