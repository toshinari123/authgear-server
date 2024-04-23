package declarative

import (
	authflow "github.com/authgear/authgear-server/pkg/lib/authenticationflow"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/authn/sso"
)

func init() {
	authflow.RegisterNode(&NodeUseAccountLinkingIdentification{})
}

type NodeUseAccountLinkingIdentification struct {
	Option   AccountLinkingIdentificationOption `json:"option,omitempty"`
	Identity *identity.Info                     `json:"identity,omitempty"`

	// oauth
	RedirectURI  string           `json:"redirect_uri,omitempty"`
	ResponseMode sso.ResponseMode `json:"response_mode,omitempty"`
}

var _ authflow.NodeSimple = &NodeUseAccountLinkingIdentification{}
var _ authflow.Milestone = &NodeUseAccountLinkingIdentification{}
var _ MilestoneUseAccountLinkingIdentification = &NodeUseAccountLinkingIdentification{}

func (*NodeUseAccountLinkingIdentification) Kind() string {
	return "NodeUseAccountLinkingIdentificationOption"
}

func (*NodeUseAccountLinkingIdentification) Milestone() {}
func (n *NodeUseAccountLinkingIdentification) MilestoneUseAccountLinkingIdentification() *identity.Info {
	return n.Identity
}
func (n *NodeUseAccountLinkingIdentification) MilestoneUseAccountLinkingIdentificationSelectedOption() AccountLinkingIdentificationOption {
	return n.Option
}
func (n *NodeUseAccountLinkingIdentification) MilestoneUseAccountLinkingIdentificationRedirectURI() string {
	return n.RedirectURI
}
func (n *NodeUseAccountLinkingIdentification) MilestoneUseAccountLinkingIdentificationResponseMode() sso.ResponseMode {
	return n.ResponseMode
}
