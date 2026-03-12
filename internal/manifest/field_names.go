package manifest

import "strings"

// TerraformAttributeName returns the Terraform-facing attribute name for an AWX field.
func TerraformAttributeName(objectName string, fieldName string) string {
	if objectName == "settings" {
		return strings.ToLower(fieldName)
	}
	return fieldName
}

// TerraformAttributeNameForField returns the Terraform-facing attribute name
// for a manifest field, including reference-link suffixing rules.
func TerraformAttributeNameForField(objectName string, field FieldSpec) string {
	if explicit := strings.TrimSpace(field.TerraformName); explicit != "" {
		return explicit
	}

	name := TerraformAttributeName(objectName, field.Name)
	if !field.Reference {
		return name
	}
	if strings.EqualFold(name, "id") || strings.HasSuffix(name, "_id") {
		return name
	}
	return name + "_id"
}
