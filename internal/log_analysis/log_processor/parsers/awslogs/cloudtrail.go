package awslogs

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
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/pantherlog"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
)

// CloudTrail is a record from the Records[*] JSON of an AWS CloudTrail API log.
// nolint:lll
type CloudTrail struct {
	AdditionalEventData pantherlog.RawMessage   `json:"additionalEventData" description:"Additional data about the event that was not part of the request or response."`
	APIVersion          pantherlog.String       `json:"apiVersion" description:"Identifies the API version associated with the AwsApiCall eventType value."`
	AWSRegion           pantherlog.String       `json:"awsRegion" validate:"required" description:"The AWS region that the request was made to, such as us-east-2."`
	ErrorCode           pantherlog.String       `json:"errorCode" description:"The AWS service error if the request returns an error."`
	ErrorMessage        pantherlog.String       `json:"errorMessage" description:"If the request returns an error, the description of the error. This message includes messages for authorization failures. CloudTrail captures the message logged by the service in its exception handling."`
	EventID             pantherlog.String       `json:"eventID" validate:"required" description:"GUID generated by CloudTrail to uniquely identify each event. You can use this value to identify a single event. For example, you can use the ID as a primary key to retrieve log data from a searchable database."`
	EventName           pantherlog.String       `json:"eventName" validate:"required" description:"The requested action, which is one of the actions in the API for that service."`
	EventSource         pantherlog.String       `json:"eventSource" validate:"required" description:"The service that the request was made to. This name is typically a short form of the service name without spaces plus .amazonaws.com."`
	EventTime           pantherlog.Time         `json:"eventTime" validate:"required" event_time:"true" tcodec:"rfc3339" description:"The date and time the request was made, in coordinated universal time (UTC)."`
	EventType           pantherlog.String       `json:"eventType" validate:"required" description:"Identifies the type of event that generated the event record. This can be the one of the following values: AwsApiCall, AwsServiceEvent, AwsConsoleSignIn"`
	EventVersion        pantherlog.String       `json:"eventVersion" validate:"required" description:"The version of the log event format."`
	ManagementEvent     pantherlog.Bool         `json:"managementEvent" description:"A Boolean value that identifies whether the event is a management event. managementEvent is shown in an event record if eventVersion is 1.06 or higher, and the event type is one of the following: AwsApiCall, AwsConsoleAction, AwsConsoleSignIn,  AwsServiceEvent"`
	ReadOnly            pantherlog.Bool         `json:"readOnly" description:"Identifies whether this operation is a read-only operation."`
	RecipientAccountID  pantherlog.String       `json:"recipientAccountId" panther:"aws_account_id" validate:"omitempty,len=12,numeric" description:"Represents the account ID that received this event. The recipientAccountID may be different from the CloudTrail userIdentity Element accountId. This can occur in cross-account resource access."`
	RequestID           pantherlog.String       `json:"requestID" description:"The value that identifies the request. The service being called generates this value."`
	RequestParameters   pantherlog.RawMessage   `json:"requestParameters" description:"The parameters, if any, that were sent with the request. These parameters are documented in the API reference documentation for the appropriate AWS service."`
	Resources           []CloudTrailResources   `json:"resources" description:"A list of resources accessed in the event."`
	ResponseElements    pantherlog.RawMessage   `json:"responseElements" description:"The response element for actions that make changes (create, update, or delete actions). If an action does not change state (for example, a request to get or list objects), this element is omitted. These actions are documented in the API reference documentation for the appropriate AWS service."`
	ServiceEventDetails pantherlog.RawMessage   `json:"serviceEventDetails" description:"Identifies the service event, including what triggered the event and the result."`
	SharedEventID       pantherlog.String       `json:"sharedEventID" description:"GUID generated by CloudTrail to uniquely identify CloudTrail events from the same AWS action that is sent to different AWS accounts."`
	SourceIPAddress     pantherlog.String       `json:"sourceIPAddress" panther:"ip" validate:"required" description:"The IP address that the request was made from. For actions that originate from the service console, the address reported is for the underlying customer resource, not the console web server. For services in AWS, only the DNS name is displayed."`
	UserAgent           pantherlog.String       `json:"userAgent" description:"The agent through which the request was made, such as the AWS Management Console, an AWS service, the AWS SDKs or the AWS CLI."`
	UserIdentity        *CloudTrailUserIdentity `json:"userIdentity" validate:"required" description:"Information about the user that made a request."`
	VPCEndpointID       pantherlog.String       `json:"vpcEndpointId" description:"Identifies the VPC endpoint in which requests were made from a VPC to another AWS service, such as Amazon S3."`
}

// CloudTrailResources are the AWS resources used in the API call.
type CloudTrailResources struct {
	ARN       pantherlog.String `json:"arn" panther:"aws_arn"`
	AccountID pantherlog.String `json:"accountId" panther:"aws_account_id"`
	Type      pantherlog.String `json:"type"`
}

// CloudTrailUserIdentity contains details about the type of IAM identity that made the request.
type CloudTrailUserIdentity struct {
	Type             pantherlog.String         `json:"type"`
	PrincipalID      pantherlog.String         `json:"principalId"`
	ARN              pantherlog.String         `json:"arn" panther:"aws_arn"`
	AccountID        pantherlog.String         `json:"accountId" panther:"aws_account_id"`
	AccessKeyID      pantherlog.String         `json:"accessKeyId"`
	Username         pantherlog.String         `json:"userName"`
	SessionContext   *CloudTrailSessionContext `json:"sessionContext"`
	InvokedBy        pantherlog.String         `json:"invokedBy"`
	IdentityProvider pantherlog.String         `json:"identityProvider"`
}

// CloudTrailSessionContext provides information about a session created for temporary credentials.
type CloudTrailSessionContext struct {
	Attributes          *CloudTrailSessionContextAttributes          `json:"attributes"`
	SessionIssuer       *CloudTrailSessionContextSessionIssuer       `json:"sessionIssuer"`
	WebIDFederationData *CloudTrailSessionContextWebIDFederationData `json:"webIdFederationData"`
}

// CloudTrailSessionContextAttributes  contains the attributes of the Session context object
type CloudTrailSessionContextAttributes struct {
	MfaAuthenticated pantherlog.String `json:"mfaAuthenticated"`
	CreationDate     pantherlog.String `json:"creationDate"`
}

// CloudTrailSessionContextSessionIssuer contains information for the SessionContextSessionIssuer
type CloudTrailSessionContextSessionIssuer struct {
	Type        pantherlog.String `json:"type"`
	PrincipalID pantherlog.String `json:"principalId"`
	Arn         pantherlog.String `json:"arn" panther:"aws_arn"`
	AccountID   pantherlog.String `json:"accountId" panther:"aws_account_id"`
	Username    pantherlog.String `json:"userName"`
}

// CloudTrailSessionContextWebIDFederationData contains Web ID federation data
type CloudTrailSessionContextWebIDFederationData struct {
	FederatedProvider pantherlog.String     `json:"federatedProvider"`
	Attributes        pantherlog.RawMessage `json:"attributes"`
}

// CloudTrailParser parses CloudTrail logs
type CloudTrailParser struct {
	builder pantherlog.ResultBuilder
}

var _ parsers.Interface = (*CloudTrailParser)(nil)

// Parse returns the parsed events or nil if parsing failed
func (p *CloudTrailParser) ParseLog(log string) (results []*parsers.Result, err error) {
	// Use strings.Reader to avoid duplicate allocation of `log` as bytes
	const bufferSize = 8192
	iter := jsoniter.Parse(jsoniter.ConfigDefault, strings.NewReader(log), bufferSize)
	// CloudTrail has all events in a single line inside an array at key `Records`
	// Seek to Records key
	const fieldNameRecords = `Records`
	for key := iter.ReadObject(); key != ""; key = iter.ReadObject() {
		if key != fieldNameRecords {
			iter.Skip()
			continue
		}
		// Pre-allocate some results to avoid multiple slice expansions
		const minResultSize = 1000
		results = make([]*parsers.Result, 0, minResultSize)
		// Go over all records parsing results
		for iter.ReadArray() {
			event := CloudTrail{}
			iter.ReadVal(&event)
			if err := iter.Error; err != nil {
				return nil, err
			}
			if err := pantherlog.ValidateStruct(&event); err != nil {
				return nil, err
			}
			result, err := p.builder.BuildResult(TypeCloudTrail, &event)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
		return results, nil
	}
	if err := iter.Error; err != nil {
		return nil, err
	}
	return nil, errors.New(`missing 'Records' field`)
}

// LogType returns the log type supported by this parser
func (p *CloudTrailParser) LogType() string {
	return TypeCloudTrail
}

var _ pantherlog.ValueWriterTo = (*CloudTrail)(nil)

func (event *CloudTrail) WriteValuesTo(w pantherlog.ValueWriter) {
	ExtractRawMessageIndicators(w, event.AdditionalEventData)
	ExtractRawMessageIndicators(w, event.RequestParameters)
	ExtractRawMessageIndicators(w, event.ResponseElements)
	ExtractRawMessageIndicators(w, event.ServiceEventDetails)
}

var _ pantherlog.ValueWriterTo = (*CloudTrailSessionContextWebIDFederationData)(nil)

func (d *CloudTrailSessionContextWebIDFederationData) WriteValuesTo(w pantherlog.ValueWriter) {
	if d != nil {
		ExtractRawMessageIndicators(w, d.Attributes)
	}
}
