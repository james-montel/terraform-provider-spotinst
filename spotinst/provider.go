package spotinst

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			string(commons.ProviderToken): {
				Type:     schema.TypeString,
				Optional: true,
				//DefaultFunc: schema.EnvDefaultFunc(credentials.EnvCredentialsVarToken, ""),
				Description: "Spotinst Personal API Access Token",
			},

			string(commons.ProviderAccount): {
				Type:     schema.TypeString,
				Optional: true,
				//DefaultFunc: schema.EnvDefaultFunc(credentials.EnvCredentialsVarAccount, ""),
				Description: "Spotinst Account ID",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			// Elastigroup.
			string(commons.ElastigroupAWSResourceName):          resourceSpotinstElastigroupAWS(),
			string(commons.GroupAWSResourceName):                resourceSpotinstElastigroupAWS(),
			string(commons.ElastigroupGCPResourceName):          resourceSpotinstElastigroupGCP(),
			string(commons.ElastigroupGKEResourceName):          resourceSpotinstElastigroupGKE(),
			string(commons.ElastigroupAWSBeanstalkResourceName): resourceSpotinstElastigroupAWSBeanstalk(),
			string(commons.ElastigroupAzureResourceName):        resourceSpotinstElastigroupAzure(),
			string(commons.SubscriptionResourceName):            resourceSpotinstSubscription(),
			string(commons.MRScalerAWSResourceName):             resourceSpotinstMRScalerAWS(),

			// Ocean.
			string(commons.OceanAWSResourceName):                 resourceSpotinstOceanAWS(),
			string(commons.OceanAWSLaunchSpecResourceName):       resourceSpotinstOceanAWSLaunchSpec(),
			string(commons.OceanGKEImportResourceName):           resourceSpotinstOceanGKEImport(),
			string(commons.OceanGKELaunchSpecResourceName):       resourceSpotinstOceanGKELaunchSpec(),
			string(commons.OceanGKELaunchSpecImportResourceName): resourceSpotinstOceanGKELaunchSpecImport(),
			string(commons.OceanECSResourceName):                 resourceSpotinstOceanECS(),
			string(commons.OceanECSLaunchSpecResourceName):       resourceSpotinstOceanECSLaunchSpec(),

			// Multai.
			string(commons.MultaiBalancerResourceName):    resourceSpotinstMultaiBalancer(),
			string(commons.MultaiDeploymentResourceName):  resourceSpotinstMultaiDeployment(),
			string(commons.MultaiListenerResourceName):    resourceSpotinstMultaiListener(),
			string(commons.MultaiRoutingRuleResourceName): resourceSpotinstMultaiRoutingRule(),
			string(commons.MultaiTargetResourceName):      resourceSpotinstMultaiTarget(),
			string(commons.MultaiTargetSetResourceName):   resourceSpotinstMultaiTargetSet(),

			// Managed Instance.
			string(commons.ManagedInstanceAWSResourceName): resourceSpotinstMangedInstanceAWS(),

			// HealthCheck
			string(commons.HealthCheckResourceName): resourceSpotinstHealthCheck(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:   d.Get(string(commons.ProviderToken)).(string),
		Account: d.Get(string(commons.ProviderAccount)).(string),
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config.Client()
}
