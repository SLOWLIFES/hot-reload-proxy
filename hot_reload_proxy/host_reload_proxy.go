package hot_reload_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type hotReloadProxy struct {
	reverseProxy *httputil.ReverseProxy
}

func New() *hotReloadProxy {
	return &hotReloadProxy{}
}

func (h *hotReloadProxy) SwitchHost(addr string) error {
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		writer.WriteHeader(404)
		_, _ = writer.Write([]byte("服务不可用"))
	}
	h.reverseProxy = proxy
	return nil
}

func (h *hotReloadProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.reverseProxy == nil {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("暂无服务可用"))
	} else {
		h.reverseProxy.ServeHTTP(w, r)
	}
}
