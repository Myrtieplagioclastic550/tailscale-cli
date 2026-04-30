# Tailscale API - Schemas / Modeles de donnees (OpenAPI)

Extraction exhaustive de tous les schemas definis dans `components.schemas` du fichier OpenAPI Tailscale.

---

## Table des matieres

1. [Device](#1-device)
2. [Error](#2-error)
3. [DeviceRoutes](#3-deviceroutes)
4. [DevicePostureAttributes](#4-devicepostureattributes)
5. [DeviceInvite](#5-deviceinvite)
6. [UserInvite](#6-userinvite)
7. [ConfigurationAuditLog](#7-configurationauditlog)
8. [ConnectionCounts](#8-connectioncounts)
9. [NetworkFlowLog](#9-networkflowlog)
10. [LogType](#10-logtype)
11. [LogstreamEndpointPublishingStatus](#11-logstreamendpointpublishingstatus)
12. [LogstreamEndpointConfiguration](#12-logstreamendpointconfiguration)
13. [AwsExternalId](#13-awsexternalid)
14. [DnsPreferences](#14-dnspreferences)
15. [DnsSearchPaths](#15-dnssearchpaths)
16. [SplitDns](#16-splitdns)
17. [DnsConfigurationResolver](#17-dnsconfigurationresolver)
18. [DnsConfigurationPreferences](#18-dnsconfigurationpreferences)
19. [DnsConfiguration](#19-dnsconfiguration)
20. [KeyCapabilities](#20-keycapabilities)
21. [Key](#21-key)
22. [PostureIntegration](#22-postureintegration)
23. [User](#23-user)
24. [Contact](#24-contact)
25. [Webhook](#25-webhook)
26. [providerType](#26-providertype)
27. [subscriptions](#27-subscriptions)
28. [TailnetSettings](#28-tailnetsettings)
29. [VIPServiceInfo](#29-vipserviceinfo)
30. [VIPServiceInfoPut](#30-vipserviceinfoput)
31. [ServiceHostInfo](#31-servicehostinfo)
32. [VIPServiceApproval](#32-vipserviceapproval)

---

## 1. Device

- **Type** : `object`
- **Description** : A Tailscale device (sometimes referred to as *node* or *machine*), is any computer or mobile device that joins a tailnet. Each device has a unique ID (`nodeId`) that is used to identify the device in API calls.
- **Champs requis** : aucun (pas de `required` explicite)

### Proprietes

| Propriete | Type | Format | Description |
|-----------|------|--------|-------------|
| `addresses` | `array` of `string` | - | A list of Tailscale IP addresses for the device, including both IPv4 (100.x.y.z) and IPv6 (fd7a:115c:a1e0:...) addresses. |
| `id` | `string` | - | The legacy identifier for a device; `nodeId` is preferred. |
| `nodeId` | `string` | - | The preferred identifier for a device; supply this value wherever {deviceId} is indicated in the endpoint. |
| `user` | `string` | - | The user who registered the node. For untagged nodes, this user is the device owner. |
| `name` | `string` | - | The MagicDNS name of the device. |
| `hostname` | `string` | - | The machine name in the admin console. |
| `clientVersion` | `string` | - | The version of the Tailscale client software; empty for external devices. |
| `updateAvailable` | `boolean` | - | `true` if a Tailscale client version upgrade is available. Empty for external devices. |
| `os` | `string` | - | The operating system that the device is running. |
| `created` | `string` | `date-time` | The date on which the device was added to the tailnet; empty for external devices. |
| `connectedToControl` | `boolean` | - | Indicates if the device recently maintained a TCP connection to the Tailscale control server. |
| `lastSeen` | `string` | `date-time` | When the device was last connected to the Tailscale control server. Omitted if the device has never been online or `connectedToControl` is true. |
| `keyExpiryDisabled` | `boolean` | - | `true` if the keys for the device will not expire. |
| `expires` | `string` | `date-time` | The expiration date of the device's auth key. |
| `authorized` | `boolean` | - | `true` if the device has been authorized to join the tailnet; otherwise, `false`. |
| `isExternal` | `boolean` | - | `true` indicates that a device is not a member of the tailnet, but is shared in to the tailnet. |
| `multipleConnections` | `boolean` | - | `true` indicates that multiple devices are currently connected using the same node key, which is usually a sign of node state being copied between machines. Omitted if only one device is connected. |
| `machineKey` | `string` | - | For internal use and is not required for any API operations. Empty for external devices. |
| `nodeKey` | `string` | - | Mostly for internal use, required for select operations, such as adding a node to a locked tailnet. |
| `blocksIncomingConnections` | `boolean` | - | `true` if the device is not allowed to accept any connections over Tailscale, including pings. |
| `enabledRoutes` | `array` of `string` | - | The subnet routes for this device that have been approved by a tailnet admin. |
| `advertisedRoutes` | `array` of `string` | - | The subnets this device requests to expose. |
| `clientConnectivity` | `object` (inline) | - | Provides a report on the device's current physical network conditions. Voir sous-objet ci-dessous. |
| `tags` | `array` of `string` | - | Tags assigned to the device. Once tagged, the tag is the owner of that device. Empty for external devices. |
| `tailnetLockError` | `string` | - | Indicates an issue with the tailnet lock node-key signature on this device. Only populated when tailnet lock is enabled. |
| `tailnetLockKey` | `string` | - | The node's tailnet lock key. Present even if tailnet lock is not enabled. |
| `sshEnabled` | `boolean` | - | `true` if Tailscale SSH is enabled on this device. |
| `postureIdentity` | `object` (inline) | - | Contains extra identifiers from the device when the tailnet has device posture identification collection enabled. Voir sous-objet ci-dessous. |
| `isEphemeral` | `boolean` | - | `true` if the device is ephemeral. |
| `distro` | `object` (inline) | - | Provides details of the operating system distribution, if available. Voir sous-objet ci-dessous. |

### Sous-objet : `clientConnectivity`

| Propriete | Type | Description |
|-----------|------|-------------|
| `endpoints` | `array` of `string` | Client's magicsock UDP IP:port endpoints (IPv4 or IPv6). |
| `mappingVariesByDestIP` | `boolean` | `true` if the host's NAT mappings vary based on the destination IP. |
| `latency` | `object` (additionalProperties) | Map of DERP server locations (string keys) to objects with `preferred` (boolean) and `latencyMs` (number/float). |
| `clientSupports` | `object` | Identifies features supported by the client. |

#### Sous-objet : `clientConnectivity.clientSupports`

| Propriete | Type | Description |
|-----------|------|-------------|
| `hairPinning` | `boolean` ou `null` | This information is no longer tracked and will always be null. |
| `ipv6` | `boolean` ou `null` | `true` if the device OS supports IPv6, regardless of whether IPv6 internet connectivity is available. |
| `pcp` | `boolean` ou `null` | `true` if PCP port-mapping service exists on your router. |
| `pmp` | `boolean` ou `null` | `true` if NAT-PMP port-mapping service exists on your router. |
| `udp` | `boolean` ou `null` | `true` if UDP traffic is enabled on the current network. |
| `upnp` | `boolean` ou `null` | `true` if UPnP port-mapping service exists on your router. |

### Sous-objet : `postureIdentity`

| Propriete | Type | Description |
|-----------|------|-------------|
| `serialNumbers` | `array` of `string` | Serial numbers of the device. |
| `disabled` | `boolean` | Indicates if posture identification collection is disabled. |

### Sous-objet : `distro`

| Propriete | Type | Description |
|-----------|------|-------------|
| `name` | `string` | The operating system distribution (e.g., "ubuntu"). |
| `version` | `string` | The OS distribution version (e.g., "25.04"). |
| `codeName` | `string` | The OS distribution code name (e.g., "Plucky Puffin"). |

---

## 2. Error

- **Type** : `object`
- **Champs requis** : `message`

### Proprietes

| Propriete | Type | Requis | Description |
|-----------|------|--------|-------------|
| `message` | `string` | Oui | Message d'erreur. |

---

## 3. DeviceRoutes

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `advertisedRoutes` | `array` of `string` | The subnets this device requests to expose. |
| `enabledRoutes` | `array` of `string` | The subnet routes for this device that have been approved by a tailnet admin. |

---

## 4. DevicePostureAttributes

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `attributes` | `object` (additionalProperties: `string` ou `number` ou `boolean`) | Contains all the posture attributes assigned to a node. Attribute values can be strings, numbers or booleans. |
| `expiries` | `object` (additionalProperties: `string`, format `date-time`) | Contains the expiry time for each posture attribute, if set. |

---

## 5. DeviceInvite

- **Type** : `object`
- **Description** : A device invite is an invitation that shares a device with an external user (a user not in the device's tailnet).
- **Champs requis** : aucun (pas de `required` explicite)

### Proprietes

| Propriete | Type | Format | Description |
|-----------|------|--------|-------------|
| `id` | `string` | - | The unique identifier for the invite. Supply this value wherever {deviceInviteId} is indicated in the endpoint. |
| `created` | `string` | `date-time` | The creation time of the invite. |
| `tailnetId` | `integer` | `int64` | The ID of the tailnet to which the shared device belongs. |
| `deviceId` | `integer` | `int64` | The ID of the device being shared. |
| `sharerId` | `integer` | `int64` | The ID of the user who created the share invite. |
| `multiUse` | `boolean` | - | Specifies whether this device invite can be accepted more than once. |
| `allowExitNode` | `boolean` | - | Specifies whether the invited user is able to use the device as an exit node when the device is advertising as one. |
| `email` | `string` | - | The email to which the invite was sent. If empty, the invite was not emailed to anyone, but the inviteUrl can be shared manually. |
| `lastEmailSentAt` | `string` | `date-time` | The last time the invite was attempted to be sent to Email. Only ever set if `email` is not empty. |
| `inviteUrl` | `string` | - | The link to accept the invite. Anyone with this link can accept the invite. |
| `accepted` | `boolean` | - | `true` when the share invite has been accepted. |
| `acceptedBy` | `object` (inline) | - | Set when the invite has been accepted. Holds information about the user who accepted the share invite. |

### Sous-objet : `acceptedBy`

| Propriete | Type | Format | Description |
|-----------|------|--------|-------------|
| `id` | `integer` | `int64` | The ID of the user who accepted the share invite. |
| `loginName` | `string` | - | The login name of the user who accepted the share invite. |
| `profilePicUrl` | `string` | - | The profile pic URL for the user who accepted the share invite. |

---

## 6. UserInvite

- **Type** : `object`
- **Description** : A user invite is an active invitation that lets a user join a tailnet with a preassigned user role.
- **Champs requis** : `id`, `role`, `tailnetId`, `inviterId`

### Proprietes

| Propriete | Type | Format | Requis | Description |
|-----------|------|--------|--------|-------------|
| `id` | `string` | - | Oui | The unique identifier for the invite. |
| `role` | `string` (enum) | - | Oui | The tailnet user role to assign to the invited user upon accepting the invite. |
| `tailnetId` | `integer` | `int64` | Oui | The ID of the tailnet to which the user was invited. |
| `inviterId` | `integer` | `int64` | Oui | The ID of the user who created the invite. |
| `email` | `string` | - | Non | The email to which the invite was sent. If empty, the invite was not emailed. |
| `lastEmailSentAt` | `string` | `date-time` | Non | The last time the invite was attempted to be sent. Only set if `email` is not empty. |
| `inviteUrl` | `string` | - | Non | Included when `email` is not part of the tailnet's domain, or when `email` is empty. Link to accept the invite. |

### Enum `role`

- `member`
- `admin`
- `it-admin`
- `network-admin`
- `billing-admin`
- `auditor`

---

## 7. ConfigurationAuditLog

- **Type** : `object`
- **Champs requis** : `eventTime`, `type`, `eventGroupID`, `origin`, `actor`, `target`, `action`

### Proprietes

| Propriete | Type | Format | Requis | Description |
|-----------|------|--------|--------|-------------|
| `eventTime` | `string` | - | Oui | Timestamp of the audit log event, in RFC 3339 format. |
| `type` | `string` (enum) | - | Oui | The type of log (always "CONFIG"). |
| `deferredAt` | `string` | - | Non | Timestamp recording the time that the audit log rate limiter enqueued the record. |
| `eventGroupID` | `string` | - | Oui | Identifier assigned to one or more audit log events, all of which are the result of a single operation. |
| `origin` | `string` (enum) | - | Oui | The initiator of the action that generated the event. |
| `actor` | `object` (inline) | - | Oui | The person who caused the action related to this event. |
| `target` | `object` (inline) | - | Oui | The object of this event's action. |
| `action` | `string` (enum) | - | Oui | The type of change attempted against the `target`. |
| `old` | `anyOf` (string, number, integer, boolean, array, object) | - | Non | The value of `target.property` prior to the event. |
| `new` | `anyOf` (string, number, integer, boolean, array, object) | - | Non | The value of `target.property` after the event. |
| `actionDetails` | `string` | - | Non | Additional information about the event, such as a client-provided reason. |
| `error` | `string` | - | Non | Provided when the configuration change failed to be completed. User-presentable reason for the failure. |

### Enum `type`

- `CONFIG`

### Enum `origin`

- `ADMIN_CONSOLE`
- `CONFIG_API`
- `CONTROL`
- `IDENTITY_PROVIDER`
- `NODE`
- `SUPPORT_REQUEST`
- `STRIPE`
- `SECURITY_NOTIFICATION`
- `LEGAL_NOTIFICATION`

### Enum `action`

- `LOGIN`
- `LOGOUT`
- `CREATE`
- `UPDATE`
- `DELETE`
- `CANCEL`
- `REVOKE`
- `APPROVE`
- `SUSPEND`
- `RESTORE`
- `ENABLE`
- `DISABLE`
- `ACCEPT`
- `EXPIRED`
- `PUSH_USER`
- `PUSH_GROUP`
- `VERIFY`
- `JOIN_WAITLIST`
- `INVITE`
- `JOIN`
- `LEAVE`
- `RESEND`
- `MIGRATE_AUTH_PROVIDER`

### Sous-objet : `actor`

| Propriete | Type | Description |
|-----------|------|-------------|
| `id` | `string` | The ID (user ID or node ID) of the actor. |
| `type` | `string` (enum) | The entity type of the actor. |
| `loginName` | `string` | The login name of the actor at time of the action. |
| `displayName` | `string` | The display name of the actor at time of the action. |
| `tags` | `array` of `string` | Indicates the tags owning a node. Only set if `type` is `NODE`. |

#### Enum `actor.type`

- `USER`
- `NODE`
- `AUTOMATED_WORKER`
- `OAUTH_CLIENT`
- `SCIM`
- `MULLVAD`
- `LOGSTREAM`
- `SECRET_SCANNER`

### Sous-objet : `target`

| Propriete | Type | Description |
|-----------|------|-------------|
| `id` | `string` | The unique ID (user id, tailnet SID, or node id) of the target. |
| `name` | `string` | Name of the entity at time of the action. |
| `type` | `string` (enum) | The entity type of Target. |
| `isEphemeral` | `boolean` | Indicates whether the target is ephemeral. Only set if `type` is `NODE`. |
| `property` | `string` (enum) | The property name on this target which was updated by the event. Empty if the event didn't update any fields. |

#### Enum `target.type`

- `TAILNET`
- `USER`
- `GROUP`
- `NODE`
- `API_KEY`
- `INVITE`
- `SHARE`
- `BILLING`
- `ADMIN_CONSOLE`
- `WEB_INTERFACE`
- `WEBHOOK_ENDPOINT`
- `FAILED_REQUEST`

#### Enum `target.property`

- `ACL`
- `ACL_TAGS`
- `ACCOUNT_EMAIL`
- `ADDRESS`
- `ALLOWED_IPS`
- `AUTO_APPROVED_ROUTES`
- `ATTRIBUTES`
- `BILLING_OWNER`
- `COLLECT_SERVICES`
- `COLLECT_POSTURE_IDENTITY`
- `MULLVAD_VPN`
- `DNS_CONFIG`
- `EMAIL`
- `EXIT_NODE`
- `FEATURE`
- `FILE_SHARING`
- `HTTPS`
- `KEY_EXPIRY_TIME`
- `KEY_EXPIRY`
- `LOG_EXIT_FLOWS`
- `LOGSTREAM_ENDPOINT`
- `MAGIC_DNS`
- `MACHINE_AUTH_NEEDED`
- `MACHINE_APPROVAL_NEEDED`
- `USER_APPROVAL_REQUIRED`
- `MACHINE_NAME`
- `MAX_KEY_DURATION`
- `NETWORK_FLOW_LOGGING`
- `GEOSTEERING`
- `NODE_SHARE`
- `TAILNET_INVITE`
- `PAYMENT_INFO`
- `POSTURE_IDENTITY`
- `POSTURE_INTEGRATION`
- `USER_ROLE`
- `SCIM`
- `SECURITY_EMAIL`
- `STRIPE_CUSTOMER_ID`
- `SUBSCRIPTION`
- `SUBSCRIBED_EVENTS`
- `SUPPORT_EMAIL`
- `SECRET`
- `TCD`
- `TKA`
- `AUTH_PROVIDER`

---

## 8. ConnectionCounts

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `proto` | `string` (enum) | IP protocol name (or number if no name used). |
| `src` | `string` | Source addr:port. |
| `dst` | `string` | Destination addr:port. |
| `txPkts` | `integer` | Number of packets sent. |
| `txBytes` | `integer` | Number of bytes sent. |
| `rxPkts` | `integer` | Number of packets received. |
| `rxBytes` | `integer` | Number of bytes received. |

### Enum `proto`

- `ah`
- `dccp`
- `egp`
- `esp`
- `gre`
- `icmp`
- `igmp`
- `igp`
- `ipv4`
- `ipv6-icmp`
- `sctp`
- `tcp`
- `udp`

---

## 9. NetworkFlowLog

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description | Reference ($ref) |
|-----------|------|-------------|-----------------|
| `logged` | `string` | Timestamp of the flow log, in RFC 3339 format. | - |
| `nodeId` | `string` | Identifier of the node. | - |
| `start` | `string` | Time at which flow started, in RFC 3339 format. | - |
| `end` | `string` | Time at which flow ended, in RFC 3339 format. | - |
| `virtualTraffic` | `array` | Virtual traffic connection counts. | Items: `$ref: '#/components/schemas/ConnectionCounts'` |
| `subnetTraffic` | `array` | Subnet traffic connection counts. | Items: `$ref: '#/components/schemas/ConnectionCounts'` |
| `exitTraffic` | `array` | Exit traffic connection counts. | Items: `$ref: '#/components/schemas/ConnectionCounts'` |
| `physicalTraffic` | `array` | Physical traffic connection counts. | Items: `$ref: '#/components/schemas/ConnectionCounts'` |

---

## 10. LogType

- **Type** : `string`
- **Description** : The type of log for logging endpoints.

### Enum

- `configuration`
- `network`

---

## 11. LogstreamEndpointPublishingStatus

- **Type** : `object`
- **Description** : Latest status of log stream publishing for a specific type of log.
- **Champs requis** : `lastActivity`, `lastError`, `maxBodySize`, `numBytesSent`, `numEntriesSent`, `numSpoofedEntries`, `numTotalRequests`, `numFailedRequests`, `rateBytesSent`, `rateEntriesSent`, `rateTotalRequests`, `rateFailedRequests`

### Proprietes

| Propriete | Type | Requis | Description |
|-----------|------|--------|-------------|
| `lastActivity` | `string` | Oui | Timestamp of the most recent publishing activity, in RFC 3339 format. |
| `lastError` | `string` | Oui | The most recent error (if any). |
| `maxBodySize` | `integer` | Oui | The size of the largest single request body. |
| `numBytesSent` | `integer` | Oui | Total bytes published across all requests. |
| `numEntriesSent` | `integer` | Oui | The total number of entries published. |
| `numSpoofedEntries` | `integer` | Oui | The number of spoofed entries published. A spoofed entry is one that failed to validate because we did not see a matching flow log from the other side of the connection. |
| `numTotalRequests` | `integer` | Oui | The total number of requests made to the streaming endpoint. |
| `numFailedRequests` | `integer` | Oui | The total number of requests to the streaming endpoint that have failed. |
| `rateBytesSent` | `number` | Oui | The exponentially weighted moving average rate at which data is being streamed, in bytes per second. |
| `rateEntriesSent` | `number` | Oui | The exponentially weighted moving average rate at which entries are being sent, in entries per second. |
| `rateTotalRequests` | `number` | Oui | The exponentially weighted moving average rate at which requests are being made, in requests per second. |
| `rateFailedRequests` | `number` | Oui | The exponentially weighted moving average rate at which requests are failing, in requests per second. |

---

## 12. LogstreamEndpointConfiguration

- **Type** : `object`
- **Description** : The current configuration of a log streaming endpoint.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | readOnly/writeOnly | Description | Reference ($ref) |
|-----------|------|-------------------|-------------|-----------------|
| `logType` | - | readOnly | The type of log that is streamed to this endpoint. | `$ref: '#/components/schemas/LogType'` |
| `destinationType` | `string` (enum) | - | The type of system to which logs are being streamed. | - |
| `url` | `string` | - | The URL to which log streams are being posted. If DestinationType is `s3`, the URL may be empty to use the official Amazon S3 endpoint. | - |
| `user` | `string` | - | The username with which log streams to this endpoint are authenticated. | - |
| `uploadPeriodMinutes` | `integer` (maximum: 1440) | - | An optional number of minutes to wait in between uploading new logs. | - |
| `compressionFormat` | `string` (enum) | - | The compression algorithm with which to compress logs. Defaults to `none`. | - |
| `token` | `string` | writeOnly | The token/password with which log streams should be authenticated. | - |
| `s3Bucket` | `string` | - | The S3 bucket name. Required if destinationType is `s3`. | - |
| `s3Region` | `string` | - | The region in which the S3 bucket is located. Required if destinationType is `s3`. | - |
| `s3KeyPrefix` | `string` | - | An optional S3 key prefix to prepend to the auto-generated S3 key name. | - |
| `s3AuthenticationType` | `string` (enum) | - | What type of authentication to use for S3. Required if destinationType is `s3`. Tailscale recommends `rolearn`. | - |
| `s3AccessKeyId` | `string` | - | The S3 access key ID. Required if destinationType is `s3` and authenticationType is `accesskey`. | - |
| `s3SecretAccessKey` | `string` | writeOnly | The S3 secret access key. Required if destinationType is `s3` and authenticationType is `accesskey`. | - |
| `s3RoleArn` | `string` | - | The Role ARN for AWS role-based authentication. Required if destinationType is `s3` and authenticationType is `rolearn`. | - |
| `s3ExternalId` | `string` | readOnly | The AWS external id for role-based authentication. Populated if destinationType is `s3` and authenticationType is `rolearn`. | - |
| `gcsBucket` | `string` | - | The GCS bucket name. Required if destinationType is `gcs`. | - |
| `gcsKeyPrefix` | `string` | - | An optional GCS key prefix to append to the GCS bucket name. | - |
| `gcsScopes` | `array` of `string` | - | The GCS scopes needed to be able to write to the GCS bucket. | - |
| `gcsCredentials` | `string` | - | The JSON workload identity credentials from GCS needed for accessing the GCS account. | - |

### Enum `destinationType`

- `splunk`
- `elastic`
- `panther`
- `cribl`
- `datadog`
- `axiom`
- `s3`

### Enum `compressionFormat`

- `zstd`
- `gzip`
- `none`

### Enum `s3AuthenticationType`

- `accesskey`
- `rolearn`

---

## 13. AwsExternalId

- **Type** : `object`
- **Description** : An external ID for use in authenticating to AWS using role-based authentication.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `externalId` | `string` | The external id that Tailscale will supply to AWS when authenticating using role-based authentication. |
| `tailscaleAwsAccountId` | `string` | The AWS account id that Tailscale will supply to AWS when authenticating using role-based authentication. |

---

## 14. DnsPreferences

- **Type** : `object`
- **Champs requis** : `magicDNS`

### Proprietes

| Propriete | Type | Requis | Description |
|-----------|------|--------|-------------|
| `magicDNS` | `boolean` | Oui | Whether MagicDNS is active for this tailnet. |

---

## 15. DnsSearchPaths

- **Type** : `object`
- **Champs requis** : `searchPaths`

### Proprietes

| Propriete | Type | Requis | Description |
|-----------|------|--------|-------------|
| `searchPaths` | `array` of `string` | Oui | The search domains for the given tailnet. |

---

## 16. SplitDns

- **Type** : `object`
- **Description** : Map of domain names to lists of nameservers or to `null`.
- **additionalProperties** : type `array` of `string` ou `null` (cle = nom de domaine)

Ce schema n'a pas de proprietes fixes. Il utilise `additionalProperties` pour representer un dictionnaire de noms de domaine vers des listes de serveurs DNS.

---

## 17. DnsConfigurationResolver

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `address` | `string` | IPv4 or IPv6 address of the DNS resolver. |
| `useWithExitNode` | `boolean` | If true, this resolver should still be used when a device is configured to use a Tailscale exit node. Requires Tailscale v1.88.1 or later. |

---

## 18. DnsConfigurationPreferences

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `overrideLocalDNS` | `boolean` | If true, resolvers in `nameservers` override the local OS DNS configuration; if false, local resolvers are used. |
| `magicDNS` | `boolean` | Whether MagicDNS is enabled for this tailnet. |

---

## 19. DnsConfiguration

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description | Reference ($ref) |
|-----------|------|-------------|-----------------|
| `nameservers` | `array` | Global DNS resolvers to use. If `preferences.overrideLocalDNS` is true, these override the local OS configuration; otherwise they are used as fallback resolvers. | Items: `$ref: '#/components/schemas/DnsConfigurationResolver'` |
| `splitDNS` | `object` (additionalProperties) | Map of DNS name suffixes (domains) to lists of resolvers for Split DNS and advanced routing overlays. | additionalProperties items: `$ref: '#/components/schemas/DnsConfigurationResolver'` |
| `searchPaths` | `array` of `string` | Search domain paths to apply. | - |
| `preferences` | - | DNS configuration preferences. | `$ref: '#/components/schemas/DnsConfigurationPreferences'` |

---

## 20. KeyCapabilities

- **Type** : `object`
- **Description** : `capabilities` is a mapping of resources to permissible actions.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `devices` | `object` (inline) | Specifies the key's permissions over devices. This field is only populated for auth keys. |

### Sous-objet : `devices`

| Propriete | Type | Description |
|-----------|------|-------------|
| `create` | `object` (inline) | Specifies the key's permissions when creating devices. |

### Sous-objet : `devices.create`

| Propriete | Type | Description |
|-----------|------|-------------|
| `reusable` | `boolean` | Reusable auth keys can be used multiple times to register different devices. |
| `ephemeral` | `boolean` | Ephemeral keys are used to connect and then clean up short-lived devices. |
| `preauthorized` | `boolean` | Pre-approved keys. `true` means devices registered with this key won't require additional approval from a tailnet admin. |
| `tags` | `array` of `string` | Tags that will be set on devices registered with this key. When creating an auth key owned by the tailnet (using OAuth), it must have tags that exactly match the tags on the OAuth client. When creating an auth key owned by a user (using a user's access token), tags are optional. |

---

## 21. Key

- **Type** : `object`
- **Description** : An API access token or Auth Key.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Format | Description | Reference ($ref) |
|-----------|------|--------|-------------|-----------------|
| `id` | `string` | - | Key identifier. | - |
| `key` | `string` | - | The secret key material (only populated at creation time). | - |
| `keyType` | `string` (enum) | - | The type of key. | - |
| `expirySeconds` | `integer` | `int64` | Duration in seconds until the key expires. Only applies to auth keys. | - |
| `created` | `string` | `date-time` | Creation timestamp. | - |
| `updated` | `string` | `date-time` | Last update timestamp. | - |
| `expires` | `string` | `date-time` | Expiration timestamp. | - |
| `revoked` | `string` | `date-time` | Revocation timestamp. | - |
| `capabilities` | - | - | Key capabilities. | `$ref: '#/components/schemas/KeyCapabilities'` |
| `scopes` | `array` of `string` | - | A list of scopes granted to the key. Only applies to OAuth clients, API access tokens, and federated identities. | - |
| `tags` | `array` of `string` | - | A list of tags associated to the trust credential. Auth keys created with this client must have these exact tags, or tags owned by the client's tags. Mandatory if the scopes include "devices:core" or "auth_keys". Only applies to OAuth clients and federated identities. | - |
| `description` | `string` | - | Key description. | - |
| `invalid` | `boolean` | - | Response for a revoked (deleted) or expired key will have `invalid` set to true. | - |
| `userId` | `string` | - | ID of the user who created this key, empty for keys created by trust credentials. | - |
| `audience` | `string` | - | The value used when matching against the `aud` claim from an OIDC identity token. Only applies to federated identities. | - |
| `issuer` | `string` | `uri` | The issuer of the OIDC identity token used in the token exchange. Must be a valid and publicly reachable https:// URL. Only applies to federated identities. | - |
| `subject` | `string` | - | The pattern used when matching against the `sub` claim from an OIDC identity token. Supports `*` wildcards. Only applies to federated identities. | - |
| `customClaimRules` | `object` (additionalProperties: `string`) | - | A map of claim names to pattern strings used to match against arbitrary claims in the OIDC identity token. Patterns can include `*` wildcards. Only applies to federated identities. | - |

### Enum `keyType`

- `auth` -- machine auth keys
- `client` -- OAuth clients
- `api` -- personal API access tokens or tokens generated using an OAuth client or federated identity
- `federated` -- federated identities

---

## 22. PostureIntegration

- **Type** : `object`
- **Description** : A configured PostureIntegration.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | readOnly/writeOnly | Description |
|-----------|------|-------------------|-------------|
| `provider` | `string` (enum) | - | The device posture provider. Required on POST requests, ignored on PATCH requests. |
| `cloudId` | `string` | - | Identifies which of the provider's clouds to integrate with. Varies by provider. |
| `clientId` | `string` | - | Unique identifier for your client. Varies by provider. |
| `tenantId` | `string` | - | The Microsoft Intune directory (tenant) ID. For other providers, this is left blank. |
| `clientSecret` | `string` | writeOnly | The secret (auth key, token, etc.) used to authenticate with the provider. Required when creating a new integration. |
| `id` | `string` | readOnly | A unique identifier for the integration (generated by the system). |
| `configUpdated` | `string` | readOnly | Timestamp of the last time this configuration was updated, in RFC 3339 format. |
| `status` | `object` (inline) | readOnly | Most recent status for this integration. |

### Enum `provider`

- `falcon`
- `intune`
- `jamfpro`
- `kandji`
- `kolide`
- `sentinelone`

### Sous-objet : `status`

| Propriete | Type | Description |
|-----------|------|-------------|
| `lastSync` | `string` | Timestamp of the last synchronization with the device posture provider, in RFC 3339 format. |
| `error` | `string` | If the last synchronization failed, this shows the error message associated with the failed synchronization. |
| `providerHostCount` | `integer` | The number of devices known to the provider. |
| `matchedCount` | `integer` | The number of Tailscale nodes that were matched with provider. |
| `possibleMatchedCount` | `integer` | The number of Tailscale nodes with identifiers for matching. |

---

## 23. User

- **Type** : `object`
- **Description** : Representation of a user within a tailnet.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Format | Description |
|-----------|------|--------|-------------|
| `id` | `string` | - | The unique identifier for the user. Supply this value wherever {userId} is indicated. |
| `displayName` | `string` | - | The name of the user. |
| `loginName` | `string` | - | The emailish login name of the user. |
| `profilePicUrl` | `string` | - | The profile pic URL for the user. |
| `tailnetId` | `string` | - | The tailnet that owns the user. |
| `created` | `string` | `date-time` | The time the user joined their tailnet. |
| `type` | `string` (enum) | - | The type of relation this user has to the tailnet associated with the request. |
| `role` | `string` (enum) | - | The role of the user. |
| `status` | `string` (enum) | - | The status of the user. |
| `deviceCount` | `integer` | - | Number of devices the user owns. |
| `lastSeen` | `string` | `date-time` | The later of: last time any of the user's nodes were connected, or last time the user authenticated to any tailscale service. |
| `currentlyConnected` | `boolean` | - | `true` when the user has a node currently connected to the control server. |

### Enum `type`

- `member`
- `shared`

### Enum `role`

- `owner`
- `member`
- `admin`
- `it-admin`
- `network-admin`
- `billing-admin`
- `auditor`

### Enum `status`

| Valeur | Description |
|--------|-------------|
| `active` | Last seen within 28 days. |
| `idle` | Last seen longer than 28 days. |
| `suspended` | Suspended from accessing the tailnet. |
| `needs-approval` | Unable to join tailnet until approved. |
| `over-billing-limit` | Unable to join tailnet until billing count increased. |

---

## 24. Contact

- **Type** : `object`
- **Description** : A tailnet contact.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `email` | `string` | The contact's email address. |
| `fallbackEmail` | `string` | The email address used when contact's email address has not been verified. |
| `needsVerification` | `boolean` | Indicates whether the contact's email address still needs to be verified. |

---

## 25. Webhook

- **Type** : `object`
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Format | Description |
|-----------|------|--------|-------------|
| `endpointId` | `string` | - | ID of the webhook endpoint. |
| `endpointUrl` | `string` | - | The endpoint that events are sent to from Tailscale via POST requests. |
| `providerType` | `string` (enum) | - | The provider type for the webhook destination, or an empty string if none are applicable. |
| `creatorLoginName` | `string` | - | The login name for the creator of the webhook endpoint. Can be blank for webhooks created with an OAuth client. |
| `created` | `string` | `date-time` | The time the webhook endpoint was created. |
| `lastModified` | `string` | `date-time` | The time the webhook endpoint was last modified. |
| `subscriptions` | `array` of `string` (enum) | - | The list of subscribed events that trigger POST requests to the configured endpoint URL. |
| `secret` | `string` | `password` | The webhook secret associated with the endpoint. Only populated on creation or when the secret is rotated. Used for generating the `Tailscale-Webhook-Signature` header. |

### Enum `providerType`

- `slack`
- `mattermost`
- `googlechat`
- `discord`

### Enum `subscriptions` (items)

- `nodeCreated`
- `nodeNeedsApproval`
- `nodeApproved`
- `nodeKeyExpiringInOneDay`
- `nodeKeyExpired`
- `nodeDeleted`
- `nodeSigned`
- `nodeNeedsSignature`
- `policyUpdate`
- `userCreated`
- `userNeedsApproval`
- `userSuspended`
- `userRestored`
- `userDeleted`
- `userApproved`
- `userRoleUpdated`
- `subnetIPForwardingNotEnabled`
- `exitNodeIPForwardingNotEnabled`

---

## 26. providerType

- **Type** : `string`
- **Description** : The provider type for the webhook destination, or an empty string if none are applicable. Outgoing webhook events are sent in the format expected by the provider type if non-empty.

### Enum

- `slack`
- `mattermost`
- `googlechat`
- `discord`

> Note : Ce schema est defini au meme niveau que les autres schemas mais sert de type reutilisable. Il est identique au champ `providerType` dans le schema `Webhook`.

---

## 27. subscriptions

- **Type** : `array` of `string` (enum)
- **Description** : The list of subscribed events that trigger POST requests to the configured endpoint URL.

### Enum (items)

- `nodeCreated`
- `nodeNeedsApproval`
- `nodeApproved`
- `nodeKeyExpiringInOneDay`
- `nodeKeyExpired`
- `nodeDeleted`
- `nodeSigned`
- `nodeNeedsSignature`
- `policyUpdate`
- `userCreated`
- `userNeedsApproval`
- `userSuspended`
- `userRestored`
- `userDeleted`
- `userApproved`
- `userRoleUpdated`
- `subnetIPForwardingNotEnabled`
- `exitNodeIPForwardingNotEnabled`

> Note : Ce schema est defini au meme niveau que les autres schemas mais sert de type reutilisable. Il est identique au champ `subscriptions` dans le schema `Webhook`.

---

## 28. TailnetSettings

- **Type** : `object`
- **Description** : Settings for a tailnet.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Contraintes | Description |
|-----------|------|-------------|-------------|
| `aclsExternallyManagedOn` | `boolean` ou `null` | - | Prevents users from editing policies in the admin console to avoid conflicts with external management workflows like GitOps or Terraform. |
| `aclsExternalLink` | `string` (format: `uri`) | - | Link to the external tailnet policy definition or management solution for this tailnet. |
| `devicesApprovalOn` | `boolean` ou `null` | - | Whether device approval is enabled for the tailnet. |
| `devicesAutoUpdatesOn` | `boolean` ou `null` | - | Whether auto updates are enabled for devices that belong to this tailnet. |
| `devicesKeyDurationDays` | `integer` | min: 1, max: 180 | The key expiry duration for devices on this tailnet. |
| `usersApprovalOn` | `boolean` ou `null` | - | Whether user approval is enabled for this tailnet. |
| `usersRoleAllowedToJoinExternalTailnets` | `string` (enum) | - | Which user roles are allowed to join external tailnets. |
| `networkFlowLoggingOn` | `boolean` ou `null` | - | Whether network flow logs are enabled for the tailnet. |
| `regionalRoutingOn` | `boolean` ou `null` | - | Whether regional routing is enabled for the tailnet. |
| `postureIdentityCollectionOn` | `boolean` ou `null` | - | Whether identity collection is enabled for device posture integrations for the tailnet. |
| `httpsEnabled` | `boolean` ou `null` | - | Whether provisioning of HTTPS certificates is enabled for this tailnet. |

### Enum `usersRoleAllowedToJoinExternalTailnets`

- `none`
- `admin`
- `member`

---

## 29. VIPServiceInfo

- **Type** : `object`
- **Description** : An information summary for a Service. Each Service has a unique name within the tailnet, one IPv4 and one IPv6 address, optional comment, list of ports, and optional tags.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `name` | `string` | The unique name of the Service (e.g., "svc:example"). |
| `addrs` | `array` of `string` | The IP addresses assigned to the Service: the IPv4 followed by the IPv6. |
| `comment` | `string` | An optional comment for the Service. |
| `ports` | `array` of `string` | A list of protocol:port pairs to be exposed by the Service. The only supported protocol is "tcp" at this time. "do-not-validate" can be used to skip validation. |
| `tags` | `array` of `string` | A list of optional tags associated with the Service. |

---

## 30. VIPServiceInfoPut

- **Type** : `allOf` (composition)
- **Description** : Schema utilise pour les operations PUT sur les Services. Compose de :
  1. Un objet inline avec une propriete `addrs` redefinissant la description
  2. Le schema `VIPServiceInfo` (via `$ref`)

### Composition

```
allOf:
  - type: object
    properties:
      addrs:
        type: array of string
        description: >
          The IP addresses assigned to the Service.
          - For new Services: either unset or a single IPv4 to assign the Service.
          - For existing Services: an IPv4 and an IPv6. The IPv4 can be updated, but not the IPv6.
  - $ref: '#/components/schemas/VIPServiceInfo'
```

### Proprietes heritees de VIPServiceInfo

Toutes les proprietes de [VIPServiceInfo](#29-vipserviceinfo) plus la redefinition de `addrs` avec une description specifique au PUT :
- Pour les nouveaux Services : non defini ou un seul IPv4.
- Pour les Services existants : un IPv4 et un IPv6. L'IPv4 peut etre mis a jour, mais pas l'IPv6.

---

## 31. ServiceHostInfo

- **Type** : `object`
- **Description** : An information summary for a device hosting a Service.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `stableNodeID` | `string` | The preferred identifier for a device. |
| `approvalLevel` | `string` (enum) | The approval level of the device hosting the Service. |
| `configured` | `string` | The configuration status of the device hosting the Service. |

### Enum `approvalLevel`

- `not-approved`
- `approved:auto`
- `approved:manual`

---

## 32. VIPServiceApproval

- **Type** : `object`
- **Description** : The approval status of a Service on a specific device.
- **Champs requis** : aucun

### Proprietes

| Propriete | Type | Description |
|-----------|------|-------------|
| `approved` | `boolean` | Indicates whether the Service is approved on the device. |
| `autoApproved` | `boolean` | Indicates whether the Service was auto-approved by an auto-approver. |

---

## Diagramme des relations entre schemas ($ref)

```
Key
 └──> capabilities: KeyCapabilities

NetworkFlowLog
 ├──> virtualTraffic[]: ConnectionCounts
 ├──> subnetTraffic[]: ConnectionCounts
 ├──> exitTraffic[]: ConnectionCounts
 └──> physicalTraffic[]: ConnectionCounts

LogstreamEndpointConfiguration
 └──> logType: LogType

DnsConfiguration
 ├──> nameservers[]: DnsConfigurationResolver
 ├──> splitDNS (additionalProperties items): DnsConfigurationResolver
 └──> preferences: DnsConfigurationPreferences

VIPServiceInfoPut
 └──> allOf: VIPServiceInfo
```

---

## Schemas reutilises dans les reponses (components.responses)

Les reponses d'erreur HTTP standard referencent toutes le schema `Error` :

| Code HTTP | Description | Schema |
|-----------|-------------|--------|
| 400 | Bad request | `$ref: '#/components/schemas/Error'` |
| 403 | Forbidden | `$ref: '#/components/schemas/Error'` |
| 404 | Not found | `$ref: '#/components/schemas/Error'` |
| 409 | Conflict | `$ref: '#/components/schemas/Error'` |
| 429 | Too Many Requests | `$ref: '#/components/schemas/Error'` |
| 500 | Internal server error | `$ref: '#/components/schemas/Error'` |
| 501 | Not implemented | `$ref: '#/components/schemas/Error'` |
| 502 | Bad gateway | `$ref: '#/components/schemas/Error'` |
| 504 | Gateway timeout | `$ref: '#/components/schemas/Error'` |
