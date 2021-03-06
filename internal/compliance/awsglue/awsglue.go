package awsglue

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/glue/glueiface"
	"github.com/pkg/errors"

	"github.com/panther-labs/panther/pkg/awsutils"
)

const (
	CloudSecurityDatabase            = "panther_cloudsecurity"
	CloudSecurityDatabaseDescription = "Hold tables related to Panther cloud security scanning"

	// https://github.com/awslabs/aws-athena-query-federation/tree/master/athena-dynamodb

	// FIXME: Update the description when the DDB connector is GA
	ResourcesTableDDB         = "panther-resources"
	ResourcesTable            = "resources"
	ResourcesTableDescription = "(ddb.panther_cloudsecurity.resources) The resources discovered by Panther scanning"

	ComplianceTableDDB         = "panther-compliance"
	ComplianceTable            = "compliance"
	ComplianceTableDescription = "(ddb.panther_cloudsecurity.compliance) The policies and statuses from Panther scanning"
)

var (
	// FIXME: Remove when the DDB connector is GA
	// Available Regions – The Athena federated query feature is available in preview in the US East (N. Virginia),
	//                     Asia Pacific (Mumbai), Europe (Ireland), and US West (Oregon) Regions.
	anthenaDDBConnectorRegions = map[string]struct{}{
		"us-east-1":  {},
		"ap-south-1": {},
		"eu-west-1":  {},
		"us-west-2":  {},
	}
)

func CreateOrUpdateCloudSecurityDatabase(glueClient glueiface.GlueAPI) error {
	dbInput := &glue.DatabaseInput{
		Description: aws.String(CloudSecurityDatabaseDescription),
		LocationUri: aws.String("dynamo-db-flag"),
		Name:        aws.String(CloudSecurityDatabase),
	}

	_, err := glueClient.CreateDatabase(&glue.CreateDatabaseInput{
		CatalogId:     nil,
		DatabaseInput: dbInput,
	})
	if awsutils.IsAnyError(err, glue.ErrCodeAlreadyExistsException) {
		return nil // nothing to do
	}
	return errors.Wrap(err, "could not create cloud security database")
}

func CreateOrUpdateResourcesTable(glueClient glueiface.GlueAPI, locationARN string) error {
	// FIXME: Remove when the DDB connector is GA
	parsedARN, err := arn.Parse(locationARN)
	if err != nil {
		return err
	}
	if _, found := anthenaDDBConnectorRegions[parsedARN.Region]; !found {
		return nil // not supported
	}

	tableInput := &glue.TableInput{
		Name:        aws.String(ResourcesTable),
		Description: aws.String(ResourcesTableDescription),
		Parameters: map[string]*string{
			// per https://github.com/awslabs/aws-athena-query-federation/tree/master/athena-dynamodb
			"classification": aws.String("dynamodb"),
			"sourceTable":    aws.String(ResourcesTableDDB),
			// for attrs with upper case
			// nolint:lll
			"columnMapping": aws.String(`expiresat=expiresAt,lastmodified=lastModified,integrationid=integrationId,integrationtype=integrationType`),
		},
		StorageDescriptor: &glue.StorageDescriptor{
			Location: &locationARN,

			Columns: []*glue.Column{
				/* Commenting out for now: always 'aws'
				{
					Name:    aws.String("integrationtype"),
					Type:    aws.String("string"),
					Comment: aws.String("Indicates what type of integration this resource came from"),
				},
				*/
				{
					Name:    aws.String("deleted"),
					Type:    aws.String("boolean"),
					Comment: aws.String("True if this is the snapshot of a deleted resource."),
				},
				{
					Name:    aws.String("integrationid"),
					Type:    aws.String("string"),
					Comment: aws.String("The unique ID indicating of the source integration."),
				},
				{
					Name:    aws.String("attributes"),
					Type:    aws.String("string"),
					Comment: aws.String("The JSON representation of the resource."),
				},
				/* Commenting out: not useful
				{
					Name:    aws.String("lowerid"),
					Type:    aws.String("string"),
					Comment: aws.String("The resource ID converted to all lower case letters."),
				},
				*/
				{
					Name:    aws.String("lastmodified"),
					Type:    aws.String("string"),
					Comment: aws.String("Timestamp of the most recent scan of this resource occurred."),
				},
				{
					Name:    aws.String("id"),
					Type:    aws.String("string"),
					Comment: aws.String("The panther wide unique identifier of the resource."),
				},
				{
					Name:    aws.String("type"),
					Type:    aws.String("string"),
					Comment: aws.String(" The type of resource (see https://docs.runpanther.io/cloud-security/resources)."),
				},
				/* Commenting out: not useful
				{
					Name:    aws.String("expiresat"),
					Type:    aws.String("bigint"),
					Comment: aws.String("Unix timestamp representing when this resource will age out of the resources table."),
				},
				*/
			},
		},
		TableType: aws.String("EXTERNAL_TABLE"),
	}

	createTableInput := &glue.CreateTableInput{
		DatabaseName: aws.String(CloudSecurityDatabase),
		TableInput:   tableInput,
	}

	_, err = glueClient.CreateTable(createTableInput)
	if err != nil {
		if awsutils.IsAnyError(err, glue.ErrCodeAlreadyExistsException) {
			// need to do an update
			updateTableInput := &glue.UpdateTableInput{
				DatabaseName: aws.String(CloudSecurityDatabase),
				TableInput:   tableInput,
			}
			_, err := glueClient.UpdateTable(updateTableInput)
			return errors.Wrapf(err, "failed to update table %s.%s", CloudSecurityDatabase, ResourcesTable)
		}
		return errors.Wrapf(err, "failed to create table %s.%s", CloudSecurityDatabase, ResourcesTable)
	}

	return nil
}

func CreateOrUpdateComplianceTable(glueClient glueiface.GlueAPI, locationARN string) error {
	// FIXME: Remove when the DDB connector is GA
	parsedARN, err := arn.Parse(locationARN)
	if err != nil {
		return err
	}
	if _, found := anthenaDDBConnectorRegions[parsedARN.Region]; !found {
		return nil // not supported
	}

	tableInput := &glue.TableInput{
		Name:        aws.String(ComplianceTable),
		Description: aws.String(ComplianceTableDescription),
		Parameters: map[string]*string{
			// per https://github.com/awslabs/aws-athena-query-federation/tree/master/athena-dynamodb
			"classification": aws.String("dynamodb"),
			"sourceTable":    aws.String(ComplianceTableDDB),
			// for attrs with upper case
			// nolint:lll
			"columnMapping": aws.String(`policyseverity=policySeverity,errormessage=errorMessage,expiresat=expiresAt,lastupdated=lastUpdated,policyid=policyId,resourceid=resourceId,resourcetype=resourceType,integrationid=integrationId`),
		},
		StorageDescriptor: &glue.StorageDescriptor{
			Location: &locationARN,

			Columns: []*glue.Column{
				{
					Name:    aws.String("lastupdated"),
					Type:    aws.String("string"),
					Comment: aws.String("That last date the specified policy was evaluated against the specified resource."),
				},
				{
					Name:    aws.String("resourceid"),
					Type:    aws.String("string"),
					Comment: aws.String("The panther wide unique identifier of the resource being evaluated."),
				},
				{
					Name:    aws.String("policyseverity"),
					Type:    aws.String("string"),
					Comment: aws.String("The severity of the policy being evaluated."),
				},
				{
					Name:    aws.String("policyid"),
					Type:    aws.String("string"),
					Comment: aws.String("The unique identifier of the policy being evaluated."),
				},
				{
					Name:    aws.String("integrationid"),
					Type:    aws.String("string"),
					Comment: aws.String("The unique ID indicating of the source integration."),
				},
				{
					Name:    aws.String("suppressed"),
					Type:    aws.String("boolean"),
					Comment: aws.String("True if this compliance status is currently being omitted from compliance findings."),
				},
				/* Commenting out: not useful
				{
					Name:    aws.String("expiresat"),
					Type:    aws.String("bigint"),
					Comment: aws.String("Unix timestamp representing when this resource will age out of the resources table."),
				},
				*/
				{
					Name:    aws.String("resourcetype"),
					Type:    aws.String("string"),
					Comment: aws.String("The type of the specified resource."),
				},
				{
					Name:    aws.String("status"),
					Type:    aws.String("string"),
					Comment: aws.String("Whether the policy evaluation of this resource resulted in a PASS, FAIL, or ERROR state."),
				},
				{
					Name:    aws.String("errormessage"),
					Type:    aws.String("string"),
					Comment: aws.String("If an error occurred, the associated error message."),
				},
			},
		},
		TableType: aws.String("EXTERNAL_TABLE"),
	}

	createTableInput := &glue.CreateTableInput{
		DatabaseName: aws.String(CloudSecurityDatabase),
		TableInput:   tableInput,
	}

	_, err = glueClient.CreateTable(createTableInput)
	if err != nil {
		if awsutils.IsAnyError(err, glue.ErrCodeAlreadyExistsException) {
			// need to do an update
			updateTableInput := &glue.UpdateTableInput{
				DatabaseName: aws.String(CloudSecurityDatabase),
				TableInput:   tableInput,
			}
			_, err := glueClient.UpdateTable(updateTableInput)
			return errors.Wrapf(err, "failed to update table %s.%s", CloudSecurityDatabase, ResourcesTable)
		}
		return errors.Wrapf(err, "failed to create table %s.%s", CloudSecurityDatabase, ResourcesTable)
	}

	return nil
}
