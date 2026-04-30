# Tailscale API - Endpoints Logging & AWS External ID

Extraction exhaustive des endpoints des lignes 1207 a 1485 du fichier `tailscale-openapi.yaml`.

---

## 1. GET /tailnet/{tailnet}/logging/configuration

**Summary :** List configuration audit logs

**Description :**
List all configuration audit logs for a tailnet.
OAuth Scope : `logs:configuration:read`.

**operationId :** `listConfigurationAuditLogs`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Peut etre un dash (`-`) pour referencer le tailnet par defaut du token. Exemple : `example.com` |

#### Query Parameters

| Nom | Type | Requis | Description | Exemple |
|-----|------|--------|-------------|---------|
| `start` | string | Oui | The start of the time window for which to retrieve logs, in RFC 3339 format. | `2023-12-19T16:39:57-08:00` |
| `end` | string | Oui | The end of the time window for which to retrieve logs, in RFC 3339 format. | `2023-12-22T02:15:23-08:00` |
| `actor` | array of string | Non | List of filters on actors, either exact actor IDs or a wildcard search on login name or display name indicated as `~search`. | `["uc4p8fRHvJ11DEVEL", "~bob"]` |
| `target` | array of string | Non | List of target elements for which to filter, attempts to match any part of any of the targets to any of the given strings. | `["mytarget1", "sometarget2"]` |
| `event` | array of string | Non | List of events for which to filter. Valeurs possibles (enum) : `ADMIN_CONSOLE.LOGIN`, `ADMIN_CONSOLE.LOGOUT`, `API_KEY.CREATE`, `API_KEY.EXPIRED`, `API_KEY.REVOKE`, `BILLING.CANCEL.SUBSCRIPTION`, `BILLING.CREATE.SUBSCRIPTION`, `BILLING.UPDATE.ADDRESS`, `BILLING.UPDATE.BILLING_OWNER`, `BILLING.UPDATE.EMAIL`, `BILLING.UPDATE.PAYMENT_INFO`, `BILLING.UPDATE.STRIPE_CUSTOMER_ID`, `BILLING.UPDATE.SUBSCRIPTION`, `FAILED_REQUEST.UPDATE`, `GROUP.PUSH_GROUP.ATTRIBUTES`, `INVITE.ACCEPT.FEATURE`, `INVITE.ACCEPT.NODE_SHARE`, `INVITE.ACCEPT.TAILNET_INVITE`, `INVITE.CREATE.FEATURE`, `INVITE.CREATE.NODE_SHARE`, `INVITE.CREATE.TAILNET_INVITE`, `INVITE.DELETE.NODE_SHARE`, `INVITE.DELETE.TAILNET_INVITE`, `INVITE.RESEND.NODE_SHARE`, `INVITE.RESEND.TAILNET_INVITE`, `NODE.APPROVE`, `NODE.CREATE`, `NODE.CREATE.ATTRIBUTES`, `NODE.DELETE`, `NODE.DELETE.ATTRIBUTES`, `NODE.DISABLE.KEY_EXPIRY`, `NODE.DISCONNECT_NODE.CLIENT_LOG`, `NODE.ENABLE.KEY_EXPIRY`, `NODE.EXPIRED.KEY_EXPIRY_TIME`, `NODE.LOGIN`, `NODE.LOGOUT`, `NODE.REVOKE`, `NODE.UPDATE.ACL_TAGS`, `NODE.UPDATE.ALLOWED_IPS`, `NODE.UPDATE.ATTRIBUTES`, `NODE.UPDATE.AUTO_APPROVED_ROUTES`, `NODE.UPDATE.EXIT_NODE`, `NODE.UPDATE.KEY_EXPIRY_TIME`, `NODE.UPDATE.MACHINE_NAME`, `NODE.UPDATE.POSTURE_IDENTITY`, `NODE.UPDATE.TKA`, `SHARE.CREATE`, `SHARE.DELETE`, `SHARE.UPDATE`, `TAILNET.ACCEPT.FEATURE`, `TAILNET.CREATE`, `TAILNET.CREATE.LOGSTREAM_ENDPOINT`, `TAILNET.CREATE.POSTURE_INTEGRATION`, `TAILNET.CREATE.TKA`, `TAILNET.DELETE.LOGSTREAM_ENDPOINT`, `TAILNET.DELETE.POSTURE_INTEGRATION`, `TAILNET.DELETE.TKA`, `TAILNET.DISABLE.COLLECT_POSTURE_IDENTITY`, `TAILNET.DISABLE.COLLECT_SERVICES`, `TAILNET.DISABLE.FILE_SHARING`, `TAILNET.DISABLE.GEOSTEERING`, `TAILNET.DISABLE.HTTPS`, `TAILNET.DISABLE.LOG_EXIT_FLOWS`, `TAILNET.DISABLE.MACHINE_APPROVAL_NEEDED`, `TAILNET.DISABLE.MAGIC_DNS`, `TAILNET.DISABLE.MULLVAD_VPN`, `TAILNET.DISABLE.NETWORK_FLOW_LOGGING`, `TAILNET.DISABLE.SCIM`, `TAILNET.DISABLE.TKA`, `TAILNET.DISABLE.USER_APPROVAL_REQUIRED`, `TAILNET.ENABLE.COLLECT_POSTURE_IDENTITY`, `TAILNET.ENABLE.COLLECT_SERVICES`, `TAILNET.ENABLE.FILE_SHARING`, `TAILNET.ENABLE.GEOSTEERING`, `TAILNET.ENABLE.HTTPS`, `TAILNET.ENABLE.LOG_EXIT_FLOWS`, `TAILNET.ENABLE.MACHINE_APPROVAL_NEEDED`, `TAILNET.ENABLE.MAGIC_DNS`, `TAILNET.ENABLE.MULLVAD_VPN`, `TAILNET.ENABLE.NETWORK_FLOW_LOGGING`, `TAILNET.ENABLE.SCIM`, `TAILNET.ENABLE.TKA`, `TAILNET.ENABLE.USER_APPROVAL_REQUIRED`, `TAILNET.JOIN`, `TAILNET.JOIN_WAITLIST.FEATURE`, `TAILNET.LEAVE`, `TAILNET.UPDATE.ACCOUNT_EMAIL`, `TAILNET.UPDATE.ACL`, `TAILNET.UPDATE.DNS_CONFIG`, `TAILNET.UPDATE.LOGSTREAM_ENDPOINT`, `TAILNET.UPDATE.MAX_KEY_DURATION`, `TAILNET.UPDATE.POSTURE_INTEGRATION`, `TAILNET.UPDATE.SECURITY_EMAIL`, `TAILNET.UPDATE.SUPPORT_EMAIL`, `TAILNET.UPDATE.TCD`, `TAILNET.UPDATE.TKA`, `TAILNET.VERIFY.ACCOUNT_EMAIL`, `TAILNET.VERIFY.SECURITY_EMAIL`, `TAILNET.VERIFY.SUPPORT_EMAIL`, `USER.APPROVE`, `USER.CREATE`, `USER.DELETE`, `USER.INVITE`, `USER.PUSH_USER.ATTRIBUTES`, `USER.RESEND.TAILNET_INVITE`, `USER.RESTORE`, `USER.RESTORE_GLOBAL`, `USER.SUSPEND`, `USER.SUSPEND_GLOBAL`, `USER.UPDATE.USER_ROLE`, `WEBHOOK_ENDPOINT.CREATE`, `WEBHOOK_ENDPOINT.DELETE`, `WEBHOOK_ENDPOINT.UPDATE.SECRET`, `WEBHOOK_ENDPOINT.UPDATE.SUBSCRIBED_EVENTS`, `WEB_INTERFACE.LOGIN`, `WEB_INTERFACE.LOGOUT` | `["USER.CREATE", "NODE.CREATE"]` |

### Reponse 200 - Successful operation

```json
{
  "version": "string",       // Version of audit logs response. Exemple : "1.1"
  "tailnet": "string",       // The tailnet on which the logged configuration changes were made. Exemple : "example.com"
  "logs": [                  // Matching log entries, ordered chronologically.
    {
      // Schema: ConfigurationAuditLog (voir detail ci-dessous)
    }
  ]
}
```

#### Schema ConfigurationAuditLog

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `eventTime` | string | Oui | Timestamp of the audit log event, in RFC 3339 format. Exemple : `2024-06-06T15:25:26.583893Z` |
| `type` | string (enum: `CONFIG`) | Oui | The type of log (always "CONFIG"). |
| `deferredAt` | string | Non | Timestamp recording the time that the audit log rate limiter enqueued the record to be logged at a future time, in RFC 3339 format. |
| `eventGroupID` | string | Oui | Identifier assigned to one or more audit log events, all of which are the result of a single operation. Exemple : `0378d8f57300d172ef7ae3826e097ef0` |
| `origin` | string | Oui | The initiator of the action. Enum : `ADMIN_CONSOLE`, `CONFIG_API`, `CONTROL`, `IDENTITY_PROVIDER`, `NODE`, `SUPPORT_REQUEST`, `STRIPE`, `SECURITY_NOTIFICATION`, `LEGAL_NOTIFICATION` |
| `actor` | object | Oui | The person who caused the action related to this event. |
| `actor.id` | string | - | The ID (user ID or node ID) of the actor. Exemple : `uZKk3KSfrH11DEVEL` |
| `actor.type` | string | - | Enum : `USER`, `NODE`, `AUTOMATED_WORKER`, `OAUTH_CLIENT`, `SCIM`, `MULLVAD`, `LOGSTREAM`, `SECRET_SCANNER` |
| `actor.loginName` | string | - | The login name of the actor at time of the action. |
| `actor.displayName` | string | - | The display name of the actor at time of the action. |
| `actor.tags` | array of string | - | Indicates the tags owning a node. Only set if `type` is `NODE`. |
| `target` | object | Oui | The object of this event's action. |
| `target.id` | string | - | The unique ID (user id, tailnet SID, or node id) of the target. |
| `target.name` | string | - | Name of the entity at time of the action. |
| `target.type` | string | - | Enum : `TAILNET`, `USER`, `GROUP`, `NODE`, `API_KEY`, `INVITE`, `SHARE`, `BILLING`, `ADMIN_CONSOLE`, `WEB_INTERFACE`, `WEBHOOK_ENDPOINT`, `FAILED_REQUEST` |
| `target.isEphemeral` | boolean | - | Indicates whether the target is ephemeral. Only set if `type` is `NODE`. |
| `target.property` | string | - | The property name on this target which was updated by the event. Enum : `ACL`, `ACL_TAGS`, `ACCOUNT_EMAIL`, `ADDRESS`, `ALLOWED_IPS`, `AUTO_APPROVED_ROUTES`, `ATTRIBUTES`, `BILLING_OWNER`, `COLLECT_SERVICES`, `COLLECT_POSTURE_IDENTITY`, `MULLVAD_VPN`, `DNS_CONFIG`, `EMAIL`, `EXIT_NODE`, `FEATURE`, `FILE_SHARING`, `HTTPS`, `KEY_EXPIRY_TIME`, `KEY_EXPIRY`, `LOG_EXIT_FLOWS`, `LOGSTREAM_ENDPOINT`, `MAGIC_DNS`, `MACHINE_AUTH_NEEDED`, `MACHINE_APPROVAL_NEEDED`, `USER_APPROVAL_REQUIRED`, `MACHINE_NAME`, `MAX_KEY_DURATION`, `NETWORK_FLOW_LOGGING`, `GEOSTEERING`, `NODE_SHARE`, `TAILNET_INVITE`, `PAYMENT_INFO`, `POSTURE_IDENTITY`, `POSTURE_INTEGRATION`, `USER_ROLE`, `SCIM`, `SECURITY_EMAIL`, `STRIPE_CUSTOMER_ID`, `SUBSCRIPTION`, `SUBSCRIBED_EVENTS`, `SUPPORT_EMAIL`, `SECRET`, `TCD`, `TKA`, `AUTH_PROVIDER` |
| `action` | string | Oui | The type of change attempted against the target. Enum : `LOGIN`, `LOGOUT`, `CREATE`, `UPDATE`, `DELETE`, `CANCEL`, `REVOKE`, `APPROVE`, `SUSPEND`, `RESTORE`, `ENABLE`, `DISABLE`, `ACCEPT`, `EXPIRED`, `PUSH_USER`, `PUSH_GROUP`, `VERIFY`, `JOIN_WAITLIST`, `INVITE`, `JOIN`, `LEAVE`, `RESEND`, `MIGRATE_AUTH_PROVIDER` |
| `old` | string / number / integer / boolean / array / object | Non | The value of `target.property` prior to the event. |
| `new` | string / number / integer / boolean / array / object | Non | The value of `target.property` after the event. |
| `actionDetails` | string | Non | Additional information about the event, such as a client-provided reason, if it exists. |
| `error` | string | Non | Provided when the configuration change failed to be completed. It is a user-presentable reason for the failure. |

**Champs requis de ConfigurationAuditLog :** `eventTime`, `type`, `eventGroupID`, `origin`, `actor`, `target`, `action`

### Autres reponses

| Code | Description |
|------|-------------|
| 400 | Request has missing or invalid parameter(s). |
| 403 | User does not have sufficient access to view configuration audit logs. |
| 404 | Logging is not supported on this deployment of Tailscale. |

---

## 2. GET /tailnet/{tailnet}/logging/network

**Summary :** List network flow logs

**Description :**
List all network flow logs for a tailnet.
OAuth Scope : `logs:network:read`.

**operationId :** `listNetworkFlowLogs`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |

#### Query Parameters

| Nom | Type | Requis | Description | Exemple |
|-----|------|--------|-------------|---------|
| `start` | string | Oui | The start of the time window for which to retrieve logs, in RFC 3339 format. | `2023-12-19T16:39:57-08:00` |
| `end` | string | Oui | The end of the time window for which to retrieve logs, in RFC 3339 format. | `2023-12-22T02:15:23-08:00` |

### Reponse 200 - Successful operation

```json
{
  "logs": [                  // Matching log entries, ordered chronologically.
    {
      // Schema: NetworkFlowLog (voir detail ci-dessous)
    }
  ]
}
```

#### Schema NetworkFlowLog

| Champ | Type | Description | Exemple |
|-------|------|-------------|---------|
| `logged` | string | Timestamp of the flow log, in RFC 3339 format. | `2024-06-06T15:27:26.583893Z` |
| `nodeId` | string | Identifier of the node. | `nBLYviWLGB21DEVEL` |
| `start` | string | Time at which flow started, in RFC 3339 format. | `2024-06-06T15:25:26.583893Z` |
| `end` | string | Time at which flow ended, in RFC 3339 format. | `2024-06-06T15:26:26.583893Z` |
| `virtualTraffic` | array of ConnectionCounts | Traffic virtuel. | |
| `subnetTraffic` | array of ConnectionCounts | Traffic sous-reseau. | |
| `exitTraffic` | array of ConnectionCounts | Traffic de sortie. | |
| `physicalTraffic` | array of ConnectionCounts | Traffic physique. | |

#### Schema ConnectionCounts

| Champ | Type | Description | Exemple |
|-------|------|-------------|---------|
| `proto` | string | IP protocol name (or number if no name used). Enum : `ah`, `dccp`, `egp`, `esp`, `gre`, `icmp`, `igmp`, `igp`, `ipv4`, `ipv6-icmp`, `sctp`, `tcp`, `udp` | `ipv4` |
| `src` | string | Source addr:port. | `108.86.185.125:52343` |
| `dst` | string | Destination addr:port. | `108.86.185.126:443` |
| `txPkts` | integer | Number of packets sent. | `10` |
| `txBytes` | integer | Number of bytes sent. | `1000` |
| `rxPkts` | integer | Number of packets received. | `10` |
| `rxBytes` | integer | Number of bytes received. | `1000` |

### Autres reponses

| Code | Description |
|------|-------------|
| 400 | Request has missing or invalid parameter(s). |
| 403 | User does not have sufficient access to view network flow logs. |
| 404 | Logging is not supported on this deployment of Tailscale. |
| 502 | The system was unable to communicate with logging server. |

---

## 3. GET /tailnet/{tailnet}/logging/{logType}/stream/status

**Summary :** Get log streaming status

**Description :**
Retrieve the log streaming status for the provided log type.
OAuth Scope : `log_streaming:read`.

**operationId :** `getLogStreamingStatus`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |
| `logType` | string | Oui | The type of log. Enum : `configuration`, `network`. Exemple : `configuration` |

### Reponse 200 - Successful operation

Schema : **LogstreamEndpointPublishingStatus**

```json
{
  "lastActivity": "string",       // Timestamp of the most recent publishing activity, in RFC 3339 format.
  "lastError": "string",          // The most recent error (if any).
  "maxBodySize": 524288,          // integer - The size of the largest single request body.
  "numBytesSent": 17238983,       // integer - Total bytes published across all requests.
  "numEntriesSent": 8363,         // integer - The total number of entries published.
  "numSpoofedEntries": 0,         // integer - The number of spoofed entries published.
  "numTotalRequests": 10610,      // integer - The total number of requests made to the streaming endpoint.
  "numFailedRequests": 5434,      // integer - The total number of requests that have failed.
  "rateBytesSent": 3.524,         // number - EWMA rate of data streamed, in bytes per second.
  "rateEntriesSent": 0.0086,      // number - EWMA rate of entries sent, in entries per second.
  "rateTotalRequests": 0.0037,    // number - EWMA rate of requests, in requests per second.
  "rateFailedRequests": 4.14e-157 // number - EWMA rate of failed requests, in requests per second.
}
```

**Tous les champs sont requis :** `lastActivity`, `lastError`, `maxBodySize`, `numBytesSent`, `numEntriesSent`, `numSpoofedEntries`, `numTotalRequests`, `numFailedRequests`, `rateBytesSent`, `rateEntriesSent`, `rateTotalRequests`, `rateFailedRequests`

### Autres reponses

| Code | Description |
|------|-------------|
| 404 | Log streaming has not been configured, this `logType` is not supported, or user does not have sufficient access to view log streaming status. |
| 502 | The system was unable to communicate with logging server. |

---

## 4. GET /tailnet/{tailnet}/logging/{logType}/stream

**Summary :** Get log streaming configuration

**Description :**
Retrieve the log streaming configuration for the provided log type.
OAuth Scope : `log_streaming:read`.

**operationId :** `getLogStreamingConfiguration`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |
| `logType` | string | Oui | The type of log. Enum : `configuration`, `network`. Exemple : `configuration` |

### Reponse 200 - Successful operation

Schema : **LogstreamEndpointConfiguration** (voir detail complet dans le PUT ci-dessous)

### Autres reponses

| Code | Description |
|------|-------------|
| 404 | Log streaming has not been configured, this `logType` is not supported, or user does not have sufficient access to view log streaming configuration. |

---

## 5. PUT /tailnet/{tailnet}/logging/{logType}/stream

**Summary :** Set log streaming configuration

**Description :**
Set the log streaming configuration for the provided log type.
OAuth Scope : `log_streaming`. `device_invites` and `policy_file` are also required if streaming to a private endpoint.

**operationId :** `setLogStreamingConfiguration`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |
| `logType` | string | Oui | The type of log. Enum : `configuration`, `network`. Exemple : `configuration` |

### Request Body (application/json)

Schema : **LogstreamEndpointConfiguration**

Description : The LogstreamEndpointConfiguration to set. `logType` is specified in the request URL rather than the body.

Exemple de body :
```json
{
  "destinationType": "elastic",
  "url": "http://100.71.134.73:80/config-log-datastream",
  "user": "myusername",
  "token": "mytoken"
}
```

#### Schema LogstreamEndpointConfiguration (complet)

| Champ | Type | Description | Exemple |
|-------|------|-------------|---------|
| `logType` | string (readOnly) | The type of log that is streamed to this endpoint. Enum : `configuration`, `network` | `configuration` |
| `destinationType` | string | The type of system to which logs are being streamed. Enum : `splunk`, `elastic`, `panther`, `cribl`, `datadog`, `axiom`, `s3` | `elastic` |
| `url` | string | The URL to which log streams are being posted. If the DestinationType is `s3`, the URL may be empty to use the official Amazon S3 endpoint. | `http://100.71.134.73:80/config-log-datastream` |
| `user` | string | The username with which log streams to this endpoint are authenticated. | `myusername` |
| `uploadPeriodMinutes` | integer (max: 1440) | An optional number of minutes to wait in between uploading new logs. If the quantity of logs does not fit within a single upload, multiple uploads will be made. | `5` |
| `compressionFormat` | string | The compression algorithm with which to compress logs. `none` disables compression. Defaults to `none`. Enum : `zstd`, `gzip`, `none` | `zstd` |
| `token` | string (writeOnly) | The token/password with which log streams to this endpoint should be authenticated. | `mytoken` |
| `s3Bucket` | string | The S3 bucket name. Required if the destinationType is `s3`. | `mycompany-mybucket` |
| `s3Region` | string | The region in which the S3 bucket is located. Required if the destinationType is `s3`. | `us-east-1` |
| `s3KeyPrefix` | string | An optional S3 key prefix to prepend to the auto-generated S3 key name. | |
| `s3AuthenticationType` | string | What type of authentication to use for S3. Required if the destinationType is `s3`. Enum : `accesskey`, `rolearn` | |
| `s3AccessKeyId` | string | The S3 access key ID. Required if the destinationType is `s3` and `authenticationType` is `accesskey`. | |
| `s3SecretAccessKey` | string (writeOnly) | The S3 secret access key. Required if the destinationType is `s3` and `authenticationType` is `accesskey`. | |
| `s3RoleArn` | string | The Role ARN that Tailscale should supply to AWS when authenticating using role-based authentication. Required if the destinationType is `s3` and `authenticationType` is `rolearn`. | |
| `s3ExternalId` | string (readOnly) | The AWS external id that Tailscale supplies when authenticating using role-based authentication. Populated if the destinationType is `s3` and `authenticationType` is `rolearn`. | |
| `gcsBucket` | string | The GCS bucket name. Required if the destinationType is `gcs`. | |
| `gcsKeyPrefix` | string | An optional GCS key prefix to append to the GCS bucket name. | |
| `gcsScopes` | array of string | The GCS scopes needed to be able to write to the GCS bucket. | |
| `gcsCredentials` | string | The JSON workload identity credentials from GCS needed for accessing the GCS account. | |

### Reponse 200 - Successful operation

Pas de corps de reponse specifique (confirmation de succes).

### Autres reponses

| Code | Description |
|------|-------------|
| 400 | Request has missing or invalid parameter(s). |
| 403 | User does not have sufficient access to update log streaming configuration. |
| 404 | Tailnet not found, this `logType` is not supported, or user does not have sufficient access to view log streaming configuration. |

---

## 6. DELETE /tailnet/{tailnet}/logging/{logType}/stream

**Summary :** Disable log streaming

**Description :**
Delete the log streaming configuration for the provided log type.
OAuth Scope : `log_streaming`.

**operationId :** `disableLogStreaming`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |
| `logType` | string | Oui | The type of log. Enum : `configuration`, `network`. Exemple : `configuration` |

### Reponse 200 - Successful operation

Pas de corps de reponse specifique (confirmation de succes).

### Autres reponses

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to update log streaming configuration. |
| 404 | Log streaming has not been configured, this `logType` is not supported, or user does not have sufficient access to view log streaming configuration. |

---

## 7. POST /tailnet/{tailnet}/aws-external-id

**Summary :** Create or get AWS external id

**Description :**
Get an AWS external id to use for streaming tailnet logs to S3 using role-based authentication, creating a new one for this tailnet when necessary.
OAuth Scope : `log_streaming`.

**operationId :** `getAwsExternalId`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description |
|-----|------|--------|-------------|
| `tailnet` | string | Oui | The tailnet ID. Exemple : `example.com` |

### Request Body (application/json)

```json
{
  "reusable": true   // boolean - If set to true, this same AWS external id will be returned
                     // on future calls to this endpoint, if and only if those calls also
                     // mark `reusable` as true, and the ID has not yet been linked with an AWS account.
}
```

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `reusable` | boolean | Non | If set to true, this same AWS external id will be returned on future calls to this endpoint, if and only if those calls also mark `reusable` as true, and the ID has not yet been linked with an AWS account. |

### Reponse 200 - Successful operation

Schema : **AwsExternalId**

```json
{
  "externalId": "string",            // The external id that Tailscale will supply to AWS
                                     // when authenticating using role-based authentication.
                                     // Exemple : "60fe9ce7-7791-4ab3-ab34-4294f5972725"
  "tailscaleAwsAccountId": "string"  // The AWS account id that Tailscale will supply to AWS
                                     // when authenticating using role-based authentication.
                                     // Exemple : "001234567890"
}
```

| Champ | Type | Description | Exemple |
|-------|------|-------------|---------|
| `externalId` | string | The external id that Tailscale will supply to AWS when authenticating using role-based authentication. | `60fe9ce7-7791-4ab3-ab34-4294f5972725` |
| `tailscaleAwsAccountId` | string | The AWS account id that Tailscale will supply to AWS when authenticating using role-based authentication. | `001234567890` |

### Autres reponses

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to obtain an AWS external id. |
| 404 | Tailnet not found. |

---

## 8. POST /tailnet/{tailnet}/aws-external-id/{id}/validate-aws-trust-policy

**Summary :** Validate external ID integration with IAM role trust policy

**Description :**
Validate that Tailscale can assume your IAM role with (and only with) this external ID.
OAuth Scope : `log_streaming`.

**operationId :** `validateAwsExternalId`

**Tags :** `Logging`

### Parametres

#### Path Parameters

| Nom | Type | Requis | Description | Exemple |
|-----|------|--------|-------------|---------|
| `tailnet` | string | Oui | The tailnet ID. | `example.com` |
| `id` | string | Oui | The AWS external ID to validate. | `60fe9ce7-7791-4ab3-ab34-4294f5972725` |

### Request Body (application/json)

```json
{
  "roleArn": "arn:aws:iam::000000000000:role/tailscale-log-writer"
}
```

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `roleArn` | string | Non (implicitement requis pour la validation) | ARN of the AWS IAM role to validate with this external ID. |

### Reponse 200 - Validation succeeded

Pas de corps de reponse specifique. Validation succeeded for this external ID and IAM role.

### Autres reponses

| Code | Description | Corps de reponse |
|------|-------------|------------------|
| 403 | User does not have sufficient access for this tailnet. | Standard error response |
| 404 | Tailnet or external ID not found. | Standard error response |
| 422 | Validation failed for this external ID and IAM role. | `{"message": "string"}` -- The reason for validation failure. |

---

## Resume des endpoints

| # | Methode | Path | operationId | Tag |
|---|---------|------|-------------|-----|
| 1 | GET | `/tailnet/{tailnet}/logging/configuration` | `listConfigurationAuditLogs` | Logging |
| 2 | GET | `/tailnet/{tailnet}/logging/network` | `listNetworkFlowLogs` | Logging |
| 3 | GET | `/tailnet/{tailnet}/logging/{logType}/stream/status` | `getLogStreamingStatus` | Logging |
| 4 | GET | `/tailnet/{tailnet}/logging/{logType}/stream` | `getLogStreamingConfiguration` | Logging |
| 5 | PUT | `/tailnet/{tailnet}/logging/{logType}/stream` | `setLogStreamingConfiguration` | Logging |
| 6 | DELETE | `/tailnet/{tailnet}/logging/{logType}/stream` | `disableLogStreaming` | Logging |
| 7 | POST | `/tailnet/{tailnet}/aws-external-id` | `getAwsExternalId` | Logging |
| 8 | POST | `/tailnet/{tailnet}/aws-external-id/{id}/validate-aws-trust-policy` | `validateAwsExternalId` | Logging |
