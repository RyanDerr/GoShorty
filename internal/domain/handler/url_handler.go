package handler

import (
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/RyanDerr/GoShorty/pkg/mapper"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	urlService service.IUrlService
}

func NewUrlHandler(urlService service.IUrlService) *UrlHandler {
	return &UrlHandler{
		urlService: urlService,
	}
}

// ShortenUrl godoc
//
//	@Summary		Shorten a URL
//	@Description	Shorten a given URL and return the shortened version
//	@Tags			URL
//	@Accept			json
//	@Produce		json
//	@Param			shortenUrl	body		request.ShortenUrlRequest	true	"URL to be shortened"
//	@Success		201			{object}	response.ShortenUrlResponse	"Shortened URL"
//	@Failure		400			{object}	response.ResponseErrorModel	"Bad Request"
//	@Failure		409			{object}	response.ResponseErrorModel	"Conflict"
//	@Failure		500			{object}	response.ResponseErrorModel	"Internal Server Error"
//	@Router			/url/shorten [post]
func (h *UrlHandler) ShortenUrl(ctx *gin.Context) {
	short := new(request.ShortenUrlRequest)
	err := ctx.BindJSON(short)
	if err != nil {
		log.Printf("Error binding JSON: %s", err.Error())
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	err = short.Validate()
	if err != nil {
		log.Printf("Error validating request: %s", err.Error())
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	shortEntity, err := mapper.MapShortenUrlRequestToEntity(short)
	if err != nil {
		log.Printf("Error mapping request to entity: %s", err.Error())
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	res, status, err := h.urlService.ShortenUrl(ctx, shortEntity)
	if err != nil {
		log.Printf("Error shortening URL: %s", err.Error())
		response.ResponseError(ctx, err.Error(), status)
		return
	}

	log.Printf("Shortened URL: %s", res.BaseUrl)
	response.ResponseCreatedWithData(ctx, mapper.MapShortenUrlEntityToResponse(res))
}

// ResolveUrl godoc
//
//	@Summary		Resolve a shortened URL
//	@Description	Resolve a shortened URL to its original URL
//	@Tags			URL
//	@Accept			json
//	@Produce		json
//	@Param			short	path		string						true	"Short to resolve"
//	@Success		301		{string}	string						"Redirect to original URL"
//	@Failure		400		{object}	response.ResponseErrorModel	"Bad Request"
//	@Failure		404		{object}	response.ResponseErrorModel	"Not Found"
//	@Failure		500		{object}	response.ResponseErrorModel	"Internal Server Error"
//	@Router			/url/{short} [get]
func (h *UrlHandler) ResolveUrl(ctx *gin.Context) {
	shortUrl := ctx.Param("short")
	log.Printf("Resolving short URL: %s", shortUrl)
	res, status, err := h.urlService.ResolveUrl(ctx, shortUrl)
	if err != nil {
		log.Printf("Error resolving URL: %s", err.Error())
		response.ResponseError(ctx, err.Error(), status)
		return
	}

	log.Printf("Redirecting to: %s", res)
	response.ResponseRedirect(ctx, res)
}
