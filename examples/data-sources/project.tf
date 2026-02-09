data "awx_project" "existing" {
  name = "automation-project"
}

output "project_id" {
  value = data.awx_project.existing.id
}
