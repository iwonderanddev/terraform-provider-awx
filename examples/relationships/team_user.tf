resource "awx_team_user_association" "membership" {
  team_id = 12 # team ID
  user_id = 34 # user ID
}

# Relationship resource imports use composite IDs.
# terraform import awx_team_user_association.membership 12:34
