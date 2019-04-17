package eltee

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"

	"github.com/tomseago/go-eltee/api"
)

type apiServer struct {
	server *Server

	grpc *grpc.Server
}

func NewApiServer(server *Server) *apiServer {
	asrv := &apiServer{
		server: server,
	}

	return asrv
}

func (asrv *apiServer) Start() {
	// TODO: More options for where to start the server
	lis, err := net.Listen("tcp", ":3434")
	if err != nil {
		log.Errorf("Unable to start api server: %v", err)
		return
	}

	asrv.grpc = grpc.NewServer()
	api.RegisterElTeeServer(asrv.grpc, asrv)

	// More config

	log.Info("gRPC server started")
	asrv.grpc.Serve(lis)

	log.Info("gRPC server stopped")
}

//////////////////
func (asrv *apiServer) Ping(ctx context.Context, req *api.StringMsg) (*api.StringMsg, error) {

	peer, _ := peer.FromContext(ctx)

	log.Infof("Ping %v from %v", req.Val, peer.Addr)

	rsp := &api.StringMsg{
		Val: "Pong!",
	}

	return rsp, nil
}

func (asrv *apiServer) ProfileLibrary(ctx context.Context, req *api.Void) (*api.ProfilesResponse, error) {

	rsp := &api.ProfilesResponse{
		Profiles: make(map[string]*api.Profile),
	}

	for k, v := range asrv.server.library.Profiles {
		rsp.Profiles[k] = v.ToAPI()
	}

	return rsp, nil
}

//////////////////
