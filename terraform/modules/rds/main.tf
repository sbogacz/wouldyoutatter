data "aws_kms_alias" "rds" {
  name = "alias/aws/rds"
}

resource "aws_db_parameter_group" "aurora_settings" {
  family = "aurora-mysql5.7"
  name   = "${var.db_name}-aurora-settings"

  parameter {
    name  = "max_connections"
    value = 2101
  }

  tags = "${var.tags}"
}

resource "aws_db_subnet_group" "db-subnet-group" {
  name = "${var.db_name}-db-subnet-group_${var.environment}"

  subnet_ids = ["${var.subnet_ids}"]

  tags = "${var.tags}"
}

resource "aws_rds_cluster_instance" "aurora_instance" {
  count             = "${length(var.availability_zones)}"
  availability_zone = "${element(var.availability_zones, count.index)}"

  identifier              = "${var.db_name}-db-${var.environment}-${count.index}"
  cluster_identifier      = "${aws_rds_cluster.db_cluster.id}"
  db_subnet_group_name    = "${aws_db_subnet_group.db-subnet-group.id}"
  instance_class          = "${var.instance_class}"
  db_parameter_group_name = "${aws_db_parameter_group.aurora_settings.id}"
  tags                    = "${var.tags}"
  engine                  = "aurora-mysql"
  engine_version          = "5.7.12"
}

resource "aws_rds_cluster" "db_cluster" {
  cluster_identifier = "${var.db_name}-db-${var.environment}"

  availability_zones   = "${var.availability_zones}"
  db_subnet_group_name = "${aws_db_subnet_group.db-subnet-group.id}"

  vpc_security_group_ids = ["${var.vpc_sg_ids}"]

  apply_immediately         = true
  final_snapshot_identifier = "${var.rds_final_snapshot_id}"
  skip_final_snapshot       = "${var.rds_skip_final_snapshot}"

  engine         = "aurora-mysql"
  engine_version = "5.7.12"

  database_name   = "${var.db_name}"
  master_username = "${var.db_username}"
  master_password = "${var.db_password}"
  port            = 3306

  tags = "${var.tags}"

  # TODO: Can we have/should we have cross-region aurora replicas sharing
  # single-region keys? Find someone who can drop some knowledge
  # encrypt at rest
  storage_encrypted = true

  # KMS keys are part of IAM, and as such are always in the same region
  kms_key_id = "${data.aws_kms_alias.rds.target_key_arn}"
}

# NB: some ideas taken from here: https://github.com/hashicorp/terraform/issues/5333

