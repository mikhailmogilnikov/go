
package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	LedgerService_AddTransaction_FullMethodName  = "/ledger.v1.LedgerService/AddTransaction"
	LedgerService_GetTransactions_FullMethodName = "/ledger.v1.LedgerService/GetTransactions"
	LedgerService_SetBudget_FullMethodName       = "/ledger.v1.LedgerService/SetBudget"
	LedgerService_GetBudgets_FullMethodName      = "/ledger.v1.LedgerService/GetBudgets"
	LedgerService_GetReport_FullMethodName       = "/ledger.v1.LedgerService/GetReport"
	LedgerService_ImportCSV_FullMethodName       = "/ledger.v1.LedgerService/ImportCSV"
	LedgerService_ExportCSV_FullMethodName       = "/ledger.v1.LedgerService/ExportCSV"
)

type LedgerServiceClient interface {
	AddTransaction(ctx context.Context, in *AddTransactionRequest, opts ...grpc.CallOption) (*AddTransactionResponse, error)
	GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error)
	SetBudget(ctx context.Context, in *SetBudgetRequest, opts ...grpc.CallOption) (*SetBudgetResponse, error)
	GetBudgets(ctx context.Context, in *GetBudgetsRequest, opts ...grpc.CallOption) (*GetBudgetsResponse, error)
	GetReport(ctx context.Context, in *GetReportRequest, opts ...grpc.CallOption) (*GetReportResponse, error)
	ImportCSV(ctx context.Context, in *ImportCSVRequest, opts ...grpc.CallOption) (*ImportCSVResponse, error)
	ExportCSV(ctx context.Context, in *ExportCSVRequest, opts ...grpc.CallOption) (*ExportCSVResponse, error)
}

type ledgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerServiceClient(cc grpc.ClientConnInterface) LedgerServiceClient {
	return &ledgerServiceClient{cc}
}

func (c *ledgerServiceClient) AddTransaction(ctx context.Context, in *AddTransactionRequest, opts ...grpc.CallOption) (*AddTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddTransactionResponse)
	err := c.cc.Invoke(ctx, LedgerService_AddTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTransactionsResponse)
	err := c.cc.Invoke(ctx, LedgerService_GetTransactions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) SetBudget(ctx context.Context, in *SetBudgetRequest, opts ...grpc.CallOption) (*SetBudgetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetBudgetResponse)
	err := c.cc.Invoke(ctx, LedgerService_SetBudget_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) GetBudgets(ctx context.Context, in *GetBudgetsRequest, opts ...grpc.CallOption) (*GetBudgetsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBudgetsResponse)
	err := c.cc.Invoke(ctx, LedgerService_GetBudgets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) GetReport(ctx context.Context, in *GetReportRequest, opts ...grpc.CallOption) (*GetReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReportResponse)
	err := c.cc.Invoke(ctx, LedgerService_GetReport_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) ImportCSV(ctx context.Context, in *ImportCSVRequest, opts ...grpc.CallOption) (*ImportCSVResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ImportCSVResponse)
	err := c.cc.Invoke(ctx, LedgerService_ImportCSV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) ExportCSV(ctx context.Context, in *ExportCSVRequest, opts ...grpc.CallOption) (*ExportCSVResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExportCSVResponse)
	err := c.cc.Invoke(ctx, LedgerService_ExportCSV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type LedgerServiceServer interface {
	AddTransaction(context.Context, *AddTransactionRequest) (*AddTransactionResponse, error)
	GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error)
	SetBudget(context.Context, *SetBudgetRequest) (*SetBudgetResponse, error)
	GetBudgets(context.Context, *GetBudgetsRequest) (*GetBudgetsResponse, error)
	GetReport(context.Context, *GetReportRequest) (*GetReportResponse, error)
	ImportCSV(context.Context, *ImportCSVRequest) (*ImportCSVResponse, error)
	ExportCSV(context.Context, *ExportCSVRequest) (*ExportCSVResponse, error)
	mustEmbedUnimplementedLedgerServiceServer()
}

type UnimplementedLedgerServiceServer struct{}

func (UnimplementedLedgerServiceServer) AddTransaction(context.Context, *AddTransactionRequest) (*AddTransactionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method AddTransaction not implemented")
}
func (UnimplementedLedgerServiceServer) GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetTransactions not implemented")
}
func (UnimplementedLedgerServiceServer) SetBudget(context.Context, *SetBudgetRequest) (*SetBudgetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method SetBudget not implemented")
}
func (UnimplementedLedgerServiceServer) GetBudgets(context.Context, *GetBudgetsRequest) (*GetBudgetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetBudgets not implemented")
}
func (UnimplementedLedgerServiceServer) GetReport(context.Context, *GetReportRequest) (*GetReportResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetReport not implemented")
}
func (UnimplementedLedgerServiceServer) ImportCSV(context.Context, *ImportCSVRequest) (*ImportCSVResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ImportCSV not implemented")
}
func (UnimplementedLedgerServiceServer) ExportCSV(context.Context, *ExportCSVRequest) (*ExportCSVResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ExportCSV not implemented")
}
func (UnimplementedLedgerServiceServer) mustEmbedUnimplementedLedgerServiceServer() {}
func (UnimplementedLedgerServiceServer) testEmbeddedByValue()                       {}

type UnsafeLedgerServiceServer interface {
	mustEmbedUnimplementedLedgerServiceServer()
}

func RegisterLedgerServiceServer(s grpc.ServiceRegistrar, srv LedgerServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LedgerService_ServiceDesc, srv)
}

func _LedgerService_AddTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).AddTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_AddTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).AddTransaction(ctx, req.(*AddTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_GetTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).GetTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_GetTransactions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).GetTransactions(ctx, req.(*GetTransactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_SetBudget_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetBudgetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).SetBudget(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_SetBudget_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).SetBudget(ctx, req.(*SetBudgetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_GetBudgets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBudgetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).GetBudgets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_GetBudgets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).GetBudgets(ctx, req.(*GetBudgetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_GetReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).GetReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_GetReport_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).GetReport(ctx, req.(*GetReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_ImportCSV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportCSVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).ImportCSV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_ImportCSV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).ImportCSV(ctx, req.(*ImportCSVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_ExportCSV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportCSVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).ExportCSV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_ExportCSV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).ExportCSV(ctx, req.(*ExportCSVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var LedgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ledger.v1.LedgerService",
	HandlerType: (*LedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTransaction",
			Handler:    _LedgerService_AddTransaction_Handler,
		},
		{
			MethodName: "GetTransactions",
			Handler:    _LedgerService_GetTransactions_Handler,
		},
		{
			MethodName: "SetBudget",
			Handler:    _LedgerService_SetBudget_Handler,
		},
		{
			MethodName: "GetBudgets",
			Handler:    _LedgerService_GetBudgets_Handler,
		},
		{
			MethodName: "GetReport",
			Handler:    _LedgerService_GetReport_Handler,
		},
		{
			MethodName: "ImportCSV",
			Handler:    _LedgerService_ImportCSV_Handler,
		},
		{
			MethodName: "ExportCSV",
			Handler:    _LedgerService_ExportCSV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ledger.proto",
}
