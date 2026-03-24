# Resource: awx_group

Manages AWX `groups` objects.

## Example Usage

### Basic configuration

```hcl
resource "awx_group" "example" {
  inventory_id = awx_inventory.example.id
  name = "example"
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

- `inventory_id` (Number, Required) Numeric ID of the related AWX inventory object.
- `name` (String, Required) AWX value stored in `name`.

### Optional

- `description` (String, Optional) AWX value stored in `description`.
- `variables` (String, Optional) Group variables in JSON or YAML format.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object identifier.

## Import

```bash
terraform import awx_group.example 42
```

## Further Reading

- [AWX Inventories](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/inventories.html)
