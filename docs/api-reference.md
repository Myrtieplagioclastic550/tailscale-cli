# Tailscale API v2 - Reference complete

**Base URL** : `https://api.tailscale.com/api/v2/`

**Authentification** : Token API (`tskey-api-xxxxx`) via Basic Auth (username = token, password vide) ou Bearer Token.
Les OAuth clients (`tskey-client-xxxxx`) peuvent creer des tokens a duree limitee avec des scopes specifiques.

**Pagination** : Non supportee. Tous les resultats sont retournes en une seule reponse.

**Format** : JSON (requetes et reponses). HuJSON supporte pour les policy files.

---

## Sommaire des endpoints (85 endpoints au total)

### 1. Devices (15 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 1 | `GET` | `/tailnet/{tailnet}/devices` | `listTailnetDevices` | Lister les devices d'un tailnet |
| 2 | `PATCH` | `/tailnet/{tailnet}/device-attributes` | `batchUpdateCustomDevicePostureAttributes` | Mise a jour batch des attributs posture |
| 3 | `GET` | `/device/{deviceId}` | `getDevice` | Obtenir un device |
| 4 | `DELETE` | `/device/{deviceId}` | `deleteDevice` | Supprimer un device |
| 5 | `POST` | `/device/{deviceId}/expire` | `expireDeviceKey` | Expirer la cle d'un device |
| 6 | `GET` | `/device/{deviceId}/routes` | `listDeviceRoutes` | Lister les routes d'un device |
| 7 | `POST` | `/device/{deviceId}/routes` | `setDeviceRoutes` | Definir les routes d'un device |
| 8 | `POST` | `/device/{deviceId}/authorized` | `authorizeDevice` | Autoriser/desautoriser un device |
| 9 | `POST` | `/device/{deviceId}/name` | `setDeviceName` | Definir le nom d'un device |
| 10 | `POST` | `/device/{deviceId}/tags` | `setDeviceTags` | Definir les tags d'un device |
| 11 | `POST` | `/device/{deviceId}/key` | `updateDeviceKey` | Activer/desactiver l'expiration de cle |
| 12 | `POST` | `/device/{deviceId}/ip` | `setDeviceIp` | Definir l'adresse IPv4 |
| 13 | `GET` | `/device/{deviceId}/attributes` | `getDevicePostureAttributes` | Obtenir les attributs posture |
| 14 | `POST` | `/device/{deviceId}/attributes/{attributeKey}` | `setCustomDevicePostureAttributes` | Definir un attribut posture custom |
| 15 | `DELETE` | `/device/{deviceId}/attributes/{attributeKey}` | `deleteCustomDevicePostureAttributes` | Supprimer un attribut posture custom |

### 2. Policy File / ACL (4 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 16 | `GET` | `/tailnet/{tailnet}/acl` | `getPolicyFile` | Recuperer le policy file |
| 17 | `POST` | `/tailnet/{tailnet}/acl` | `setPolicyFile` | Definir le policy file |
| 18 | `POST` | `/tailnet/{tailnet}/acl/preview` | `previewRuleMatches` | Previsualiser les regles |
| 19 | `POST` | `/tailnet/{tailnet}/acl/validate` | `validateAndTestPolicyFile` | Valider et tester un policy file |

### 3. Keys (5 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 20 | `GET` | `/tailnet/{tailnet}/keys` | `listTailnetKeys` | Lister les cles du tailnet |
| 21 | `POST` | `/tailnet/{tailnet}/keys` | `createKey` | Creer un auth key ou trust credential |
| 22 | `GET` | `/tailnet/{tailnet}/keys/{keyId}` | `getKey` | Obtenir les details d'une cle |
| 23 | `DELETE` | `/tailnet/{tailnet}/keys/{keyId}` | `deleteKey` | Supprimer une cle |
| 24 | `PUT` | `/tailnet/{tailnet}/keys/{keyId}` | `setKey` | Configurer un OAuth client/federated identity |

### 4. DNS (11 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 25 | `GET` | `/tailnet/{tailnet}/dns/nameservers` | `listDnsNameservers` | Lister les nameservers DNS |
| 26 | `POST` | `/tailnet/{tailnet}/dns/nameservers` | `setDnsNameservers` | Definir les nameservers DNS |
| 27 | `GET` | `/tailnet/{tailnet}/dns/preferences` | `getDnsPreferences` | Obtenir les preferences DNS |
| 28 | `POST` | `/tailnet/{tailnet}/dns/preferences` | `setDnsPreferences` | Definir les preferences DNS (MagicDNS) |
| 29 | `GET` | `/tailnet/{tailnet}/dns/searchpaths` | `listDnsSearchPaths` | Lister les search paths DNS |
| 30 | `POST` | `/tailnet/{tailnet}/dns/searchpaths` | `setDnsSearchPaths` | Definir les search paths DNS |
| 31 | `GET` | `/tailnet/{tailnet}/dns/split-dns` | `getSplitDns` | Obtenir la config split DNS |
| 32 | `PATCH` | `/tailnet/{tailnet}/dns/split-dns` | `updateSplitDns` | Mise a jour partielle split DNS |
| 33 | `PUT` | `/tailnet/{tailnet}/dns/split-dns` | `setSplitDns` | Remplacer la config split DNS |
| 34 | `GET` | `/tailnet/{tailnet}/dns/configuration` | `getDnsConfiguration` | Obtenir la config DNS complete |
| 35 | `POST` | `/tailnet/{tailnet}/dns/configuration` | `setDnsConfiguration` | Definir la config DNS complete |

### 5. Logging (8 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 36 | `GET` | `/tailnet/{tailnet}/logging/configuration` | `listConfigurationAuditLogs` | Lister les logs d'audit |
| 37 | `GET` | `/tailnet/{tailnet}/logging/network` | `listNetworkFlowLogs` | Lister les logs de flux reseau |
| 38 | `GET` | `/tailnet/{tailnet}/logging/{logType}/stream/status` | `getLogStreamingStatus` | Statut du streaming de logs |
| 39 | `GET` | `/tailnet/{tailnet}/logging/{logType}/stream` | `getLogStreamingConfiguration` | Config du streaming de logs |
| 40 | `PUT` | `/tailnet/{tailnet}/logging/{logType}/stream` | `setLogStreamingConfiguration` | Definir la config de streaming |
| 41 | `DELETE` | `/tailnet/{tailnet}/logging/{logType}/stream` | `disableLogStreaming` | Desactiver le streaming de logs |
| 42 | `POST` | `/tailnet/{tailnet}/aws-external-id` | `getAwsExternalId` | Creer/obtenir un external ID AWS |
| 43 | `POST` | `/tailnet/{tailnet}/aws-external-id/{id}/validate-aws-trust-policy` | `validateAwsExternalId` | Valider l'integration AWS IAM |

### 6. Users (7 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 44 | `GET` | `/tailnet/{tailnet}/users` | `listUsers` | Lister les utilisateurs |
| 45 | `GET` | `/users/{userId}` | `getUser` | Obtenir un utilisateur |
| 46 | `POST` | `/users/{userId}/role` | `updateUserRole` | Modifier le role d'un utilisateur |
| 47 | `POST` | `/users/{userId}/approve` | `approveUser` | Approuver un utilisateur |
| 48 | `POST` | `/users/{userId}/suspend` | `suspendUser` | Suspendre un utilisateur |
| 49 | `POST` | `/users/{userId}/restore` | `restoreUser` | Restaurer un utilisateur |
| 50 | `POST` | `/users/{userId}/delete` | `deleteUser` | Supprimer un utilisateur |

### 7. User Invites (5 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 51 | `GET` | `/tailnet/{tailnet}/user-invites` | `listUserInvites` | Lister les invitations utilisateur |
| 52 | `POST` | `/tailnet/{tailnet}/user-invites` | `createUserInvites` | Creer des invitations utilisateur |
| 53 | `GET` | `/user-invites/{userInviteId}` | `getUserInvite` | Obtenir une invitation |
| 54 | `DELETE` | `/user-invites/{userInviteId}` | `deleteUserInvite` | Supprimer une invitation |
| 55 | `POST` | `/user-invites/{userInviteId}/resend` | `resendUserInvite` | Renvoyer une invitation |

### 8. Device Invites (6 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 56 | `GET` | `/device/{deviceId}/device-invites` | `listDeviceInvites` | Lister les invites de partage |
| 57 | `POST` | `/device/{deviceId}/device-invites` | `createDeviceInvites` | Creer des invites de partage |
| 58 | `GET` | `/device-invites/{deviceInviteId}` | `getDeviceInvite` | Obtenir une invite |
| 59 | `DELETE` | `/device-invites/{deviceInviteId}` | `deleteDeviceInvite` | Supprimer une invite |
| 60 | `POST` | `/device-invites/{deviceInviteId}/resend` | `resendDeviceInvite` | Renvoyer une invite |
| 61 | `POST` | `/device-invites/-/accept` | `acceptDeviceInvite` | Accepter une invite |

### 9. Device Posture (5 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 62 | `GET` | `/tailnet/{tailnet}/posture/integrations` | `getPostureIntegrations` | Lister les integrations posture |
| 63 | `POST` | `/tailnet/{tailnet}/posture/integrations` | `createPostureIntegration` | Creer une integration posture |
| 64 | `GET` | `/posture/integrations/{id}` | `getPostureIntegration` | Obtenir une integration |
| 65 | `PATCH` | `/posture/integrations/{id}` | `updatePostureIntegration` | Mettre a jour une integration |
| 66 | `DELETE` | `/posture/integrations/{id}` | `deletePostureIntegration` | Supprimer une integration |

### 10. Contacts (3 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 67 | `GET` | `/tailnet/{tailnet}/contacts` | `getContacts` | Obtenir les contacts |
| 68 | `PATCH` | `/tailnet/{tailnet}/contacts/{contactType}` | `updateContact` | Modifier un contact |
| 69 | `POST` | `/tailnet/{tailnet}/contacts/{contactType}/resend-verification-email` | `resendContactVerificationEmail` | Renvoyer l'email de verification |

### 11. Webhooks (7 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 70 | `GET` | `/tailnet/{tailnet}/webhooks` | `listWebhooks` | Lister les webhooks |
| 71 | `POST` | `/tailnet/{tailnet}/webhooks` | `createWebhook` | Creer un webhook |
| 72 | `GET` | `/webhooks/{endpointId}` | `getWebhook` | Obtenir un webhook |
| 73 | `PATCH` | `/webhooks/{endpointId}` | `updateWebhook` | Modifier un webhook |
| 74 | `DELETE` | `/webhooks/{endpointId}` | `deleteWebhook` | Supprimer un webhook |
| 75 | `POST` | `/webhooks/{endpointId}/test` | `testWebhook` | Tester un webhook |
| 76 | `POST` | `/webhooks/{endpointId}/rotate` | `rotateWebhookSecret` | Rotation du secret webhook |

### 12. Tailnet Settings (2 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 77 | `GET` | `/tailnet/{tailnet}/settings` | `getTailnetSettings` | Obtenir les parametres du tailnet |
| 78 | `PATCH` | `/tailnet/{tailnet}/settings` | `updateTailnetSettings` | Modifier les parametres du tailnet |

### 13. Services (7 endpoints)

| # | Methode | Path | operationId | Description |
|---|---------|------|-------------|-------------|
| 79 | `GET` | `/tailnet/{tailnet}/services` | `listServices` | Lister les Services |
| 80 | `GET` | `/tailnet/{tailnet}/services/{serviceName}` | `getService` | Obtenir un Service |
| 81 | `PUT` | `/tailnet/{tailnet}/services/{serviceName}` | `updateService` | Creer/modifier un Service |
| 82 | `DELETE` | `/tailnet/{tailnet}/services/{serviceName}` | `deleteService` | Supprimer un Service |
| 83 | `GET` | `/tailnet/{tailnet}/services/{serviceName}/devices` | `listServiceHosts` | Lister les devices d'un Service |
| 84 | `GET` | `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved` | `getServiceDeviceApproval` | Statut approbation d'un device |
| 85 | `POST` | `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved` | `updateServiceDeviceApproval` | Modifier approbation d'un device |

---

## Parametres globaux

### `tailnet` (path)
- **Type** : string
- **Requis** : oui
- **Description** : L'ID du tailnet. Utiliser `-` pour le tailnet par defaut du token, ou le Tailnet ID (ex: `T1234CNTRL`).

### `deviceId` (path)
- **Type** : string
- **Requis** : oui
- **Description** : ID du device. Preferer `nodeId` (ex: `n292kg92CNTRL`), mais l'`id` numerique est accepte.

### `fields` (query) - Devices uniquement
- **Type** : string (enum: `all`, `default`)
- **Description** : `all` retourne tous les champs, `default` retourne un sous-ensemble.

---

## Modeles de donnees principaux (32 schemas)

### Device
Objet complet avec 27+ proprietes : `addresses`, `id`, `nodeId`, `user`, `name`, `hostname`, `clientVersion`, `updateAvailable`, `os`, `created`, `connectedToControl`, `lastSeen`, `keyExpiryDisabled`, `expires`, `authorized`, `isExternal`, `multipleConnections`, `machineKey`, `nodeKey`, `blocksIncomingConnections`, `enabledRoutes`, `advertisedRoutes`, `clientConnectivity` (object), `tags`, `tailnetLockError`, `tailnetLockKey`, `sshEnabled`, `postureIdentity` (object), `isEphemeral`, `distro` (object).

### Key
`id`, `key` (secret, uniquement a la creation), `keyType` (enum: `auth`, `client`, `api`, `federated`), `expirySeconds`, `created`, `updated`, `expires`, `revoked`, `capabilities` (KeyCapabilities), `scopes`, `tags`, `description`, `invalid`, `userId`, `audience`, `issuer`, `subject`, `customClaimRules`.

### User
`id`, `displayName`, `loginName`, `profilePicUrl`, `tailnetId`, `created`, `type` (member/shared), `role` (owner/member/admin/it-admin/network-admin/billing-admin/auditor), `status` (active/idle/suspended/needs-approval/over-billing-limit), `deviceCount`, `lastSeen`, `currentlyConnected`.

### TailnetSettings
`aclsExternallyManagedOn`, `aclsExternalLink`, `devicesApprovalOn`, `devicesAutoUpdatesOn`, `devicesKeyDurationDays` (1-180), `usersApprovalOn`, `usersRoleAllowedToJoinExternalTailnets` (none/admin/member), `networkFlowLoggingOn`, `regionalRoutingOn`, `postureIdentityCollectionOn`, `httpsEnabled`.

### Webhook
`endpointId`, `endpointUrl`, `providerType` (slack/mattermost/googlechat/discord), `creatorLoginName`, `created`, `lastModified`, `subscriptions` (19 evenements possibles), `secret`.

### VIPServiceInfo
`name` (prefixe `svc:`), `addrs` (IPv4+IPv6), `comment`, `ports` (protocol:port), `tags`.

### PostureIntegration
`id`, `provider` (falcon/intune/jamfpro/kandji/kolide/sentinelone), `cloudId`, `clientId`, `tenantId`, `clientSecret` (writeOnly), `configUpdated`, `status` (lastSync/error/counts).

### Autres schemas
- `Error` : `{ message: string }`
- `DeviceRoutes` : `advertisedRoutes`, `enabledRoutes`
- `DevicePostureAttributes` : `attributes` (map), `expiries` (map)
- `DeviceInvite` / `UserInvite`
- `ConfigurationAuditLog` / `NetworkFlowLog` / `ConnectionCounts`
- `LogstreamEndpointConfiguration` / `LogstreamEndpointPublishingStatus`
- `DnsConfiguration` / `DnsPreferences` / `DnsSearchPaths` / `SplitDns`
- `KeyCapabilities`
- `Contact`

---

## Codes de reponse HTTP standards

| Code | Description |
|------|-------------|
| 200 | Succes |
| 202 | Accepted (async, webhooks test) |
| 400 | Parametres invalides |
| 403 | Acces insuffisant |
| 404 | Ressource non trouvee |
| 409 | Conflit (ex: integration posture dupliquee) |
| 412 | Precondition failed (If-Match ETag mismatch) |
| 422 | Validation echouee |
| 429 | Rate limited |
| 500 | Erreur serveur |
| 501 | Non implemente |
| 502 | Bad gateway |
| 504 | Timeout |

---

## OAuth Scopes

| Scope | Description |
|-------|-------------|
| `devices:core` | Gestion des devices |
| `devices:core:read` | Lecture des devices |
| `devices:routes` | Gestion des routes |
| `devices:routes:read` | Lecture des routes |
| `devices:posture_attributes` | Gestion des attributs posture |
| `devices:posture_attributes:read` | Lecture des attributs posture |
| `policy_file` | Gestion du policy file |
| `policy_file:read` | Lecture du policy file |
| `auth_keys` | Gestion des auth keys |
| `auth_keys:read` | Lecture des auth keys |
| `api_access_tokens` | Gestion des API tokens |
| `api_access_tokens:read` | Lecture des API tokens |
| `oauth_keys` | Gestion des OAuth clients |
| `oauth_keys:read` | Lecture des OAuth clients |
| `federated_keys` | Gestion des federated identities |
| `federated_keys:read` | Lecture des federated identities |
| `users` | Gestion des utilisateurs |
| `users:read` | Lecture des utilisateurs |
| `device_invites` | Gestion des invites device |
| `device_invites:read` | Lecture des invites device |
| `webhooks` | Gestion des webhooks |
| `webhooks:read` | Lecture des webhooks |
| `log_streaming` | Gestion du streaming de logs |
| `log_streaming:read` | Lecture du streaming de logs |
| `logs:configuration:read` | Lecture des logs d'audit |
| `logs:network:read` | Lecture des logs reseau |
| `logs:network` | Gestion des logs reseau |
| `account_settings` | Gestion des contacts/parametres compte |
| `account_settings:read` | Lecture des contacts/parametres compte |
| `feature_settings` | Gestion des parametres tailnet |
| `feature_settings:read` | Lecture des parametres tailnet |
| `networking_settings` | Gestion des parametres reseau (HTTPS) |
| `networking_settings:read` | Lecture des parametres reseau |
| `services` | Gestion des Services |
| `services:read` | Lecture des Services |
