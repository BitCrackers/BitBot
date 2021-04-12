package responses

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type CustomResponseHandler struct {
	Responses []Response
}

type Response struct {
	Name string
	Send func(s *discordgo.Session, m *discordgo.Message, reply bool) error
}

func New() CustomResponseHandler {
	return CustomResponseHandler{Responses: []Response{
		BuildResponse,
	}}
}

func (c *CustomResponseHandler) GetCustomResponse(name string) (Response, error) {
	for _, response := range c.Responses {
		if response.Name == name {
			return response, nil
		}
	}
	return Response{}, fmt.Errorf("no custom response found with name %s", name)
}
