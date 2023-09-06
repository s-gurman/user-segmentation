package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/logger"

	"github.com/gorilla/mux"
)

type segmentHandler struct {
	uc SegmentationUseCase
	l  logger.Logger
}

func newSegmentHandler(uc SegmentationUseCase, l logger.Logger) routeHandler {
	return segmentHandler{uc: uc, l: l}
}

func (h segmentHandler) addRoutes(r *mux.Router) {
	r.HandleFunc("/segment", h.createSegment).Methods(http.MethodPost)
	r.HandleFunc("/segment", h.deleteSegment).Methods(http.MethodDelete)
}

type (
	createSegmentRequest struct {
		SegmentName string `json:"name" example:"AVITO_VOICE_MESSAGES"`
	}
	deleteSegmentRequest struct {
		SegmentName string `json:"name" example:"AVITO_VOICE_MESSAGES"`
	}
)

// @Tags           segments
// @Summary        Creates segment
// @Description    Сreates a new segment with input name.
// @Router         /segment [post]
// @Accept         json
// @Produce        json
// @Param          body    body     createSegmentRequest true "Segment name"
// @Success        200     {object} successResponse{result=string}
// @Failure        400,500 {object} failedResponse
func (h segmentHandler) createSegment(w http.ResponseWriter, r *http.Request) {
	var req createSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := failedResponse{Msg: "invalid request, check swagger file", Code: 400, err: err}
		writeAndLogError(w, resp, h.l, "httpapi - decode request body" /*log*/)
		return
	}

	id, err := h.uc.CreateSegment(r.Context(), req.SegmentName)
	if err != nil {
		var (
			custom e.CustomError
			resp   = failedResponse{Msg: "internal error", Code: 500}
		)
		if errors.As(err, &custom) {
			resp = failedResponse{Msg: custom.Message(), Code: custom.Code()}
		}
		resp.err = err
		writeAndLogError(w, resp, h.l, "httpapi - create segment")
		return
	}

	msg := fmt.Sprintf("created segment '%s' with id=%d", req.SegmentName, id)
	resp := successResponse{Value: msg}
	writeAndLogValue(w, resp, h.l, msg)
}

// @Tags           segments
// @Summary        Deletes segment
// @Description    Deletes an existing segment by input name.
// @Router         /segment [delete]
// @Accept         json
// @Produce        json
// @Param          body        body     deleteSegmentRequest true "Segment name"
// @Success        200         {object} successResponse{result=string}
// @Failure        400,404,500 {object} failedResponse
func (h segmentHandler) deleteSegment(w http.ResponseWriter, r *http.Request) {
	var req deleteSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := failedResponse{Msg: "invalid request, check swagger file", Code: 400, err: err}
		writeAndLogError(w, resp, h.l, "httpapi - decode request body" /*log*/)
		return
	}

	err := h.uc.DeleteSegment(r.Context(), req.SegmentName)
	if err != nil {
		var (
			custom e.CustomError
			resp   = failedResponse{Msg: "internal error", Code: 500}
		)
		if errors.As(err, &custom) {
			resp = failedResponse{Msg: custom.Message(), Code: custom.Code()}
		}
		resp.err = err
		writeAndLogError(w, resp, h.l, "httpapi - delete segment")
		return
	}

	msg := fmt.Sprintf("segment '%s' deleted", req.SegmentName)
	resp := successResponse{Value: msg}
	writeAndLogValue(w, resp, h.l, msg)
}
