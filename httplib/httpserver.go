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
	handles     map[string]*Handle
	handlesList []string

	handlesFile     map[string]*Handle
	handlesListFile []string

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
		if strings.LastIndex(urlPath, "/") == len(urlPath)-1 {
			urlPath = strings.TrimRight(urlPath, "/")
		}
		handle, ok := h.handles[urlPath]
		if !ok {
			handle, ok = h.handlesFile[urlPath]
			if !ok {
				for _, value := range h.handlesListFile {
					if strings.HasPrefix(urlPath, value) {
						handle = h.handlesFile[value]
						break
					}
				}
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

//sort sort
func (h *HTTPServer) sort(key string, handlesList []string) []string {

	n := len(handlesList)
	i := sort.Search(n, func(i int) bool {
		return len(handlesList[i]) < len(key)
	})
	if i == n {
		handlesList = append(handlesList, key)
		return handlesList
	}
	var newHandleList []string
	newHandleList = append(newHandleList, handlesList[0:i]...)
	newHandleList = append(newHandleList, key)
	newHandleList = append(newHandleList, handlesList[i:]...)
	return newHandleList
}

func (h *HTTPServer) registerHandle(url string, callBackHandle CallBackHandle, needLogin bool) {
	if h.handles == nil {
		h.handles = make(map[string]*Handle)
	}
	if strings.LastIndex(url, "/") == len(url)-1 {
		url = strings.TrimRight(url, "/")
	}

	handle := new(Handle)
	handle.NeedLogin = needLogin
	handle.CallBackHandle = callBackHandle
	h.handles[strings.ToLower(url)] = handle
	h.handlesList = h.sort(url, h.handlesList)
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
	if h.handles == nil {
		h.handles = make(map[string]*Handle)
	}
	if strings.LastIndex(url, "/") == len(url)-1 {
		url = strings.TrimRight(url, "/")
	}
	handleVar := new(Handle)
	handleVar.NeedLogin = false
	handleVar.Handler = handle
	h.handles[strings.ToLower(url)] = handleVar
	h.handlesList = h.sort(url, h.handlesList)
}

func (h *HTTPServer) registerVirtualPath(virtualPath, localPath string, needLogin bool) {
	if h.handlesFile == nil {
		h.handlesFile = make(map[string]*Handle)
	}

	fileHandle := http.StripPrefix(virtualPath, http.FileServer(http.Dir(localPath)))

	handle := new(Handle)
	handle.NeedLogin = needLogin
	handle.Handler = fileHandle
	h.handlesFile[strings.ToLower(virtualPath)] = handle
	h.handlesListFile = h.sort(virtualPath, h.handlesListFile)
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
	if h.handlesFile == nil {
		h.handlesFile = make(map[string]*Handle)
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

	handleVar := new(Handle)
	handleVar.NeedLogin = false
	handleVar.Handler = fileHandle
	h.handlesFile[strings.ToLower(virtualPath)] = handleVar
	h.handlesListFile = h.sort(virtualPath, h.handlesListFile)
}
