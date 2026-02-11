package manifest

import "strings"

// RelationshipObjectIDAttribute returns the canonical Terraform attribute name
// for a relationship object identifier input.
func RelationshipObjectIDAttribute(collectionName string) string {
	base := SingularizeCollectionName(collectionName)
	if strings.TrimSpace(base) == "" {
		return "id"
	}
	if strings.HasSuffix(base, "_id") {
		return base
	}
	return base + "_id"
}

// RelationshipParentIDAttribute returns the canonical Terraform attribute name
// for the parent relationship object identifier.
func RelationshipParentIDAttribute(rel Relationship) string {
	if strings.TrimSpace(rel.ParentIDAttribute) != "" {
		return rel.ParentIDAttribute
	}
	return RelationshipObjectIDAttribute(rel.ParentObject)
}

// RelationshipChildIDAttribute returns the canonical Terraform attribute name
// for the child relationship object identifier.
func RelationshipChildIDAttribute(rel Relationship) string {
	if strings.TrimSpace(rel.ChildIDAttribute) != "" {
		return rel.ChildIDAttribute
	}
	return RelationshipObjectIDAttribute(rel.ChildObject)
}
