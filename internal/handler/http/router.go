package httpapi

import (
	"context"
	"net/http"

	_ "github.com/s-gurman/user-segmentation/docs"
	"github.com/s-gurman/user-segmentation/pkg/http/middleware"
	"github.com/s-gurman/user-segmentation/pkg/logger"

	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger/v2"
)

type SegmentationUseCase interface {
	CreateSegment(ctx context.Context, name string) (int, error)
	DeleteSegment(ctx context.Context, name string) error
	UpdateExperiments(ctx context.Context, userID int, segmentsToDel, segmentsToAdd []string) error
	GetUserExperiments(ctx context.Context, userID int) ([]string, error)
}

type muxRouter struct {
	*mux.Router
}

func NewRouter(uc SegmentationUseCase, l logger.Logger) http.Handler {
	segHandler := newSegmentHandler(uc, l)
	expHandler := newExperimentHandler(uc, l)

	router := muxRouter{mux.NewRouter()}.
		WithHandler(segHandler).
		WithHandler(expHandler).
		WithSwagger().
		WithMiddleware(l)

	return router
}

func (r muxRouter) WithSwagger() muxRouter {
	swaggerHandler := swagger.Handler(
		swagger.DeepLinking(true),
		swagger.DocExpansion("full"),
		swagger.DomID("swagger-ui"),
	)
	r.PathPrefix("/swagger/").Handler(swaggerHandler).Methods(http.MethodGet)
	return r
}

type routeHandler interface {
	addRoutes(r *mux.Router)
}

func (r muxRouter) WithHandler(h routeHandler) muxRouter {
	api := r.PathPrefix("/api").Subrouter()
	h.addRoutes(api)
	return r
}

func (r muxRouter) WithMiddleware(l logger.Logger) http.Handler {
	router := http.Handler(r.Router)
	router = middleware.AccessLog(router, l)
	router = middleware.PanicRecovery(router, l)
	return router
}
