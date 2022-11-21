data "infra_destination" "my_cluster" {
  name = "my-cluster"
}

data "infra_credential" "my_credential" {}

provider "kubernetes" {
  host                   = data.infra_destination.my_cluster.kubernetes.endpoint
  cluster_ca_certificate = base64decode(data.infra_destination.my_cluster.kubernetes.certificate_authority[0].data)
  token                  = data.infra_credential.my_cluster.token
}
