package ticket

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/heroticket/pkg/contracts/heroticket"
)

type Service interface {
	TbaByAddress(ctx context.Context, owner common.Address) (*common.Address, error)
	IsIssuedTicket(ctx context.Context, contractAddress common.Address) (bool, error)
	HasTicket(ctx context.Context, contractAddress, owner common.Address) (bool, error)
	IsWhitelisted(ctx context.Context, contractAddress, to common.Address) (bool, error)
	OnChainTicketInfo(ctx context.Context, contractAddress common.Address) (*OnchainTicketInfo, error)
	TicketsByOwner(ctx context.Context, owner common.Address) ([]common.Address, error)
	TokenBalanceOf(ctx context.Context, owner common.Address) (*big.Int, error)
	UpdateWhitelist(ctx context.Context, contractAddress, to common.Address) error
	CreateTBA(ctx context.Context, to common.Address, tokenURI string) (*heroticket.HeroticketTBACreated, error)
	IssueTicket(ctx context.Context, params IssueTicketParams) (*heroticket.HeroticketTicketIssued, error)
	BuyTicketByToken(ctx context.Context, contractAddress, buyerAddress common.Address) (*heroticket.HeroticketTicketSold, error)

	// TODO: repo에 저장하는 메서드 추가
	CreateTicketCollection(ctx context.Context, params CreateTicketCollectionParams) (*TicketCollection, error)
	GetOwnedNFT(ctx context.Context, owner common.Address) (OwnedNFT, error)
	FindTicketCollectionByContractAddress(ctx context.Context, contractAddress string) (*TicketCollection, error)
	FindTicketCollections(ctx context.Context, filter TicketCollectionFilter) ([]*TicketCollection, error)
}

type TicketService struct {
	client *ethclient.Client
	hero   *heroticket.Heroticket
	pvk    *ecdsa.PrivateKey
	repo   Repository
}

type OwnedNFT struct {
	Status string `json:"status"`
	NFTs   []NFT  `json:"nfts`
}

type NFT struct {
	tokenId      string `json:"token_id"`
	tokenAddress string `json:"token_address"`
	name         string `json:"name"`
	symbol       string `json:"symbol"`
	tokenUri     string `json:"token_uri"`
}

type OwnedNFTResponse struct {
	Status string `json:"status"`
	Result []struct {
		TokenId      string `json:"token_id"`
		TokenAddress string `json:"token_address"`
		Name         string `json:"name"`
		Symbol       string `json:"symbol"`
		TokenURI     string `json:"token_uri"`
	} `json:"result"`
}

func New(client *ethclient.Client, hero *heroticket.Heroticket, pvk *ecdsa.PrivateKey, repo Repository) Service {
	return &TicketService{
		client: client,
		hero:   hero,
		pvk:    pvk,
		repo:   repo,
	}
}

func (s *TicketService) TbaByAddress(ctx context.Context, owner common.Address) (*common.Address, error) {
	tba, err := s.hero.TbaAddress(&bind.CallOpts{Context: ctx}, owner)
	if err != nil {
		return nil, err
	}

	return &tba, nil

}

func (s *TicketService) IsIssuedTicket(ctx context.Context, contractAddress common.Address) (bool, error) {
	return s.hero.IssuedTicket(&bind.CallOpts{Context: ctx}, contractAddress)
}

func (s *TicketService) HasTicket(ctx context.Context, contractAddress, owner common.Address) (bool, error) {
	return s.hero.HasTicket(&bind.CallOpts{Context: ctx}, contractAddress, owner)
}

func (s *TicketService) IsWhitelisted(ctx context.Context, contractAddress, to common.Address) (bool, error) {
	return s.hero.IsWhiteListed(&bind.CallOpts{Context: ctx}, contractAddress, to)
}

func (s *TicketService) OnChainTicketInfo(ctx context.Context, contractAddress common.Address) (*OnchainTicketInfo, error) {
	issuer, remain, ethPrice, tokenPrice, saleStartAt, saleEndAt, err := s.hero.TicketInfo(&bind.CallOpts{Context: ctx}, contractAddress)
	if err != nil {
		return nil, err
	}

	return &OnchainTicketInfo{
		ContractAddress: contractAddress,
		Issuer:          issuer,
		Remaining:       remain,
		EthPrice:        ethPrice,
		TokenPrice:      tokenPrice,
		SaleStartAt:     saleStartAt,
		SaleEndAt:       saleEndAt,
	}, nil
}

func (s *TicketService) TicketsByOwner(ctx context.Context, owner common.Address) ([]common.Address, error) {
	ownedTickets, err := s.hero.TicketsByOwner(&bind.CallOpts{Context: ctx}, owner)
	if err != nil {
		return nil, err
	}

	// TODO: repo에 저장된 콜렉션 정보와 합하여 반환

	return ownedTickets, nil
}

func (s *TicketService) TokenBalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return s.hero.TokenBalanceOf(&bind.CallOpts{Context: ctx}, owner)
}

func (s *TicketService) UpdateWhitelist(ctx context.Context, contractAddress, to common.Address) error {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return err
	}

	tx, err := s.hero.UpdateWhiteList(auth, contractAddress, to)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, s.client, tx)
	return err
}

func (s *TicketService) CreateTBA(ctx context.Context, to common.Address, tokenURI string) (*heroticket.HeroticketTBACreated, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.hero.CreateTBA(auth, to, tokenURI)
	if err != nil {
		return nil, err
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, err
	}

	// 이 방법 안되면 filter 사용이나 tbaAddress를 직접 받아서 사용
	var tbaCreated *heroticket.HeroticketTBACreated

	for _, log := range receipt.Logs {

		tbaCreated, err = s.hero.ParseTBACreated(*log)
		if err == nil {
			break
		}
	}

	if tbaCreated == nil {
		return nil, err
	}

	return tbaCreated, nil
}

func (s *TicketService) IssueTicket(ctx context.Context, params IssueTicketParams) (*heroticket.HeroticketTicketIssued, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.hero.IssueTicket(auth, params.TicketName, params.TicketSymbol, params.TicketUri, params.Issuer,
		params.TicketAmount, params.TicketEthPrice, params.TicketTokenPrice, params.SaleDuration)
	if err != nil {
		return nil, err
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, err
	}

	var ticketIssued *heroticket.HeroticketTicketIssued

	for _, log := range receipt.Logs {
		ticketIssued, err = s.hero.ParseTicketIssued(*log)
		if err == nil {
			break
		}
	}

	if ticketIssued == nil {
		return nil, err
	}

	return ticketIssued, nil
}

func (s *TicketService) BuyTicketByToken(ctx context.Context, contractAddress, buyerAddress common.Address) (*heroticket.HeroticketTicketSold, error) {
	auth, err := s.txOpts(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.hero.BuyTicketByToken(auth, contractAddress, buyerAddress)
	if err != nil {
		return nil, err
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, err
	}

	var ticketSold *heroticket.HeroticketTicketSold

	for _, log := range receipt.Logs {
		ticketSold, err = s.hero.ParseTicketSold(*log)
		if err == nil {
			break
		}
	}

	if ticketSold == nil {
		return nil, err
	}

	return ticketSold, nil
}

func (s *TicketService) GetOwnedNFT(ctx context.Context, owner common.Address) (OwnedNFT, error) {
	tbaAddress, err := s.hero.TbaAddress(&bind.CallOpts{Context: ctx}, owner)
	if err != nil {
		return OwnedNFT{}, err
	}

	url := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/%s/nft?chain=mumbai&format=decimal&media_items=false", tbaAddress.String())

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjEzNGNhN2Y2LTNkODMtNGIxZC1iOGIwLWE0NmVhMTllNmM4NiIsIm9yZ0lkIjoiMzY2NDMzIiwidXNlcklkIjoiMzc2NTk0IiwidHlwZUlkIjoiYTRmNGMzNTQtM2Y3Zi00YmU5LWI4ZjItZDkzOTM1MmJjZmVkIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE3MDEzNDMxOTYsImV4cCI6NDg1NzEwMzE5Nn0.KxfO8preWRqP1BvMTkW_FvPzH6cuQSTwzxz8DvBhZjc")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return OwnedNFT{}, err
	}
	defer res.Body.Close()

	var result OwnedNFTResponse
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&result); err != nil {
		return OwnedNFT{}, err
	}

	ownedNFT := OwnedNFT{
		Status: result.Status,
		NFTs:   make([]NFT, len(result.Result)),
	}

	for i, nft := range result.Result {
		ownedNFT.NFTs[i] = NFT{
			tokenId:      nft.TokenId,
			tokenAddress: nft.TokenAddress,
			name:         nft.Name,
			symbol:       nft.Symbol,
			tokenUri:     nft.TokenURI,
		}
	}

	return ownedNFT, nil
}

func (s *TicketService) txOpts(ctx context.Context) (*bind.TransactOpts, error) {
	address := crypto.PubkeyToAddress(s.pvk.PublicKey)

	nonce, err := s.client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := s.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.pvk, chainID)
	if err != nil {
		return nil, err
	}

	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 3000000

	return auth, nil
}

func (s *TicketService) CreateTicketCollection(ctx context.Context, params CreateTicketCollectionParams) (*TicketCollection, error) {
	return s.repo.CreateTicketCollection(ctx, params)
}

func (s *TicketService) FindTicketCollections(ctx context.Context, filter TicketCollectionFilter) ([]*TicketCollection, error) {
	return s.repo.FindTicketCollections(ctx, filter)
}

func (s *TicketService) FindTicketCollectionByContractAddress(ctx context.Context, contractAddress string) (*TicketCollection, error) {
	return s.repo.FindTicketCollectionByContractAddress(ctx, contractAddress)
}
