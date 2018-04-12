output "db_endpoint" {
  description = "the endpoint for the created Aurora cluster"
  value       = "${aws_rds_cluster.db_cluster.endpoint}"
}

output "db_ro_endpoint" {
  description = "the read-only endpoint for the created Aurora cluster"
  value       = "${aws_rds_cluster.db_cluster.reader_endpoint}"
}
