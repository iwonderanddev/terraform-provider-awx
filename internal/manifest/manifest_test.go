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

func TestScheduleExtraDataIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("schedules")
	if !ok {
		t.Fatalf("expected schedules object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "extra_data" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected schedules.extra_data type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected extra_data field on schedules object")
}

func TestWorkflowJobTemplateNodeExtraDataIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("workflow_job_template_nodes")
	if !ok {
		t.Fatalf("expected workflow_job_template_nodes object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "extra_data" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected workflow_job_template_nodes.extra_data type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected extra_data field on workflow_job_template_nodes object")
}

func TestWorkflowJobNodeExtraDataIsObjectField(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	object, ok := catalog.ObjectByName("workflow_job_nodes")
	if !ok {
		t.Fatalf("expected workflow_job_nodes object in catalog")
	}

	for _, field := range object.Fields {
		if field.Name != "extra_data" {
			continue
		}
		if field.Type != FieldTypeObject {
			t.Fatalf("expected workflow_job_nodes.extra_data type=%q, got=%q", FieldTypeObject, field.Type)
		}
		return
	}

	t.Fatalf("expected extra_data field on workflow_job_nodes object")
}

func TestRoleAssignmentObjectIDFieldsAreNumeric(t *testing.T) {
	t.Parallel()

	catalog := MustLoad()
	objectNames := []string{"role_team_assignments", "role_user_assignments"}

	for _, objectName := range objectNames {
		object, ok := catalog.ObjectByName(objectName)
		if !ok {
			t.Fatalf("expected %s object in catalog", objectName)
		}

		found := false
		for _, field := range object.Fields {
			if field.Name != "object_id" {
				continue
			}
			found = true
			if field.Type != FieldTypeInt {
				t.Fatalf("expected %s.object_id type=%q, got=%q", objectName, FieldTypeInt, field.Type)
			}
		}

		if !found {
			t.Fatalf("expected object_id field on %s", objectName)
		}
	}
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

	if got := TerraformAttributeNameForField("teams", FieldSpec{Name: "organization", Type: FieldTypeInt, Reference: true}); got != "organization_id" {
		t.Fatalf("unexpected reference field mapping: got=%q want=%q", got, "organization_id")
	}

	if got := TerraformAttributeNameForField("inventories", FieldSpec{Name: "organization_id", Type: FieldTypeInt, Reference: true}); got != "organization_id" {
		t.Fatalf("unexpected pre-suffixed reference field mapping: got=%q want=%q", got, "organization_id")
	}

	if got := TerraformAttributeNameForField("settings", FieldSpec{Name: "AUTH_LDAP_4_SERVER_URI", Type: FieldTypeString}); got != "auth_ldap_4_server_uri" {
		t.Fatalf("unexpected settings field mapping with field helper: got=%q want=%q", got, "auth_ldap_4_server_uri")
	}

	if got := TerraformAttributeNameForField("projects", FieldSpec{Name: "credential", Type: FieldTypeInt, Reference: true, TerraformName: "scm_credential_id"}); got != "scm_credential_id" {
		t.Fatalf("unexpected explicit Terraform name mapping: got=%q want=%q", got, "scm_credential_id")
	}
}

func TestRelationshipObjectIDAttribute(t *testing.T) {
	t.Parallel()

	if got := RelationshipObjectIDAttribute("job_templates"); got != "job_template_id" {
		t.Fatalf("unexpected relationship object ID attribute: got=%q want=%q", got, "job_template_id")
	}
	if got := RelationshipObjectIDAttribute("users"); got != "user_id" {
		t.Fatalf("unexpected relationship object ID attribute: got=%q want=%q", got, "user_id")
	}
	if got := RelationshipObjectIDAttribute("custom_id"); got != "custom_id" {
		t.Fatalf("unexpected pre-suffixed relationship object ID attribute: got=%q want=%q", got, "custom_id")
	}
}

func TestRelationshipIDAttributesFallbackToCollectionNames(t *testing.T) {
	t.Parallel()

	rel := Relationship{
		ParentObject: "workflow_job_templates",
		ChildObject:  "credentials",
	}
	if got := RelationshipParentIDAttribute(rel); got != "workflow_job_template_id" {
		t.Fatalf("unexpected parent ID attribute fallback: got=%q want=%q", got, "workflow_job_template_id")
	}
	if got := RelationshipChildIDAttribute(rel); got != "credential_id" {
		t.Fatalf("unexpected child ID attribute fallback: got=%q want=%q", got, "credential_id")
	}
}

func TestRelationshipIDAttributesUseManifestOverrides(t *testing.T) {
	t.Parallel()

	rel := Relationship{
		ParentObject:      "job_templates",
		ChildObject:       "credentials",
		ParentIDAttribute: "job_template_id",
		ChildIDAttribute:  "credential_id",
	}
	if got := RelationshipParentIDAttribute(rel); got != "job_template_id" {
		t.Fatalf("unexpected parent ID attribute override: got=%q want=%q", got, "job_template_id")
	}
	if got := RelationshipChildIDAttribute(rel); got != "credential_id" {
		t.Fatalf("unexpected child ID attribute override: got=%q want=%q", got, "credential_id")
	}
}

func TestSingularizeCollectionName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want string
	}{
		{in: "job_templates", want: "job_template"},
		{in: "credentials", want: "credential"},
		{in: "classes", want: "class"},
		{in: "policies", want: "policy"},
		{in: "settings", want: "setting"},
		{in: "glass", want: "glass"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			if got := SingularizeCollectionName(tc.in); got != tc.want {
				t.Fatalf("unexpected singularization: input=%q got=%q want=%q", tc.in, got, tc.want)
			}
		})
	}
}
