package httplib

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

//HTTPServer HTTPServer
type HTTPServer struct {
	AllowOrigin []string
	handles     []*Handle

	server *http.Server

	controlContext context.Context
	controlCancel  context.CancelFunc

	authHandle AuthHandle
}

//AddCORS AddCORS
func (h *HTTPServer) AddCORS(urls ...string) {
	for _, url := range urls {
		h.AllowOrigin = append(h.AllowOrigin, url)
	}
}

//StartServer StartServer
func (h *HTTPServer) StartServer(port int) error {
	var handlesNoStart []*Handle
	var handlesHasStart []*Handle
	for _, item := range h.handles {
		if item.MatchStart {
			handlesHasStart = append(handlesHasStart, item)
		} else {
			handlesNoStart = append(handlesNoStart, item)
		}
	}

	sort.Slice(handlesNoStart, func(i, j int) bool {
		return len(handlesNoStart[i].Path) > len(handlesNoStart[j].Path) // 按字符串长度 降序 排列
	})

	sort.Slice(handlesHasStart, func(i, j int) bool {
		return len(handlesHasStart[i].Path) > len(handlesHasStart[j].Path) // 按字符串长度 降序 排列
	})

	h.handles = append(handlesNoStart, handlesHasStart...)

	h.controlContext, h.controlCancel = context.WithCancel(context.Background())
	h.server = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: h}

	return h.server.ListenAndServe()
}

func (h *HTTPServer) SetAuthHandle(authHandle AuthHandle) {
	h.authHandle = authHandle
}

//StartServer StartServer
func (h *HTTPServer) StartServerTLS(port int, certFile string, keyFile string) error {
	h.controlContext, h.controlCancel = context.WithCancel(context.Background())
	h.server = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: h}

	return h.server.ListenAndServeTLS(certFile, keyFile)
}

func (h *HTTPServer) StopServer() {
	h.server.Shutdown(h.controlContext)
}

func (h *HTTPServer) checkCORS(w http.ResponseWriter, r *http.Request) bool {
	isAccess := false
	origin := r.Header.Get("Origin")
	origin = strings.Trim(origin, "http://")
	if origin != "" && origin != r.Host {
		for _, item := range h.AllowOrigin {
			if item == "*" || item == origin {
				w.Header().Add("Access-Control-Allow-Origin", "http://"+origin)
				w.Header().Add("Access-Control-Allow-Credentials", "true")
				if r.Method == "OPTIONS" {
					w.Header().Add("Access-Control-Allow-Headers", "Content-Type, ws, accept, Authorization")
					w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, UPDATE")
					isAccess = false
				} else {
					isAccess = true
				}

				break
			}
		}
	} else {
		return true
	}
	return isAccess
}

//ServeHTTP ServeHTTP
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.checkCORS(w, r) {
		urlPath := strings.ToLower(r.URL.Path)
		var handle *Handle

		for _, item := range h.handles {
			if item.Path == urlPath || (strings.Index(urlPath, item.Path) == 0 && item.MatchStart) {
				handle = item
				break
			}
		}
		if handle != nil {
			var authInfo interface{}
			var err error
			if handle.NeedLogin {
				if h.authHandle != nil {
					authInfo, err = h.authHandle(r)
				} else {
					err = errors.New("")
				}
			}
			if err == nil {

				handle.handle(w, r, authInfo)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		} else {
			http.NotFound(w, r)
		}
	}
}

func (h *HTTPServer) registerHandle(url string, callBackHandle CallBackHandle, needLogin bool) {
	if url == "" {
		return
	}
	handle := new(Handle)
	handle.NeedLogin = needLogin
	handle.CallBackHandle = callBackHandle
	if url[len(url)-1] == '*' {
		url = url[:len(url)-1]
		handle.MatchStart = true
	}
	handle.Path = strings.ToLower(url)
	h.handles = append(h.handles, handle)
}

//RegisterHandle RegisterHandle
func (h *HTTPServer) RegisterHandle(url string, callBackHandle CallBackHandle) {
	h.registerHandle(url, callBackHandle, false)
}

//RegisterHandleAuth RegisterHandleAuth
func (h *HTTPServer) RegisterHandleAuth(url string, callBackHandle CallBackHandle) {
	h.registerHandle(url, callBackHandle, true)
}

//RegisterWebSocket RegisterWebSocket
func (h *HTTPServer) RegisterWebSocket(url string, handle websocket.Handler) {
	handleVar := new(Handle)
	handleVar.NeedLogin = false
	handleVar.Handler = handle
	handleVar.Path = strings.ToLower(url)
	h.handles = append(h.handles, handleVar)
}

func (h *HTTPServer) registerVirtualPath(virtualPath, localPath string, needLogin bool) {
	handle := new(Handle)
	if virtualPath[len(virtualPath)-1] == '*' {
		virtualPath = virtualPath[:len(virtualPath)-1]
		handle.MatchStart = true
	}

	fileHandle := http.StripPrefix(virtualPath, http.FileServer(http.Dir(localPath)))
	handle.NeedLogin = needLogin
	handle.Handler = fileHandle
	handle.Path = strings.ToLower(virtualPath)
	h.handles = append(h.handles, handle)
}

//RegisterVirtualPath RegisterVirtualPath
func (h *HTTPServer) RegisterVirtualPath(virtualPath, localPath string) {
	h.registerVirtualPath(virtualPath, localPath, false)
}

//RegisterVirtualPath RegisterVirtualPath
func (h *HTTPServer) RegisterVirtualPathAuth(virtualPath, localPath string) {
	h.registerVirtualPath(virtualPath, localPath, true)
}

//RegisterVirtualPath RegisterVirtualPath
func (h *HTTPServer) RegisterFS(virtualPath string, embedPath string, f fs.FS) {
	handleVar := new(Handle)
	if virtualPath[len(virtualPath)-1] == '*' {
		virtualPath = virtualPath[:len(virtualPath)-1]
		handleVar.MatchStart = true
	}

	handle := http.FileServer(http.FS(f))
	fileHandle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, virtualPath)
		rp := strings.TrimPrefix(r.URL.RawPath, virtualPath)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = embedPath + p
			r2.URL.RawPath = rp
			handle.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})

	handleVar.NeedLogin = false
	handleVar.Handler = fileHandle
	handleVar.Path = strings.ToLower(virtualPath)
	h.handles = append(h.handles, handleVar)
}
