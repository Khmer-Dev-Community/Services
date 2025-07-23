package utils

type Message string

const (
	SuccessMessage  Message = "Operation successful"
	WarningMessage  Message = "Warning: Operation may have unexpected results"
	ErrorMessage    Message = "An error occurred"
	NotFoundMessage Message = "record not found"
	ExecuteMessage  Message = "Execute"
	RequestMessage  Message = "Request"
	ResponseMessage Message = "Response"
)
