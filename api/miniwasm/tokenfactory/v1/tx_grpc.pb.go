// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: miniwasm/tokenfactory/v1/tx.proto

package tokenfactoryv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Msg_CreateDenom_FullMethodName       = "/miniwasm.tokenfactory.v1.Msg/CreateDenom"
	Msg_Mint_FullMethodName              = "/miniwasm.tokenfactory.v1.Msg/Mint"
	Msg_Burn_FullMethodName              = "/miniwasm.tokenfactory.v1.Msg/Burn"
	Msg_ChangeAdmin_FullMethodName       = "/miniwasm.tokenfactory.v1.Msg/ChangeAdmin"
	Msg_SetDenomMetadata_FullMethodName  = "/miniwasm.tokenfactory.v1.Msg/SetDenomMetadata"
	Msg_SetBeforeSendHook_FullMethodName = "/miniwasm.tokenfactory.v1.Msg/SetBeforeSendHook"
	Msg_ForceTransfer_FullMethodName     = "/miniwasm.tokenfactory.v1.Msg/ForceTransfer"
	Msg_UpdateParams_FullMethodName      = "/miniwasm.tokenfactory.v1.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	CreateDenom(ctx context.Context, in *MsgCreateDenom, opts ...grpc.CallOption) (*MsgCreateDenomResponse, error)
	Mint(ctx context.Context, in *MsgMint, opts ...grpc.CallOption) (*MsgMintResponse, error)
	Burn(ctx context.Context, in *MsgBurn, opts ...grpc.CallOption) (*MsgBurnResponse, error)
	ChangeAdmin(ctx context.Context, in *MsgChangeAdmin, opts ...grpc.CallOption) (*MsgChangeAdminResponse, error)
	SetDenomMetadata(ctx context.Context, in *MsgSetDenomMetadata, opts ...grpc.CallOption) (*MsgSetDenomMetadataResponse, error)
	SetBeforeSendHook(ctx context.Context, in *MsgSetBeforeSendHook, opts ...grpc.CallOption) (*MsgSetBeforeSendHookResponse, error)
	ForceTransfer(ctx context.Context, in *MsgForceTransfer, opts ...grpc.CallOption) (*MsgForceTransferResponse, error)
	// UpdateParams defines an operation for updating the x/tokenfactory module
	// parameters.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateDenom(ctx context.Context, in *MsgCreateDenom, opts ...grpc.CallOption) (*MsgCreateDenomResponse, error) {
	out := new(MsgCreateDenomResponse)
	err := c.cc.Invoke(ctx, Msg_CreateDenom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Mint(ctx context.Context, in *MsgMint, opts ...grpc.CallOption) (*MsgMintResponse, error) {
	out := new(MsgMintResponse)
	err := c.cc.Invoke(ctx, Msg_Mint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Burn(ctx context.Context, in *MsgBurn, opts ...grpc.CallOption) (*MsgBurnResponse, error) {
	out := new(MsgBurnResponse)
	err := c.cc.Invoke(ctx, Msg_Burn_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ChangeAdmin(ctx context.Context, in *MsgChangeAdmin, opts ...grpc.CallOption) (*MsgChangeAdminResponse, error) {
	out := new(MsgChangeAdminResponse)
	err := c.cc.Invoke(ctx, Msg_ChangeAdmin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetDenomMetadata(ctx context.Context, in *MsgSetDenomMetadata, opts ...grpc.CallOption) (*MsgSetDenomMetadataResponse, error) {
	out := new(MsgSetDenomMetadataResponse)
	err := c.cc.Invoke(ctx, Msg_SetDenomMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetBeforeSendHook(ctx context.Context, in *MsgSetBeforeSendHook, opts ...grpc.CallOption) (*MsgSetBeforeSendHookResponse, error) {
	out := new(MsgSetBeforeSendHookResponse)
	err := c.cc.Invoke(ctx, Msg_SetBeforeSendHook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ForceTransfer(ctx context.Context, in *MsgForceTransfer, opts ...grpc.CallOption) (*MsgForceTransferResponse, error) {
	out := new(MsgForceTransferResponse)
	err := c.cc.Invoke(ctx, Msg_ForceTransfer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	CreateDenom(context.Context, *MsgCreateDenom) (*MsgCreateDenomResponse, error)
	Mint(context.Context, *MsgMint) (*MsgMintResponse, error)
	Burn(context.Context, *MsgBurn) (*MsgBurnResponse, error)
	ChangeAdmin(context.Context, *MsgChangeAdmin) (*MsgChangeAdminResponse, error)
	SetDenomMetadata(context.Context, *MsgSetDenomMetadata) (*MsgSetDenomMetadataResponse, error)
	SetBeforeSendHook(context.Context, *MsgSetBeforeSendHook) (*MsgSetBeforeSendHookResponse, error)
	ForceTransfer(context.Context, *MsgForceTransfer) (*MsgForceTransferResponse, error)
	// UpdateParams defines an operation for updating the x/tokenfactory module
	// parameters.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) CreateDenom(context.Context, *MsgCreateDenom) (*MsgCreateDenomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDenom not implemented")
}
func (UnimplementedMsgServer) Mint(context.Context, *MsgMint) (*MsgMintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mint not implemented")
}
func (UnimplementedMsgServer) Burn(context.Context, *MsgBurn) (*MsgBurnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Burn not implemented")
}
func (UnimplementedMsgServer) ChangeAdmin(context.Context, *MsgChangeAdmin) (*MsgChangeAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAdmin not implemented")
}
func (UnimplementedMsgServer) SetDenomMetadata(context.Context, *MsgSetDenomMetadata) (*MsgSetDenomMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDenomMetadata not implemented")
}
func (UnimplementedMsgServer) SetBeforeSendHook(context.Context, *MsgSetBeforeSendHook) (*MsgSetBeforeSendHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetBeforeSendHook not implemented")
}
func (UnimplementedMsgServer) ForceTransfer(context.Context, *MsgForceTransfer) (*MsgForceTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForceTransfer not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_CreateDenom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateDenom)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateDenom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateDenom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateDenom(ctx, req.(*MsgCreateDenom))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Mint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Mint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Mint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Mint(ctx, req.(*MsgMint))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Burn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBurn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Burn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Burn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Burn(ctx, req.(*MsgBurn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ChangeAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgChangeAdmin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ChangeAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ChangeAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ChangeAdmin(ctx, req.(*MsgChangeAdmin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetDenomMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetDenomMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetDenomMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetDenomMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetDenomMetadata(ctx, req.(*MsgSetDenomMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetBeforeSendHook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetBeforeSendHook)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetBeforeSendHook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetBeforeSendHook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetBeforeSendHook(ctx, req.(*MsgSetBeforeSendHook))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ForceTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgForceTransfer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ForceTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ForceTransfer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ForceTransfer(ctx, req.(*MsgForceTransfer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "miniwasm.tokenfactory.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDenom",
			Handler:    _Msg_CreateDenom_Handler,
		},
		{
			MethodName: "Mint",
			Handler:    _Msg_Mint_Handler,
		},
		{
			MethodName: "Burn",
			Handler:    _Msg_Burn_Handler,
		},
		{
			MethodName: "ChangeAdmin",
			Handler:    _Msg_ChangeAdmin_Handler,
		},
		{
			MethodName: "SetDenomMetadata",
			Handler:    _Msg_SetDenomMetadata_Handler,
		},
		{
			MethodName: "SetBeforeSendHook",
			Handler:    _Msg_SetBeforeSendHook_Handler,
		},
		{
			MethodName: "ForceTransfer",
			Handler:    _Msg_ForceTransfer_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "miniwasm/tokenfactory/v1/tx.proto",
}
