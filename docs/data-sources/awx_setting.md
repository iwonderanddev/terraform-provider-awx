# Data Source: awx_setting

Reads AWX `settings` objects.

Use `id = "all"` as the default and recommended settings scope.
Category-scoped IDs (for example `system`, `authentication`, and `bulk`) remain supported for optional advanced scoping.
Avoid overlapping ownership of the same setting key across `id = "all"` and category-scoped resources, because overlaps can cause configuration conflicts.

## Example Usage

```hcl
data "awx_setting" "example" {
  id = "all"
}
```

## Schema

### Optional

- `id` (String, Optional) AWX object identifier used in the detail endpoint path.

### Read-Only

- `id` (String, Read-Only) AWX detail-path identifier for this object.
- `activity_stream_enabled` (Boolean, Read-Only) Enable capturing activity for the activity stream.
- `activity_stream_enabled_for_inventory_sync` (Boolean, Read-Only) Enable capturing activity for the activity stream when running inventory sync.
- `ad_hoc_commands` (String, Read-Only) List of modules allowed to be used by ad-hoc jobs.
- `allow_jinja_in_extra_vars` (String, Read-Only) Ansible allows variable substitution via the Jinja2 templating language for --extra-vars. This poses a potential security risk where users with the ability to specify extra vars at job launch time can use Jinja2 templates to run arbitrary Python.  It is recommended that this value be set to "template" or "never".
  - `always` - Always
  - `never` - Never
  - `template` - Only On Job Template Definitions
- `allow_metrics_for_anonymous_users` (Boolean, Read-Only) If true, anonymous users are allowed to poll metrics.
- `allow_oauth2_for_external_users` (Boolean, Read-Only) For security reasons, users from external auth providers (LDAP, SAML, SSO, Radius, and others) are not allowed to create OAuth2 tokens. To change this behavior, enable this setting. Existing tokens will not be deleted when this setting is toggled off.
- `ansible_fact_cache_timeout` (Number, Read-Only) Maximum time, in seconds, that stored Ansible facts are considered valid since the last time they were modified. Only valid, non-stale, facts will be accessible by a playbook. Note, this does not influence the deletion of ansible_facts from the database. Use a value of 0 to indicate that no timeout should be imposed.
- `api_400_error_log_format` (String, Read-Only) The format of logged messages when an API 4XX error occurs, the following variables will be substituted:
  status_code - The HTTP status code of the error
  user_name - The user name attempting to use the API
  url_path - The URL path to the API endpoint called
  remote_addr - The remote address seen for the user
  error - The error set by the api endpoint
  Variables need to be in the format {&lt;variable name&gt;}.
- `authentication_backends` (String, Read-Only) List of authentication backends that are enabled based on license features and other authentication settings.
- `auth_basic_enabled` (Boolean, Read-Only) Enable HTTP Basic Auth for the API Browser.
- `auth_ldap_1_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_1_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_1_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_1_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_1_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_1_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_1_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_1_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_1_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_1_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_1_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_1_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_1_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_1_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_1_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_1_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `auth_ldap_2_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_2_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_2_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_2_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_2_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_2_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_2_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_2_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_2_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_2_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_2_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_2_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_2_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_2_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_2_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_2_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `auth_ldap_3_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_3_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_3_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_3_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_3_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_3_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_3_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_3_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_3_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_3_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_3_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_3_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_3_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_3_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_3_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_3_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `auth_ldap_4_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_4_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_4_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_4_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_4_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_4_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_4_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_4_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_4_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_4_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_4_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_4_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_4_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_4_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_4_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_4_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `auth_ldap_5_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_5_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_5_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_5_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_5_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_5_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_5_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_5_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_5_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_5_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_5_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_5_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_5_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_5_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_5_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_5_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `auth_ldap_bind_dn` (String, Read-Only) DN (Distinguished Name) of user to bind for all search queries. This is the system user account we will use to login to query LDAP for other user information. Refer to the documentation for example syntax.
- `auth_ldap_bind_password` (String, Read-Only, Sensitive) Password used to bind LDAP user account.
- `auth_ldap_connection_options` (Object, Read-Only) Additional options to set for the LDAP connection. LDAP referrals are disabled by default (to prevent certain LDAP queries from hanging with AD). Option names should be strings (e.g. "OPT_REFERRALS"). Refer to https://www.python-ldap.org/doc/html/ldap.html#options for possible options and values that can be set.
- `auth_ldap_deny_group` (String, Read-Only) Group DN denied from login. If specified, user will not be allowed to login if a member of this group. Only one deny group is supported.
- `auth_ldap_group_search` (String, Read-Only) Users are mapped to organizations based on their membership in LDAP groups. This setting defines the LDAP search query to find groups. Unlike the user search, group search does not support LDAPSearchUnion.
- `auth_ldap_group_type` (String, Read-Only) The group type may need to be changed based on the type of the LDAP server. Values are listed at: https://django-auth-ldap.readthedocs.io/en/stable/groups.html#types-of-groups
- `auth_ldap_group_type_params` (Object, Read-Only) Key value parameters to send the chosen group type init method.
- `auth_ldap_organization_map` (Object, Read-Only) Mapping between organization admins/users and LDAP groups. This controls which users are placed into which organizations relative to their LDAP group memberships. Configuration details are available in the documentation.
- `auth_ldap_require_group` (String, Read-Only) Group DN required to login. If specified, user must be a member of this group to login via LDAP. If not set, everyone in LDAP that matches the user search will be able to login to the service. Only one require group is supported.
- `auth_ldap_server_uri` (String, Read-Only) URI to connect to LDAP server, such as "ldap://ldap.example.com:389" (non-SSL) or "ldaps://ldap.example.com:636" (SSL). Multiple LDAP servers may be specified by separating with spaces or commas. LDAP authentication is disabled if this parameter is empty.
- `auth_ldap_start_tls` (Boolean, Read-Only) Whether to enable TLS when the LDAP connection is not using SSL.
- `auth_ldap_team_map` (Object, Read-Only) Mapping between team members (users) and LDAP groups. Configuration details are available in the documentation.
- `auth_ldap_user_attr_map` (Object, Read-Only) Mapping of LDAP user schema to API user attributes. The default setting is valid for ActiveDirectory but users with other LDAP configurations may need to change the values. Refer to the documentation for additional details.
- `auth_ldap_user_dn_template` (String, Read-Only) Alternative to user search, if user DNs are all of the same format. This approach is more efficient for user lookups than searching if it is usable in your organizational environment. If this setting has a value it will be used instead of AUTH_LDAP_USER_SEARCH.
- `auth_ldap_user_flags_by_group` (Object, Read-Only) Retrieve users from a given group. At this time, superuser and system auditors are the only groups supported. Refer to the documentation for more detail.
- `auth_ldap_user_search` (String, Read-Only) LDAP search query to find users. Any user that matches the given pattern will be able to login to the service. The user should also be mapped into an organization (as defined in the AUTH_LDAP_ORGANIZATION_MAP setting). If multiple search queries need to be supported use of "LDAPUnion" is possible. See the documentation for details.
- `automation_analytics_gather_interval` (Number, Read-Only) Interval (in seconds) between data gathering.
- `automation_analytics_last_entries` (String, Read-Only) AWX value stored in `automation_analytics_last_entries`.
- `automation_analytics_last_gather` (String, Read-Only) AWX value stored in `automation_analytics_last_gather`.
- `automation_analytics_url` (String, Read-Only) This setting is used to to configure the upload URL for data collection for Automation Analytics.
- `awx_ansible_callback_plugins` (String, Read-Only) List of paths to search for extra callback plugins to be used when running jobs. Enter one path per line.
- `awx_cleanup_paths` (Boolean, Read-Only) Enable or Disable TMP Dir cleanup
- `awx_collections_enabled` (Boolean, Read-Only) Allows collections to be dynamically downloaded from a requirements.yml file for SCM projects.
- `awx_isolation_base_path` (String, Read-Only) The directory in which the service will create new temporary directories for job execution and isolation (such as credential files).
- `awx_isolation_show_paths` (String, Read-Only) List of paths that would otherwise be hidden to expose to isolated jobs. Enter one path per line. Volumes will be mounted from the execution node to the container. The supported format is HOST-DIR[:CONTAINER-DIR[:OPTIONS]].
- `awx_mount_isolated_paths_on_k8s` (Boolean, Read-Only) Expose paths via hostPath for the Pods created by a Container Group. HostPath volumes present many security risks, and it is a best practice to avoid the use of HostPaths when possible.
- `awx_request_profile` (Boolean, Read-Only) Debug web request python timing
- `awx_roles_enabled` (Boolean, Read-Only) Allows roles to be dynamically downloaded from a requirements.yml file for SCM projects.
- `awx_runner_keepalive_seconds` (Number, Read-Only) Only applies to jobs running in a Container Group. If not 0, send a message every so-many seconds to keep connection open.
- `awx_show_playbook_links` (Boolean, Read-Only) Follow symbolic links when scanning for playbooks. Be aware that setting this to True can lead to infinite recursion if a link points to a parent directory of itself.
- `awx_task_env` (Object, Read-Only) Additional environment variables set for playbook runs, inventory updates, project updates, and notification sending.
- `bulk_host_max_create` (Number, Read-Only) Max number of hosts to allow to be created in a single bulk action
- `bulk_host_max_delete` (Number, Read-Only) Max number of hosts to allow to be deleted in a single bulk action
- `bulk_job_max_launch` (Number, Read-Only) Max jobs to allow bulk jobs to launch
- `cleanup_host_metrics_last_ts` (String, Read-Only) AWX value stored in `cleanup_host_metrics_last_ts`.
- `csrf_trusted_origins` (String, Read-Only) If the service is behind a reverse proxy/load balancer, use this setting to configure the schema://addresses from which the service should trust Origin header values.
- `custom_login_info` (String, Read-Only) If needed, you can add specific information (such as a legal notice or a disclaimer) to a text box in the login modal using this setting. Any content added must be in plain text or an HTML fragment, as other markup languages are not supported.
- `custom_logo` (String, Read-Only) To set up a custom logo, provide a file that you create. For the custom logo to look its best, use a .png file with a transparent background. GIF, PNG and JPEG formats are supported.
- `custom_venv_paths` (String, Read-Only) Paths where Tower will look for custom virtual environments (in addition to /var/lib/awx/venv/). Enter one path per line.
- `default_container_run_options` (String, Read-Only) List of options to pass to podman run example: ['--network', 'slirp4netns:enable_ipv6=true', '--log-level', 'debug']
- `default_control_plane_queue_name` (String, Read-Only) AWX value stored in `default_control_plane_queue_name`.
- `default_execution_environment_id` (Number, Read-Only) The Execution Environment to be used when one has not been configured for a job template.
- `default_execution_queue_name` (String, Read-Only) AWX value stored in `default_execution_queue_name`.
- `default_inventory_update_timeout` (Number, Read-Only) Maximum time in seconds to allow inventory updates to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual inventory source will override this.
- `default_job_idle_timeout` (Number, Read-Only) If no output is detected from ansible in this number of seconds the execution will be terminated. Use value of 0 to indicate that no idle timeout should be imposed.
- `default_job_timeout` (Number, Read-Only) Maximum time in seconds to allow jobs to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual job template will override this.
- `default_project_update_timeout` (Number, Read-Only) Maximum time in seconds to allow project updates to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual project will override this.
- `disable_local_auth` (Boolean, Read-Only) Controls whether users are prevented from using the built-in authentication system.
- `event_stdout_max_bytes_display` (Number, Read-Only) Maximum Size of Standard Output in bytes to display for a single job or ad hoc command event. `stdout` will end with `…` when truncated.
- `galaxy_ignore_certs` (Boolean, Read-Only) If set to true, certificate validation will not be done when installing content from any Galaxy server.
- `galaxy_task_env` (Object, Read-Only) Additional environment variables set for invocations of ansible-galaxy within project updates. Useful if you must use a proxy server for ansible-galaxy but not git.
- `host_metric_summary_task_last_ts` (String, Read-Only) AWX value stored in `host_metric_summary_task_last_ts`.
- `insights_tracking_state` (Boolean, Read-Only) Enables the service to gather data on automation and send it to Automation Analytics.
- `install_uuid` (String, Read-Only) AWX value stored in `install_uuid`.
- `is_k8s` (Boolean, Read-Only) Indicates whether the instance is part of a kubernetes-based deployment.
- `license` (Object, Read-Only) The license controls which features and functionality are enabled. Use /api/v2/config/ to update or change the license.
- `local_password_min_digits` (Number, Read-Only, Sensitive) Minimum number of digit characters required in a local password. 0 means no minimum
- `local_password_min_length` (Number, Read-Only, Sensitive) Minimum number of characters required in a local password. 0 means no minimum
- `local_password_min_special` (Number, Read-Only, Sensitive) Minimum number of special characters required in a local password. 0 means no minimum
- `local_password_min_upper` (Number, Read-Only, Sensitive) Minimum number of uppercase characters required in a local password. 0 means no minimum
- `login_redirect_override` (String, Read-Only) URL to which unauthorized users will be redirected to log in.  If blank, users will be sent to the login page.
- `log_aggregator_action_max_disk_usage_gb` (Number, Read-Only) Amount of data to store (in gigabytes) if an rsyslog action takes time to process an incoming message (defaults to 1). Equivalent to the rsyslogd queue.maxdiskspace setting on the action (e.g. omhttp). It stores files in the directory specified by LOG_AGGREGATOR_MAX_DISK_USAGE_PATH.
- `log_aggregator_action_queue_size` (Number, Read-Only) Defines how large the rsyslog action queue can grow in number of messages stored. This can have an impact on memory utilization. When the queue reaches 75% of this number, the queue will start writing to disk (queue.highWatermark in rsyslog). When it reaches 90%, NOTICE, INFO, and DEBUG messages will start to be discarded (queue.discardMark with queue.discardSeverity=5).
- `log_aggregator_enabled` (Boolean, Read-Only) Enable sending logs to external log aggregator.
- `log_aggregator_host` (String, Read-Only) Hostname/IP where external logs will be sent to.
- `log_aggregator_individual_facts` (Boolean, Read-Only) If set, system tracking facts will be sent for each package, service, or other item found in a scan, allowing for greater search query granularity. If unset, facts will be sent as a single dictionary, allowing for greater efficiency in fact processing.
- `log_aggregator_level` (String, Read-Only) Level threshold used by log handler. Severities from lowest to highest are DEBUG, INFO, WARNING, ERROR, CRITICAL. Messages less severe than the threshold will be ignored by log handler. (messages under category awx.anlytics ignore this setting)
  - `DEBUG` - DEBUG
  - `INFO` - INFO
  - `WARNING` - WARNING
  - `ERROR` - ERROR
  - `CRITICAL` - CRITICAL
- `log_aggregator_loggers` (String, Read-Only) List of loggers that will send HTTP logs to the collector, these can include any or all of:
  awx - service logs
  activity_stream - activity stream records
  job_events - callback data from Ansible job events
  system_tracking - facts gathered from scan jobs
  broadcast_websocket - errors pertaining to websockets broadcast metrics
  job_lifecycle - logs related to processing of a job
- `log_aggregator_max_disk_usage_path` (String, Read-Only) Location to persist logs that should be retried after an outage of the external log aggregator (defaults to /var/lib/awx). Equivalent to the rsyslogd queue.spoolDirectory setting.
- `log_aggregator_password` (String, Read-Only, Sensitive) Password or authentication token for external log aggregator (if required; HTTP/s only).
- `log_aggregator_port` (Number, Read-Only) Port on Logging Aggregator to send logs to (if required and not provided in Logging Aggregator).
- `log_aggregator_protocol` (String, Read-Only) Protocol used to communicate with log aggregator.  HTTPS/HTTP assumes HTTPS unless http:// is explicitly used in the Logging Aggregator hostname.
  - `https` - HTTPS/HTTP
  - `tcp` - TCP
  - `udp` - UDP
- `log_aggregator_rsyslogd_debug` (Boolean, Read-Only) Enabled high verbosity debugging for rsyslogd.  Useful for debugging connection issues for external log aggregation.
- `log_aggregator_tcp_timeout` (Number, Read-Only) Number of seconds for a TCP connection to external log aggregator to timeout. Applies to HTTPS and TCP log aggregator protocols.
- `log_aggregator_tower_uuid` (String, Read-Only) Useful to uniquely identify instances.
- `log_aggregator_type` (String, Read-Only) Format messages for the chosen log aggregator.
  - `logstash` - logstash
  - `splunk` - splunk
  - `loggly` - loggly
  - `sumologic` - sumologic
  - `other` - other
- `log_aggregator_username` (String, Read-Only) Username for external log aggregator (if required; HTTP/s only).
- `log_aggregator_verify_cert` (Boolean, Read-Only) Flag to control enable/disable of certificate verification when LOG_AGGREGATOR_PROTOCOL is "https". If enabled, the log handler will verify certificate sent by external log aggregator before establishing connection.
- `manage_organization_auth` (Boolean, Read-Only) Controls whether any Organization Admin has the privileges to create and manage users and teams.
- `max_forks` (Number, Read-Only) Saving a Job Template with more than this number of forks will result in an error. When set to 0, no limit is applied.
- `max_ui_job_events` (Number, Read-Only) Maximum number of job events for the UI to retrieve within a single request.
- `max_websocket_event_rate` (Number, Read-Only) Maximum number of messages to update the UI live job output with per second. Value of 0 means no limit.
- `named_url_formats` (Object, Read-Only) Read-only list of key-value pairs that shows the standard format of all available named URLs.
- `named_url_graph_nodes` (Object, Read-Only) Read-only list of key-value pairs that exposes named URL graph topology. Use this list to programmatically generate named URLs for resources
- `oauth2_provider` (Object, Read-Only) Dictionary for customizing OAuth 2 timeouts, available items are `ACCESS_TOKEN_EXPIRE_SECONDS`, the duration of access tokens in the number of seconds, `AUTHORIZATION_CODE_EXPIRE_SECONDS`, the duration of authorization codes in the number of seconds, and `REFRESH_TOKEN_EXPIRE_SECONDS`, the duration of refresh tokens, after expired access tokens, in the number of seconds.
- `opa_auth_ca_cert` (String, Read-Only) The content of the CA certificate for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `opa_auth_client_cert` (String, Read-Only) The content of the client certificate file for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `opa_auth_client_key` (String, Read-Only) The content of the client key for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `opa_auth_custom_headers` (Object, Read-Only) Optional custom headers included in requests to the OPA server. Defaults to empty dictionary ({}).
- `opa_auth_token` (String, Read-Only, Sensitive) The token for authentication to the OPA server. Required when OPA_AUTH_TYPE is "Token". If an authorization header is defined in OPA_AUTH_CUSTOM_HEADERS, it will be overridden by OPA_AUTH_TOKEN.
- `opa_auth_type` (String, Read-Only) The authentication type that will be used to connect to the OPA server: "None", "Token", or "Certificate".
  - `None` - None
  - `Token` - Token
  - `Certificate` - Certificate
- `opa_host` (String, Read-Only) The hostname used to connect to the OPA server. If empty, policy enforcement will be disabled.
- `opa_port` (Number, Read-Only) The port used to connect to the OPA server. Defaults to 8181.
- `opa_request_retries` (Number, Read-Only) The number of retry attempts for connecting to the OPA server. Default is 2.
- `opa_request_timeout` (Number, Read-Only) The number of seconds after which the connection to the OPA server will time out. Defaults to 1.5 seconds.
- `opa_ssl` (Boolean, Read-Only) Enable or disable the use of SSL to connect to the OPA server. Defaults to false.
- `org_admins_can_see_all_users` (Boolean, Read-Only) Controls whether any Organization Admin can view all users and teams, even those not associated with their Organization.
- `pendo_tracking_state` (String, Read-Only) Enable or Disable User Analytics Tracking.
- `project_update_vvv` (Boolean, Read-Only) Adds the CLI -vvv flag to ansible-playbook runs of project_update.yml used for project updates.
- `proxy_ip_allowed_list` (String, Read-Only) If the service is behind a reverse proxy/load balancer, use this setting to configure the proxy IP addresses from which the service should trust custom REMOTE_HOST_HEADERS header values. If this setting is an empty list (the default), the headers specified by REMOTE_HOST_HEADERS will be trusted unconditionally')
- `radius_port` (Number, Read-Only) Port of RADIUS server.
- `radius_secret` (String, Read-Only, Sensitive) Shared secret for authenticating to RADIUS server.
- `radius_server` (String, Read-Only) Hostname/IP of RADIUS server. RADIUS authentication is disabled if this setting is empty.
- `receptor_keep_work_on_error` (Boolean, Read-Only) Prevent receptor work from being released on when error is detected
- `receptor_release_work` (Boolean, Read-Only) Release receptor work
- `redhat_password` (String, Read-Only, Sensitive) Client secret used to send data to Automation Analytics
- `redhat_username` (String, Read-Only) Client ID used to send data to Automation Analytics
- `remote_host_headers` (String, Read-Only) HTTP headers and meta keys to search to determine remote host name or IP. Add additional items to this list, such as "HTTP_X_FORWARDED_FOR", if behind a reverse proxy. See the "Proxy Support" section of the AAP Installation guide for more details.
- `saml_auto_create_objects` (Boolean, Read-Only) When enabled (the default), mapped Organizations and Teams will be created automatically on successful SAML login.
- `schedule_max_jobs` (Number, Read-Only) Maximum number of the same job template that can be waiting to run when launching from a schedule before no more are created.
- `sessions_per_user` (Number, Read-Only) Maximum number of simultaneous logged in sessions a user may have. To disable enter -1.
- `session_cookie_age` (Number, Read-Only) Number of seconds that a user is inactive before they will need to login again.
- `social_auth_azuread_oauth2_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_azuread_oauth2_key` (String, Read-Only) The OAuth2 key (Client ID) from your Azure AD application.
- `social_auth_azuread_oauth2_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_azuread_oauth2_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your Azure AD application.
- `social_auth_azuread_oauth2_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_github_enterprise_api_url` (String, Read-Only) The API URL for your GitHub Enterprise instance, e.g.: http(s)://hostname/api/v3/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_enterprise_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_github_enterprise_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub Enterprise developer application.
- `social_auth_github_enterprise_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_enterprise_org_api_url` (String, Read-Only) The API URL for your GitHub Enterprise instance, e.g.: http(s)://hostname/api/v3/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_enterprise_org_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_github_enterprise_org_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub Enterprise organization application.
- `social_auth_github_enterprise_org_name` (String, Read-Only) The name of your GitHub Enterprise organization, as used in your organization's URL: https://github.com/&lt;yourorg&gt;/.
- `social_auth_github_enterprise_org_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_enterprise_org_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub Enterprise organization application.
- `social_auth_github_enterprise_org_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_enterprise_org_url` (String, Read-Only) The URL for your Github Enterprise instance, e.g.: http(s)://hostname/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_enterprise_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub Enterprise developer application.
- `social_auth_github_enterprise_team_api_url` (String, Read-Only) The API URL for your GitHub Enterprise instance, e.g.: http(s)://hostname/api/v3/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_enterprise_team_callback_url` (String, Read-Only) Create an organization-owned application at https://github.com/organizations/&lt;yourorg&gt;/settings/applications and obtain an OAuth2 key (Client ID) and secret (Client Secret). Provide this URL as the callback URL for your application.
- `social_auth_github_enterprise_team_id` (String, Read-Only) Find the numeric team ID using the Github Enterprise API: http://fabian-kostadinov.github.io/2015/01/16/how-to-find-a-github-team-id/.
- `social_auth_github_enterprise_team_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub Enterprise organization application.
- `social_auth_github_enterprise_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_enterprise_team_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_enterprise_team_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub Enterprise organization application.
- `social_auth_github_enterprise_team_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_enterprise_team_url` (String, Read-Only) The URL for your Github Enterprise instance, e.g.: http(s)://hostname/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_enterprise_url` (String, Read-Only) The URL for your Github Enterprise instance, e.g.: http(s)://hostname/. Refer to Github Enterprise documentation for more details.
- `social_auth_github_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub developer application.
- `social_auth_github_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_org_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_github_org_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub organization application.
- `social_auth_github_org_name` (String, Read-Only) The name of your GitHub organization, as used in your organization's URL: https://github.com/&lt;yourorg&gt;/.
- `social_auth_github_org_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_org_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub organization application.
- `social_auth_github_org_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub developer application.
- `social_auth_github_team_callback_url` (String, Read-Only) Create an organization-owned application at https://github.com/organizations/&lt;yourorg&gt;/settings/applications and obtain an OAuth2 key (Client ID) and secret (Client Secret). Provide this URL as the callback URL for your application.
- `social_auth_github_team_id` (String, Read-Only) Find the numeric team ID using the Github API: http://fabian-kostadinov.github.io/2015/01/16/how-to-find-a-github-team-id/.
- `social_auth_github_team_key` (String, Read-Only) The OAuth2 key (Client ID) from your GitHub organization application.
- `social_auth_github_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_github_team_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_github_team_secret` (String, Read-Only, Sensitive) The OAuth2 secret (Client Secret) from your GitHub organization application.
- `social_auth_github_team_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_google_oauth2_auth_extra_arguments` (Object, Read-Only) Extra arguments for Google OAuth2 login. You can restrict it to only allow a single domain to authenticate, even if the user is logged in with multple Google accounts. Refer to the documentation for more detail.
- `social_auth_google_oauth2_callback_url` (String, Read-Only) Provide this URL as the callback URL for your application as part of your registration process. Refer to the documentation for more detail.
- `social_auth_google_oauth2_key` (String, Read-Only) The OAuth2 key from your web application.
- `social_auth_google_oauth2_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_google_oauth2_secret` (String, Read-Only, Sensitive) The OAuth2 secret from your web application.
- `social_auth_google_oauth2_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_google_oauth2_whitelisted_domains` (String, Read-Only) Update this setting to restrict the domains who are allowed to login using Google OAuth2.
- `social_auth_oidc_key` (String, Read-Only) The OIDC key (Client ID) from your IDP.
- `social_auth_oidc_oidc_endpoint` (String, Read-Only) The URL for your OIDC provider including the path up to /.well-known/openid-configuration
- `social_auth_oidc_secret` (String, Read-Only, Sensitive) The OIDC secret (Client Secret) from your IDP.
- `social_auth_oidc_verify_ssl` (Boolean, Read-Only) Verify the OIDC provider ssl certificate.
- `social_auth_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_saml_callback_url` (String, Read-Only) Register the service as a service provider (SP) with each identity provider (IdP) you have configured. Provide your SP Entity ID and this ACS URL for your application.
- `social_auth_saml_enabled_idps` (Object, Read-Only) Configure the Entity ID, SSO URL and certificate for each identity provider (IdP) in use. Multiple SAML IdPs are supported. Some IdPs may provide user data using attribute names that differ from the default OIDs. Attribute names may be overridden for each IdP. Refer to the Ansible documentation for additional details and syntax.
- `social_auth_saml_extra_data` (String, Read-Only) A list of tuples that maps IDP attributes to extra_attributes. Each attribute will be a list of values, even if only 1 value.
- `social_auth_saml_metadata_url` (String, Read-Only) If your identity provider (IdP) allows uploading an XML metadata file, you can download one from this URL.
- `social_auth_saml_organization_attr` (Object, Read-Only) Used to translate user organization membership.
- `social_auth_saml_organization_map` (Object, Read-Only) Mapping to organization admins/users from social auth accounts. This setting controls which users are placed into which organizations based on their username and email address. Configuration details are available in the documentation.
- `social_auth_saml_org_info` (Object, Read-Only) Provide the URL, display name, and the name of your app. Refer to the documentation for example syntax.
- `social_auth_saml_security_config` (Object, Read-Only) A dict of key value pairs that are passed to the underlying python-saml security setting https://github.com/onelogin/python-saml#settings
- `social_auth_saml_sp_entity_id` (String, Read-Only) The application-defined unique identifier used as the audience of the SAML service provider (SP) configuration. This is usually the URL for the service.
- `social_auth_saml_sp_extra` (Object, Read-Only) A dict of key value pairs to be passed to the underlying python-saml Service Provider configuration setting.
- `social_auth_saml_sp_private_key` (String, Read-Only, Sensitive) Create a keypair to use as a service provider (SP) and include the private key content here.
- `social_auth_saml_sp_public_cert` (String, Read-Only) Create a keypair to use as a service provider (SP) and include the certificate content here.
- `social_auth_saml_support_contact` (Object, Read-Only) Provide the name and email address of the support contact for your service provider. Refer to the documentation for example syntax.
- `social_auth_saml_team_attr` (Object, Read-Only) Used to translate user team membership.
- `social_auth_saml_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_saml_technical_contact` (Object, Read-Only) Provide the name and email address of the technical contact for your service provider. Refer to the documentation for example syntax.
- `social_auth_saml_user_flags_by_attr` (Object, Read-Only) Used to map super users and system auditors from SAML.
- `social_auth_team_map` (Object, Read-Only) Mapping of team members (users) from social auth accounts. Configuration details are available in the documentation.
- `social_auth_username_is_full_email` (Boolean, Read-Only) Enabling this setting will tell social auth to use the full Email as username instead of the full name
- `social_auth_user_fields` (String, Read-Only) When set to an empty list `[]`, this setting prevents new user accounts from being created. Only users who have previously logged in using social auth or have a user account with a matching email address will be able to login.
- `stdout_max_bytes_display` (Number, Read-Only) Maximum Size of Standard Output in bytes to display before requiring the output be downloaded.
- `subscriptions_client_id` (String, Read-Only) Client ID used to retrieve subscription and content information
- `subscriptions_client_secret` (String, Read-Only, Sensitive) Client secret used to retrieve subscription and content information
- `subscriptions_password` (String, Read-Only, Sensitive) Password used to retrieve subscription and content information
- `subscriptions_username` (String, Read-Only) Username used to retrieve subscription and content information
- `subscription_usage_model` (String, Read-Only) Allowed values:
  - `` - No subscription. Deletion of host_metrics will not be considered for purposes of managed host counting
  - `unique_managed_hosts` - Usage based on unique managed nodes in a large historical time frame and delete functionality for no longer used managed nodes
- `tacacsplus_auth_protocol` (String, Read-Only) Choose the authentication protocol used by TACACS+ client.
- `tacacsplus_host` (String, Read-Only) Hostname of TACACS+ server.
- `tacacsplus_port` (Number, Read-Only) Port number of TACACS+ server.
- `tacacsplus_rem_addr` (Boolean, Read-Only) Enable the client address sending by TACACS+ client.
- `tacacsplus_secret` (String, Read-Only, Sensitive) Shared secret for authenticating to TACACS+ server.
- `tacacsplus_session_timeout` (Number, Read-Only) TACACS+ session timeout value in seconds, 0 disables timeout.
- `tower_url_base` (String, Read-Only) This setting is used by services like notifications to render a valid url to the service.
- `ui_live_updates_enabled` (Boolean, Read-Only) If disabled, the page will not refresh when events are received. Reloading the page will be required to get the latest details.
- `ui_next` (Boolean, Read-Only) Enable preview of new user interface.

## Further Reading

- [AWX Configuration](https://docs.ansible.com/projects/awx/en/24.6.1/administration/configure_awx.html)
