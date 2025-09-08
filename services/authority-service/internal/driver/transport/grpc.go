package transport

import (
	"log"
	"net"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/controller"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/presenter"
	adaptersec "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/security"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/driver/security"
	inport "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/input"
	outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/usecase"
    pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/proto/authority/v1"
	"google.golang.org/grpc"
)

// Bootstrap wires everything and starts gRPC server.
func Bootstrap(addr string,
	accountRepo outgateway.AccountRepository,
	idRepo outgateway.IdentityRepository,
	roleRepo outgateway.RoleRepository,
	verifier outgateway.TokenVerifier,
	idp outgateway.IdentityProvider,
	clock outgateway.Clock,
) error {
	// presenter instance
	p := &presenter.GRPCPresenter{}
	// claims provider via context
	claimsProvider := &adaptersec.ContextClaimsProvider{}

	// usecase
	var uc inport.AuthorityInputPort = usecase.NewAuthorityInteractor(
		accountRepo, idRepo, roleRepo, verifier, claimsProvider, idp, clock, p,
	)

	// handler
	h := controller.NewHandler(uc, p)

	// server with auth interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(security.UnaryAuthInterceptor(verifier)),
	)
	pb.RegisterAuthorityServiceServer(s, h)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("authority-service gRPC listening on %s", addr)
	return s.Serve(lis)
}
