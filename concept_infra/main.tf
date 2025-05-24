provider "aws" {
  region = "eu-central-1"
}

resource "aws_vpc" "agroport_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_security_group" "agroport_rds_sg" {
  name_prefix = "agroport-rds-sg"
  vpc_id = aws_vpc.agroport_vpc.id

  ingress {
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_rds_cluster" "agroport_db_cluster" {
  cluster_identifier = "agroport-db-cluster"
  engine = "aurora-postgresql"
  engine_mode = "serverless"
  engine_version = null
  database_name = "agroport_db"
  master_username = "admin"
  master_password = "example123"
  vpc_security_group_ids = [aws_security_group.agroport_sg.id]
  db_subnet_group_name = "agroport-subnet-group"
  scaling_configuration {
    auto_pause = true
    max_capacity = 16
    min_capacity = 2
    seconds_until_auto_pause = 300
    timeout_action = "ForceApplyCapacityChange"
  }
}

resource "aws_rds_cluster_instance" "agroport_db_instance" {
  count = 1
  identifier = "agroport-db-instance-${count.index}"
  instance_class = "db.t3.small"
  engine = "aurora-postgresql"
  db_subnet_group_name = "agroport-subnet-group"
  apply_immediately = true
  auto_minor_version_upgrade = true
  db_parameter_group_name = "default.aurora-mysql5.7"
  db_cluster_identifier = aws_rds_cluster.agroport_db_cluster.id
  tags = {
    Name = "agroport-db-instance-${count.index}"
  }
}

resource "aws_db_subnet_group" "agroport_subnet_group" {
  name = "agroport-subnet-group"
  subnet_ids = aws_vpc.agroport_vpc.private_subnets
}
