resource "awx_team_user_association" "membership" {
  parent_id = 12 # team ID
  child_id  = 34 # user ID
}

# Relationship resource imports use composite IDs.
# terraform import awx_team_user_association.membership 12:34
