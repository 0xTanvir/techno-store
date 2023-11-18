package dto

type Error struct {
	Message string `json:"message"`
}

func Builder() Error {
	return Error{}
}

func (error Error) SetMessage(message string) Error {
	error.Message = message
	return error
}
