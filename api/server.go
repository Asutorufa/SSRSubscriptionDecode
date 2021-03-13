package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Asutorufa/yuhaiin/app"
	"github.com/Asutorufa/yuhaiin/config"
	"github.com/Asutorufa/yuhaiin/net/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Process struct {
	UnimplementedProcessInitServer

	singleInstance chan bool
	message        chan string
	m              *manager
}

func NewProcess(e *app.Entrance) (*Process, error) {
	p := &Process{}
	p.m = newManager(e)
	err := p.m.Start()
	return p, err
}

func (s *Process) Host() string {
	return s.m.Host()
}

func (s *Process) CreateLockFile(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if !s.m.lockfile() {
		return &emptypb.Empty{}, errors.New("create lock file false")
	}

	if !s.m.initApp() {
		return &emptypb.Empty{}, errors.New("init Process Failed")
	}

	s.m.connect()
	return &emptypb.Empty{}, nil
}

func (s *Process) ProcessInit(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Process) GetRunningHost(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	host, err := app.ReadLockFile()
	if err != nil {
		return &wrapperspb.StringValue{}, err
	}
	return &wrapperspb.StringValue{Value: host}, nil
}

func (s *Process) ClientOn(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if s.singleInstance != nil {
		select {
		case <-s.singleInstance:
			break
		default:
			s.message <- "on"
			return &emptypb.Empty{}, nil
		}
	}
	return &emptypb.Empty{}, errors.New("no client")
}

func (s *Process) ProcessExit(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, app.LockFileClose()
}

func (s *Process) SingleInstance(srv ProcessInit_SingleInstanceServer) error {
	if s.singleInstance != nil {
		select {
		case <-s.singleInstance:
			break
		default:
			return errors.New("already exist one client")
		}
	}

	s.singleInstance = make(chan bool)
	s.message = make(chan string, 1)
	ctx := srv.Context()

	for {
		select {
		case m := <-s.message:
			err := srv.Send(&wrapperspb.StringValue{Value: m})
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Call Client Open Window.")
		case <-ctx.Done():
			close(s.message)
			close(s.singleInstance)
			if s.m.killWDC {
				panic("client exit")
			}
			return ctx.Err()
		}
	}
}

func (s *Process) GetKernelPid(context.Context, *emptypb.Empty) (*wrapperspb.UInt32Value, error) {
	return &wrapperspb.UInt32Value{Value: uint32(os.Getpid())}, nil
}

func (s *Process) StopKernel(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	defer os.Exit(0)
	return &emptypb.Empty{}, nil
}

type Config struct {
	UnimplementedConfigServer
	entrance *app.Entrance
}

func NewConfig(e *app.Entrance) *Config {
	return &Config{
		entrance: e,
	}
}

func (c *Config) GetConfig(context.Context, *emptypb.Empty) (*config.Setting, error) {
	conf, err := c.entrance.GetConfig()
	return conf, err
}

func (c *Config) SetConfig(_ context.Context, req *config.Setting) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.entrance.SetConFig(req)
}

func (c *Config) ReimportRule(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.entrance.RefreshMapping()
}

func (c *Config) GetRate(_ *emptypb.Empty, srv Config_GetRateServer) error {
	fmt.Println("Start Send Flow Message to Client.")
	//TODO deprecated string
	da, ua := c.entrance.GetDownload(), c.entrance.GetUpload()
	var dr string
	var ur string
	ctx := srv.Context()
	for {
		dr = utils.ReducedUnitStr(float64(c.entrance.GetDownload()-da)) + "/S"
		ur = utils.ReducedUnitStr(float64(c.entrance.GetUpload()-ua)) + "/S"
		da, ua = c.entrance.GetDownload(), c.entrance.GetUpload()

		err := srv.Send(&DaUaDrUr{
			Download: utils.ReducedUnitStr(float64(da)),
			Upload:   utils.ReducedUnitStr(float64(ua)),
			DownRate: dr,
			UpRate:   ur,
		})
		if err != nil {
			log.Println(err)
		}
		select {
		case <-ctx.Done():
			fmt.Println("Client is Hidden, Close Stream.")
			return ctx.Err()
		case <-time.After(time.Second):
			continue
		}
	}
}

type Node struct {
	UnimplementedNodeServer
	entrance *app.Entrance
}

func NewNode(e *app.Entrance) *Node {
	return &Node{
		entrance: e,
	}
}

func (n *Node) GetNodes(context.Context, *emptypb.Empty) (*Nodes, error) {
	nodes := &Nodes{Value: map[string]*AllGroupOrNode{}}
	nods := n.entrance.GetANodes()
	for key := range nods {
		nodes.Value[key] = &AllGroupOrNode{Value: nods[key]}
	}
	return nodes, nil
}

func (n *Node) GetGroup(context.Context, *emptypb.Empty) (*AllGroupOrNode, error) {
	groups, err := n.entrance.GetGroups()
	return &AllGroupOrNode{Value: groups}, err
}

func (n *Node) GetNode(_ context.Context, req *wrapperspb.StringValue) (*AllGroupOrNode, error) {
	nodes, err := n.entrance.GetNodes(req.Value)
	return &AllGroupOrNode{Value: nodes}, err
}

func (n *Node) GetNowGroupAndName(context.Context, *emptypb.Empty) (*GroupAndNode, error) {
	node, group := n.entrance.GetNNodeAndNGroup()
	return &GroupAndNode{Node: node, Group: group}, nil
}

func (n *Node) AddNode(_ context.Context, req *NodeMap) (*emptypb.Empty, error) {
	// TODO add node
	return &emptypb.Empty{}, nil
}

func (n *Node) ModifyNode(context.Context, *NodeMap) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (n *Node) DeleteNode(_ context.Context, req *GroupAndNode) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, n.entrance.DeleteNode(req.Group, req.Node)
}

func (n *Node) ChangeNowNode(_ context.Context, req *GroupAndNode) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, n.entrance.ChangeNNode(req.Group, req.Node)
}

func (n *Node) Latency(_ context.Context, req *GroupAndNode) (*wrapperspb.StringValue, error) {
	latency, err := n.entrance.Latency(req.Group, req.Node)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: latency.String()}, err
}

type Subscribe struct {
	UnimplementedSubscribeServer
	entrance *app.Entrance
}

func NewSubscribe(e *app.Entrance) *Subscribe {
	return &Subscribe{
		entrance: e,
	}
}

func (s *Subscribe) UpdateSub(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.entrance.UpdateSub()
}

func (s *Subscribe) GetSubLinks(context.Context, *emptypb.Empty) (*Links, error) {
	links, err := s.entrance.GetLinks()
	if err != nil {
		return nil, err
	}
	l := &Links{}
	l.Value = map[string]*Link{}
	for key := range links {
		l.Value[key] = &Link{
			Type: links[key].Type,
			Url:  links[key].Url,
		}
	}
	return l, nil
}

func (s *Subscribe) AddSubLink(ctx context.Context, req *Link) (*Links, error) {
	err := s.entrance.AddLink(req.Name, req.Type, req.Url)
	if err != nil {
		return nil, fmt.Errorf("api:AddSubLink -> %v", err)
	}
	return s.GetSubLinks(ctx, &emptypb.Empty{})
}

func (s *Subscribe) DeleteSubLink(ctx context.Context, req *wrapperspb.StringValue) (*Links, error) {
	err := s.entrance.DeleteLink(req.Value)
	if err != nil {
		return nil, err
	}
	return s.GetSubLinks(ctx, &emptypb.Empty{})
}
