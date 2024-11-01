variable "cache_name" {
  description = "Name of the cache cluster"
  type        = string
  default     = "GoShorty-${var.env}-Cache"
}

resource "aws_elasticache_cluster" "cache" {
  cluster_id           = var.cache_name
  engine               = "redis"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis6.x" # Redis OSS
}