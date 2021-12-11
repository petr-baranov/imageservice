package services

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/draw"
)

type ImageService interface {
	Encode(w io.Writer, r io.Reader) error
	Scale(w io.Writer, r io.Reader, scale string) error
	Scales() []string
}

type imageServiceImpl struct {
	config []ScaleConfig
}

func NewImageService(config []ScaleConfig) ImageService {
	return &imageServiceImpl{config: config}
}

func encode(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg":
		if err := jpeg.Encode(w, img, nil); err != nil {
			return err
		}
	case "png":
		if err := png.Encode(w, img); err != nil {
			return err
		}
	default:
		return errors.New("unsupported image format")
	}
	return nil
}

func (o *imageServiceImpl) Encode(w io.Writer, r io.Reader) error {
	img, format, err := image.Decode(r)
	if err != nil {
		return err
	}
	if err := encode(w, img, format); err != nil {
		return err
	}
	return nil
}
func (o *imageServiceImpl) getScaleConfig(scale string) *ScaleConfig {
	for i := range o.config {
		if strings.ToLower(scale) == strings.ToLower(o.config[i].Name) {
			return &o.config[i]
		}
	}
	return nil
}
func (o *imageServiceImpl) Scales() []string {
	r := []string{}
	for i := range o.config {
		r = append(r, strings.ToLower(o.config[i].Name))
	}
	return r
}
func (o *imageServiceImpl) Scale(w io.Writer, r io.Reader, scale string) error {
	scaleConfig := o.getScaleConfig(scale)
	if scaleConfig == nil {
		return errors.New("not supported scale")
	}
	img, format, err := image.Decode(r)
	if err != nil {
		return err
	}
	kernel := draw.BiLinear
	rect := image.Rect(0, 0, int(float32(img.Bounds().Max.X)*scaleConfig.Factor), int(float32(img.Bounds().Max.Y)*scaleConfig.Factor))
	dstImg := image.NewRGBA(rect)
	kernel.Scale(dstImg, rect, img, img.Bounds(), draw.Over, nil)
	if err := encode(w, dstImg, format); err != nil {
		return err
	}
	return nil
}

type ScaleConfig struct {
	Name   string
	Factor float32
}
