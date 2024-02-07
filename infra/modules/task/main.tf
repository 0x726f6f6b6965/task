resource "kubernetes_deployment" "task_deployment" {
  metadata {
    name = var.deployment_name
    labels = {
      app = var.label_app
    }
  }
  spec {
    replicas = 3
    selector {
      match_labels = {
        app = var.label_app
      }
    }
    template {
      metadata {
        labels = {
          app = var.label_app
        }
      }
      spec {
        container {
          image = var.task_image
          name  = "task-container"
          volume_mount {
            mount_path = "/app/application.yaml"
            name       = "task-cfg"
          }
          env {
            name  = "CONFIG"
            value = var.env_config
          }
        }
        volume {
          name = "task-cfg"
          host_path {
            path = var.task_cfg
            type = "File"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "task_service" {
  metadata {
    name = var.service_name
  }
  spec {
    selector = {
      app = kubernetes_deployment.task_deployment.metadata.0.labels.app
    }
    port {
      port        = 64530
      target_port = 64530
    }
    type = "ClusterIP"
  }
  depends_on = [kubernetes_deployment.task_deployment]
}
