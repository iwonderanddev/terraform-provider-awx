package manifest

import "testing"

func TestOrganizationsMaxHostsComputed(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("organizations")
	if !ok {
		t.Fatalf("expected organizations object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "max_hosts" {
			continue
		}
		if !field.Computed {
			t.Fatalf("expected organizations.max_hosts to be computed")
		}
		return
	}

	t.Fatalf("expected max_hosts field on organizations object")
}

func TestInventoryPreventInstanceGroupFallbackComputed(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("inventories")
	if !ok {
		t.Fatalf("expected inventories object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "prevent_instance_group_fallback" {
			continue
		}
		if !field.Computed {
			t.Fatalf("expected inventories.prevent_instance_group_fallback to be computed")
		}
		return
	}

	t.Fatalf("expected prevent_instance_group_fallback field on inventories object")
}

func TestCredentialsDescriptionNotComputed(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("credentials")
	if !ok {
		t.Fatalf("expected credentials object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "description" {
			continue
		}
		if field.Computed {
			t.Fatalf("expected credentials.description to remain optional (not computed)")
		}
		return
	}

	t.Fatalf("expected description field on credentials object")
}

func TestTerraformAttributeName(t *testing.T) {
	t.Parallel()

	if got := TerraformAttributeName("settings", "AUTH_LDAP_4_SERVER_URI"); got != "auth_ldap_4_server_uri" {
		t.Fatalf("unexpected settings field mapping: got=%q want=%q", got, "auth_ldap_4_server_uri")
	}

	if got := TerraformAttributeName("inventories", "name"); got != "name" {
		t.Fatalf("unexpected non-settings field mapping: got=%q want=%q", got, "name")
	}
}
