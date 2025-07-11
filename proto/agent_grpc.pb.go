// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/agent.proto

package agentpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AgentReporter_ReportData_FullMethodName    = "/agent.AgentReporter/ReportData"
	AgentReporter_StreamPodLogs_FullMethodName = "/agent.AgentReporter/StreamPodLogs"
)

// AgentReporterClient is the client API for AgentReporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// gRPC service for sending agent data
type AgentReporterClient interface {
	ReportData(ctx context.Context, in *AgentData, opts ...grpc.CallOption) (*ReportResponse, error)
	StreamPodLogs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LogStream], error)
}

type agentReporterClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentReporterClient(cc grpc.ClientConnInterface) AgentReporterClient {
	return &agentReporterClient{cc}
}

func (c *agentReporterClient) ReportData(ctx context.Context, in *AgentData, opts ...grpc.CallOption) (*ReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, AgentReporter_ReportData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentReporterClient) StreamPodLogs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LogStream], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AgentReporter_ServiceDesc.Streams[0], AgentReporter_StreamPodLogs_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[LogRequest, LogStream]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AgentReporter_StreamPodLogsClient = grpc.ServerStreamingClient[LogStream]

// AgentReporterServer is the server API for AgentReporter service.
// All implementations must embed UnimplementedAgentReporterServer
// for forward compatibility.
//
// gRPC service for sending agent data
type AgentReporterServer interface {
	ReportData(context.Context, *AgentData) (*ReportResponse, error)
	StreamPodLogs(*LogRequest, grpc.ServerStreamingServer[LogStream]) error
	mustEmbedUnimplementedAgentReporterServer()
}

// UnimplementedAgentReporterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAgentReporterServer struct{}

func (UnimplementedAgentReporterServer) ReportData(context.Context, *AgentData) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportData not implemented")
}
func (UnimplementedAgentReporterServer) StreamPodLogs(*LogRequest, grpc.ServerStreamingServer[LogStream]) error {
	return status.Errorf(codes.Unimplemented, "method StreamPodLogs not implemented")
}
func (UnimplementedAgentReporterServer) mustEmbedUnimplementedAgentReporterServer() {}
func (UnimplementedAgentReporterServer) testEmbeddedByValue()                       {}

// UnsafeAgentReporterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentReporterServer will
// result in compilation errors.
type UnsafeAgentReporterServer interface {
	mustEmbedUnimplementedAgentReporterServer()
}

func RegisterAgentReporterServer(s grpc.ServiceRegistrar, srv AgentReporterServer) {
	// If the following call pancis, it indicates UnimplementedAgentReporterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AgentReporter_ServiceDesc, srv)
}

func _AgentReporter_ReportData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentReporterServer).ReportData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentReporter_ReportData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentReporterServer).ReportData(ctx, req.(*AgentData))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentReporter_StreamPodLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LogRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentReporterServer).StreamPodLogs(m, &grpc.GenericServerStream[LogRequest, LogStream]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AgentReporter_StreamPodLogsServer = grpc.ServerStreamingServer[LogStream]

// AgentReporter_ServiceDesc is the grpc.ServiceDesc for AgentReporter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentReporter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "agent.AgentReporter",
	HandlerType: (*AgentReporterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportData",
			Handler:    _AgentReporter_ReportData_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamPodLogs",
			Handler:       _AgentReporter_StreamPodLogs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/agent.proto",
}
