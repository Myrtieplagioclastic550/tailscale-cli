# Tailscale API - Endpoints DeviceInvites & UserInvites

Extraction exhaustive de tous les endpoints des sections DeviceInvites et UserInvites (lignes 728 a 1210 du fichier OpenAPI).

---

## Schemas de reference

### Schema `DeviceInvite`

Une invitation de partage de device avec un utilisateur externe (hors du tailnet du device).

| Champ           | Type             | Exemple                                              | Description                                                                                                  |
|-----------------|------------------|------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `id`            | string           | `"12346"`                                            | Identifiant unique de l'invite. A fournir pour `{deviceInviteId}`.                                           |
| `created`       | string (date-time) | `"2024-04-03T21:38:49.333829261Z"`                | Date de creation de l'invite.                                                                                |
| `tailnetId`     | integer (int64)  | `59954`                                              | ID du tailnet auquel le device partage appartient.                                                           |
| `deviceId`      | integer (int64)  | `11055`                                              | ID du device partage.                                                                                        |
| `sharerId`      | integer (int64)  | `22012`                                              | ID de l'utilisateur qui a cree l'invite de partage.                                                          |
| `multiUse`      | boolean          | `false`                                              | Indique si l'invite peut etre acceptee plus d'une fois.                                                      |
| `allowExitNode` | boolean          | `false`                                              | Indique si l'utilisateur invite peut utiliser le device comme exit node quand celui-ci l'annonce.             |
| `email`         | string           | `"user@example.com"`                                 | Email a laquelle l'invite a ete envoyee. Si vide, l'invite n'a pas ete envoyee par email mais l'URL peut etre partagee manuellement. |
| `lastEmailSentAt` | string (date-time) | `"2024-04-03T21:38:49.333829261Z"`              | Derniere tentative d'envoi de l'invite par email. Defini uniquement si `email` n'est pas vide.               |
| `inviteUrl`     | string           | `"https://login.tailscale.com/admin/invite/<code>"`  | Lien pour accepter l'invite. N'importe qui avec ce lien peut accepter l'invite (pas restreint a l'email).    |
| `accepted`      | boolean          | `false`                                              | `true` quand l'invite de partage a ete acceptee.                                                             |
| `acceptedBy`    | object           | -                                                    | Defini quand l'invite a ete acceptee. Contient les infos de l'utilisateur qui a accepte.                     |
| `acceptedBy.id`          | integer (int64) | `33223`                  | ID de l'utilisateur qui a accepte.             |
| `acceptedBy.loginName`   | string          | `"someone@example.com"`  | Login de l'utilisateur qui a accepte.          |
| `acceptedBy.profilePicUrl` | string        | `""`                     | URL de la photo de profil de l'utilisateur.    |

### Schema `UserInvite`

Une invitation permettant a un utilisateur de rejoindre un tailnet avec un role preassigne.

Champs requis : `id`, `role`, `tailnetId`, `inviterId`.

| Champ           | Type             | Exemple                                              | Description                                                                                                  |
|-----------------|------------------|------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `id`            | string           | `"12346"`                                            | Identifiant unique de l'invite. A fournir pour `{userInviteId}`.                                             |
| `role`          | string (enum)    | `"admin"`                                            | Role a assigner a l'utilisateur invite. Valeurs : `member`, `admin`, `it-admin`, `network-admin`, `billing-admin`, `auditor`. |
| `tailnetId`     | integer (int64)  | `59954`                                              | ID du tailnet auquel l'utilisateur est invite.                                                               |
| `inviterId`     | integer (int64)  | `22012`                                              | ID de l'utilisateur qui a cree l'invite.                                                                     |
| `email`         | string           | `"user@example.com"`                                 | Email a laquelle l'invite a ete envoyee. Si vide, l'invite n'a pas ete envoyee par email mais l'URL peut etre partagee manuellement. |
| `lastEmailSentAt` | string (date-time) | `"2024-04-03T21:38:49.333829261Z"`              | Derniere tentative d'envoi de l'invite par email. Defini uniquement si `email` n'est pas vide.               |
| `inviteUrl`     | string           | `"https://login.tailscale.com/admin/invite/<code>"`  | Inclus quand `email` n'est pas dans le domaine du tailnet, ou quand `email` est vide. Lien pour accepter l'invite. N'importe qui avec ce lien peut accepter. Quand `email` est dans le domaine du tailnet, l'utilisateur peut rejoindre automatiquement via https://login.tailscale.com/start. |

### Parametres de chemin references

| Parametre         | In   | Type   | Requis | Description                                                                                                                    |
|-------------------|------|--------|--------|--------------------------------------------------------------------------------------------------------------------------------|
| `deviceId`        | path | string | oui    | ID du device. L'utilisation du `nodeId` du device est preferee, mais sa valeur numerique `id` peut aussi etre utilisee.        |
| `tailnet`         | path | string | oui    | L'ID du tailnet. On peut fournir `-` pour referencer le tailnet par defaut du token, ou l'ID du tailnet (ex: `T1234CNTRL`).    |
| `userInviteId`    | path | string | oui    | ID de l'invite utilisateur.                                                                                                    |
| `deviceInviteId`  | path | string | oui    | ID de l'invite device.                                                                                                         |

---

## Endpoint 1 : List device invites

- **Chemin** : `/device/{deviceId}/device-invites`
- **Methode HTTP** : `GET`
- **operationId** : `listDeviceInvites`
- **Tags** : `DeviceInvites`
- **Summary** : List device invites
- **Description** : List all share invites for a device. OAuth Scope: `device_invites:read`.

### Parametres

| Nom        | In   | Type   | Requis | Description                                                                                                             |
|------------|------|--------|--------|-------------------------------------------------------------------------------------------------------------------------|
| `deviceId` | path | string | oui    | ID du device. L'utilisation du `nodeId` est preferee, mais sa valeur numerique `id` peut aussi etre utilisee.            |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Type : `array` d'objets `DeviceInvite`

```json
[
  {
    "id": "12345",
    "created": "2024-05-08T20:19:51.777861756Z",
    "tailnetId": 59954,
    "deviceId": 11055,
    "sharerId": 22011,
    "allowExitNode": true,
    "email": "user@example.com",
    "lastEmailSentAt": "2024-05-08T20:19:51.777861756Z",
    "inviteUrl": "https://login.tailscale.com/admin/invite/<code>",
    "accepted": false
  },
  {
    "id": "12346",
    "created": "2024-04-03T21:38:49.333829261Z",
    "tailnetId": 59954,
    "deviceId": 11055,
    "sharerId": 22012,
    "inviteUrl": "https://login.tailscale.com/admin/invite/<code>",
    "accepted": true,
    "acceptedBy": {
      "id": 33223,
      "loginName": "someone@example.com",
      "profilePicUrl": ""
    }
  }
]
```

#### 404 - Device not found

#### 500 - Internal server error

#### 504 - Gateway timeout

---

## Endpoint 2 : Create device invites

- **Chemin** : `/device/{deviceId}/device-invites`
- **Methode HTTP** : `POST`
- **operationId** : `createDeviceInvites`
- **Tags** : `DeviceInvites`
- **Summary** : Create device invites
- **Description** : Create new share invites for a device. Note that device invites cannot be created using an API access token generated from an OAuth client as the shared device is scoped to a user.

### Parametres

| Nom        | In   | Type   | Requis | Description                                                                                                             |
|------------|------|--------|--------|-------------------------------------------------------------------------------------------------------------------------|
| `deviceId` | path | string | oui    | ID du device. L'utilisation du `nodeId` est preferee, mais sa valeur numerique `id` peut aussi etre utilisee.            |

### Corps de la requete (requestBody)

Content-Type : `application/json`

Type : `array` d'objets avec les proprietes suivantes :

| Champ           | Type    | Requis | Exemple              | Description                                                                                                                                   |
|-----------------|---------|--------|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `multiUse`      | boolean | non    | `false`              | Si l'invite peut etre acceptee plus d'une fois. Quand `true`, l'invite peut etre acceptee jusqu'a 1 000 fois.                                 |
| `allowExitNode` | boolean | non    | `false`              | Si l'utilisateur invite peut utiliser le device comme exit node quand celui-ci l'annonce.                                                      |
| `email`         | string  | non    | `"user@example.com"` | Email a laquelle envoyer l'invite creee. Si non defini, l'endpoint genere et retourne une URL d'invite (sans l'envoyer).                      |

### Reponses

#### 200 - Successful operation

Type : `array` d'objets `DeviceInvite`

```json
[
  {
    "id": "12345",
    "created": "2024-05-08T20:19:51.777861756Z",
    "tailnetId": 59954,
    "deviceId": 11055,
    "sharerId": 22011,
    "allowExitNode": true,
    "email": "user@example.com",
    "lastEmailSentAt": "2024-05-08T20:19:51.777861756Z",
    "inviteUrl": "https://login.tailscale.com/admin/invite/<code>",
    "accepted": false
  },
  {
    "id": "12346",
    "created": "2024-04-03T21:38:49.333829261Z",
    "tailnetId": 59954,
    "deviceId": 11055,
    "sharerId": 22012,
    "inviteUrl": "https://login.tailscale.com/admin/invite/<code>",
    "accepted": false
  }
]
```

#### 404 - Device not found

#### 500 - Internal server error

#### 504 - Gateway timeout

---

## Endpoint 3 : List user invites

- **Chemin** : `/tailnet/{tailnet}/user-invites`
- **Methode HTTP** : `GET`
- **operationId** : `listUserInvites`
- **Tags** : `UserInvites`
- **Summary** : List user invites
- **Description** : List all open (not yet accepted) user invites to the tailnet.

### Parametres

| Nom       | In   | Type   | Requis | Description                                                                                                                 |
|-----------|------|--------|--------|-----------------------------------------------------------------------------------------------------------------------------|
| `tailnet` | path | string | oui    | L'ID du tailnet. Peut etre `-` pour le tailnet par defaut du token, ou l'ID du tailnet (ex: `T1234CNTRL`).                  |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Type : `array` d'objets `UserInvite`

```json
[
  {
    "id": "29214",
    "role": "admin",
    "tailnetId": 12345,
    "inviterId": 34567,
    "email": "user@example.com",
    "lastEmailSentAt": "2024-05-09T16:23:26.91778771Z",
    "inviteUrl": "https://login.tailscale.com/uinv/<code>"
  },
  {
    "id": "29215",
    "role": "admin",
    "tailnetId": 12345,
    "inviterId": 34567,
    "email": "someoneelse@example.com",
    "lastEmailSentAt": "2024-05-09T17:23:30.91778771Z",
    "inviteUrl": "https://login.tailscale.com/uinv/<code>"
  }
]
```

#### 404 - Tailnet not found

#### 500 - Internal server error

---

## Endpoint 4 : Create user invites

- **Chemin** : `/tailnet/{tailnet}/user-invites`
- **Methode HTTP** : `POST`
- **operationId** : `createUserInvites`
- **Tags** : `UserInvites`
- **Summary** : Create user invites
- **Description** : Create, and optionally email out, new user invites to join the tailnet. Only permitted for user-owned keys, because invites require an inviting user.

### Parametres

| Nom       | In   | Type   | Requis | Description                                                                                                                 |
|-----------|------|--------|--------|-----------------------------------------------------------------------------------------------------------------------------|
| `tailnet` | path | string | oui    | L'ID du tailnet. Peut etre `-` pour le tailnet par defaut du token, ou l'ID du tailnet (ex: `T1234CNTRL`).                  |

### Corps de la requete (requestBody)

Content-Type : `application/json`

Type : `array` d'objets avec les proprietes suivantes :

| Champ  | Type          | Requis | Defaut    | Exemple              | Description                                                                                                                        |
|--------|---------------|--------|-----------|----------------------|------------------------------------------------------------------------------------------------------------------------------------|
| `role` | string (enum) | non    | `member`  | `"admin"`            | Role a assigner a l'utilisateur invite. Valeurs possibles : `member`, `admin`, `it-admin`, `network-admin`, `billing-admin`, `auditor`. |
| `email`| string        | non    | -         | `"user@example.com"` | Email a laquelle envoyer l'invite. Si non defini, l'endpoint genere et retourne une URL d'invite sans l'envoyer par email.         |

### Reponses

#### 200 - Successful operation

Type : `array` d'objets `UserInvite`

```json
[
  {
    "id": "29214",
    "role": "admin",
    "tailnetId": 12345,
    "inviterId": 34567,
    "email": "user@example.com",
    "lastEmailSentAt": "2024-05-09T16:23:26.91778771Z",
    "inviteUrl": "https://login.tailscale.com/uinv/<code>"
  }
]
```

#### 404 - Tailnet not found

#### 500 - Internal server error

---

## Endpoint 5 : Get a user invite

- **Chemin** : `/user-invites/{userInviteId}`
- **Methode HTTP** : `GET`
- **operationId** : `getUserInvite`
- **Tags** : `UserInvites`
- **Summary** : Get a user invite
- **Description** : Retrieve a specific user invite.

### Parametres

| Nom            | In   | Type   | Requis | Description                 |
|----------------|------|--------|--------|-----------------------------|
| `userInviteId` | path | string | oui    | ID de l'invite utilisateur. |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Type : objet `UserInvite` (voir schema ci-dessus)

#### 404 - User invite not found

#### 500 - Internal server error

---

## Endpoint 6 : Delete a user invite

- **Chemin** : `/user-invites/{userInviteId}`
- **Methode HTTP** : `DELETE`
- **operationId** : `deleteUserInvite`
- **Tags** : `UserInvites`
- **Summary** : Delete a user invite
- **Description** : Deletes a specific user invite. Only permitted for user-owned keys, because invites require an inviting user.

### Parametres

| Nom            | In   | Type   | Requis | Description                 |
|----------------|------|--------|--------|-----------------------------|
| `userInviteId` | path | string | oui    | ID de l'invite utilisateur. |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Pas de corps de reponse.

#### 404 - User invite not found

#### 500 - Internal server error

---

## Endpoint 7 : Resend a user invite

- **Chemin** : `/user-invites/{userInviteId}/resend`
- **Methode HTTP** : `POST`
- **operationId** : `resendUserInvite`
- **Tags** : `UserInvites`
- **Summary** : Resend a user invite
- **Description** : Resend a user invite by email. You can only use this if the specified invite was originally created with an email specified. Refer to creating user invites for a tailnet. Note: Invite resends are rate limited to one per minute. Only permitted for user-owned keys, because invites require an inviting user.

### Parametres

| Nom            | In   | Type   | Requis | Description                 |
|----------------|------|--------|--------|-----------------------------|
| `userInviteId` | path | string | oui    | ID de l'invite utilisateur. |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Pas de corps de reponse.

#### 404 - User invite not found

#### 500 - Internal server error

---

## Endpoint 8 : Get a device invite

- **Chemin** : `/device-invites/{deviceInviteId}`
- **Methode HTTP** : `GET`
- **operationId** : `getDeviceInvite`
- **Tags** : `DeviceInvites`
- **Summary** : Get a device invite
- **Description** : Retrieve a specific device invite. OAuth Scope: `device_invites:read`.

### Parametres

| Nom              | In   | Type   | Requis | Description              |
|------------------|------|--------|--------|--------------------------|
| `deviceInviteId` | path | string | oui    | ID de l'invite device.   |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Type : objet `DeviceInvite` (voir schema ci-dessus)

#### 404 - Device invite not found

#### 500 - Internal server error

---

## Endpoint 9 : Delete a device invite

- **Chemin** : `/device-invites/{deviceInviteId}`
- **Methode HTTP** : `DELETE`
- **operationId** : `deleteDeviceInvite`
- **Tags** : `DeviceInvites`
- **Summary** : Delete a device invite
- **Description** : Delete a specific device invite. OAuth Scope: `device_invites`.

### Parametres

| Nom              | In   | Type   | Requis | Description              |
|------------------|------|--------|--------|--------------------------|
| `deviceInviteId` | path | string | oui    | ID de l'invite device.   |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Pas de corps de reponse.

#### 404 - Device invite not found

#### 500 - Internal server error

---

## Endpoint 10 : Resend a device invite

- **Chemin** : `/device-invites/{deviceInviteId}/resend`
- **Methode HTTP** : `POST`
- **operationId** : `resendDeviceInvite`
- **Tags** : `DeviceInvites`
- **Summary** : Resend a device invite
- **Description** : Resend a device invite by email. You can only use this if the specified invite was originally created with an email specified. Refer to creating device invites for a device. Note: Invite resends are rate limited to one per minute. Note that device invites cannot be resent using an API access token generated from an OAuth client as the shared device is scoped to a user.

### Parametres

| Nom              | In   | Type   | Requis | Description              |
|------------------|------|--------|--------|--------------------------|
| `deviceInviteId` | path | string | oui    | ID de l'invite device.   |

Pas de query params. Pas de body.

### Reponses

#### 200 - Successful operation

Pas de corps de reponse.

#### 404 - Device invite not found

#### 500 - Internal server error

---

## Endpoint 11 : Accept a device invite

- **Chemin** : `/device-invites/-/accept`
- **Methode HTTP** : `POST`
- **operationId** : `acceptDeviceInvite`
- **Tags** : `DeviceInvites`
- **Summary** : Accept a device invite
- **Description** : Accepts the invitation to share a device into the requesting user's tailnet. Note that device invites cannot be accepted using an API access token generated from an OAuth client as the shared device is scoped to a user.

### Parametres

Pas de path params. Pas de query params.

### Corps de la requete (requestBody)

Content-Type : `application/json`

Type : `object`

| Champ    | Type   | Requis | Exemple                                                | Description                                                                                                                                             |
|----------|--------|--------|---------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| `invite` | string | **oui** | `"https://login.tailscale.com/admin/invite/xxxxxx"`    | L'URL de l'invite (sous la forme `https://login.tailscale.com/admin/invite/{code}`) ou le composant `{code}` seul de l'URL.                             |

Exemples de body :

```json
// Avec l'URL complete
{
  "invite": "https://login.tailscale.com/admin/invite/xxxxxx"
}

// Avec le code uniquement
{
  "invite": "xxxxxx"
}
```

### Reponses

#### 200 - Successful operation

Type : `object` avec les proprietes suivantes :

**`device`** (object) - Informations sur le device partage :

| Champ              | Type    | Exemple                          | Description                                                                                  |
|--------------------|---------|----------------------------------|----------------------------------------------------------------------------------------------|
| `device.id`        | string  | `"12346"`                        | Le `nodeId` du device.                                                                       |
| `device.os`        | string  | `"iOS"`                          | Le systeme d'exploitation du device.                                                         |
| `device.name`      | string  | `"my-phone"`                     | Le nom du device.                                                                            |
| `device.fqdn`      | string  | `"my-phone.something.ts.net"`    | Le nom MagicDNS du device.                                                                   |
| `device.ipv4`      | string  | `"100.x.y.z"`                    | L'adresse IPv4 du device.                                                                    |
| `device.ipv6`      | string  | `"fd7a:115c:x::y:z"`            | L'adresse IPv6 du device.                                                                    |
| `device.includeExitNode` | boolean | `false`                    | Indique si l'utilisateur invite peut utiliser le device comme exit node.                     |

**`sharer`** (object) - L'utilisateur qui a cree l'invite de partage :

| Champ                  | Type   | Exemple                    | Description                                              |
|------------------------|--------|----------------------------|----------------------------------------------------------|
| `sharer.id`            | string | `"22012"`                  | ID de l'utilisateur qui a cree l'invite.                 |
| `sharer.displayName`   | string | `"Some User"`              | Nom d'affichage de l'utilisateur qui a cree l'invite.    |
| `sharer.loginName`     | string | `"someuser@example.com"`   | Adresse email de l'utilisateur qui a cree l'invite.      |
| `sharer.profilePicURL` | string | `""`                       | URL de la photo de profil du createur de l'invite.       |

**`acceptedBy`** (object) - L'utilisateur qui accepte l'invite :

| Champ                      | Type   | Exemple                        | Description                                              |
|----------------------------|--------|--------------------------------|----------------------------------------------------------|
| `acceptedBy.id`            | string | `"33233"`                      | ID de l'utilisateur qui a accepte l'invite.              |
| `acceptedBy.displayName`   | string | `"Another User"`               | Nom d'affichage de l'utilisateur qui a accepte.          |
| `acceptedBy.loginName`     | string | `"anotheruser@example2.com"`   | Adresse email de l'utilisateur qui a accepte.            |
| `acceptedBy.profilePicURL` | string | `""`                           | URL de la photo de profil de l'utilisateur qui a accepte.|

#### 400 - Bad request

#### 500 - Internal server error

---

## Resume des endpoints

| #  | Methode  | Chemin                                       | operationId           | Tag            |
|----|----------|----------------------------------------------|-----------------------|----------------|
| 1  | GET      | `/device/{deviceId}/device-invites`          | listDeviceInvites     | DeviceInvites  |
| 2  | POST     | `/device/{deviceId}/device-invites`          | createDeviceInvites   | DeviceInvites  |
| 3  | GET      | `/tailnet/{tailnet}/user-invites`            | listUserInvites       | UserInvites    |
| 4  | POST     | `/tailnet/{tailnet}/user-invites`            | createUserInvites     | UserInvites    |
| 5  | GET      | `/user-invites/{userInviteId}`               | getUserInvite         | UserInvites    |
| 6  | DELETE   | `/user-invites/{userInviteId}`               | deleteUserInvite      | UserInvites    |
| 7  | POST     | `/user-invites/{userInviteId}/resend`        | resendUserInvite      | UserInvites    |
| 8  | GET      | `/device-invites/{deviceInviteId}`           | getDeviceInvite       | DeviceInvites  |
| 9  | DELETE   | `/device-invites/{deviceInviteId}`           | deleteDeviceInvite    | DeviceInvites  |
| 10 | POST     | `/device-invites/{deviceInviteId}/resend`    | resendDeviceInvite    | DeviceInvites  |
| 11 | POST     | `/device-invites/-/accept`                   | acceptDeviceInvite    | DeviceInvites  |
