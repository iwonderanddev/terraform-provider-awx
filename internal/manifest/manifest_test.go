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

func TestCredentialTypeInputsAndInjectorsAreObjectFields(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("credential_types")
	if !ok {
		t.Fatalf("expected credential_types object in catalog")
	}

	var inputsType FieldType
	var injectorsType FieldType
	for _, field := range object.Fields {
		switch field.Name {
		case "inputs":
			inputsType = field.Type
		case "injectors":
			injectorsType = field.Type
		}
	}

	if inputsType != FieldTypeObject {
		t.Fatalf("expected credential_types.inputs type=%q, got=%q", FieldTypeObject, inputsType)
	}
	if injectorsType != FieldTypeObject {
		t.Fatalf("expected credential_types.injectors type=%q, got=%q", FieldTypeObject, injectorsType)
	}
}

func TestNotificationTemplateMessagesIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("notification_templates")
	if !ok {
		t.Fatalf("expected notification_templates object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "messages" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected notification_templates.messages type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected messages field on notification_templates object")
}

func TestNotificationTemplateConfigurationIsWriteOnlySensitiveObject(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("notification_templates")
	if !ok {
		t.Fatalf("expected notification_templates object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "notification_configuration" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected notification_templates.notification_configuration type=%q, got=%q", FieldTypeObject, field.Type)
		}
		if !field.Sensitive {
			t.Fatalf("expected notification_templates.notification_configuration to be sensitive")
		}
		if !field.WriteOnly {
			t.Fatalf("expected notification_templates.notification_configuration to be write-only")
		}
		return
	}

	t.Fatalf("expected notification_configuration field on notification_templates object")
}

func TestJobTemplateExtraVarsIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("job_templates")
	if !ok {
		t.Fatalf("expected job_templates object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "extra_vars" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected job_templates.extra_vars type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected extra_vars field on job_templates object")
}

func TestWorkflowJobTemplateExtraVarsIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("workflow_job_templates")
	if !ok {
		t.Fatalf("expected workflow_job_templates object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "extra_vars" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected workflow_job_templates.extra_vars type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected extra_vars field on workflow_job_templates object")
}

func TestSettingsHostMetricTimestampsAreComputed(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("settings")
	if !ok {
		t.Fatalf("expected settings object in catalog")
	}

	requiredComputed := map[string]bool{
		"CLEANUP_HOST_METRICS_LAST_TS":      false,
		"HOST_METRIC_SUMMARY_TASK_LAST_TS":  false,
		"AUTOMATION_ANALYTICS_LAST_ENTRIES": false,
		"AUTOMATION_ANALYTICS_LAST_GATHER":  false,
	}

	for _, field := range object.Fields {
		if _, tracked := requiredComputed[field.Name]; !tracked {
			continue
		}
		requiredComputed[field.Name] = field.Computed
	}

	for name, computed := range requiredComputed {
		if !computed {
			t.Fatalf("expected settings.%s to be computed", name)
		}
	}
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
