package main

import (
	"image"
	"net/http"
	"strconv"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/effect"
)

type Params struct {
	brightness float64
	contrast   float64
	invert     float64
}

var defaultParams = Params{
	brightness: 0,
	contrast:   0,
	invert:     0,
}

func transform(image image.Image, query *http.Request) image.Image {
	final := image

	brightness := parseQueryParamNum(query, "brightness", defaultParams.brightness)
	contrast := parseQueryParamNum(query, "contrast", defaultParams.contrast)
	invert := parseQueryParamNum(query, "invert", defaultParams.invert)

	params := Params{
		brightness: brightness,
		contrast:   contrast,
		invert:     invert,
	}

	if params.invert != defaultParams.invert {
		final = effect.Invert(final)
	}

	if params.brightness != defaultParams.brightness {
		final = adjust.Brightness(final, params.brightness)
	}

	if params.contrast != defaultParams.contrast {
		final = adjust.Contrast(final, params.contrast)
	}

	return final
}

func parseQueryParamNum(query *http.Request, key string, defaultValue float64) float64 {
	queryParams := query.URL.Query()

	val := queryParams[key]

	if len(val) == 0 {
		return defaultValue
	}

	value, err := strconv.ParseFloat(val[0], 64)

	if err != nil {
		value = defaultValue
	}

	return value
}
