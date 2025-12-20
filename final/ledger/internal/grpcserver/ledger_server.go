package grpcserver

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mikhailmogilnikov/go/final/ledger/internal/csv"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/domain"
	pb "github.com/mikhailmogilnikov/go/final/ledger/internal/pb/ledger/v1"
	"github.com/mikhailmogilnikov/go/final/ledger/internal/service"
)

// LedgerServer реализует gRPC сервер
type LedgerServer struct {
	pb.UnimplementedLedgerServiceServer
	ledgerService *service.LedgerService
}

// NewLedgerServer создаёт новый сервер
func NewLedgerServer(ledgerService *service.LedgerService) *LedgerServer {
	return &LedgerServer{ledgerService: ledgerService}
}

// AddTransaction добавляет транзакцию
func (s *LedgerServer) AddTransaction(ctx context.Context, req *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.GetCategory() == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}

	tx := &domain.Transaction{
		UserID:      req.GetUserId(),
		Amount:      req.GetAmount(),
		Category:    req.GetCategory(),
		Description: req.GetDescription(),
	}
	if req.GetDate() != nil {
		tx.Date = req.GetDate().AsTime()
	}

	budgetWarning, err := s.ledgerService.AddTransaction(ctx, tx)
	if err != nil {
		// Превышение бюджета - возвращаем конфликт (транзакция отклонена)
		if errors.Is(err, service.ErrBudgetExceeded) {
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to add transaction: %v", err)
	}

	return &pb.AddTransactionResponse{
		Transaction:    toProtoTransaction(tx),
		BudgetExceeded: false,
		BudgetWarning:  budgetWarning,
	}, nil
}

// GetTransactions возвращает транзакции
func (s *LedgerServer) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	transactions, err := s.ledgerService.GetTransactions(ctx, req.GetUserId(),
		timeFromProto(req.GetFrom()), timeFromProto(req.GetTo()), req.GetCategory())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	protoTxs := make([]*pb.Transaction, 0, len(transactions))
	for i := range transactions {
		protoTxs = append(protoTxs, toProtoTransaction(&transactions[i]))
	}

	return &pb.GetTransactionsResponse{
		Transactions: protoTxs,
	}, nil
}

// SetBudget устанавливает бюджет
func (s *LedgerServer) SetBudget(ctx context.Context, req *pb.SetBudgetRequest) (*pb.SetBudgetResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetCategory() == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}
	if req.GetLimitAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit_amount must be positive")
	}

	budget := &domain.Budget{
		UserID:      req.GetUserId(),
		Category:    req.GetCategory(),
		LimitAmount: req.GetLimitAmount(),
		Period:      req.GetPeriod(),
	}

	if err := s.ledgerService.SetBudget(ctx, budget); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set budget: %v", err)
	}

	return &pb.SetBudgetResponse{
		Budget: toProtoBudget(budget),
	}, nil
}

// GetBudgets возвращает бюджеты
func (s *LedgerServer) GetBudgets(ctx context.Context, req *pb.GetBudgetsRequest) (*pb.GetBudgetsResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	budgets, err := s.ledgerService.GetBudgets(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budgets: %v", err)
	}

	protoBudgets := make([]*pb.Budget, 0, len(budgets))
	for i := range budgets {
		protoBudgets = append(protoBudgets, toProtoBudget(&budgets[i]))
	}

	return &pb.GetBudgetsResponse{
		Budgets: protoBudgets,
	}, nil
}

// GetReport возвращает отчёт
func (s *LedgerServer) GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetFrom() == nil || req.GetTo() == nil {
		return nil, status.Error(codes.InvalidArgument, "from and to dates are required")
	}

	from := req.GetFrom().AsTime()
	to := req.GetTo().AsTime()

	categories, totalExpenses, err := s.ledgerService.GetReport(ctx, req.GetUserId(), from, to)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get report: %v", err)
	}

	protoCategories := make([]*pb.CategorySummary, 0, len(categories))
	for _, c := range categories {
		protoCategories = append(protoCategories, &pb.CategorySummary{
			Category:         c.Category,
			Total:            c.Total,
			BudgetLimit:      c.BudgetLimit,
			BudgetPercentage: c.BudgetPercentage,
		})
	}

	return &pb.GetReportResponse{
		Categories:    protoCategories,
		TotalExpenses: totalExpenses,
		From:          req.GetFrom(),
		To:            req.GetTo(),
	}, nil
}

// ImportCSV импортирует транзакции из CSV
func (s *LedgerServer) ImportCSV(ctx context.Context, req *pb.ImportCSVRequest) (*pb.ImportCSVResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.GetCsvData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "csv_data is required")
	}

	transactions, parseErrors, err := csv.ParseCSV(req.GetCsvData(), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse CSV: %v", err)
	}

	var importedCount int32
	var skippedCount int32
	var importErrors []string

	for i := range transactions {
		_, err := s.ledgerService.AddTransaction(ctx, &transactions[i])
		if err != nil {
			skippedCount++
			importErrors = append(importErrors, err.Error())
		} else {
			importedCount++
		}
	}

	importErrors = append(parseErrors, importErrors...)

	return &pb.ImportCSVResponse{
		ImportedCount: importedCount,
		SkippedCount:  skippedCount + int32(len(parseErrors)),
		Errors:        importErrors,
	}, nil
}

// ExportCSV экспортирует транзакции в CSV
func (s *LedgerServer) ExportCSV(ctx context.Context, req *pb.ExportCSVRequest) (*pb.ExportCSVResponse, error) {
	if req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	transactions, err := s.ledgerService.GetTransactions(ctx, req.GetUserId(),
		timeFromProto(req.GetFrom()), timeFromProto(req.GetTo()), "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	csvData, err := csv.GenerateCSV(transactions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate CSV: %v", err)
	}

	return &pb.ExportCSVResponse{
		CsvData:   csvData,
		RowsCount: int32(len(transactions)),
	}, nil
}

// Хелперы для конвертации

func toProtoTransaction(tx *domain.Transaction) *pb.Transaction {
	return &pb.Transaction{
		Id:          tx.ID,
		UserId:      tx.UserID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        timestamppb.New(tx.Date),
		CreatedAt:   timestamppb.New(tx.CreatedAt),
	}
}

func toProtoBudget(b *domain.Budget) *pb.Budget {
	return &pb.Budget{
		Id:          b.ID,
		UserId:      b.UserID,
		Category:    b.Category,
		LimitAmount: b.LimitAmount,
		Period:      b.Period,
	}
}

func timeFromProto(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
