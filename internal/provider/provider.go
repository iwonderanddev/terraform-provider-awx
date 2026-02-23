package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/damien/terraform-provider-awx-iwd/internal/client"
	"github.com/damien/terraform-provider-awx-iwd/internal/manifest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*awxProvider)(nil)

type awxProvider struct {
	version    string
	catalog    *manifest.Catalog
	catalogErr error
}

type providerModel struct {
	Hostname              types.String `tfsdk:"hostname"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	InsecureSkipTLSVerify types.Bool   `tfsdk:"insecure_skip_tls_verify"`
	CACertPEM             types.String `tfsdk:"ca_cert_pem"`
	RequestTimeoutSeconds types.Int64  `tfsdk:"request_timeout_seconds"`
	RetryMaxAttempts      types.Int64  `tfsdk:"retry_max_attempts"`
	RetryBackoffMillis    types.Int64  `tfsdk:"retry_backoff_millis"`
}

type configuredProvider struct {
	client  *client.Client
	catalog *manifest.Catalog
}

// New returns the AWX provider implementation.
func New(version string) func() provider.Provider {
	catalog, err := manifest.Load()
	return func() provider.Provider {
		return &awxProvider{
			version:    version,
			catalog:    catalog,
			catalogErr: err,
		}
	}
}

func (p *awxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awx"
	resp.Version = p.version
}

func (p *awxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{
		Description: "Terraform provider for AWX API v2.",
		Attributes: map[string]providerschema.Attribute{
			"hostname": providerschema.StringAttribute{
				Description: "Base URL for AWX, for example https://awx.example.com.",
				Required:    true,
			},
			"username": providerschema.StringAttribute{
				Description: "HTTP Basic username for AWX authentication.",
				Required:    true,
			},
			"password": providerschema.StringAttribute{
				Description: "HTTP Basic password for AWX authentication.",
				Required:    true,
				Sensitive:   true,
			},
			"insecure_skip_tls_verify": providerschema.BoolAttribute{
				Description: "Skip TLS certificate verification for HTTPS connections.",
				Optional:    true,
			},
			"ca_cert_pem": providerschema.StringAttribute{
				Description: "Optional PEM-encoded CA bundle used to verify AWX server TLS certificates.",
				Optional:    true,
				Sensitive:   true,
			},
			"request_timeout_seconds": providerschema.Int64Attribute{
				Description: "HTTP request timeout in seconds.",
				Optional:    true,
			},
			"retry_max_attempts": providerschema.Int64Attribute{
				Description: "Maximum retry attempts for retryable API failures.",
				Optional:    true,
			},
			"retry_backoff_millis": providerschema.Int64Attribute{
				Description: "Initial retry backoff in milliseconds for retryable API failures.",
				Optional:    true,
			},
		},
	}
}

func (p *awxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.catalogErr != nil {
		resp.Diagnostics.AddError("Failed to load generated manifest assets", p.catalogErr.Error())
		return
	}

	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Hostname.IsUnknown() || config.Hostname.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("hostname"), "Missing AWX hostname", "Set hostname to the AWX server URL.")
	}
	if config.Username.IsUnknown() || config.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Missing username", "Set username for HTTP Basic auth.")
	}
	if config.Password.IsUnknown() || config.Password.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing password", "Set password for HTTP Basic auth.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	validateConfig(config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clientConfig := client.Config{
		BaseURL:               config.Hostname.ValueString(),
		Username:              config.Username.ValueString(),
		Password:              config.Password.ValueString(),
		InsecureSkipTLSVerify: !config.InsecureSkipTLSVerify.IsNull() && config.InsecureSkipTLSVerify.ValueBool(),
		CACertPEM:             strings.TrimSpace(config.CACertPEM.ValueString()),
		Timeout:               durationFromSeconds(config.RequestTimeoutSeconds, 30*time.Second),
		RetryMaxAttempts:      int(valueOrDefault(config.RetryMaxAttempts, 3)),
		RetryInitialBackoff:   time.Duration(valueOrDefault(config.RetryBackoffMillis, 500)) * time.Millisecond,
		UserAgent:             fmt.Sprintf("terraform-provider-awx-iwd/%s", p.version),
	}

	apiClient, err := client.New(clientConfig)
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure AWX API client", err.Error())
		return
	}

	if err := apiClient.Ping(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Unable to connect to AWX API",
			fmt.Sprintf("Connection validation failed against /api/v2/ with configured credentials: %s", err.Error()),
		)
		return
	}

	configured := &configuredProvider{client: apiClient, catalog: p.catalog}
	resp.ResourceData = configured
	resp.DataSourceData = configured
}

func (p *awxProvider) Resources(_ context.Context) []func() resource.Resource {
	resources := make([]func() resource.Resource, 0)
	if p.catalog == nil {
		return resources
	}

	for _, obj := range p.catalog.ManagedResourceObjects() {
		objCopy := obj
		resources = append(resources, func() resource.Resource {
			return NewObjectResource(objCopy)
		})
	}
	for _, rel := range p.catalog.Relationships {
		relCopy := rel
		resources = append(resources, func() resource.Resource {
			return NewRelationshipResource(relCopy)
		})
	}
	return resources
}

func (p *awxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	dataSources := make([]func() datasource.DataSource, 0)
	if p.catalog == nil {
		return dataSources
	}

	for _, obj := range p.catalog.ManagedDataSourceObjects() {
		objCopy := obj
		dataSources = append(dataSources, func() datasource.DataSource {
			return NewObjectDataSource(objCopy)
		})
	}
	return dataSources
}

func validateConfig(config providerModel, diags *diag.Diagnostics) {
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() {
		return
	}

	parsedURL, err := url.Parse(strings.TrimSpace(config.Hostname.ValueString()))
	if err != nil {
		diags.AddAttributeError(path.Root("hostname"), "Invalid AWX hostname", err.Error())
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		diags.AddAttributeError(path.Root("hostname"), "Invalid AWX hostname", "hostname must use http or https.")
	}
	if parsedURL.Host == "" {
		diags.AddAttributeError(path.Root("hostname"), "Invalid AWX hostname", "hostname must include a host.")
	}

	if !config.RequestTimeoutSeconds.IsNull() && !config.RequestTimeoutSeconds.IsUnknown() && config.RequestTimeoutSeconds.ValueInt64() <= 0 {
		diags.AddAttributeError(path.Root("request_timeout_seconds"), "Invalid timeout", "request_timeout_seconds must be > 0")
	}
	if !config.RetryMaxAttempts.IsNull() && !config.RetryMaxAttempts.IsUnknown() && config.RetryMaxAttempts.ValueInt64() <= 0 {
		diags.AddAttributeError(path.Root("retry_max_attempts"), "Invalid retry attempts", "retry_max_attempts must be > 0")
	}
	if !config.RetryBackoffMillis.IsNull() && !config.RetryBackoffMillis.IsUnknown() && config.RetryBackoffMillis.ValueInt64() <= 0 {
		diags.AddAttributeError(path.Root("retry_backoff_millis"), "Invalid retry backoff", "retry_backoff_millis must be > 0")
	}
}

func durationFromSeconds(value types.Int64, defaultDuration time.Duration) time.Duration {
	if value.IsNull() || value.IsUnknown() {
		return defaultDuration
	}
	return time.Duration(value.ValueInt64()) * time.Second
}

func valueOrDefault(value types.Int64, fallback int64) int64 {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}
	return value.ValueInt64()
}
