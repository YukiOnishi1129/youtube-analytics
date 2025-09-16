package controller

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/presenter"
	inport "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/input"
	outpresenter "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/presenter"
	pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/authority/v1"
)

// Handler wires gRPC to use cases and presenter.
type Handler struct {
	uc        inport.AuthorityInputPort
	presenter outpresenter.AuthorityPresenter
	pb.UnimplementedAuthorityServiceServer
}

func NewHandler(uc inport.AuthorityInputPort, p outpresenter.AuthorityPresenter) *Handler {
	return &Handler{uc: uc, presenter: p}
}

func (h *Handler) GetAccount(ctx context.Context, _ *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	if err := h.uc.GetAccount(ctx); err != nil {
		return nil, err
	}
	// Cast to gRPC presenter to extract the built response
	if gp, ok := h.presenter.(*presenter.GRPCPresenter); ok && gp.AccountResp != nil {
		return gp.AccountResp, nil
	}
	return &pb.GetAccountResponse{}, nil
}

func (h *Handler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if err := h.uc.SignUp(ctx, req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}
	if gp, ok := h.presenter.(*presenter.GRPCPresenter); ok && gp.SignUpResp != nil {
		return gp.SignUpResp, nil
	}
	return &pb.SignUpResponse{}, nil
}

func (h *Handler) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	if err := h.uc.SignIn(ctx, req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}
	if gp, ok := h.presenter.(*presenter.GRPCPresenter); ok && gp.SignInResp != nil {
		return gp.SignInResp, nil
	}
	return &pb.SignInResponse{}, nil
}

func (h *Handler) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	if err := h.uc.SignOut(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}
	if gp, ok := h.presenter.(*presenter.GRPCPresenter); ok && gp.SignOutResp != nil {
		return gp.SignOutResp, nil
	}
	return &pb.SignOutResponse{}, nil
}

func (h *Handler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := h.uc.ResetPassword(ctx, req.GetEmail()); err != nil {
		return nil, err
	}
	if gp, ok := h.presenter.(*presenter.GRPCPresenter); ok && gp.ResetPassResp != nil {
		return gp.ResetPassResp, nil
	}
	return &pb.ResetPasswordResponse{}, nil
}
