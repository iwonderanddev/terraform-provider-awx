package manifest

import "strings"

// SingularizeCollectionName converts a plural collection token into a singular
// resource token used by Terraform-facing names.
func SingularizeCollectionName(collectionName string) string {
	name := strings.TrimSpace(collectionName)
	if strings.HasSuffix(name, "ies") && len(name) > 3 {
		return strings.TrimSuffix(name, "ies") + "y"
	}
	if strings.HasSuffix(name, "sses") {
		return strings.TrimSuffix(name, "es")
	}
	if strings.HasSuffix(name, "ses") && len(name) > 3 {
		return strings.TrimSuffix(name, "es")
	}
	if strings.HasSuffix(name, "s") && !strings.HasSuffix(name, "ss") && len(name) > 1 {
		return strings.TrimSuffix(name, "s")
	}
	return name
}
