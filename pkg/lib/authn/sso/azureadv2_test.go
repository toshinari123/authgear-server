package sso

import (
	"net/http"
	"testing"

	"github.com/authgear/authgear-server/pkg/lib/config"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestAzureadv2Impl(t *testing.T) {
	Convey("Azureadv2Impl", t, func() {
		gock.InterceptClient(http.DefaultClient)
		defer gock.Off()

		g := &Azureadv2Impl{
			ProviderConfig: config.OAuthSSOProviderConfig{
				ClientID: "client_id",
				Type:     config.OAuthSSOProviderTypeAzureADv2,
				Tenant:   "common",
			},
		}

		gock.New("https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration").
			Reply(200).
			BodyString(`
{
  "token_endpoint": "https://login.microsoftonline.com/common/oauth2/v2.0/token",
  "token_endpoint_auth_methods_supported": [
    "client_secret_post",
    "private_key_jwt",
    "client_secret_basic"
  ],
  "jwks_uri": "https://login.microsoftonline.com/common/discovery/v2.0/keys",
  "response_modes_supported": [
    "query",
    "fragment",
    "form_post"
  ],
  "subject_types_supported": [
    "pairwise"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "response_types_supported": [
    "code",
    "id_token",
    "code id_token",
    "id_token token"
  ],
  "scopes_supported": [
    "openid",
    "profile",
    "email",
    "offline_access"
  ],
  "issuer": "https://login.microsoftonline.com/{tenantid}/v2.0",
  "request_uri_parameter_supported": false,
  "userinfo_endpoint": "https://graph.microsoft.com/oidc/userinfo",
  "authorization_endpoint": "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
  "device_authorization_endpoint": "https://login.microsoftonline.com/common/oauth2/v2.0/devicecode",
  "http_logout_supported": true,
  "frontchannel_logout_supported": true,
  "end_session_endpoint": "https://login.microsoftonline.com/common/oauth2/v2.0/logout",
  "claims_supported": [
    "sub",
    "iss",
    "cloud_instance_name",
    "cloud_instance_host_name",
    "cloud_graph_host_name",
    "msgraph_host",
    "aud",
    "exp",
    "iat",
    "auth_time",
    "acr",
    "nonce",
    "preferred_username",
    "name",
    "tid",
    "ver",
    "at_hash",
    "c_hash",
    "email"
  ],
  "kerberos_endpoint": "https://login.microsoftonline.com/common/kerberos",
  "tenant_region_scope": null,
  "cloud_instance_name": "microsoftonline.com",
  "cloud_graph_host_name": "graph.windows.net",
  "msgraph_host": "graph.microsoft.com",
  "rbac_url": "https://pas.windows.net"
}
			`)
		defer func() { gock.Flush() }()

		u, err := g.GetAuthURL(GetAuthURLParam{
			RedirectURI:  "https://localhost/",
			ResponseMode: ResponseModeFormPost,
			Nonce:        "nonce",
			State:        "state",
			Prompt:       []string{"login"},
		})
		So(err, ShouldBeNil)
		So(u, ShouldEqual, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=client_id&nonce=nonce&prompt=login&redirect_uri=https%3A%2F%2Flocalhost%2F&response_mode=form_post&response_type=code&scope=openid+profile+email&state=state")
	})
}