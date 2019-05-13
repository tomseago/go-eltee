package eltee

import (
	"context"
	"fmt"
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

	log.Infof("Ping %v from %v", req.GetVal(), peer.Addr)

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

func (asrv *apiServer) StateNames(ctx context.Context, req *api.Void) (*api.StringMsg, error) {
	rsp := &api.StringMsg{}

	rsp.List = asrv.server.stateJuggler.StateNames()

	return rsp, nil
}

func (asrv *apiServer) ControlPoints(ctx context.Context, req *api.StringMsg) (*api.ControlPointList, error) {

	juggler := asrv.server.stateJuggler
	stateName := req.GetVal()
	state := juggler.State(stateName)
	if state == nil {
		return nil, fmt.Errorf("Invalid request. No state named %v", stateName)
	}

	rsp := &api.ControlPointList{
		Cps:   make([]*api.ControlPoint, 0, len(state.ControlPoints())),
		State: stateName,
	}

	for _, cp := range state.ControlPoints() {
		rsp.Cps = append(rsp.Cps, cp.ToApi())
	}

	return rsp, nil
}

func (asrv *apiServer) SetControlPoints(ctx context.Context, req *api.ControlPointList) (*api.Void, error) {

	state := asrv.server.stateJuggler.State(req.GetState())
	if state == nil {
		return nil, fmt.Errorf("Could not find state %v", req.GetState())
	}

	for _, apiCP := range req.GetCps() {
		cp := state.ControlPoint(apiCP.GetName())

		if cp == nil {
			if !req.GetUpsert() {
				return nil, fmt.Errorf("State does not have a control point named '%v' and upsert was false", apiCP.GetName())
			}
			// Else, we make it happen
			cp = CreateControlPointFromApi(apiCP)
			if cp == nil {
				return nil, fmt.Errorf("Was not able to create control point '%v'", apiCP.GetName())
			}

			state.AddControlPoint(cp.Name(), cp)
		}

		//log.Debugf("Old: %v", cp)
		cp.SetFromApi(apiCP)
		//log.Infof("New: %v", cp)
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) RemoveControlPoints(ctx context.Context, req *api.ControlPointList) (*api.Void, error) {

	state := asrv.server.stateJuggler.State(req.GetState())
	if state == nil {
		return nil, fmt.Errorf("Could not find state %v", req.GetState())
	}

	for _, apiCP := range req.GetCps() {
		state.RemoveControlPoint(apiCP.GetName())
	}

	return &api.Void{}, nil
}

//////////////////

func (asrv *apiServer) ApplyState(ctx context.Context, req *api.StringMsg) (*api.Void, error) {

	err := asrv.server.stateJuggler.ApplyState(req.GetVal())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) LoadableStateNames(ctx context.Context, req *api.Void) (*api.StringMsg, error) {

	rsp := &api.StringMsg{}

	var err error
	rsp.List, err = asrv.server.stateJuggler.LoadableStateNames(asrv.server.loadableStatesDir)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (asrv *apiServer) LoadLoadableState(ctx context.Context, req *api.StringMsg) (*api.Void, error) {

	err := asrv.server.stateJuggler.LoadLoadableState(asrv.server.loadableStatesDir, req.GetVal())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) SaveLoadableState(ctx context.Context, req *api.StringMsg) (*api.Void, error) {

	err := asrv.server.stateJuggler.SaveLoadableState(asrv.server.loadableStatesDir, req.GetVal())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) SaveAllStates(ctx context.Context, req *api.Void) (*api.Void, error) {

	asrv.server.stateJuggler.SaveAll(asrv.server.loadableStatesDir)

	return &api.Void{}, nil
}

//////////////////

func (asrv *apiServer) AddState(ctx context.Context, req *api.StringMsg) (*api.Void, error) {

	err := asrv.server.stateJuggler.AddState(req.GetVal())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) RemoveState(ctx context.Context, req *api.StringMsg) (*api.Void, error) {

	err := asrv.server.stateJuggler.RemoveState(req.GetVal())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) CopyStateTo(ctx context.Context, req *api.SrcDest) (*api.Void, error) {

	err := asrv.server.stateJuggler.CopyStateTo(req.GetSrc(), req.GetDest())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) MoveStateTo(ctx context.Context, req *api.SrcDest) (*api.Void, error) {

	err := asrv.server.stateJuggler.MoveStateTo(req.GetSrc(), req.GetDest())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) ApplyStateTo(ctx context.Context, req *api.SrcDest) (*api.Void, error) {

	err := asrv.server.stateJuggler.ApplyStateTo(req.GetSrc(), req.GetDest())
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}
