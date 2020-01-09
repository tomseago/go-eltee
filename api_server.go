package eltee

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"

	"github.com/tomseago/go-eltee/api"
)

type cpListener struct {
    stream api.ElTee_ControlPointChangesServer
    c chan *api.ControlPointList

    cleanup func(l *cpListener)
}

func NewCpListener(stream api.ElTee_ControlPointChangesServer, cleanup func(l *cpListener)) *cpListener {
    out := &cpListener{
        stream: stream,
        c: make(chan *api.ControlPointList, 4),
        cleanup: cleanup,
    }

    return out
}

func (l *cpListener) SendList(list *api.ControlPointList) {
    log.Debugf("cpListener.SendList list=%v",list)
    l.c <- list
}

func (l *cpListener) Run() {
    log.Debugf("cpListener.Run")


    for done := false ; !done ; {
        select {
        case list := <-l.c:
            log.Debugf("from channel list=%v",list)
            err := l.stream.Send(list)
            if err != nil {
                log.Debugf("cpListener.stream.Send err=%v", err);
                // l.done <- true
                done = true
            }
        }
    }

    log.Debugf("cpListener calling cleanup");
    l.cleanup(l)
}


type apiServer struct {
	server *Server

	grpc *grpc.Server

	// List of client channels that want to receive cp changes
	cpListeners []*cpListener
	deadListeners chan *cpListener
}

func NewApiServer(server *Server) *apiServer {
	asrv := &apiServer{
		server: server,
		cpListeners: make([]*cpListener,0),
		deadListeners: make(chan *cpListener, 2),
	}

	return asrv
}

func (asrv *apiServer) SendCPChanges(list *api.ControlPointList) {
    listeners := asrv.cpListeners

    for _, v := range listeners {
        v.SendList(list)
    }
}

func (asrv *apiServer) deadListenerRemover() {
    for {
        dead := <- asrv.deadListeners

        ix := -1
        for i, v := range asrv.cpListeners {
            if v == dead {
                ix = i
                break
            }
        }
        if ix != -1 {
            // Copy last element into this position
            asrv.cpListeners[ix] = asrv.cpListeners[len(asrv.cpListeners) - 1]
            // Remove the last element
            asrv.cpListeners = asrv.cpListeners[:len(asrv.cpListeners) - 1]
        }
    }
}

////////

func (asrv *apiServer) Start() {
	// TODO: More options for where to start the server
	lis, err := net.Listen("tcp", ":3434")
	if err != nil {
		log.Errorf("Unable to start api server: %v", err)
		return
	}

	asrv.grpc = grpc.NewServer()
	api.RegisterElTeeServer(asrv.grpc, asrv)

	// Need to get rid of dead listeners
	go asrv.deadListenerRemover()

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

func (asrv *apiServer) FixtureList(ctx context.Context, req *api.Void) (*api.FixtureListResponse, error) {

    rsp := &api.FixtureListResponse{
        Fixtures: make([]*api.Fixture, 0),
    }

    for _, v := range asrv.server.fixtures {
        rsp.Fixtures = append(rsp.Fixtures, v.ToAPI())
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

	// Make sure all the control points exist before we set any of them
    for _, apiCP := range req.GetCps() {
        cp := state.ControlPoint(apiCP.GetName())
        if cp == nil {
            if !req.GetUpsert() {
                return nil, fmt.Errorf("State does not have a control point named '%v' and upsert was false", apiCP.GetName())
            }
        }
    }

    // Okay everything should be fine

	for _, apiCP := range req.GetCps() {
		cp := state.ControlPoint(apiCP.GetName())

		if cp == nil {
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

	// And now tell other listeners about this
	asrv.SendCPChanges(req)

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

func (asrv *apiServer) ControlPointChanges(req *api.Void, stream api.ElTee_ControlPointChangesServer ) (error) {

    done := make(chan bool)

    listener := NewCpListener(stream, func(l *cpListener) {
        asrv.deadListeners <- l
        done <- true
    })
    asrv.cpListeners = append(asrv.cpListeners, listener)

    go listener.Run()

    <-done

    return nil
}

//////////////////

func (asrv *apiServer) SetFixturePatches(ctx context.Context, req *api.FixturePatchMap) (*api.Void, error) {

	state := asrv.server.stateJuggler.State(req.GetState())
	if state == nil {
		return nil, fmt.Errorf("Could not find state %v", req.GetState())
	}

	for fixtureName, apiFixPatch := range req.GetByFixture() {
		fp := state.FixturePatch(fixtureName)

		if fp == nil {
			if !req.GetUpsert() {
				return nil, fmt.Errorf("State does not have a fixture patch for '%v' and upsert was false", fixtureName)
			}

			// // Else, we make it happen
			// cp = CreateControlPointFromApi(apiCP)
			// if cp == nil {
			// 	return nil, fmt.Errorf("Was not able to create control point '%v'", apiCP.GetName())
			// }

			// state.AddControlPoint(cp.Name(), cp)
		}

		//log.Debugf("Old: %v", cp)
		fp.SetFromApi(apiFixPatch)
		//log.Infof("New: %v", cp)
	}

	return &api.Void{}, nil
}

func (asrv *apiServer) RemoveFixturePatches(ctx context.Context, req *api.FixturePatchMap) (*api.Void, error) {

	state := asrv.server.stateJuggler.State(req.GetState())
	if state == nil {
		return nil, fmt.Errorf("Could not find state %v", req.GetState())
	}

	for fixtureName, apiFixPatch := range req.GetByFixture() {
		fp := state.FixturePatch(fixtureName)

		if fp == nil {
			continue
		}

		fp.RemoveFromApi(apiFixPatch)

		// TODO: Check for empty and remove the whole thing??
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

//////////////////

func (asrv *apiServer) FixturePatches(ctx context.Context, req *api.StringMsg) (*api.FixturePatchMap, error) {

	state := asrv.server.stateJuggler.State(req.GetVal())
	if state == nil {
		return nil, fmt.Errorf("Could not find state %v", req.GetVal())
	}

	return state.FixturePatchMap(), nil
}
