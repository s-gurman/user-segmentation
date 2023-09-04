package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/logger"

	"github.com/gorilla/mux"
)

type experimentHandler struct {
	uc SegmentationUseCase
	l  logger.Logger
}

func newExperimentHandler(uc SegmentationUseCase, l logger.Logger) experimentHandler {
	return experimentHandler{uc: uc, l: l}
}

func (h experimentHandler) addRoutes(r *mux.Router) {
	r.HandleFunc("/experiments/user/{user_id:[0-9]+}", h.updateExperiments).Methods(http.MethodPost)
	r.HandleFunc("/experiments/user/{user_id:[0-9]+}", h.getExperiments).Methods(http.MethodGet)
}

type updateExperimentsRequest struct {
	SegmentsToDel []string `json:"delete" example:"AVITO_PERFORMANCE_VAS,AVITO_DISCOUNT_30"`
	SegmentsToAdd []string `json:"add"    example:"AVITO_VOICE_MESSAGES,AVITO_DISCOUNT_50"`
}

// @Tags           experiments
// @Summary        Updates user experiments
// @Description    Deletes user's active segments and adds new segments from existing ones.
// @Router         /experiments/user/{user_id} [post]
// @Accept         json
// @Produce        json
// @Param          user_id     path     int true "User ID"
// @Param          body        body     updateExperimentsRequest true "Lists of deleting and adding active user segments"
// @Success        200         {object} successResponse{result=string}
// @Failure        400,404,500 {object} failedResponse
func (h experimentHandler) updateExperiments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		resp := failedResponse{Msg: "internal error", Code: 500, err: err}
		writeAndLogError(w, resp, h.l, "httpapi - strconv id err" /*log*/)
		return
	}

	var req updateExperimentsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := failedResponse{Msg: "invalid json in request body", Code: 400, err: err}
		writeAndLogError(w, resp, h.l, "httpapi - decode json err" /*log*/)
		return
	}

	err = h.uc.UpdateExperiments(r.Context(), userID, req.SegmentsToDel, req.SegmentsToAdd)
	if err != nil {
		var (
			custom e.CustomError
			resp   = failedResponse{Msg: "internal error", Code: 500}
		)
		if errors.As(err, &custom) {
			resp = failedResponse{Msg: custom.Message(), Code: custom.Code()}
		}
		resp.err = err
		writeAndLogError(w, resp, h.l, "httpapi - update experiments" /*log*/)
		return
	}

	msg := fmt.Sprintf(
		"added %d segments to user '%d' and deleted %d active ones",
		len(req.SegmentsToAdd), userID, len(req.SegmentsToDel),
	)
	resp := successResponse{Value: msg}
	writeAndLogValue(w, resp, h.l, msg)
}

// @Tags           experiments
// @Summary        Gets user experiments
// @Description    Gets the user's active segments.
// @Router         /experiments/user/{user_id} [get]
// @Accept         json
// @Produce        json
// @Param          user_id path     int true "User ID"
// @Success        200     {object} successResponse{result=string}
// @Failure        400,500 {object} failedResponse
func (h experimentHandler) getExperiments(w http.ResponseWriter, r *http.Request) {
	// id, err := h.uc.CreateSegment(r.Context(), req.SegmentName)
	// if err != nil {
	// 	var (
	// 		custom e.CustomError
	// 		resp = failedResponse{Msg: "internal error", Code: 500, err: err}
	// 	)
	// 	if errors.As(err, &custom) {
	// 		resp = failedResponse{Msg: custom.Message(), Code: custom.Code(), err: err}
	// 	}
	// 	writeAndLogError(w, resp, h.l, "httpapi - create segment")
	// 	return
	// }

	// msg := fmt.Sprintf("created segment '%s' with id=%d", req.SegmentName, id)
	// resp := successResponse{Value: msg}
	// writeAndLogValue(w, resp, h.l, msg)
}
