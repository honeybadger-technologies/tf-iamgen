package parser

import (
	"strings"
)

// Resource type constants for AWS services
const (
	// S3
	ResourceS3Bucket = "aws_s3_bucket"

	// EC2
	ResourceEC2Instance   = "aws_instance"
	ResourceSecurityGroup = "aws_security_group"
	ResourceVPC           = "aws_vpc"
	ResourceSubnet        = "aws_subnet"

	// RDS
	ResourceDBInstance = "aws_db_instance"

	// IAM
	ResourceIAMRole   = "aws_iam_role"
	ResourceIAMPolicy = "aws_iam_policy"
	ResourceIAMUser   = "aws_iam_user"

	// Lambda
	ResourceLambdaFunction = "aws_lambda_function"

	// DynamoDB
	ResourceDynamoDBTable = "aws_dynamodb_table"

	// SNS
	ResourceSNSTopic = "aws_sns_topic"

	// SQS
	ResourceSQSQueue = "aws_sqs_queue"

	// CloudWatch
	ResourceCloudWatchLogGroup = "aws_cloudwatch_log_group"
)

// SupportedResourceTypes contains all AWS resource types that can be analyzed
var SupportedResourceTypes = map[string]bool{
	// S3
	"aws_s3_bucket":                    true,
	"aws_s3_bucket_versioning":         true,
	"aws_s3_bucket_encryption":         true,
	"aws_s3_bucket_policy":             true,
	"aws_s3_bucket_acl":                true,
	"aws_s3_bucket_cors_configuration": true,
	"aws_s3_object":                    true,

	// EC2
	"aws_instance":            true,
	"aws_security_group":      true,
	"aws_security_group_rule": true,
	"aws_vpc":                 true,
	"aws_subnet":              true,
	"aws_route_table":         true,
	"aws_route":               true,
	"aws_internet_gateway":    true,
	"aws_nat_gateway":         true,
	"aws_eip":                 true,
	"aws_network_interface":   true,
	"aws_key_pair":            true,
	"aws_volume":              true,

	// RDS
	"aws_db_instance":        true,
	"aws_db_parameter_group": true,
	"aws_db_subnet_group":    true,
	"aws_db_security_group":  true,

	// IAM
	"aws_iam_role":                    true,
	"aws_iam_policy":                  true,
	"aws_iam_user":                    true,
	"aws_iam_group":                   true,
	"aws_iam_role_policy":             true,
	"aws_iam_user_policy":             true,
	"aws_iam_group_policy":            true,
	"aws_iam_role_policy_attachment":  true,
	"aws_iam_user_policy_attachment":  true,
	"aws_iam_group_policy_attachment": true,
	"aws_iam_instance_profile":        true,

	// Lambda
	"aws_lambda_function":      true,
	"aws_lambda_alias":         true,
	"aws_lambda_permission":    true,
	"aws_lambda_layer_version": true,

	// DynamoDB
	"aws_dynamodb_table":      true,
	"aws_dynamodb_table_item": true,

	// SNS
	"aws_sns_topic":              true,
	"aws_sns_topic_policy":       true,
	"aws_sns_topic_subscription": true,

	// SQS
	"aws_sqs_queue":        true,
	"aws_sqs_queue_policy": true,

	// CloudWatch
	"aws_cloudwatch_log_group":           true,
	"aws_cloudwatch_log_stream":          true,
	"aws_cloudwatch_log_resource_policy": true,

	// KMS
	"aws_kms_key":   true,
	"aws_kms_alias": true,

	// Data sources (also important for understanding dependencies)
	"aws_ami":                true,
	"aws_availability_zones": true,
}

// IsAWSResource checks if a resource type is an AWS resource
func IsAWSResource(resourceType string) bool {
	return strings.HasPrefix(resourceType, "aws_")
}

// IsKnownResource checks if a resource type is known and supported
func IsKnownResource(resourceType string) bool {
	return SupportedResourceTypes[resourceType]
}

// GetServiceFromResourceType extracts the AWS service name from a resource type
// e.g., "aws_s3_bucket" -> "s3"
func GetServiceFromResourceType(resourceType string) string {
	parts := strings.Split(resourceType, "_")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

// GetResourceCategoryFromResourceType gets the resource category
// e.g., "aws_s3_bucket" -> "bucket", "aws_security_group_rule" -> "security_group_rule"
func GetResourceCategoryFromResourceType(resourceType string) string {
	parts := strings.Split(resourceType, "_")
	if len(parts) >= 3 {
		// Join everything after "aws_"
		return strings.Join(parts[2:], "_")
	}
	if len(parts) == 2 {
		return parts[1]
	}
	return "unknown"
}

// ResourceMetadata contains information about a resource type
type ResourceMetadata struct {
	Type        string // Resource type (e.g., "aws_s3_bucket")
	Service     string // AWS service (e.g., "s3")
	Category    string // Resource category (e.g., "bucket")
	IsSupported bool   // Whether this resource is in the supported list
	Description string // Human-readable description
}

// GetResourceMetadata returns metadata about a resource type
func GetResourceMetadata(resourceType string) ResourceMetadata {
	return ResourceMetadata{
		Type:        resourceType,
		Service:     GetServiceFromResourceType(resourceType),
		Category:    GetResourceCategoryFromResourceType(resourceType),
		IsSupported: IsKnownResource(resourceType),
		Description: getResourceDescription(resourceType),
	}
}

// getResourceDescription returns a human-readable description of a resource
func getResourceDescription(resourceType string) string {
	descriptions := map[string]string{
		"aws_s3_bucket":       "Amazon S3 Bucket",
		"aws_instance":        "EC2 Instance",
		"aws_security_group":  "VPC Security Group",
		"aws_vpc":             "Virtual Private Cloud",
		"aws_subnet":          "VPC Subnet",
		"aws_db_instance":     "RDS Database Instance",
		"aws_iam_role":        "IAM Role",
		"aws_iam_policy":      "IAM Policy",
		"aws_lambda_function": "Lambda Function",
		"aws_dynamodb_table":  "DynamoDB Table",
	}

	if desc, ok := descriptions[resourceType]; ok {
		return desc
	}
	return "AWS Resource"
}

// ServicePriorityMap defines the priority order for displaying services
// Higher priority services appear first in analysis output
var ServicePriorityMap = map[string]int{
	"iam":        1,
	"s3":         2,
	"ec2":        3,
	"rds":        4,
	"lambda":     5,
	"dynamodb":   6,
	"sns":        7,
	"sqs":        8,
	"kms":        9,
	"cloudwatch": 10,
}
