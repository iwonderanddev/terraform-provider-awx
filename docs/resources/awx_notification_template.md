# Resource: awx_notification_template

Manages AWX `notification_templates` objects.

## Example Usage

```hcl
resource "awx_notification_template" "example" {
  name = "example"
  notification_type = "example"
  organization_id = awx_organization.example.id
  messages = { key = "value" }
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `messages` (Optional, Computed) Optional custom messages as a Terraform object.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `notification_configuration` (Optional, Computed, Sensitive) Notification transport configuration as a write-only sensitive Terraform object.
- `notification_type` (Required) * `awssns` - AWS SNS
  - `email` - Email
  - `grafana` - Grafana
  - `irc` - IRC
  - `mattermost` - Mattermost
  - `pagerduty` - Pagerduty
  - `rocketchat` - Rocket.Chat
  - `slack` - Slack
  - `twilio` - Twilio
  - `webhook` - Webhook
- `organization_id` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_notification_template.example 42
```
