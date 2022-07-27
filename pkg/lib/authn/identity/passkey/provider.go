package passkey

import (
	"sort"

	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/webauthn"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/uuid"
)

type Provider struct {
	Store *Store
	Clock clock.Clock
}

func (p *Provider) List(userID string) ([]*identity.Passkey, error) {
	is, err := p.Store.List(userID)
	if err != nil {
		return nil, err
	}

	sortIdentities(is)
	return is, nil
}

func (p *Provider) ListByClaim(name string, value string) ([]*identity.Passkey, error) {
	is, err := p.Store.ListByClaim(name, value)
	if err != nil {
		return nil, err
	}

	sortIdentities(is)
	return is, nil
}

func (p *Provider) Get(userID, id string) (*identity.Passkey, error) {
	return p.Store.Get(userID, id)
}

func (p *Provider) GetByAssertionResponse(assertionResponse []byte) (*identity.Passkey, error) {
	credentialID, _, err := webauthn.ParseAssertionResponse(assertionResponse)
	if err != nil {
		return nil, err
	}
	return p.Store.GetByCredentialID(credentialID)
}

func (p *Provider) GetMany(ids []string) ([]*identity.Passkey, error) {
	return p.Store.GetMany(ids)
}

func (p *Provider) New(
	userID string,
	creationOptions *webauthn.CreationOptions,
	attestationResponse []byte,
) (*identity.Passkey, error) {
	credentialID, _, err := webauthn.VerifyAttestationResponse(attestationResponse)
	if err != nil {
		return nil, err
	}

	i := &identity.Passkey{
		ID:                  uuid.New(),
		UserID:              userID,
		CredentialID:        credentialID,
		CreationOptions:     creationOptions,
		AttestationResponse: attestationResponse,
	}
	return i, nil
}

func (p *Provider) Create(i *identity.Passkey) error {
	now := p.Clock.NowUTC()
	i.CreatedAt = now
	i.UpdatedAt = now
	return p.Store.Create(i)
}

func (p *Provider) Delete(i *identity.Passkey) error {
	return p.Store.Delete(i)
}

func sortIdentities(is []*identity.Passkey) {
	sort.Slice(is, func(i, j int) bool {
		return is[i].CreatedAt.Before(is[j].CreatedAt)
	})
}
