package rpcserver

import (
	"log"
	"net/http"
	"strings"
	"path"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"github.com/elazarl/go-bindata-assetfs"
	"wps_store/pkg/swagger"

	gw "wps_store/rpc"
)

var (
	ServerHttpPort string
	HttpEndPoint   string
)

func RunHttpServer() (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// grpc服务地址
	EndPoint = ":" + ServerPort
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// HTTP转grpc
	err = gw.RegisterStoreApiServiceHandlerFromEndpoint(ctx, gwmux, EndPoint, opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}


	HttpEndPoint = ":" + ServerHttpPort
	
	// swagger ui
	mux := http.NewServeMux()
    mux.Handle("/", gwmux)
    mux.HandleFunc("/swagger/", serveSwaggerFile)
	serveSwaggerUI(mux)
	
	log.Println("HTTP Listen success:", HttpEndPoint)

    s := &http.Server{
        Addr:      HttpEndPoint,
        Handler:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            mux.ServeHTTP(w, r)
        }),
    }
	s.ListenAndServe()
	return err
}

func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if ! strings.HasSuffix(r.URL.Path, "swagger.json") {
	  log.Printf("Not Found: %s", r.URL.Path)
	  http.NotFound(w, r)
	  return
  }

  p := strings.TrimPrefix(r.URL.Path, "/swagger/")
  p = path.Join("rpc", p)

  log.Printf("Serving swagger-file: %s", p)

  http.ServeFile(w, r, p)
}

func serveSwaggerUI(mux *http.ServeMux) {
  fileServer := http.FileServer(&assetfs.AssetFS{
	  Asset:    swagger.Asset,
	  AssetDir: swagger.AssetDir,
	  Prefix:   "swagger-ui",
  })
  prefix := "/swagger-ui/"
  mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
