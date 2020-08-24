package user

import (
	"time"

	"github.com/authgear/authgear-server/pkg/api/model"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
)

type IdentityService interface {
	ListByUser(userID string) ([]*identity.Info, error)
}

type VerificationService interface {
	IsUserVerified(identities []*identity.Info, userID string) (bool, error)
	IsVerified(identities []*identity.Info, authenticators []*authenticator.Info) bool
}

type Queries struct {
	Store        store
	Identities   IdentityService
	Verification VerificationService
}

func (p *Queries) Get(id string) (*model.User, error) {
	user, err := p.Store.Get(id)
	if err != nil {
		return nil, err
	}

	identities, err := p.Identities.ListByUser(id)
	if err != nil {
		return nil, err
	}

	isVerified, err := p.Verification.IsUserVerified(identities, id)
	if err != nil {
		return nil, err
	}

	return newUserModel(user, identities, isVerified), nil
}

func (p *Queries) GetManyRaw(ids []string) ([]*User, error) {
	return p.Store.GetByIDs(ids)
}

func (p *Queries) Count() (uint64, error) {
	return p.Store.Count()
}

func (p *Queries) QueryPage(after, before model.PageCursor, first, last *uint64) ([]model.PageItem, error) {
	users, err := p.Store.QueryPage(after, before, first, last)
	if err != nil {
		return nil, err
	}

	var models = make([]model.PageItem, len(users))
	for i, u := range users {
		pageKey := db.PageKey{
			Key: u.CreatedAt.Format(time.RFC3339Nano),
			ID:  u.ID,
		}
		cursor, err := pageKey.ToPageCursor()
		if err != nil {
			return nil, err
		}

		models[i] = model.PageItem{Value: u, Cursor: cursor}
	}
	return models, nil
}
