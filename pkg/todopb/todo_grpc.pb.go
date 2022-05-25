// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: todo/todo.proto

package todopb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ToDoServiceClient is the client API for ToDoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ToDoServiceClient interface {
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AllProjects(ctx context.Context, in *AllProjectsRequest, opts ...grpc.CallOption) (*AllProjectsResponse, error)
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReplaceTask(ctx context.Context, in *ReplaceTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ProjectsUpdates(ctx context.Context, in *ProjectsUpdatesRequest, opts ...grpc.CallOption) (ToDoService_ProjectsUpdatesClient, error)
}

type toDoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewToDoServiceClient(cc grpc.ClientConnInterface) ToDoServiceClient {
	return &toDoServiceClient{cc}
}

func (c *toDoServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) AllProjects(ctx context.Context, in *AllProjectsRequest, opts ...grpc.CallOption) (*AllProjectsResponse, error) {
	out := new(AllProjectsResponse)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/AllProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/AddTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) ReplaceTask(ctx context.Context, in *ReplaceTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/ReplaceTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/DeleteTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/todo.ToDoService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) ProjectsUpdates(ctx context.Context, in *ProjectsUpdatesRequest, opts ...grpc.CallOption) (ToDoService_ProjectsUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ToDoService_ServiceDesc.Streams[0], "/todo.ToDoService/ProjectsUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &toDoServiceProjectsUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ToDoService_ProjectsUpdatesClient interface {
	Recv() (*Project, error)
	grpc.ClientStream
}

type toDoServiceProjectsUpdatesClient struct {
	grpc.ClientStream
}

func (x *toDoServiceProjectsUpdatesClient) Recv() (*Project, error) {
	m := new(Project)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ToDoServiceServer is the server API for ToDoService service.
// All implementations must embed UnimplementedToDoServiceServer
// for forward compatibility
type ToDoServiceServer interface {
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
	UpdateProject(context.Context, *UpdateProjectRequest) (*emptypb.Empty, error)
	AllProjects(context.Context, *AllProjectsRequest) (*AllProjectsResponse, error)
	AddTask(context.Context, *AddTaskRequest) (*emptypb.Empty, error)
	ReplaceTask(context.Context, *ReplaceTaskRequest) (*emptypb.Empty, error)
	DeleteTask(context.Context, *DeleteTaskRequest) (*emptypb.Empty, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*emptypb.Empty, error)
	ProjectsUpdates(*ProjectsUpdatesRequest, ToDoService_ProjectsUpdatesServer) error
	mustEmbedUnimplementedToDoServiceServer()
}

// UnimplementedToDoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedToDoServiceServer struct {
}

func (UnimplementedToDoServiceServer) CreateProject(context.Context, *CreateProjectRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedToDoServiceServer) UpdateProject(context.Context, *UpdateProjectRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProject not implemented")
}
func (UnimplementedToDoServiceServer) AllProjects(context.Context, *AllProjectsRequest) (*AllProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllProjects not implemented")
}
func (UnimplementedToDoServiceServer) AddTask(context.Context, *AddTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTask not implemented")
}
func (UnimplementedToDoServiceServer) ReplaceTask(context.Context, *ReplaceTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplaceTask not implemented")
}
func (UnimplementedToDoServiceServer) DeleteTask(context.Context, *DeleteTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
func (UnimplementedToDoServiceServer) DeleteProject(context.Context, *DeleteProjectRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedToDoServiceServer) ProjectsUpdates(*ProjectsUpdatesRequest, ToDoService_ProjectsUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method ProjectsUpdates not implemented")
}
func (UnimplementedToDoServiceServer) mustEmbedUnimplementedToDoServiceServer() {}

// UnsafeToDoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ToDoServiceServer will
// result in compilation errors.
type UnsafeToDoServiceServer interface {
	mustEmbedUnimplementedToDoServiceServer()
}

func RegisterToDoServiceServer(s grpc.ServiceRegistrar, srv ToDoServiceServer) {
	s.RegisterService(&ToDoService_ServiceDesc, srv)
}

func _ToDoService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).UpdateProject(ctx, req.(*UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_AllProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).AllProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/AllProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).AllProjects(ctx, req.(*AllProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/AddTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_ReplaceTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplaceTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).ReplaceTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/ReplaceTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).ReplaceTask(ctx, req.(*ReplaceTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_DeleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).DeleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/DeleteTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).DeleteTask(ctx, req.(*DeleteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo.ToDoService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_ProjectsUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProjectsUpdatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ToDoServiceServer).ProjectsUpdates(m, &toDoServiceProjectsUpdatesServer{stream})
}

type ToDoService_ProjectsUpdatesServer interface {
	Send(*Project) error
	grpc.ServerStream
}

type toDoServiceProjectsUpdatesServer struct {
	grpc.ServerStream
}

func (x *toDoServiceProjectsUpdatesServer) Send(m *Project) error {
	return x.ServerStream.SendMsg(m)
}

// ToDoService_ServiceDesc is the grpc.ServiceDesc for ToDoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ToDoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "todo.ToDoService",
	HandlerType: (*ToDoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProject",
			Handler:    _ToDoService_CreateProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _ToDoService_UpdateProject_Handler,
		},
		{
			MethodName: "AllProjects",
			Handler:    _ToDoService_AllProjects_Handler,
		},
		{
			MethodName: "AddTask",
			Handler:    _ToDoService_AddTask_Handler,
		},
		{
			MethodName: "ReplaceTask",
			Handler:    _ToDoService_ReplaceTask_Handler,
		},
		{
			MethodName: "DeleteTask",
			Handler:    _ToDoService_DeleteTask_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ToDoService_DeleteProject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ProjectsUpdates",
			Handler:       _ToDoService_ProjectsUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "todo/todo.proto",
}
