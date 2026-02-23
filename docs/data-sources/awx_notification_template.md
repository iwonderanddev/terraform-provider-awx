# Data Source: awx_notification_template

Reads AWX `notification_templates` objects.

## Example Usage

```hcl
data "awx_notification_template" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) AWX value stored in `description`.
- `messages` (Object, Read-Only) Optional custom messages as a Terraform object.
- `name` (String, Read-Only) AWX value stored in `name`.
- `notification_configuration` (Object, Read-Only, Sensitive) Notification transport configuration as a write-only sensitive Terraform object.
- `notification_type` (String, Read-Only) Allowed values:
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
- `organization_id` (Number, Read-Only) Numeric ID of the related AWX organization object.

## Further Reading

- [AWX Notifications](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/notifications.html)
