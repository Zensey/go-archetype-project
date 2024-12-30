package pow_service

import (
	"context"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

const (
	challengeTTL         = 30 * time.Second
	computeMaxIterations = 1 << 30
	hashAlgo             = "sha256"
)

type PoW struct {
	hasher              hash.Hasher
	pow                 *pow.POW
	challengeDifficulty int
	secret              string
}

func New(challengeDifficulty int, secret string) *PoW {
	hasher, err := hash.NewHasher(hashAlgo)
	if err != nil {
		panic(err)
	}
	powService := pow.New(hasher, pow.WithChallengeExpDuration(challengeTTL))

	return &PoW{
		hasher:              hasher,
		pow:                 powService,
		challengeDifficulty: challengeDifficulty,
		secret:              secret,
	}
}

func (p *PoW) GenerateChallenge(resource string) (string, error) {
	challenge, err := pow.InitHashcash(int32(p.challengeDifficulty), resource, pow.SignExt(p.secret, p.hasher))
	if err != nil {
		return "", err
	}
	return challenge.String(), nil
}

func (p *PoW) ComputeResponse(ctx context.Context, challenge *pow.Hashcach) (string, error) {
	response, err := p.pow.Compute(ctx, challenge, computeMaxIterations)
	if err != nil {
		return "", err
	}
	return response.String(), nil
}

func (p *PoW) VerifyResponse(response *pow.Hashcach, resource string) error {
	err := p.pow.Verify(response, resource)
	return err
}
