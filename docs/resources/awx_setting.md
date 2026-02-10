# Resource: awx_setting

Manages AWX `settings` objects.

## Example Usage

```hcl
resource "awx_setting" "example" {
  id = "example"
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `id` (Required) AWX detail-path identifier for this object.
- `ACTIVITY_STREAM_ENABLED` (Optional, Computed) Enable capturing activity for the activity stream.
- `ACTIVITY_STREAM_ENABLED_FOR_INVENTORY_SYNC` (Optional, Computed) Enable capturing activity for the activity stream when running inventory sync.
- `AD_HOC_COMMANDS` (Optional, Computed) List of modules allowed to be used by ad-hoc jobs.
- `ALLOW_JINJA_IN_EXTRA_VARS` (Optional, Computed) Ansible allows variable substitution via the Jinja2 templating language for --extra-vars. This poses a potential security risk where users with the ability to specify extra vars at job launch time can use Jinja2 templates to run arbitrary Python.  It is recommended that this value be set to "template" or "never".

* `always` - Always
* `never` - Never
* `template` - Only On Job Template Definitions
- `ALLOW_METRICS_FOR_ANONYMOUS_USERS` (Optional, Computed) If true, anonymous users are allowed to poll metrics.
- `ANSIBLE_FACT_CACHE_TIMEOUT` (Optional, Computed) Maximum time, in seconds, that stored Ansible facts are considered valid since the last time they were modified. Only valid, non-stale, facts will be accessible by a playbook. Note, this does not influence the deletion of ansible_facts from the database. Use a value of 0 to indicate that no timeout should be imposed.
- `API_400_ERROR_LOG_FORMAT` (Optional, Computed) The format of logged messages when an API 4XX error occurs, the following variables will be substituted: 
status_code - The HTTP status code of the error
user_name - The user name attempting to use the API
url_path - The URL path to the API endpoint called
remote_addr - The remote address seen for the user
error - The error set by the api endpoint
Variables need to be in the format {<variable name>}.
- `AUTH_BASIC_ENABLED` (Optional, Computed) Enable HTTP Basic Auth for the API Browser.
- `AUTOMATION_ANALYTICS_GATHER_INTERVAL` (Optional, Computed) Interval (in seconds) between data gathering.
- `AUTOMATION_ANALYTICS_LAST_ENTRIES` (Optional) Managed field from AWX OpenAPI schema.
- `AUTOMATION_ANALYTICS_LAST_GATHER` (Optional) Managed field from AWX OpenAPI schema.
- `AUTOMATION_ANALYTICS_URL` (Optional, Computed) This setting is used to to configure the upload URL for data collection for Automation Analytics.
- `AWX_ANSIBLE_CALLBACK_PLUGINS` (Optional, Computed) List of paths to search for extra callback plugins to be used when running jobs. Enter one path per line.
- `AWX_CLEANUP_PATHS` (Optional, Computed) Enable or Disable TMP Dir cleanup
- `AWX_COLLECTIONS_ENABLED` (Optional, Computed) Allows collections to be dynamically downloaded from a requirements.yml file for SCM projects.
- `AWX_ISOLATION_BASE_PATH` (Optional, Computed) The directory in which the service will create new temporary directories for job execution and isolation (such as credential files).
- `AWX_ISOLATION_SHOW_PATHS` (Optional, Computed) List of paths that would otherwise be hidden to expose to isolated jobs. Enter one path per line. Volumes will be mounted from the execution node to the container. The supported format is HOST-DIR[:CONTAINER-DIR[:OPTIONS]].
- `AWX_MOUNT_ISOLATED_PATHS_ON_K8S` (Optional, Computed) Expose paths via hostPath for the Pods created by a Container Group. HostPath volumes present many security risks, and it is a best practice to avoid the use of HostPaths when possible.
- `AWX_REQUEST_PROFILE` (Optional, Computed) Debug web request python timing
- `AWX_ROLES_ENABLED` (Optional, Computed) Allows roles to be dynamically downloaded from a requirements.yml file for SCM projects.
- `AWX_RUNNER_KEEPALIVE_SECONDS` (Optional, Computed) Only applies to jobs running in a Container Group. If not 0, send a message every so-many seconds to keep connection open.
- `AWX_SHOW_PLAYBOOK_LINKS` (Optional, Computed) Follow symbolic links when scanning for playbooks. Be aware that setting this to True can lead to infinite recursion if a link points to a parent directory of itself.
- `AWX_TASK_ENV` (Optional, Computed) Additional environment variables set for playbook runs, inventory updates, project updates, and notification sending.
- `BULK_HOST_MAX_CREATE` (Optional, Computed) Max number of hosts to allow to be created in a single bulk action
- `BULK_HOST_MAX_DELETE` (Optional, Computed) Max number of hosts to allow to be deleted in a single bulk action
- `BULK_JOB_MAX_LAUNCH` (Optional, Computed) Max jobs to allow bulk jobs to launch
- `CLEANUP_HOST_METRICS_LAST_TS` (Optional) Managed field from AWX OpenAPI schema.
- `CSRF_TRUSTED_ORIGINS` (Optional, Computed) If the service is behind a reverse proxy/load balancer, use this setting to configure the schema://addresses from which the service should trust Origin header values.
- `CUSTOM_LOGIN_INFO` (Optional) If needed, you can add specific information (such as a legal notice or a disclaimer) to a text box in the login modal using this setting. Any content added must be in plain text or an HTML fragment, as other markup languages are not supported.
- `CUSTOM_LOGO` (Optional) To set up a custom logo, provide a file that you create. For the custom logo to look its best, use a .png file with a transparent background. GIF, PNG and JPEG formats are supported.
- `CUSTOM_VENV_PATHS` (Optional, Computed) Paths where Tower will look for custom virtual environments (in addition to /var/lib/awx/venv/). Enter one path per line.
- `DEFAULT_CONTAINER_RUN_OPTIONS` (Optional, Computed) List of options to pass to podman run example: ['--network', 'slirp4netns:enable_ipv6=true', '--log-level', 'debug']
- `DEFAULT_EXECUTION_ENVIRONMENT` (Optional) The Execution Environment to be used when one has not been configured for a job template.
- `DEFAULT_INVENTORY_UPDATE_TIMEOUT` (Optional, Computed) Maximum time in seconds to allow inventory updates to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual inventory source will override this.
- `DEFAULT_JOB_IDLE_TIMEOUT` (Optional, Computed) If no output is detected from ansible in this number of seconds the execution will be terminated. Use value of 0 to indicate that no idle timeout should be imposed.
- `DEFAULT_JOB_TIMEOUT` (Optional, Computed) Maximum time in seconds to allow jobs to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual job template will override this.
- `DEFAULT_PROJECT_UPDATE_TIMEOUT` (Optional, Computed) Maximum time in seconds to allow project updates to run. Use value of 0 to indicate that no timeout should be imposed. A timeout set on an individual project will override this.
- `DISABLE_LOCAL_AUTH` (Optional, Computed) Controls whether users are prevented from using the built-in authentication system.
- `EVENT_STDOUT_MAX_BYTES_DISPLAY` (Optional, Computed) Maximum Size of Standard Output in bytes to display for a single job or ad hoc command event. `stdout` will end with `…` when truncated.
- `GALAXY_IGNORE_CERTS` (Optional, Computed) If set to true, certificate validation will not be done when installing content from any Galaxy server.
- `GALAXY_TASK_ENV` (Optional, Computed) Additional environment variables set for invocations of ansible-galaxy within project updates. Useful if you must use a proxy server for ansible-galaxy but not git.
- `HOST_METRIC_SUMMARY_TASK_LAST_TS` (Optional) Managed field from AWX OpenAPI schema.
- `INSIGHTS_TRACKING_STATE` (Optional, Computed) Enables the service to gather data on automation and send it to Automation Analytics.
- `LOGIN_REDIRECT_OVERRIDE` (Optional) URL to which unauthorized users will be redirected to log in.  If blank, users will be sent to the login page.
- `LOG_AGGREGATOR_ACTION_MAX_DISK_USAGE_GB` (Optional, Computed) Amount of data to store (in gigabytes) if an rsyslog action takes time to process an incoming message (defaults to 1). Equivalent to the rsyslogd queue.maxdiskspace setting on the action (e.g. omhttp). It stores files in the directory specified by LOG_AGGREGATOR_MAX_DISK_USAGE_PATH.
- `LOG_AGGREGATOR_ACTION_QUEUE_SIZE` (Optional, Computed) Defines how large the rsyslog action queue can grow in number of messages stored. This can have an impact on memory utilization. When the queue reaches 75% of this number, the queue will start writing to disk (queue.highWatermark in rsyslog). When it reaches 90%, NOTICE, INFO, and DEBUG messages will start to be discarded (queue.discardMark with queue.discardSeverity=5).
- `LOG_AGGREGATOR_ENABLED` (Optional, Computed) Enable sending logs to external log aggregator.
- `LOG_AGGREGATOR_HOST` (Optional) Hostname/IP where external logs will be sent to.
- `LOG_AGGREGATOR_INDIVIDUAL_FACTS` (Optional, Computed) If set, system tracking facts will be sent for each package, service, or other item found in a scan, allowing for greater search query granularity. If unset, facts will be sent as a single dictionary, allowing for greater efficiency in fact processing.
- `LOG_AGGREGATOR_LEVEL` (Optional, Computed) Level threshold used by log handler. Severities from lowest to highest are DEBUG, INFO, WARNING, ERROR, CRITICAL. Messages less severe than the threshold will be ignored by log handler. (messages under category awx.anlytics ignore this setting)

* `DEBUG` - DEBUG
* `INFO` - INFO
* `WARNING` - WARNING
* `ERROR` - ERROR
* `CRITICAL` - CRITICAL
- `LOG_AGGREGATOR_LOGGERS` (Optional, Computed) List of loggers that will send HTTP logs to the collector, these can include any or all of: 
awx - service logs
activity_stream - activity stream records
job_events - callback data from Ansible job events
system_tracking - facts gathered from scan jobs
broadcast_websocket - errors pertaining to websockets broadcast metrics
job_lifecycle - logs related to processing of a job
- `LOG_AGGREGATOR_MAX_DISK_USAGE_PATH` (Optional, Computed) Location to persist logs that should be retried after an outage of the external log aggregator (defaults to /var/lib/awx). Equivalent to the rsyslogd queue.spoolDirectory setting.
- `LOG_AGGREGATOR_PASSWORD` (Optional, Sensitive) Password or authentication token for external log aggregator (if required; HTTP/s only).
- `LOG_AGGREGATOR_PORT` (Optional) Port on Logging Aggregator to send logs to (if required and not provided in Logging Aggregator).
- `LOG_AGGREGATOR_PROTOCOL` (Optional, Computed) Protocol used to communicate with log aggregator.  HTTPS/HTTP assumes HTTPS unless http:// is explicitly used in the Logging Aggregator hostname.

* `https` - HTTPS/HTTP
* `tcp` - TCP
* `udp` - UDP
- `LOG_AGGREGATOR_RSYSLOGD_DEBUG` (Optional, Computed) Enabled high verbosity debugging for rsyslogd.  Useful for debugging connection issues for external log aggregation.
- `LOG_AGGREGATOR_TCP_TIMEOUT` (Optional, Computed) Number of seconds for a TCP connection to external log aggregator to timeout. Applies to HTTPS and TCP log aggregator protocols.
- `LOG_AGGREGATOR_TOWER_UUID` (Optional) Useful to uniquely identify instances.
- `LOG_AGGREGATOR_TYPE` (Optional) Format messages for the chosen log aggregator.

* `logstash` - logstash
* `splunk` - splunk
* `loggly` - loggly
* `sumologic` - sumologic
* `other` - other
- `LOG_AGGREGATOR_USERNAME` (Optional) Username for external log aggregator (if required; HTTP/s only).
- `LOG_AGGREGATOR_VERIFY_CERT` (Optional, Computed) Flag to control enable/disable of certificate verification when LOG_AGGREGATOR_PROTOCOL is "https". If enabled, the log handler will verify certificate sent by external log aggregator before establishing connection.
- `MANAGE_ORGANIZATION_AUTH` (Optional, Computed) Controls whether any Organization Admin has the privileges to create and manage users and teams.
- `MAX_FORKS` (Optional, Computed) Saving a Job Template with more than this number of forks will result in an error. When set to 0, no limit is applied.
- `MAX_UI_JOB_EVENTS` (Optional, Computed) Maximum number of job events for the UI to retrieve within a single request.
- `MAX_WEBSOCKET_EVENT_RATE` (Optional, Computed) Maximum number of messages to update the UI live job output with per second. Value of 0 means no limit.
- `OPA_AUTH_CA_CERT` (Optional) The content of the CA certificate for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `OPA_AUTH_CLIENT_CERT` (Optional) The content of the client certificate file for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `OPA_AUTH_CLIENT_KEY` (Optional) The content of the client key for mTLS authentication to the OPA server. Required when OPA_AUTH_TYPE is "Certificate".
- `OPA_AUTH_CUSTOM_HEADERS` (Optional, Computed) Optional custom headers included in requests to the OPA server. Defaults to empty dictionary ({}).
- `OPA_AUTH_TOKEN` (Optional, Sensitive) The token for authentication to the OPA server. Required when OPA_AUTH_TYPE is "Token". If an authorization header is defined in OPA_AUTH_CUSTOM_HEADERS, it will be overridden by OPA_AUTH_TOKEN.
- `OPA_AUTH_TYPE` (Optional, Computed) The authentication type that will be used to connect to the OPA server: "None", "Token", or "Certificate".

* `None` - None
* `Token` - Token
* `Certificate` - Certificate
- `OPA_HOST` (Optional) The hostname used to connect to the OPA server. If empty, policy enforcement will be disabled.
- `OPA_PORT` (Optional, Computed) The port used to connect to the OPA server. Defaults to 8181.
- `OPA_REQUEST_RETRIES` (Optional, Computed) The number of retry attempts for connecting to the OPA server. Default is 2.
- `OPA_REQUEST_TIMEOUT` (Optional, Computed) The number of seconds after which the connection to the OPA server will time out. Defaults to 1.5 seconds.
- `OPA_SSL` (Optional, Computed) Enable or disable the use of SSL to connect to the OPA server. Defaults to false.
- `ORG_ADMINS_CAN_SEE_ALL_USERS` (Optional, Computed) Controls whether any Organization Admin can view all users and teams, even those not associated with their Organization.
- `PROJECT_UPDATE_VVV` (Optional, Computed) Adds the CLI -vvv flag to ansible-playbook runs of project_update.yml used for project updates.
- `PROXY_IP_ALLOWED_LIST` (Optional, Computed) If the service is behind a reverse proxy/load balancer, use this setting to configure the proxy IP addresses from which the service should trust custom REMOTE_HOST_HEADERS header values. If this setting is an empty list (the default), the headers specified by REMOTE_HOST_HEADERS will be trusted unconditionally')
- `RECEPTOR_KEEP_WORK_ON_ERROR` (Optional, Computed) Prevent receptor work from being released on when error is detected
- `RECEPTOR_RELEASE_WORK` (Optional, Computed) Release receptor work
- `REDHAT_PASSWORD` (Optional, Sensitive) Client secret used to send data to Automation Analytics
- `REDHAT_USERNAME` (Optional) Client ID used to send data to Automation Analytics
- `REMOTE_HOST_HEADERS` (Optional, Computed) HTTP headers and meta keys to search to determine remote host name or IP. Add additional items to this list, such as "HTTP_X_FORWARDED_FOR", if behind a reverse proxy. See the "Proxy Support" section of the AAP Installation guide for more details.
- `SCHEDULE_MAX_JOBS` (Optional, Computed) Maximum number of the same job template that can be waiting to run when launching from a schedule before no more are created.
- `SESSIONS_PER_USER` (Optional, Computed) Maximum number of simultaneous logged in sessions a user may have. To disable enter -1.
- `SESSION_COOKIE_AGE` (Optional, Computed) Number of seconds that a user is inactive before they will need to login again.
- `STDOUT_MAX_BYTES_DISPLAY` (Optional, Computed) Maximum Size of Standard Output in bytes to display before requiring the output be downloaded.
- `SUBSCRIPTIONS_CLIENT_ID` (Optional) Client ID used to retrieve subscription and content information
- `SUBSCRIPTIONS_CLIENT_SECRET` (Optional, Sensitive) Client secret used to retrieve subscription and content information
- `SUBSCRIPTIONS_PASSWORD` (Optional, Sensitive) Password used to retrieve subscription and content information
- `SUBSCRIPTIONS_USERNAME` (Optional) Username used to retrieve subscription and content information
- `SUBSCRIPTION_USAGE_MODEL` (Optional) * `` - No subscription. Deletion of host_metrics will not be considered for purposes of managed host counting
* `unique_managed_hosts` - Usage based on unique managed nodes in a large historical time frame and delete functionality for no longer used managed nodes
- `TOWER_URL_BASE` (Optional, Computed) This setting is used by services like notifications to render a valid url to the service.
- `UI_LIVE_UPDATES_ENABLED` (Optional, Computed) If disabled, the page will not refresh when events are received. Reloading the page will be required to get the latest details.

## Attributes Reference

- `id` (String) AWX detail-path identifier for this object.

## Import

```bash
terraform import awx_setting.example example
```
