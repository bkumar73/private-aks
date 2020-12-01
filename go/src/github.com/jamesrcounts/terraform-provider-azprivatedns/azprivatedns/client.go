package azprivatedns

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
)

// Client contians AzureRM clients used by this provider.
type Client struct {
	StopContext context.Context

	RouteTablesClient    *network.RouteTablesClient
	SecurityGroupsClient *network.SecurityGroupsClient
	PrivateZonesClient   *privatedns.PrivateZonesClient
}

// Build creates an initialized Client struct.
func Build(config *authentication.Config) (*Client, error) {

	if config == nil {
		return nil, fmt.Errorf("error build config is nil: %v", config)
	}

	sender := sender.BuildSender("AzureRM")

	env, err := authentication.DetermineEnvironment(config.Environment)
	if err != nil {
		return nil, fmt.Errorf("error determining environment: %v", err)
	}

	oauthConfig, err := config.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error building oauth config: %v", err)
	}

	auth, err := config.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, fmt.Errorf("error retrieving auth token: %v", err)
	}

	subscriptionID := config.SubscriptionID
	routeTablesClient := network.NewRouteTablesClient(subscriptionID)
	securityGroupsClient := network.NewSecurityGroupsClient(subscriptionID)
	zonesClient := privatedns.NewPrivateZonesClient(subscriptionID)

	routeTablesClient.Authorizer = auth
	securityGroupsClient.Authorizer = auth
	zonesClient.Authorizer = auth

	return &Client{
		RouteTablesClient:    &routeTablesClient,
		SecurityGroupsClient: &securityGroupsClient,
		PrivateZonesClient:   &zonesClient,
	}, nil
}
