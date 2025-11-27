package errors

import "google.golang.org/grpc/codes"

type Response struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (r *Response) Error() string {
	return r.Message
}
