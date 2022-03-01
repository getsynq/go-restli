package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
)

var ErrorResponseIdentifier = utils.Identifier{
	Name:      "ErrorResponse",
	Namespace: utils.StdTypesPackage,
}

// ErrorResponse is manually parsed from https://github.com/linkedin/rest.li/blob/master/restli-common/src/main/pegasus/com/linkedin/restli/common/ErrorResponse.pdl
var ErrorResponse = &types.Record{
	NamedType: types.NamedType{Identifier: ErrorResponseIdentifier},
	Fields: []types.Field{
		{
			Type:       types.RestliType{Primitive: &types.Int32Primitive},
			Name:       "status",
			Doc:        "The HTTP status code.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "message",
			Doc:        "A human-readable explanation of the error.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "exceptionClass",
			Doc:        "The FQCN of the exception thrown by the server.",
			IsOptional: true,
		},
		{
			Type:       types.RestliType{Primitive: &types.StringPrimitive},
			Name:       "stackTrace",
			Doc:        "The full stack trace of the exception thrown by the server.",
			IsOptional: true,
		},
	},
}

func init() {
	utils.TypeRegistry.RegisterExternalType(
		ErrorResponseIdentifier,
		types.RecordShouldUsePointer,
		utils.StdTypesPackage,
		ErrorResponseIdentifier.Name,
	)
}
