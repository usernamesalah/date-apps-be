package api

import (
	"date-apps-be/pkg/derrors"
	"date-apps-be/pkg/logger"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ResponseFormat struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination any         `json:"pagination,omitempty"`
}

// Response converts a Go value to JSON and sends it to the client.
func Response(c echo.Context, data interface{}, status string, message string, httpCode int) error {
	return c.JSON(httpCode, ResponseFormat{Status: status, Message: message, Data: data})
}

// ResponseWithPagination converts a Go value to JSON and sends it to the client.
func ResponseWithPagination(c echo.Context, data interface{}, pagination any, status string, message string, httpCode int) error {
	return c.JSON(httpCode, ResponseFormat{Status: status, Message: message, Data: data, Pagination: &pagination})
}

// ResponseOK converts a Go value to JSON and sends it to the client.
func ResponseOK(c echo.Context, data interface{}, HTTPStatus int) error {
	return Response(c, data, StatusCodeOK, StatusMessageOK, HTTPStatus)
}

func ResponseSuccess(c echo.Context, data interface{}, message string, HTTPStatus int) error {
	return Response(c, data, StatusSuccess, message, HTTPStatus)
}

func ResponseSuccessWithPagination(c echo.Context, data interface{}, pagination any, message string, HTTPStatus int) error {
	return ResponseWithPagination(c, data, pagination, StatusSuccess, message, HTTPStatus)
}

// ResponseError sends an error reponse back to the client.
func ResponseError(c echo.Context, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(*Error); ok {
		if err := Response(c, nil, webErr.Status, webErr.MessageStatus, webErr.HTTPStatus); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	if err := Response(c, nil, StatusCodeInternalServerError, StatusMessageInternalServerError, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}

// RenderErrorResponse sends an error reponse back to the client.
// Use this to get error response using derrors package
func RenderErrorResponse(c echo.Context, r *http.Request, err error) error {
	var ierr *derrors.Error
	if errors.As(err, &ierr) {
		status := derrors.ToStatus(ierr)
		msg := "Internal Server Error"
		if status < 500 {
			msg = ierr.Error()
		}

		reqID := c.Response().Header().Get(echo.HeaderXRequestID)

		logger.GetL().Error(
			"API error",
			zap.Int("status", status),
			zap.String("msg", msg),
			zap.String("request_id", reqID),
			zap.Error(err),
		)
		if errResp := Response(c, nil, StatusError, msg, status); errResp != nil {
			return errResp
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	logger.GetL().Error(
		"API error",
		zap.Int("status", http.StatusInternalServerError),
		zap.String("orig", err.Error()),
	)
	if err := Response(c, nil, StatusCodeInternalServerError, StatusMessageInternalServerError, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
