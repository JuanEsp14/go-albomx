package dto

type HelloWorldRequest struct {
	Name string `form:"name,default=World"` //form=queryParam
} //@Name HelloWorldRequest
