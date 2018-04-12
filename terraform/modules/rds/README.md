# AWS RDS Module

## Variables
* aws_region: string
  	* where the aws region is
  	* defaults     = us-east-2
* availability_zones 
  	* list of AZs that we want to deploy to. Cardinality should match public/private/database subnet lists
  	* defaults     = [us-east-2a, us-east-2b, us-east-2c]
* account_id: string
  	* defaults     = 797387530767
  	* the dev account we have access to
* environment: string 
  	* development, production, etc.
  	* defaults     = development
* instance_class: string 
  	* db.t2.small, etc.
  	* defaults     = db.t2.small
* tags: map
  	* mandatory tags to prevent Krampus from our killing services
* db_password: string
  	* password for the database that holds the template and version tables
* db_username: string
  	* the user for the database
* db_name: string
  	* the name of the database
* db_identifier: string
  	* the identifier given to the database
* vpc_sg_ids: list
  	* a list of the security groups associated with the desired VPC
* db_subnet_ids: list
  	* a list of the subnets in the VPC designated for the DB to use
* rds_final_snapshot_id: string
	* the id to give the DB's final snapshot on deletion. Requires rds_skip_final_snapshot to be false
* rds_skip_final_snapshot: boolean
	* indicates whether AWS should store a snapshot of the DB prior to deletion

## Outputs
* db_instance_endpoint: string

