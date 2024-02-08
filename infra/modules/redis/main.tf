resource "kubernetes_deployment" "redis_deployment" {
  metadata {
    name = "redis-deployment"
    labels = {
      app = "redis-app"
    }
  }
  spec {
    selector {
      match_labels = {
        app = "redis-app"
      }
    }
    template {
      metadata {
        labels = {
          app = "redis-app"
        }
      }
      spec {
        container {
          image = var.redis_image
          name  = "redis-container"
        }
      }
    }
  }
}

resource "kubernetes_service" "redis_service" {
  metadata {
    name = "redis-service"
  }
  spec {
    selector = {
      app = kubernetes_deployment.redis_deployment.metadata.0.labels.app
    }
    port {
      port        = 6379
      target_port = 6379
    }
    type = "ClusterIP"
  }
  depends_on = [kubernetes_deployment.redis_deployment]
}
