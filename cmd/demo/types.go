package main

import (
	"image"
	"io"

	"github.com/Zensey/go-archetype-project/pkg/logger"
)

type TRequest struct {
	Urls []string `json:"imgsUrls,omitempty"`
	Imgs []string `json:"imgs,omitempty"`
}

type TResponse struct {
	thumbnails *[]image.Image
	Thumbs     []string `json:"thumbs,omitempty"`
}

func NewTResponse() TResponse {
	return TResponse{thumbnails: &[]image.Image{}}
}

type decodeAndAddThumb func(rr io.Reader, l logger.Logger)
