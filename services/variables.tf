variable "kubeconfig" {
  type    = string
  default = "~/.kube/config"
}

variable "kubectx" {
  type    = string
  default = ""
}

variable "rode_version" {
  type    = string
  default = "0.14.5"
}

variable "grafeas_version" {
  type    = string
  default = "0.8.2"
}
