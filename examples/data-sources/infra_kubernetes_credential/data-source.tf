data "infra_kubernetes_cluster" "my_cluster" {
  name = "my-cluster"
}

data "infra_kubernetes_credential" "my_credential" {}

provider "kubernetes" {
  host                   = data.infra_kubernetes_cluster.my_cluster.endpoint
  cluster_ca_certificate = base64decode(data.infra_kubernetes_cluster.my_cluster.certificate_authority[0].data)
  token                  = data.infra_kubernetes_credential.my_cluster.token
}
