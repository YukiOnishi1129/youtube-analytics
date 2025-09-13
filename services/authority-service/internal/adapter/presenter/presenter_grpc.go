package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
	outpresenter "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/presenter"
	pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/authority/v1"
)

// GRPCPresenter collects results to be returned by gRPC handlers.
// In a real setup, handlers would pass a reference to receive the built message.
type GRPCPresenter struct {
	AccountResp   *pb.GetAccountResponse
	SignUpResp    *pb.SignUpResponse
	SignInResp    *pb.SignInResponse
	SignOutResp   *pb.SignOutResponse
	ResetPassResp *pb.ResetPasswordResponse
}

var _ outpresenter.AuthorityPresenter = (*GRPCPresenter)(nil)

func (p *GRPCPresenter) PresentGetAccount(a *domain.Account) error {
	p.AccountResp = &pb.GetAccountResponse{Account: toPBAccount(a)}
	return nil
}

func (p *GRPCPresenter) PresentSignUp(a *domain.Account, idToken, refreshToken string) error {
	p.SignUpResp = &pb.SignUpResponse{
		Account:      toPBAccount(a),
		IdToken:      idToken,
		RefreshToken: refreshToken,
	}
	return nil
}

func (p *GRPCPresenter) PresentSignIn(idToken, refreshToken string, expiresIn int64) error {
	p.SignInResp = &pb.SignInResponse{
		IdToken:      idToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	return nil
}

func (p *GRPCPresenter) PresentSignOut(success bool) error {
	p.SignOutResp = &pb.SignOutResponse{Success: success}
	return nil
}

func (p *GRPCPresenter) PresentResetPassword(emailSent bool) error {
	p.ResetPassResp = &pb.ResetPasswordResponse{EmailSent: emailSent}
	return nil
}

func toPBAccount(a *domain.Account) *pb.Account {
	var roles []string
	for _, r := range a.Roles {
		roles = append(roles, r.Name)
	}
	lastLogin := ""
	if a.LastLoginAt != nil {
		lastLogin = a.LastLoginAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return &pb.Account{
		Id:            a.ID,
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
		DisplayName:   a.DisplayName,
		PhotoUrl:      a.PhotoURL,
		IsActive:      a.IsActive,
		LastLoginAt:   lastLogin,
		Roles:         roles,
	}
}
