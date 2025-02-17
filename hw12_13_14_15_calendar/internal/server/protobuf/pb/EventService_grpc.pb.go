// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/EventService.proto

package pb

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
	Calendar_AddEvent_FullMethodName    = "/event.Calendar/AddEvent"
	Calendar_UpdateEvent_FullMethodName = "/event.Calendar/UpdateEvent"
	Calendar_DeleteEvent_FullMethodName = "/event.Calendar/DeleteEvent"
	Calendar_GetEvent_FullMethodName    = "/event.Calendar/GetEvent"
	Calendar_GetEvents_FullMethodName   = "/event.Calendar/GetEvents"
)

// CalendarClient is the client API for Calendar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarClient interface {
	AddEvent(ctx context.Context, in *AddEventReq, opts ...grpc.CallOption) (*AddEventResp, error)
	UpdateEvent(ctx context.Context, in *UpdateEventReq, opts ...grpc.CallOption) (*UpdateEventResp, error)
	DeleteEvent(ctx context.Context, in *DeleteEventReq, opts ...grpc.CallOption) (*DeleteEventResp, error)
	GetEvent(ctx context.Context, in *GetEventReq, opts ...grpc.CallOption) (*GetEventResp, error)
	GetEvents(ctx context.Context, in *GetEventsReq, opts ...grpc.CallOption) (*GetEventsResp, error)
}

type calendarClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarClient(cc grpc.ClientConnInterface) CalendarClient {
	return &calendarClient{cc}
}

func (c *calendarClient) AddEvent(ctx context.Context, in *AddEventReq, opts ...grpc.CallOption) (*AddEventResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddEventResp)
	err := c.cc.Invoke(ctx, Calendar_AddEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) UpdateEvent(ctx context.Context, in *UpdateEventReq, opts ...grpc.CallOption) (*UpdateEventResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEventResp)
	err := c.cc.Invoke(ctx, Calendar_UpdateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) DeleteEvent(ctx context.Context, in *DeleteEventReq, opts ...grpc.CallOption) (*DeleteEventResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteEventResp)
	err := c.cc.Invoke(ctx, Calendar_DeleteEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEvent(ctx context.Context, in *GetEventReq, opts ...grpc.CallOption) (*GetEventResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEventResp)
	err := c.cc.Invoke(ctx, Calendar_GetEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetEvents(ctx context.Context, in *GetEventsReq, opts ...grpc.CallOption) (*GetEventsResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEventsResp)
	err := c.cc.Invoke(ctx, Calendar_GetEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServer is the server API for Calendar service.
// All implementations must embed UnimplementedCalendarServer
// for forward compatibility.
type CalendarServer interface {
	AddEvent(context.Context, *AddEventReq) (*AddEventResp, error)
	UpdateEvent(context.Context, *UpdateEventReq) (*UpdateEventResp, error)
	DeleteEvent(context.Context, *DeleteEventReq) (*DeleteEventResp, error)
	GetEvent(context.Context, *GetEventReq) (*GetEventResp, error)
	GetEvents(context.Context, *GetEventsReq) (*GetEventsResp, error)
	mustEmbedUnimplementedCalendarServer()
}

// UnimplementedCalendarServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCalendarServer struct{}

func (UnimplementedCalendarServer) AddEvent(context.Context, *AddEventReq) (*AddEventResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEvent not implemented")
}
func (UnimplementedCalendarServer) UpdateEvent(context.Context, *UpdateEventReq) (*UpdateEventResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedCalendarServer) DeleteEvent(context.Context, *DeleteEventReq) (*DeleteEventResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedCalendarServer) GetEvent(context.Context, *GetEventReq) (*GetEventResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedCalendarServer) GetEvents(context.Context, *GetEventsReq) (*GetEventsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (UnimplementedCalendarServer) mustEmbedUnimplementedCalendarServer() {}
func (UnimplementedCalendarServer) testEmbeddedByValue()                  {}

// UnsafeCalendarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServer will
// result in compilation errors.
type UnsafeCalendarServer interface {
	mustEmbedUnimplementedCalendarServer()
}

func RegisterCalendarServer(s grpc.ServiceRegistrar, srv CalendarServer) {
	// If the following call pancis, it indicates UnimplementedCalendarServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Calendar_ServiceDesc, srv)
}

func _Calendar_AddEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).AddEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calendar_AddEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).AddEvent(ctx, req.(*AddEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calendar_UpdateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).UpdateEvent(ctx, req.(*UpdateEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calendar_DeleteEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).DeleteEvent(ctx, req.(*DeleteEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calendar_GetEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEvent(ctx, req.(*GetEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calendar_GetEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetEvents(ctx, req.(*GetEventsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Calendar_ServiceDesc is the grpc.ServiceDesc for Calendar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calendar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.Calendar",
	HandlerType: (*CalendarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddEvent",
			Handler:    _Calendar_AddEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _Calendar_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _Calendar_DeleteEvent_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _Calendar_GetEvent_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _Calendar_GetEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/EventService.proto",
}
