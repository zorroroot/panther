// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

// DeleteDataModelReader is a Reader for the DeleteDataModel structure.
type DeleteDataModelReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDataModelReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteDataModelOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteDataModelBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteDataModelInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDeleteDataModelOK creates a DeleteDataModelOK with default headers values
func NewDeleteDataModelOK() *DeleteDataModelOK {
	return &DeleteDataModelOK{}
}

/*DeleteDataModelOK handles this case with default header values.

OK
*/
type DeleteDataModelOK struct {
}

func (o *DeleteDataModelOK) Error() string {
	return fmt.Sprintf("[POST /datamodel/delete][%d] deleteDataModelOK ", 200)
}

func (o *DeleteDataModelOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDataModelBadRequest creates a DeleteDataModelBadRequest with default headers values
func NewDeleteDataModelBadRequest() *DeleteDataModelBadRequest {
	return &DeleteDataModelBadRequest{}
}

/*DeleteDataModelBadRequest handles this case with default header values.

Bad request
*/
type DeleteDataModelBadRequest struct {
	Payload *models.Error
}

func (o *DeleteDataModelBadRequest) Error() string {
	return fmt.Sprintf("[POST /datamodel/delete][%d] deleteDataModelBadRequest  %+v", 400, o.Payload)
}

func (o *DeleteDataModelBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteDataModelBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteDataModelInternalServerError creates a DeleteDataModelInternalServerError with default headers values
func NewDeleteDataModelInternalServerError() *DeleteDataModelInternalServerError {
	return &DeleteDataModelInternalServerError{}
}

/*DeleteDataModelInternalServerError handles this case with default header values.

Internal server error
*/
type DeleteDataModelInternalServerError struct {
}

func (o *DeleteDataModelInternalServerError) Error() string {
	return fmt.Sprintf("[POST /datamodel/delete][%d] deleteDataModelInternalServerError ", 500)
}

func (o *DeleteDataModelInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
