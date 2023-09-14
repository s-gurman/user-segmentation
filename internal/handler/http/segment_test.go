package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/s-gurman/user-segmentation/internal/handler/http/mocks"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestDeleteSegment_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := "some slug"
	reqBodyStr := fmt.Sprintf(`{"name": %q}`, name)
	reqBody := strings.NewReader(reqBodyStr)

	req := httptest.NewRequest("DELETE", "/api/segment", reqBody)
	w := httptest.NewRecorder()

	uc := mocks.NewMockSegmentationUseCase(ctrl)
	uc.
		EXPECT().
		DeleteSegment(req.Context(), name).
		Return(nil)

	h := segmentHandler{
		uc: uc,
		l:  zap.NewNop().Sugar(),
	}

	h.deleteSegment(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("body read err: %s", err)
	}

	expectedCode := http.StatusOK
	expectedResult := successResponse{
		Value: fmt.Sprintf(`segment '%s' deleted`, name),
	}
	var result successResponse

	if resp.StatusCode != expectedCode {
		t.Errorf("unexpected status code:\n\twant: %#v\n\thave: %#v", expectedCode, resp.StatusCode)
	}
	if err = json.Unmarshal(body, &result); err != nil {
		t.Errorf("unexpected response body: %s", body)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("results don't match:\n\twant: %#v\n\thave: %#v", expectedResult, result)
	}
}
