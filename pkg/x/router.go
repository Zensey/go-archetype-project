package x

import (
	"github.com/julienschmidt/httprouter"

	"github.com/ory/x/serverx"
)

type RouterPublic struct {
	*httprouter.Router
}

func NewRouterPublic() *RouterPublic {
	router := httprouter.New()
	router.NotFound = serverx.DefaultNotFoundHandler
	return &RouterPublic{
		Router: router,
	}
}
