
package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Transaction struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        float64                `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Category      string                 `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Description   string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Date          *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=date,proto3" json:"date,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_ledger_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Transaction) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Transaction) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Transaction) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Transaction) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Transaction) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Transaction) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *Transaction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type AddTransactionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        float64                `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Category      string                 `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Date          *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddTransactionRequest) Reset() {
	*x = AddTransactionRequest{}
	mi := &file_ledger_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddTransactionRequest) ProtoMessage() {}

func (x *AddTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*AddTransactionRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{1}
}

func (x *AddTransactionRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AddTransactionRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *AddTransactionRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *AddTransactionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *AddTransactionRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

type AddTransactionResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Transaction    *Transaction           `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	BudgetExceeded bool                   `protobuf:"varint,2,opt,name=budget_exceeded,json=budgetExceeded,proto3" json:"budget_exceeded,omitempty"` 
	BudgetWarning  string                 `protobuf:"bytes,3,opt,name=budget_warning,json=budgetWarning,proto3" json:"budget_warning,omitempty"`     
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *AddTransactionResponse) Reset() {
	*x = AddTransactionResponse{}
	mi := &file_ledger_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddTransactionResponse) ProtoMessage() {}

func (x *AddTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*AddTransactionResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{2}
}

func (x *AddTransactionResponse) GetTransaction() *Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

func (x *AddTransactionResponse) GetBudgetExceeded() bool {
	if x != nil {
		return x.BudgetExceeded
	}
	return false
}

func (x *AddTransactionResponse) GetBudgetWarning() string {
	if x != nil {
		return x.BudgetWarning
	}
	return ""
}

type GetTransactionsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	From          *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`         
	To            *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`             
	Category      string                 `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"` 
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTransactionsRequest) Reset() {
	*x = GetTransactionsRequest{}
	mi := &file_ledger_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTransactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionsRequest) ProtoMessage() {}

func (x *GetTransactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetTransactionsRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{3}
}

func (x *GetTransactionsRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetTransactionsRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *GetTransactionsRequest) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *GetTransactionsRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

type GetTransactionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Transactions  []*Transaction         `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTransactionsResponse) Reset() {
	*x = GetTransactionsResponse{}
	mi := &file_ledger_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTransactionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionsResponse) ProtoMessage() {}

func (x *GetTransactionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetTransactionsResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{4}
}

func (x *GetTransactionsResponse) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

type Budget struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Category      string                 `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	LimitAmount   float64                `protobuf:"fixed64,4,opt,name=limit_amount,json=limitAmount,proto3" json:"limit_amount,omitempty"`
	Period        string                 `protobuf:"bytes,5,opt,name=period,proto3" json:"period,omitempty"` 
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Budget) Reset() {
	*x = Budget{}
	mi := &file_ledger_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Budget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Budget) ProtoMessage() {}

func (x *Budget) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Budget) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{5}
}

func (x *Budget) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Budget) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Budget) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Budget) GetLimitAmount() float64 {
	if x != nil {
		return x.LimitAmount
	}
	return 0
}

func (x *Budget) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

type SetBudgetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Category      string                 `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	LimitAmount   float64                `protobuf:"fixed64,3,opt,name=limit_amount,json=limitAmount,proto3" json:"limit_amount,omitempty"`
	Period        string                 `protobuf:"bytes,4,opt,name=period,proto3" json:"period,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetBudgetRequest) Reset() {
	*x = SetBudgetRequest{}
	mi := &file_ledger_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetBudgetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetBudgetRequest) ProtoMessage() {}

func (x *SetBudgetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetBudgetRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{6}
}

func (x *SetBudgetRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SetBudgetRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *SetBudgetRequest) GetLimitAmount() float64 {
	if x != nil {
		return x.LimitAmount
	}
	return 0
}

func (x *SetBudgetRequest) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

type SetBudgetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Budget        *Budget                `protobuf:"bytes,1,opt,name=budget,proto3" json:"budget,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetBudgetResponse) Reset() {
	*x = SetBudgetResponse{}
	mi := &file_ledger_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetBudgetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetBudgetResponse) ProtoMessage() {}

func (x *SetBudgetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetBudgetResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{7}
}

func (x *SetBudgetResponse) GetBudget() *Budget {
	if x != nil {
		return x.Budget
	}
	return nil
}

type GetBudgetsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBudgetsRequest) Reset() {
	*x = GetBudgetsRequest{}
	mi := &file_ledger_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBudgetsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBudgetsRequest) ProtoMessage() {}

func (x *GetBudgetsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetBudgetsRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{8}
}

func (x *GetBudgetsRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetBudgetsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Budgets       []*Budget              `protobuf:"bytes,1,rep,name=budgets,proto3" json:"budgets,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBudgetsResponse) Reset() {
	*x = GetBudgetsResponse{}
	mi := &file_ledger_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBudgetsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBudgetsResponse) ProtoMessage() {}

func (x *GetBudgetsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetBudgetsResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{9}
}

func (x *GetBudgetsResponse) GetBudgets() []*Budget {
	if x != nil {
		return x.Budgets
	}
	return nil
}

type CategorySummary struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Category         string                 `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	Total            float64                `protobuf:"fixed64,2,opt,name=total,proto3" json:"total,omitempty"`
	BudgetLimit      float64                `protobuf:"fixed64,3,opt,name=budget_limit,json=budgetLimit,proto3" json:"budget_limit,omitempty"`                
	BudgetPercentage float64                `protobuf:"fixed64,4,opt,name=budget_percentage,json=budgetPercentage,proto3" json:"budget_percentage,omitempty"` 
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *CategorySummary) Reset() {
	*x = CategorySummary{}
	mi := &file_ledger_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CategorySummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategorySummary) ProtoMessage() {}

func (x *CategorySummary) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CategorySummary) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{10}
}

func (x *CategorySummary) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *CategorySummary) GetTotal() float64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *CategorySummary) GetBudgetLimit() float64 {
	if x != nil {
		return x.BudgetLimit
	}
	return 0
}

func (x *CategorySummary) GetBudgetPercentage() float64 {
	if x != nil {
		return x.BudgetPercentage
	}
	return 0
}

type GetReportRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	From          *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To            *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetReportRequest) Reset() {
	*x = GetReportRequest{}
	mi := &file_ledger_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportRequest) ProtoMessage() {}

func (x *GetReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetReportRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{11}
}

func (x *GetReportRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetReportRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *GetReportRequest) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type GetReportResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Categories    []*CategorySummary     `protobuf:"bytes,1,rep,name=categories,proto3" json:"categories,omitempty"`
	TotalExpenses float64                `protobuf:"fixed64,2,opt,name=total_expenses,json=totalExpenses,proto3" json:"total_expenses,omitempty"`
	From          *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	To            *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetReportResponse) Reset() {
	*x = GetReportResponse{}
	mi := &file_ledger_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportResponse) ProtoMessage() {}

func (x *GetReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetReportResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{12}
}

func (x *GetReportResponse) GetCategories() []*CategorySummary {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *GetReportResponse) GetTotalExpenses() float64 {
	if x != nil {
		return x.TotalExpenses
	}
	return 0
}

func (x *GetReportResponse) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *GetReportResponse) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type ImportCSVRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CsvData       []byte                 `protobuf:"bytes,2,opt,name=csv_data,json=csvData,proto3" json:"csv_data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ImportCSVRequest) Reset() {
	*x = ImportCSVRequest{}
	mi := &file_ledger_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ImportCSVRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportCSVRequest) ProtoMessage() {}

func (x *ImportCSVRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ImportCSVRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{13}
}

func (x *ImportCSVRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ImportCSVRequest) GetCsvData() []byte {
	if x != nil {
		return x.CsvData
	}
	return nil
}

type ImportCSVResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ImportedCount int32                  `protobuf:"varint,1,opt,name=imported_count,json=importedCount,proto3" json:"imported_count,omitempty"`
	SkippedCount  int32                  `protobuf:"varint,2,opt,name=skipped_count,json=skippedCount,proto3" json:"skipped_count,omitempty"`
	Errors        []string               `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ImportCSVResponse) Reset() {
	*x = ImportCSVResponse{}
	mi := &file_ledger_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ImportCSVResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportCSVResponse) ProtoMessage() {}

func (x *ImportCSVResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ImportCSVResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{14}
}

func (x *ImportCSVResponse) GetImportedCount() int32 {
	if x != nil {
		return x.ImportedCount
	}
	return 0
}

func (x *ImportCSVResponse) GetSkippedCount() int32 {
	if x != nil {
		return x.SkippedCount
	}
	return 0
}

func (x *ImportCSVResponse) GetErrors() []string {
	if x != nil {
		return x.Errors
	}
	return nil
}

type ExportCSVRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	From          *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To            *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExportCSVRequest) Reset() {
	*x = ExportCSVRequest{}
	mi := &file_ledger_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExportCSVRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportCSVRequest) ProtoMessage() {}

func (x *ExportCSVRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ExportCSVRequest) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{15}
}

func (x *ExportCSVRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ExportCSVRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *ExportCSVRequest) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type ExportCSVResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CsvData       []byte                 `protobuf:"bytes,1,opt,name=csv_data,json=csvData,proto3" json:"csv_data,omitempty"`
	RowsCount     int32                  `protobuf:"varint,2,opt,name=rows_count,json=rowsCount,proto3" json:"rows_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExportCSVResponse) Reset() {
	*x = ExportCSVResponse{}
	mi := &file_ledger_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExportCSVResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportCSVResponse) ProtoMessage() {}

func (x *ExportCSVResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ExportCSVResponse) Descriptor() ([]byte, []int) {
	return file_ledger_proto_rawDescGZIP(), []int{16}
}

func (x *ExportCSVResponse) GetCsvData() []byte {
	if x != nil {
		return x.CsvData
	}
	return nil
}

func (x *ExportCSVResponse) GetRowsCount() int32 {
	if x != nil {
		return x.RowsCount
	}
	return 0
}

var File_ledger_proto protoreflect.FileDescriptor

const file_ledger_proto_rawDesc = "" +
	"\n" +
	"\fledger.proto\x12\tledger.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\xf7\x01\n" +
	"\vTransaction\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\x12\x16\n" +
	"\x06amount\x18\x03 \x01(\x01R\x06amount\x12\x1a\n" +
	"\bcategory\x18\x04 \x01(\tR\bcategory\x12 \n" +
	"\vdescription\x18\x05 \x01(\tR\vdescription\x12.\n" +
	"\x04date\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\x04date\x129\n" +
	"\n" +
	"created_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\"\xb6\x01\n" +
	"\x15AddTransactionRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12\x16\n" +
	"\x06amount\x18\x02 \x01(\x01R\x06amount\x12\x1a\n" +
	"\bcategory\x18\x03 \x01(\tR\bcategory\x12 \n" +
	"\vdescription\x18\x04 \x01(\tR\vdescription\x12.\n" +
	"\x04date\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\x04date\"\xa2\x01\n" +
	"\x16AddTransactionResponse\x128\n" +
	"\vtransaction\x18\x01 \x01(\v2\x16.ledger.v1.TransactionR\vtransaction\x12'\n" +
	"\x0fbudget_exceeded\x18\x02 \x01(\bR\x0ebudgetExceeded\x12%\n" +
	"\x0ebudget_warning\x18\x03 \x01(\tR\rbudgetWarning\"\xa9\x01\n" +
	"\x16GetTransactionsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12.\n" +
	"\x04from\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x04from\x12*\n" +
	"\x02to\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x02to\x12\x1a\n" +
	"\bcategory\x18\x04 \x01(\tR\bcategory\"U\n" +
	"\x17GetTransactionsResponse\x12:\n" +
	"\ftransactions\x18\x01 \x03(\v2\x16.ledger.v1.TransactionR\ftransactions\"\x88\x01\n" +
	"\x06Budget\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\x12\x1a\n" +
	"\bcategory\x18\x03 \x01(\tR\bcategory\x12!\n" +
	"\flimit_amount\x18\x04 \x01(\x01R\vlimitAmount\x12\x16\n" +
	"\x06period\x18\x05 \x01(\tR\x06period\"\x82\x01\n" +
	"\x10SetBudgetRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12\x1a\n" +
	"\bcategory\x18\x02 \x01(\tR\bcategory\x12!\n" +
	"\flimit_amount\x18\x03 \x01(\x01R\vlimitAmount\x12\x16\n" +
	"\x06period\x18\x04 \x01(\tR\x06period\">\n" +
	"\x11SetBudgetResponse\x12)\n" +
	"\x06budget\x18\x01 \x01(\v2\x11.ledger.v1.BudgetR\x06budget\",\n" +
	"\x11GetBudgetsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\"A\n" +
	"\x12GetBudgetsResponse\x12+\n" +
	"\abudgets\x18\x01 \x03(\v2\x11.ledger.v1.BudgetR\abudgets\"\x93\x01\n" +
	"\x0fCategorySummary\x12\x1a\n" +
	"\bcategory\x18\x01 \x01(\tR\bcategory\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x01R\x05total\x12!\n" +
	"\fbudget_limit\x18\x03 \x01(\x01R\vbudgetLimit\x12+\n" +
	"\x11budget_percentage\x18\x04 \x01(\x01R\x10budgetPercentage\"\x87\x01\n" +
	"\x10GetReportRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12.\n" +
	"\x04from\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x04from\x12*\n" +
	"\x02to\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x02to\"\xd2\x01\n" +
	"\x11GetReportResponse\x12:\n" +
	"\n" +
	"categories\x18\x01 \x03(\v2\x1a.ledger.v1.CategorySummaryR\n" +
	"categories\x12%\n" +
	"\x0etotal_expenses\x18\x02 \x01(\x01R\rtotalExpenses\x12.\n" +
	"\x04from\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x04from\x12*\n" +
	"\x02to\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\x02to\"F\n" +
	"\x10ImportCSVRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12\x19\n" +
	"\bcsv_data\x18\x02 \x01(\fR\acsvData\"w\n" +
	"\x11ImportCSVResponse\x12%\n" +
	"\x0eimported_count\x18\x01 \x01(\x05R\rimportedCount\x12#\n" +
	"\rskipped_count\x18\x02 \x01(\x05R\fskippedCount\x12\x16\n" +
	"\x06errors\x18\x03 \x03(\tR\x06errors\"\x87\x01\n" +
	"\x10ExportCSVRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12.\n" +
	"\x04from\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x04from\x12*\n" +
	"\x02to\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x02to\"M\n" +
	"\x11ExportCSVResponse\x12\x19\n" +
	"\bcsv_data\x18\x01 \x01(\fR\acsvData\x12\x1d\n" +
	"\n" +
	"rows_count\x18\x02 \x01(\x05R\trowsCount2\xab\x04\n" +
	"\rLedgerService\x12U\n" +
	"\x0eAddTransaction\x12 .ledger.v1.AddTransactionRequest\x1a!.ledger.v1.AddTransactionResponse\x12X\n" +
	"\x0fGetTransactions\x12!.ledger.v1.GetTransactionsRequest\x1a\".ledger.v1.GetTransactionsResponse\x12F\n" +
	"\tSetBudget\x12\x1b.ledger.v1.SetBudgetRequest\x1a\x1c.ledger.v1.SetBudgetResponse\x12I\n" +
	"\n" +
	"GetBudgets\x12\x1c.ledger.v1.GetBudgetsRequest\x1a\x1d.ledger.v1.GetBudgetsResponse\x12F\n" +
	"\tGetReport\x12\x1b.ledger.v1.GetReportRequest\x1a\x1c.ledger.v1.GetReportResponse\x12F\n" +
	"\tImportCSV\x12\x1b.ledger.v1.ImportCSVRequest\x1a\x1c.ledger.v1.ImportCSVResponse\x12F\n" +
	"\tExportCSV\x12\x1b.ledger.v1.ExportCSVRequest\x1a\x1c.ledger.v1.ExportCSVResponseB8Z6github.com/mikhailmogilnikov/go/final/pkg/pb/ledger/v1b\x06proto3"

var (
	file_ledger_proto_rawDescOnce sync.Once
	file_ledger_proto_rawDescData []byte
)

func file_ledger_proto_rawDescGZIP() []byte {
	file_ledger_proto_rawDescOnce.Do(func() {
		file_ledger_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_ledger_proto_rawDesc), len(file_ledger_proto_rawDesc)))
	})
	return file_ledger_proto_rawDescData
}

var file_ledger_proto_msgTypes = make([]protoimpl.MessageInfo, 17)
var file_ledger_proto_goTypes = []any{
	(*Transaction)(nil),             
	(*AddTransactionRequest)(nil),   
	(*AddTransactionResponse)(nil),  
	(*GetTransactionsRequest)(nil),  
	(*GetTransactionsResponse)(nil), 
	(*Budget)(nil),                  
	(*SetBudgetRequest)(nil),        
	(*SetBudgetResponse)(nil),       
	(*GetBudgetsRequest)(nil),       
	(*GetBudgetsResponse)(nil),      
	(*CategorySummary)(nil),         
	(*GetReportRequest)(nil),        
	(*GetReportResponse)(nil),       
	(*ImportCSVRequest)(nil),        
	(*ImportCSVResponse)(nil),       
	(*ExportCSVRequest)(nil),        
	(*ExportCSVResponse)(nil),       
	(*timestamppb.Timestamp)(nil),   
}
var file_ledger_proto_depIdxs = []int32{
	17, 
	17, 
	17, 
	0,  
	17, 
	17, 
	0,  
	5,  
	5,  
	17, 
	17, 
	10, 
	17, 
	17, 
	17, 
	17, 
	1,  
	3,  
	6,  
	8,  
	11, 
	13, 
	15, 
	2,  
	4,  
	7,  
	9,  
	12, 
	14, 
	16, 
	23, 
	16, 
	16, 
	16, 
	0,  
}

func init() { file_ledger_proto_init() }
func file_ledger_proto_init() {
	if File_ledger_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_ledger_proto_rawDesc), len(file_ledger_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   17,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ledger_proto_goTypes,
		DependencyIndexes: file_ledger_proto_depIdxs,
		MessageInfos:      file_ledger_proto_msgTypes,
	}.Build()
	File_ledger_proto = out.File
	file_ledger_proto_goTypes = nil
	file_ledger_proto_depIdxs = nil
}
