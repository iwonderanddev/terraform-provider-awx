package manifest

import "strings"

// TerraformAttributeName returns the Terraform-facing attribute name for an AWX field.
func TerraformAttributeName(objectName string, fieldName string) string {
	if objectName == "settings" {
		return strings.ToLower(fieldName)
	}
	return fieldName
}
