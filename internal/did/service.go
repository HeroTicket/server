package did

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/heroticket/internal/cache"
	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	"github.com/iden3/iden3comm/v2/protocol"
)

type Service interface {
	LoginRequest(ctx context.Context, id, audience, callbackUrl string) (protocol.AuthorizationRequestMessage, error)
	LoginCallback(ctx context.Context, id, token string) (*protocol.AuthorizationResponseMessage, error)
	// TODO: add more methods
	// CreateCredential()
	// VerifyCredential()
	// VerifyCredentialCallback()
	// RevokeCredential()

	CreateVerifier(ctx context.Context, v *Verifier) (*Verifier, error)
	DeleteVerifier(ctx context.Context, id string) error
	FindVerifierByID(ctx context.Context, id string) (*Verifier, error)
	FindMatchedVerifier(ctx context.Context, didStr, walletAddress, contractAddress string) (*Verifier, error)
}

type didService struct {
	rpcUrl string

	repo         Repository
	requestCache cache.Cache
}

func New(repo Repository, requestCache cache.Cache, rpcUrl string) Service {
	svc := &didService{
		repo:         repo,
		requestCache: requestCache,
		rpcUrl:       rpcUrl,
	}

	return svc
}

func (s *didService) LoginRequest(ctx context.Context, id string, audience string, callbackUrl string) (protocol.AuthorizationRequestMessage, error) {
	var request protocol.AuthorizationRequestMessage = auth.CreateAuthorizationRequestWithMessage(
		"Login to Hero Ticket",
		"Scan the QR code to login to Hero Ticket",
		audience,
		callbackUrl,
	)

	request.ID = id
	request.ThreadID = id

	err := s.requestCache.Set(ctx, id, request, DefaultCacheExpiry)
	if err != nil {
		return protocol.AuthorizationRequestMessage{}, err
	}

	return request, nil
}

func (s *didService) LoginCallback(ctx context.Context, id string, token string) (*protocol.AuthorizationResponseMessage, error) {
	var request protocol.AuthorizationRequestMessage

	err := s.requestCache.Get(ctx, id, &request)
	if err != nil {
		return nil, err
	}

	ipfsUrl := "https://ipfs.io"
	contractAddress := "134B1BE34911E39A8397ec6289782989729807a4"
	resolverPrefix := "polygon:mumbai"
	ketDir := "./keys"

	var verificationKeyLoader = &loaders.FSKeyLoader{
		Dir: ketDir,
	}

	resolver := state.ETHResolver{
		RPCUrl:          s.rpcUrl,
		ContractAddress: common.HexToAddress(contractAddress),
	}

	resolvers := map[string]pubsignals.StateResolver{
		resolverPrefix: &resolver,
	}

	verifier, err := auth.NewVerifier(
		verificationKeyLoader,
		resolvers,
		auth.WithIPFSGateway(ipfsUrl),
	)
	if err != nil {
		return nil, err
	}

	response, err := verifier.FullVerify(
		ctx,
		token,
		request,
		pubsignals.WithAcceptedProofGenerationDelay(time.Minute*5),
	)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.requestCache.Delete(ctx, id)
	}()

	return response, nil
}

func (s *didService) CreateVerifier(ctx context.Context, v *Verifier) (*Verifier, error) {
	return s.repo.CreateVerifier(ctx, v)
}

func (s *didService) DeleteVerifier(ctx context.Context, id string) error {
	return s.repo.DeleteVerifier(ctx, id)
}

func (s *didService) FindVerifierByID(ctx context.Context, id string) (*Verifier, error) {
	return s.repo.FindVerifierByID(ctx, id)
}

func (s *didService) FindMatchedVerifier(ctx context.Context, didStr, walletAddress, contractAddress string) (*Verifier, error) {
	return s.repo.FindMatchedVerifier(ctx, didStr, walletAddress, contractAddress)
}
