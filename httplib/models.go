package httplib

import "net/http"

//CallBackHandle CallBackHandle
type CallBackHandle func(response *Response, request *Request)

//AuthHandle AuthHandle
type AuthHandle func(request *http.Request) (interface{}, error)

type Handle struct {
	NeedLogin      bool
	AuthInfo       interface{}
	Handler        http.Handler
	CallBackHandle CallBackHandle
	Path           string
	MatchStart     bool
}

func (h *Handle) handle(w http.ResponseWriter, r *http.Request, authInfo interface{}) {
	if h.Handler != nil {
		h.Handler.ServeHTTP(w, r)
	} else if h.CallBackHandle != nil {
		response := new(Response)
		response.ResponseWriter = w

		request := new(Request)
		request.Request = r
		request.AuthInfo = authInfo

		h.CallBackHandle(response, request)
	}
}
