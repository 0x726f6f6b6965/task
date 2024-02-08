resource "kubernetes_deployment" "nginx_deployment" {
  metadata {
    name = "nginx-deployment"
    labels = {
      app = "nginx-app"
    }
  }
  spec {
    replicas = 3
    selector {
      match_labels = {
        app = "nginx-app"
      }
    }
    template {
      metadata {
        labels = {
          app = "nginx-app"
        }
      }
      spec {
        container {
          image = var.nginx_image
          name  = "nginx-container"
          volume_mount {
            mount_path = "/etc/nginx/nginx.conf"
            name       = "cfg"
          }
        }
        volume {
          name = "cfg"
          host_path {
            path = var.nginx_cfg
            type = "File"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "nginx_service" {
  metadata {
    name = "nginx-service"
  }
  spec {
    selector = {
      app = kubernetes_deployment.nginx_deployment.metadata.0.labels.app
    }
    port {
      port      = 80
      node_port = var.nginx_node_port
    }
    type = "NodePort"
  }
  depends_on = [kubernetes_deployment.nginx_deployment]
}
