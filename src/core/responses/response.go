package responses

import "github.com/gin-gonic/gin"

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func OfSucceed(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, &response{
		Data: data,
	})
}

func Ok(ctx *gin.Context, statusCode int) {
	ctx.JSON(statusCode, &response{
		Data: "ok",
	})
}

func OfError(ctx *gin.Context, statusCode int, errorCode string) {
	ctx.JSON(statusCode, &response{
		Error: errorCode,
	})
}
