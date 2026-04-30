# Tailscale API - Endpoints Detailles (DevicePosture, Users, Contacts, Webhooks, TailnetSettings, Services)

---

## 1. DevicePosture

### 1.1 List all posture integrations

- **Chemin** : `/tailnet/{tailnet}/posture/integrations`
- **Methode** : `GET`
- **Tags** : `DevicePosture`
- **operationId** : `getPostureIntegrations`
- **Summary** : List all posture integrations
- **Description** : List all of the posture integrations for a tailnet. OAuth Scope: `feature_settings:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. Peut etre `-` pour le tailnet par defaut du token. |

#### Reponse 200

```json
{
  "integrations": [
    {
      "id": "pcBEPQVMpki7DEVEL",           // string, readOnly - identifiant unique
      "provider": "falcon",                  // string, enum: falcon, intune, jamfpro, kandji, kolide, sentinelone
      "cloudId": "us-1",                     // string - cloud du provider
      "clientId": "myclientid",              // string - identifiant client
      "tenantId": "",                        // string - Microsoft Intune tenant ID
      "configUpdated": "2024-06-18T13:43:43.239839Z",  // string, readOnly - timestamp derniere MAJ
      "status": {                            // object, readOnly
        "lastSync": "2024-06-18T08:43:43.777283-05:00",    // string - timestamp derniere synchro
        "error": "...",                                      // string - message d'erreur si echec
        "providerHostCount": 0,                              // integer - nb devices connus du provider
        "matchedCount": 0,                                   // integer - nb noeuds Tailscale matches
        "possibleMatchedCount": 0                            // integer - nb noeuds avec identifiants
      }
    }
  ]
}
```

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to list posture integrations. |

---

### 1.2 Create a posture integration

- **Chemin** : `/tailnet/{tailnet}/posture/integrations`
- **Methode** : `POST`
- **Tags** : `DevicePosture`
- **operationId** : `createPostureIntegration`
- **Summary** : Create a posture integration
- **Description** : Create a posture integration, returning the resulting PostureIntegration. Must include `provider` and `clientSecret`. Currently, only one integration for each provider is supported. OAuth Scope: `feature_settings`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Corps de la requete (application/json)

Schema : `PostureIntegration`

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `provider` | string | oui (POST) | enum: `falcon`, `intune`, `jamfpro`, `kandji`, `kolide`, `sentinelone`. The device posture provider. |
| `cloudId` | string | non | Identifie le cloud du provider (ex: `us-1`, `us-2`, `eu-1`, `us-gov` pour CrowdStrike ; `global`, `us-gov` pour Intune ; FQDN pour Jamf Pro, Kandji, SentinelOne ; vide pour Kolide). |
| `clientId` | string | non | Identifiant client unique (UUID pour Intune, client id pour CrowdStrike/Jamf Pro, vide pour Kandji/Kolide/SentinelOne). |
| `tenantId` | string | non | Microsoft Intune directory (tenant) ID. Vide pour les autres providers. |
| `clientSecret` | string | oui (POST) | writeOnly. Le secret (auth key, token, etc.) pour s'authentifier aupres du provider. |

Exemple :
```json
{
  "provider": "intune",
  "cloudId": "global",
  "clientId": "93013672-b00c-4344-80ca-7ecf74f9dce1",
  "tenantId": "d1ae389b-5207-43a2-afca-2de6b03ac7e3",
  "clientSecret": "as32598d#@%asdf"
}
```

#### Reponse 200

Schema : `PostureIntegration` (meme structure que ci-dessus, avec les champs `id`, `configUpdated`, `status` remplis).

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to create posture integrations. |
| 409 | A posture integration for the same provider already exists. |

---

### 1.3 Get a posture integration

- **Chemin** : `/posture/integrations/{id}`
- **Methode** : `GET`
- **Tags** : `DevicePosture`
- **operationId** : `getPostureIntegration`
- **Summary** : Get a posture integration
- **Description** : Gets the posture integration identified by `{id}`. OAuth Scope: `feature_settings:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `id` | path | string | oui | Unique identifier for a posture integration. Exemple: `p56wQiqrn7mfDEVEL` |

#### Reponse 200

Schema : `PostureIntegration` (structure complete avec `id`, `provider`, `cloudId`, `clientId`, `tenantId`, `configUpdated`, `status`).

#### Erreurs

| Code | Description |
|------|-------------|
| 404 | Posture integration not found, or user is not authorized to read it. |

---

### 1.4 Update a posture integration

- **Chemin** : `/posture/integrations/{id}`
- **Methode** : `PATCH`
- **Tags** : `DevicePosture`
- **operationId** : `updatePostureIntegration`
- **Summary** : Update a posture integration
- **Description** : Updates the posture integration identified by `{id}`. You may omit the `clientSecret` from your request to retain the previously configured `clientSecret`. OAuth Scope: `feature_settings`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `id` | path | string | oui | Unique identifier for a posture integration. |

#### Corps de la requete (application/json)

Schema : `PostureIntegration`

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `cloudId` | string | non | Cloud du provider. |
| `clientId` | string | non | Identifiant client. |
| `tenantId` | string | non | Microsoft Intune tenant ID. |
| `clientSecret` | string | non | writeOnly. Peut etre omis pour conserver le secret existant. |

Note : `provider` est ignore sur les requetes PATCH.

Exemple :
```json
{
  "cloudId": "global",
  "clientId": "93013672-b00c-4344-80ca-7ecf74f9dce1",
  "tenantId": "d1ae389b-5207-43a2-afca-2de6b03ac7e3",
  "clientSecret": "as32598d#@%asdf"
}
```

#### Reponse 200

Schema : `PostureIntegration`

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to update this posture integration. |
| 404 | Posture integration not found. |

---

### 1.5 Delete a posture integration

- **Chemin** : `/posture/integrations/{id}`
- **Methode** : `DELETE`
- **Tags** : `DevicePosture`
- **operationId** : `deletePostureIntegration`
- **Summary** : Delete a posture integration
- **Description** : Delete a specific posture integration. OAuth Scope: `feature_settings`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `id` | path | string | oui | Unique identifier for a posture integration. |

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to delete this posture integration. |
| 404 | Posture integration not found. |

---

## 2. Users

### 2.1 List users

- **Chemin** : `/tailnet/{tailnet}/users`
- **Methode** : `GET`
- **Tags** : `Users`
- **operationId** : `listUsers`
- **Summary** : List users
- **Description** : List all users of a tailnet. OAuth Scope: `users:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Defaut | Description |
|-----|-------------|------|--------|--------|-------------|
| `tailnet` | path | string | oui | - | The tailnet ID. |
| `type` | query | string | non | `member` | Filtre par type d'utilisateur. enum: `member`, `shared`, `all`. |
| `role` | query | string | non | `all` | Filtre par role utilisateur. enum: `owner`, `member`, `admin`, `it-admin`, `network-admin`, `billing-admin`, `auditor`, `all`. |

#### Reponse 200

```json
{
  "users": [
    {
      "id": "123456",                              // string - identifiant unique
      "displayName": "Some User",                  // string - nom affiche
      "loginName": "someuser@example.com",         // string - login emailish
      "profilePicUrl": "",                          // string - URL photo de profil
      "tailnetId": "example.com",                  // string - tailnet proprietaire
      "created": "2022-12-01T05:23:30Z",           // string, date-time - date d'ajout au tailnet
      "type": "member",                            // string, enum: member, shared
      "role": "member",                            // string, enum: owner, member, admin, it-admin, network-admin, billing-admin, auditor
      "status": "active",                          // string, enum: active, idle, suspended, needs-approval, over-billing-limit
      "deviceCount": 4,                            // integer - nombre de devices possedes
      "lastSeen": "2022-12-01T05:23:30Z",          // string, date-time - derniere connexion
      "currentlyConnected": true                   // boolean - actuellement connecte
    }
  ]
}
```

**Details des valeurs `status`** :
- `active` : Last seen within 28 days.
- `idle` : Last seen longer than 28 days.
- `suspended` : Suspended from accessing the tailnet.
- `needs-approval` : Unable to join tailnet until approved.
- `over-billing-limit` : Unable to join tailnet until billing count increased.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Tailnet not found, or user does not have access to read users. |
| 500 | Internal server error. |

---

### 2.2 Get a user

- **Chemin** : `/users/{userId}`
- **Methode** : `GET`
- **Tags** : `Users`
- **operationId** : `getUser`
- **Summary** : Get a user
- **Description** : Retrieve details about the specified user. OAuth Scope: `users:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Reponse 200

Schema : `User` (meme structure que dans la liste ci-dessus).

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 2.3 Update user role

- **Chemin** : `/users/{userId}/role`
- **Methode** : `POST`
- **Tags** : `Users`
- **operationId** : `updateUserRole`
- **Summary** : Update user role
- **Description** : Update the role for the specified user. Learn more about user roles. OAuth Scope: `users`. Note: User-based access tokens cannot update their own user's role.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Corps de la requete (application/json)

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `role` | string | non (schema) | enum: `owner`, `member`, `admin`, `it-admin`, `network-admin`, `billing-admin`, `auditor`. The role of the user. |

Exemple :
```json
{
  "role": "member"
}
```

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 2.4 Approve a user

- **Chemin** : `/users/{userId}/approve`
- **Methode** : `POST`
- **Tags** : `Users`
- **operationId** : `approveUser`
- **Summary** : Approve a user
- **Description** : Approve a pending user's access to the tailnet. This is a no-op if user approval has not been enabled for the tailnet, or if the user is already approved. User approval can be managed using the tailnet settings endpoints. OAuth Scope: `users`. Note: User-based access tokens cannot approve their own user.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 2.5 Suspend a user

- **Chemin** : `/users/{userId}/suspend`
- **Methode** : `POST`
- **Tags** : `Users`
- **operationId** : `suspendUser`
- **Summary** : Suspend a user
- **Description** : Suspends a user from their tailnet. Learn more about suspending users. OAuth Scope: `users`. Note: User-based access tokens cannot suspend their own user.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 2.6 Restore a user

- **Chemin** : `/users/{userId}/restore`
- **Methode** : `POST`
- **Tags** : `Users`
- **operationId** : `restoreUser`
- **Summary** : Restore a user
- **Description** : Restores a suspended user's access to their tailnet. Learn more about restoring users. OAuth Scope: `users`. Note: User-based access tokens cannot restore their own user.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 2.7 Delete a user

- **Chemin** : `/users/{userId}/delete`
- **Methode** : `POST`
- **Tags** : `Users`
- **operationId** : `deleteUser`
- **Summary** : Delete a user
- **Description** : Delete a user from their tailnet. Learn more about deleting users. OAuth Scope: `users`. Note: User-based access tokens cannot delete their own user.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `userId` | path | string | oui | ID of the user. |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

## 3. Contacts

### 3.1 Get contacts

- **Chemin** : `/tailnet/{tailnet}/contacts`
- **Methode** : `GET`
- **Tags** : `Contacts`
- **operationId** : `getContacts`
- **Summary** : Get contacts
- **Description** : Retrieve the tailnet's current contacts. OAuth Scope: `account_settings:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Reponse 200

```json
{
  "account": {
    "email": "owner@example.com",                  // string - adresse email du contact
    "fallbackEmail": "otheruser@example.com",      // string - email de secours (si pas encore verifie)
    "needsVerification": true                       // boolean - indique si la verification est necessaire
  },
  "support": {
    "email": "support@example.com",
    "fallbackEmail": "...",
    "needsVerification": false
  },
  "security": {
    "email": "security@example.com",
    "fallbackEmail": "...",
    "needsVerification": false
  }
}
```

**Schema Contact** :

| Champ | Type | Description |
|-------|------|-------------|
| `email` | string | The contact's email address. |
| `fallbackEmail` | string | The email address used when contact's email address has not been verified. |
| `needsVerification` | boolean | Indicates whether the contact's email address still needs to be verified. |

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to view contacts on this tailnet. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 3.2 Update contact

- **Chemin** : `/tailnet/{tailnet}/contacts/{contactType}`
- **Methode** : `PATCH`
- **Tags** : `Contacts`
- **operationId** : `updateContact`
- **Summary** : Update contact
- **Description** : Update the preferences for this type of contact. If the email address has changed, the system will send a verification email to confirm the change. OAuth Scope: `account_settings`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `contactType` | path | string | oui | Type of contact. enum: `account`, `support`, `security`. |

#### Corps de la requete (application/json)

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `email` | string | **oui** | The contact's email address. |

Exemple :
```json
{
  "email": "newuser@example.com"
}
```

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 403 | User does not have sufficient access to update contacts for this tailnet. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 3.3 Resend verification email

- **Chemin** : `/tailnet/{tailnet}/contacts/{contactType}/resend-verification-email`
- **Methode** : `POST`
- **Tags** : `Contacts`
- **operationId** : `resendContactVerificationEmail`
- **Summary** : Resend verification email
- **Description** : Resends the verification email for this contact, if and only if verification is still pending. OAuth Scope: `account_settings`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `contactType` | path | string | oui | Type of contact. enum: `account`, `support`, `security`. |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Verification is not required, can't resend email. |
| 403 | User does not have sufficient access to update contacts for this tailnet. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

## 4. Webhooks

### 4.1 List webhooks

- **Chemin** : `/tailnet/{tailnet}/webhooks`
- **Methode** : `GET`
- **Tags** : `Webhooks`
- **operationId** : `listWebhooks`
- **Summary** : List webhooks
- **Description** : List all webhooks for a tailnet. OAuth Scope: `webhooks:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Reponse 200

```json
{
  "webhooks": [
    {
      "endpointId": "123456",                        // string - ID du webhook endpoint
      "endpointUrl": "https://example.com/endpoint",  // string - URL de destination
      "providerType": "slack",                         // string, enum: slack, mattermost, googlechat, discord
      "creatorLoginName": "user@example.com",          // string - login du createur
      "created": "2022-12-01T05:23:30Z",               // string, date-time
      "lastModified": "2022-12-01T05:23:30Z",          // string, date-time
      "subscriptions": [                               // array of string
        "nodeCreated",
        "userDeleted"
      ],
      "secret": "xxxxx"                                // string, password - seulement a la creation ou rotation
    }
  ]
}
```

**Valeurs possibles pour `subscriptions`** :
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

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 4.2 Create a webhook

- **Chemin** : `/tailnet/{tailnet}/webhooks`
- **Methode** : `POST`
- **Tags** : `Webhooks`
- **operationId** : `createWebhook`
- **Summary** : Create a webhook
- **Description** : Create a webhook within a tailnet. OAuth Scope: `webhooks`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Corps de la requete (application/json)

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `endpointUrl` | string | **oui** | The endpoint that events are sent to from Tailscale via POST requests. |
| `providerType` | string | non | enum: `slack`, `mattermost`, `googlechat`, `discord`. Format de sortie des evenements. |
| `subscriptions` | array of string | **oui** | Liste des evenements auxquels s'abonner (voir valeurs possibles ci-dessus). |

Exemple :
```json
{
  "endpointUrl": "https://example.com/endpoint",
  "providerType": "slack",
  "subscriptions": ["nodeCreated", "userDeleted"]
}
```

#### Reponse 200

Schema : `Webhook` (structure complete incluant `endpointId`, `secret` genere automatiquement, etc.)

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 4.3 Get webhook

- **Chemin** : `/webhooks/{endpointId}`
- **Methode** : `GET`
- **Tags** : `Webhooks`
- **operationId** : `getWebhook`
- **Summary** : Get webhook
- **Description** : Retrieve a specific webhook. OAuth Scope: `webhooks:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `endpointId` | path | string | oui | ID for the webhook endpoint. |

#### Reponse 200

Schema : `Webhook`

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Webhook not found. |
| 500 | Internal server error. |

---

### 4.4 Update webhook

- **Chemin** : `/webhooks/{endpointId}`
- **Methode** : `PATCH`
- **Tags** : `Webhooks`
- **operationId** : `updateWebhook`
- **Summary** : Update webhook
- **Description** : Update a specific webhook. OAuth Scope: `webhooks`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `endpointId` | path | string | oui | ID for the webhook endpoint. |

#### Corps de la requete (application/json)

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `subscriptions` | array of string | non | Liste des evenements auxquels s'abonner (voir valeurs possibles dans la section 4.1). |

#### Reponse 200

Schema : `Webhook`

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 4.5 Delete webhook

- **Chemin** : `/webhooks/{endpointId}`
- **Methode** : `DELETE`
- **Tags** : `Webhooks`
- **operationId** : `deleteWebhook`
- **Summary** : Delete webhook
- **Description** : Delete a specific webhook. OAuth Scope: `webhooks`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `endpointId` | path | string | oui | ID for the webhook endpoint. |

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Webhook not found. |
| 500 | Internal server error. |

---

### 4.6 Test a webhook

- **Chemin** : `/webhooks/{endpointId}/test`
- **Methode** : `POST`
- **Tags** : `Webhooks`
- **operationId** : `testWebhook`
- **Summary** : Test a webhook
- **Description** : Test a specific webhook by sending out a test event to the endpoint URL. This endpoint queues the event which is sent out asynchronously. If your webhook is configured correctly, within a few seconds your webhook endpoint should receive an event with type of "test". OAuth Scope: `webhooks`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `endpointId` | path | string | oui | ID for the webhook endpoint. |

#### Corps de la requete

Aucun.

#### Reponse 202

Pas de corps de reponse. (Successfully queued test event.)

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | User not found. |
| 500 | Internal server error. |

---

### 4.7 Rotate webhook secret

- **Chemin** : `/webhooks/{endpointId}/rotate`
- **Methode** : `POST`
- **Tags** : `Webhooks`
- **operationId** : `rotateWebhookSecret`
- **Summary** : Rotate webhook secret
- **Description** : Rotate and generate a new secret for a specific webhook. This secret is used for generating the `Tailscale-Webhook-Signature` header in requests sent to the endpoint URL. Learn more about verifying webhook event signatures. OAuth Scope: `webhooks`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `endpointId` | path | string | oui | ID for the webhook endpoint. |

#### Corps de la requete

Aucun.

#### Reponse 200

Schema : `Webhook` (la reponse inclut le nouveau `secret`).

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Webhook not found. |
| 500 | Internal server error. |

---

## 5. TailnetSettings

### 5.1 Get tailnet settings

- **Chemin** : `/tailnet/{tailnet}/settings`
- **Methode** : `GET`
- **Tags** : `TailnetSettings`
- **operationId** : `getTailnetSettings`
- **Summary** : Get tailnet settings
- **Description** : Retrieve the settings for a specific tailnet. OAuth Scopes necessaires selon les parametres :
  - `feature_settings:read` : pour tous les parametres sauf ceux ci-dessous.
  - `logs:network:read` : pour le parametre `networkFlowLoggingOn`.
  - `networking_settings:read` : pour le parametre `httpsCertificates`.
  - `policy_file:read` : pour les parametres `aclsExternallyManagedOn` et `aclsExternalLink`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Reponse 200

Schema : `TailnetSettings`

```json
{
  "aclsExternallyManagedOn": false,                        // boolean|null - empeche edition ACL dans la console admin
  "aclsExternalLink": "https://github.com/example/tailnet-policy",  // string, format uri - lien vers la gestion externe des ACL
  "devicesApprovalOn": false,                              // boolean|null - device approval active
  "devicesAutoUpdatesOn": false,                           // boolean|null - auto-updates pour les devices
  "devicesKeyDurationDays": 180,                           // integer (min: 1, max: 180) - duree d'expiration des cles en jours
  "usersApprovalOn": true,                                 // boolean|null - user approval active
  "usersRoleAllowedToJoinExternalTailnets": "admin",       // string, enum: none, admin, member
  "networkFlowLoggingOn": false,                           // boolean|null - logs de flux reseau actives
  "regionalRoutingOn": false,                              // boolean|null - routage regional actif
  "postureIdentityCollectionOn": false,                    // boolean|null - collecte d'identite posture active
  "httpsEnabled": false                                    // boolean|null - certificats HTTPS actives
}
```

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 5.2 Update tailnet settings

- **Chemin** : `/tailnet/{tailnet}/settings`
- **Methode** : `PATCH`
- **Tags** : `TailnetSettings`
- **operationId** : `updateTailnetSettings`
- **Summary** : Update tailnet settings
- **Description** : Update the settings for a specific tailnet. OAuth Scopes necessaires selon les parametres :
  - `feature_settings` : pour tous les parametres sauf ceux ci-dessous.
  - `logs:network` : pour le parametre `networkFlowLoggingOn`.
  - `networking_settings` : pour le parametre `httpsCertificates`.
  - `policy_file` : pour les parametres `aclsExternallyManagedOn` et `aclsExternalLink`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Corps de la requete (application/json)

Schema : `TailnetSettings` (tous les champs sont optionnels, seuls les champs fournis sont mis a jour).

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `aclsExternallyManagedOn` | boolean\|null | non | Empeche l'edition des ACL dans la console admin. |
| `aclsExternalLink` | string (uri) | non | Lien vers la definition/gestion externe des ACL. |
| `devicesApprovalOn` | boolean\|null | non | Active/desactive le device approval. |
| `devicesAutoUpdatesOn` | boolean\|null | non | Active/desactive les auto-updates. |
| `devicesKeyDurationDays` | integer (1-180) | non | Duree d'expiration des cles en jours. |
| `usersApprovalOn` | boolean\|null | non | Active/desactive le user approval. |
| `usersRoleAllowedToJoinExternalTailnets` | string | non | enum: `none`, `admin`, `member`. |
| `networkFlowLoggingOn` | boolean\|null | non | Active/desactive les logs de flux reseau. |
| `regionalRoutingOn` | boolean\|null | non | Active/desactive le routage regional. |
| `postureIdentityCollectionOn` | boolean\|null | non | Active/desactive la collecte d'identite posture. |
| `httpsEnabled` | boolean\|null | non | Active/desactive les certificats HTTPS. |

#### Reponse 200

Schema : `TailnetSettings` (parametres apres mise a jour).

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

## 6. Services

### 6.1 List all Services

- **Chemin** : `/tailnet/{tailnet}/services`
- **Methode** : `GET`
- **Tags** : `Services`
- **operationId** : `listServices`
- **Summary** : List all Services
- **Description** : List all Services configured for the tailnet. This includes all Services in the "advertised" tab of the Services page in the Tailscale admin console. OAuth Scope: `services:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |

#### Reponse 200

```json
{
  "vipServices": [
    {
      "name": "svc:example",                   // string - nom unique du Service
      "addrs": [                                // array of string - adresses IP (IPv4 puis IPv6)
        "100.93.49.180",
        "fd7a:115c:a1e0::3456:3cb4"
      ],
      "comment": "Example Service",             // string - commentaire optionnel
      "ports": [                                 // array of string - paires protocol:port
        "tcp:80",
        "tcp:443"
      ],
      "tags": [                                  // array of string - tags optionnels
        "tag:example"
      ]
    }
  ]
}
```

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Forbidden. |
| 404 | Tailnet not found. |
| 500 | Internal server error. |

---

### 6.2 Get a Service

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}`
- **Methode** : `GET`
- **Tags** : `Services`
- **operationId** : `getService`
- **Summary** : Get a Service
- **Description** : Retrieve the details for the specified Service. OAuth Scope: `services:read`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service. Must be prefixed with "svc:". Must be unique across the tailnet (names already used for a machine cannot be reused). Exemple: `svc:example` |

#### Reponse 200

Schema : `VIPServiceInfo`

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Access to the Service is forbidden. |
| 404 | Service not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

### 6.3 Update a Service (Create or Update)

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}`
- **Methode** : `PUT`
- **Tags** : `Services`
- **operationId** : `updateService`
- **Summary** : Update a Service
- **Description** : Update or create the specified Service. If the Service does not exist, it will create a Service with the provided details. When creating a new Service, the name in the request body must match the serviceName path parameter. When updating an existing Service, the path parameter is the current name of the Service, and the name in the request body can be used to rename the Service. OAuth Scope: `services`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service (prefixed with "svc:"). |

#### Corps de la requete (application/json) - Requis : oui

Schema : `VIPServiceInfoPut` (extends `VIPServiceInfo`)

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `name` | string | non | The unique name of the Service. Pour la creation, doit correspondre au path param. Pour la mise a jour, peut etre different pour renommer. |
| `addrs` | array of string | non | Pour les nouveaux Services : soit non defini, soit un seul IPv4 a assigner. Pour les Services existants : un IPv4 et un IPv6. L'IPv4 peut etre mis a jour, mais pas l'IPv6. |
| `comment` | string | non | Commentaire optionnel. |
| `ports` | array of string | non | Liste de paires protocol:port (seul "tcp" supporte actuellement, "do-not-validate" pour ignorer la validation). |
| `tags` | array of string | non | Tags optionnels. |

Exemple :
```json
{
  "name": "svc:example",
  "addrs": ["100.93.49.180", "fd7a:115c:a1e0::3456:3cb4"],
  "comment": "Example Service",
  "ports": ["tcp:80", "tcp:443"],
  "tags": ["tag:example"]
}
```

#### Reponse 200

Schema : `VIPServiceInfo`

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Access to modify the Service is forbidden. |
| 404 | Service not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

### 6.4 Delete a Service

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}`
- **Methode** : `DELETE`
- **Tags** : `Services`
- **operationId** : `deleteService`
- **Summary** : Delete a Service
- **Description** : Delete the specified Service from the tailnet. OAuth Scope: `services`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service (prefixed with "svc:"). |

#### Corps de la requete

Aucun.

#### Reponse 200

Pas de corps de reponse.

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Bad request. |
| 403 | Access to delete the Service is forbidden. |
| 404 | Service not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

### 6.5 List devices hosting a Service

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}/devices`
- **Methode** : `GET`
- **Tags** : `Services`
- **operationId** : `listServiceHosts`
- **Summary** : List devices hosting a Service
- **Description** : List all devices that are hosting the specified Service. OAuth Scope: `services`, `devices:core`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service (prefixed with "svc:"). |

#### Reponse 200

```json
{
  "hosts": [
    {
      "stableNodeID": "n292kg92CNTRL",      // string - identifiant prefere du device
      "approvalLevel": "approved:auto",      // string, enum: not-approved, approved:auto, approved:manual
      "configured": "ready"                  // string - statut de configuration du device
    }
  ]
}
```

**Schema ServiceHostInfo** :

| Champ | Type | Description |
|-------|------|-------------|
| `stableNodeID` | string | The preferred identifier for a device. |
| `approvalLevel` | string | enum: `not-approved`, `approved:auto`, `approved:manual`. The approval level of the device hosting the Service. |
| `configured` | string | The configuration status of the device hosting the Service. |

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Invalid parameters or no permission to Services. |
| 403 | Access to the Service or devices is forbidden. |
| 404 | Service not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

### 6.6 Get approval status of Service on a device

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved`
- **Methode** : `GET`
- **Tags** : `Services`
- **operationId** : `getServiceDeviceApproval`
- **Summary** : Get approval status of Service on a device
- **Description** : Retrieve the approval status of the specified Service on a specific device. OAuth Scope: `services`, `devices:core`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service (prefixed with "svc:"). |
| `deviceId` | path | string | oui | ID of the device. Using the device's `nodeId` is preferred, but its numeric `id` value can also be used. |

#### Reponse 200

Schema : `VIPServiceApproval`

```json
{
  "approved": true,           // boolean - indique si le Service est approuve sur le device
  "autoApproved": true        // boolean - indique si le Service a ete auto-approuve par un auto-approver
}
```

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Invalid parameters or no permission to Services. |
| 403 | Access to the Service or device is forbidden. |
| 404 | Service or device not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

### 6.7 Update approval status of Service on a device

- **Chemin** : `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved`
- **Methode** : `POST`
- **Tags** : `Services`
- **operationId** : `updateServiceDeviceApproval`
- **Summary** : Update approval status of Service on a device
- **Description** : Update the approval status of the specified Service on a specific device. OAuth Scope: `services`, `devices:core`.

#### Parametres

| Nom | Emplacement | Type | Requis | Description |
|-----|-------------|------|--------|-------------|
| `tailnet` | path | string | oui | The tailnet ID. |
| `serviceName` | path | string | oui | The unique name of a Service (prefixed with "svc:"). |
| `deviceId` | path | string | oui | ID of the device. |

#### Corps de la requete (application/json) - Requis : oui

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `approved` | boolean | non (schema) | Indicates whether to approve or revoke approval for the Service on the device. |

Exemple :
```json
{
  "approved": true
}
```

#### Reponse 200

Schema : `VIPServiceApproval`

```json
{
  "approved": true,
  "autoApproved": false
}
```

#### Erreurs

| Code | Description |
|------|-------------|
| 400 | Invalid parameters or no permission to Services. |
| 403 | Access to the Service or device is forbidden. |
| 404 | Service or device not found. |
| 500 | Internal server error. |
| 504 | Gateway timeout. |

---

## Annexe : Schemas de reference complets

### PostureIntegration

| Champ | Type | Lecture/Ecriture | Description |
|-------|------|------------------|-------------|
| `id` | string | readOnly | Identifiant unique genere par le systeme. |
| `provider` | string | write (POST uniquement) | enum: `falcon`, `intune`, `jamfpro`, `kandji`, `kolide`, `sentinelone`. |
| `cloudId` | string | read/write | Cloud du provider. |
| `clientId` | string | read/write | Identifiant client. |
| `tenantId` | string | read/write | Microsoft Intune tenant ID. |
| `clientSecret` | string | writeOnly | Secret d'authentification. |
| `configUpdated` | string | readOnly | Timestamp de la derniere MAJ de configuration. |
| `status` | object | readOnly | Contient : `lastSync` (string), `error` (string), `providerHostCount` (integer), `matchedCount` (integer), `possibleMatchedCount` (integer). |

### User

| Champ | Type | Description |
|-------|------|-------------|
| `id` | string | Identifiant unique. |
| `displayName` | string | Nom affiche. |
| `loginName` | string | Login emailish. |
| `profilePicUrl` | string | URL photo de profil. |
| `tailnetId` | string | Tailnet proprietaire. |
| `created` | string (date-time) | Date d'ajout au tailnet. |
| `type` | string | enum: `member`, `shared`. |
| `role` | string | enum: `owner`, `member`, `admin`, `it-admin`, `network-admin`, `billing-admin`, `auditor`. |
| `status` | string | enum: `active`, `idle`, `suspended`, `needs-approval`, `over-billing-limit`. |
| `deviceCount` | integer | Nombre de devices possedes. |
| `lastSeen` | string (date-time) | Derniere connexion ou authentification. |
| `currentlyConnected` | boolean | True si un node est actuellement connecte. |

### Contact

| Champ | Type | Description |
|-------|------|-------------|
| `email` | string | Adresse email du contact. |
| `fallbackEmail` | string | Email de secours utilise quand l'email principal n'est pas verifie. |
| `needsVerification` | boolean | Indique si la verification est encore necessaire. |

### Webhook

| Champ | Type | Description |
|-------|------|-------------|
| `endpointId` | string | ID du webhook endpoint. |
| `endpointUrl` | string | URL de destination des evenements. |
| `providerType` | string | enum: `slack`, `mattermost`, `googlechat`, `discord` (ou vide). |
| `creatorLoginName` | string | Login du createur. |
| `created` | string (date-time) | Date de creation. |
| `lastModified` | string (date-time) | Date de derniere modification. |
| `subscriptions` | array of string | Liste des evenements abonnes. |
| `secret` | string (password) | Secret webhook (uniquement peuple a la creation ou rotation). |

### TailnetSettings

| Champ | Type | Description |
|-------|------|-------------|
| `aclsExternallyManagedOn` | boolean\|null | Empeche edition des ACL dans la console admin. |
| `aclsExternalLink` | string (uri) | Lien vers gestion externe des ACL. |
| `devicesApprovalOn` | boolean\|null | Device approval active. |
| `devicesAutoUpdatesOn` | boolean\|null | Auto-updates actives. |
| `devicesKeyDurationDays` | integer (1-180) | Duree d'expiration des cles en jours. |
| `usersApprovalOn` | boolean\|null | User approval active. |
| `usersRoleAllowedToJoinExternalTailnets` | string | enum: `none`, `admin`, `member`. |
| `networkFlowLoggingOn` | boolean\|null | Logs de flux reseau actives. |
| `regionalRoutingOn` | boolean\|null | Routage regional actif. |
| `postureIdentityCollectionOn` | boolean\|null | Collecte d'identite posture active. |
| `httpsEnabled` | boolean\|null | Certificats HTTPS actives. |

### VIPServiceInfo

| Champ | Type | Description |
|-------|------|-------------|
| `name` | string | Nom unique du Service (prefixe "svc:"). |
| `addrs` | array of string | Adresses IP : IPv4 puis IPv6. |
| `comment` | string | Commentaire optionnel. |
| `ports` | array of string | Paires protocol:port (ex: `tcp:80`). Seul "tcp" supporte. |
| `tags` | array of string | Tags optionnels. |

### VIPServiceInfoPut (extends VIPServiceInfo)

Meme que `VIPServiceInfo`, mais le champ `addrs` a un comportement different :
- Pour les nouveaux Services : soit non defini, soit un seul IPv4.
- Pour les Services existants : un IPv4 et un IPv6. L'IPv4 peut etre mis a jour, mais pas l'IPv6.

### ServiceHostInfo

| Champ | Type | Description |
|-------|------|-------------|
| `stableNodeID` | string | Identifiant prefere du device. |
| `approvalLevel` | string | enum: `not-approved`, `approved:auto`, `approved:manual`. |
| `configured` | string | Statut de configuration du device. |

### VIPServiceApproval

| Champ | Type | Description |
|-------|------|-------------|
| `approved` | boolean | Indique si le Service est approuve sur le device. |
| `autoApproved` | boolean | Indique si le Service a ete auto-approuve. |

---

## Resume des endpoints

| # | Methode | Chemin | Tag | Summary |
|---|---------|--------|-----|---------|
| 1 | GET | `/tailnet/{tailnet}/posture/integrations` | DevicePosture | List all posture integrations |
| 2 | POST | `/tailnet/{tailnet}/posture/integrations` | DevicePosture | Create a posture integration |
| 3 | GET | `/posture/integrations/{id}` | DevicePosture | Get a posture integration |
| 4 | PATCH | `/posture/integrations/{id}` | DevicePosture | Update a posture integration |
| 5 | DELETE | `/posture/integrations/{id}` | DevicePosture | Delete a posture integration |
| 6 | GET | `/tailnet/{tailnet}/users` | Users | List users |
| 7 | GET | `/users/{userId}` | Users | Get a user |
| 8 | POST | `/users/{userId}/role` | Users | Update user role |
| 9 | POST | `/users/{userId}/approve` | Users | Approve a user |
| 10 | POST | `/users/{userId}/suspend` | Users | Suspend a user |
| 11 | POST | `/users/{userId}/restore` | Users | Restore a user |
| 12 | POST | `/users/{userId}/delete` | Users | Delete a user |
| 13 | GET | `/tailnet/{tailnet}/contacts` | Contacts | Get contacts |
| 14 | PATCH | `/tailnet/{tailnet}/contacts/{contactType}` | Contacts | Update contact |
| 15 | POST | `/tailnet/{tailnet}/contacts/{contactType}/resend-verification-email` | Contacts | Resend verification email |
| 16 | GET | `/tailnet/{tailnet}/webhooks` | Webhooks | List webhooks |
| 17 | POST | `/tailnet/{tailnet}/webhooks` | Webhooks | Create a webhook |
| 18 | GET | `/webhooks/{endpointId}` | Webhooks | Get webhook |
| 19 | PATCH | `/webhooks/{endpointId}` | Webhooks | Update webhook |
| 20 | DELETE | `/webhooks/{endpointId}` | Webhooks | Delete webhook |
| 21 | POST | `/webhooks/{endpointId}/test` | Webhooks | Test a webhook |
| 22 | POST | `/webhooks/{endpointId}/rotate` | Webhooks | Rotate webhook secret |
| 23 | GET | `/tailnet/{tailnet}/settings` | TailnetSettings | Get tailnet settings |
| 24 | PATCH | `/tailnet/{tailnet}/settings` | TailnetSettings | Update tailnet settings |
| 25 | GET | `/tailnet/{tailnet}/services` | Services | List all Services |
| 26 | GET | `/tailnet/{tailnet}/services/{serviceName}` | Services | Get a Service |
| 27 | PUT | `/tailnet/{tailnet}/services/{serviceName}` | Services | Update a Service |
| 28 | DELETE | `/tailnet/{tailnet}/services/{serviceName}` | Services | Delete a Service |
| 29 | GET | `/tailnet/{tailnet}/services/{serviceName}/devices` | Services | List devices hosting a Service |
| 30 | GET | `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved` | Services | Get approval status of Service on a device |
| 31 | POST | `/tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved` | Services | Update approval status of Service on a device |
