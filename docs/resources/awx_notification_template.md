# Resource: awx_notification_template

Manages AWX `notification_templates` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_notification_template" "example" {
  name = "example"
  notification_type = "example"
  organization_id = awx_organization.example.id
  messages = { key = "value" }
}
```

## Schema

### Qualifiers

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Computed`: AWX sets the value during create or refresh.
- `Read-Only`: Cannot be set in configuration; Terraform records the value AWX returns.
- `Sensitive`: Terraform redacts the value in normal CLI output.
- `Write-Only`: Sent to AWX during create/update and not read back.

### Required

- `name` (String, Required) AWX value stored in `name`.
- `notification_type` (String, Required) Allowed values:
  - `awssns` - AWS SNS
  - `email` - Email
  - `grafana` - Grafana
  - `irc` - IRC
  - `mattermost` - Mattermost
  - `pagerduty` - Pagerduty
  - `rocketchat` - Rocket.Chat
  - `slack` - Slack
  - `twilio` - Twilio
  - `webhook` - Webhook
- `organization_id` (Number, Required) Numeric ID of the related AWX organization object.

### Optional

- `description` (String, Optional) AWX value stored in `description`.
- `messages` (Object, Optional, Computed) Optional custom messages as a Terraform object.
- `notification_configuration` (Object, Optional, Computed, Sensitive, Write-Only) Notification transport configuration as a write-only sensitive Terraform object.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_notification_template.example 42
```

## Further Reading

- [AWX Notifications](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/notifications.html)
