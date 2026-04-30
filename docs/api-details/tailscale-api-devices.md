# Tailscale API - Section Devices (Endpoints exhaustifs)

---

## Paramètres réutilisables (composants référencés)

### `tailnet` (path parameter)
- **In :** path
- **Name :** `tailnet`
- **Type :** string
- **Requis :** oui
- **Description :** L'identifiant du tailnet. On peut fournir un tiret (`-`) pour référencer le tailnet par défaut du token d'accès, ou le tailnet ID trouvé dans les General Settings de la console admin (ex : `T1234CNTRL`).
- **Exemple :** `example.com`

### `deviceId` (path parameter)
- **In :** path
- **Name :** `deviceId`
- **Type :** string
- **Requis :** oui
- **Description :** ID du device. L'utilisation du `nodeId` du device est préférée, mais sa valeur numérique `id` peut aussi être utilisée.

### `fields` (query parameter)
- **In :** query
- **Name :** `fields`
- **Type :** string (enum : `all`, `default`)
- **Requis :** non
- **Description :** Contrôle si la réponse retourne **tous** les champs ou seulement un sous-ensemble prédéfini.
  - `all` : retourne tous les champs
  - `default` : retourne uniquement : `addresses`, `id`, `nodeId`, `user`, `name`, `hostname`, `clientVersion`, `updateAvailable`, `os`, `created`, `connectedToControl`, `lastSeen`, `keyExpiryDisabled`, `expires`, `authorized`, `isExternal`, `machineKey`, `nodeKey`, `blocksIncomingConnections`, `tailnetLockKey`, `tailnetLockError`, `tags`, `isEphemeral`
- **Exemple :** `all`

### `attributeKey` (path parameter)
- **In :** path
- **Name :** `attributeKey`
- **Type :** string
- **Requis :** oui
- **Description :** Le nom de l'attribut posture à définir. Doit être préfixé par `custom:`. Longueur max 128 caractères (namespace inclus). Ne peut contenir que lettres, chiffres, underscores et deux-points. Les clés sont sensibles à la casse mais vérifiées pour l'unicité de manière insensible à la casse. Toutes les valeurs pour une clé donnée doivent être du même type.

---

## Schemas réutilisés

### Schema `Device`

Objet représentant un appareil Tailscale (aussi appelé *node* ou *machine*).

| Champ | Type | Exemple | Description |
|---|---|---|---|
| `addresses` | array of string | `["100.87.74.78", "fd7a:115c:a1e0:ac82:4843:ca90:697d:c36e"]` | Liste des adresses IP Tailscale (IPv4 100.x.y.z et IPv6) |
| `id` | string | `"92960230385"` | Identifiant legacy du device (préférer `nodeId`) |
| `nodeId` | string | `"n292kg92CNTRL"` | Identifiant préféré du device |
| `user` | string | `"amelie@example.com"` | L'utilisateur ayant enregistré le node |
| `name` | string | `"pangolin.tailfe8c.ts.net"` | Nom MagicDNS du device |
| `hostname` | string | `"pangolin"` | Nom de machine dans la console admin |
| `clientVersion` | string | `"v1.36.0"` | Version du client Tailscale (vide pour les devices externes) |
| `updateAvailable` | boolean | `false` | `true` si une mise à jour est disponible (vide pour devices externes) |
| `os` | string | `"linux"` | Système d'exploitation du device |
| `created` | string (date-time) | `"2022-12-01T05:23:30Z"` | Date d'ajout au tailnet (vide pour devices externes) |
| `connectedToControl` | boolean | `true` | Indique si le device a récemment maintenu une connexion TCP au serveur de contrôle |
| `lastSeen` | string (date-time) | `"2022-12-01T05:23:30Z"` | Dernière connexion au serveur de contrôle. Omis si jamais en ligne ou `connectedToControl` est `true` |
| `keyExpiryDisabled` | boolean | `false` | `true` si les clés ne vont pas expirer |
| `expires` | string (date-time) | `"2023-05-30T04:44:05Z"` | Date d'expiration de la clé d'auth du device |
| `authorized` | boolean | `false` | `true` si le device est autorisé à rejoindre le tailnet |
| `isExternal` | boolean | `false` | `true` si le device est partagé dans le tailnet (pas membre direct) |
| `multipleConnections` | boolean | `true` | `true` si plusieurs devices sont connectés avec la même node key. Omis si 0 ou 1 connexion |
| `machineKey` | string | `""` | Usage interne, vide pour devices externes |
| `nodeKey` | string | `"nodekey:01234567890abcdef"` | Usage interne, requis pour certaines opérations (ex: tailnet lock) |
| `blocksIncomingConnections` | boolean | `false` | `true` si le device refuse toutes les connexions entrantes (y compris les pings) |
| `enabledRoutes` | array of string | `["10.0.0.0/16", "192.168.1.0/24"]` | Routes subnet approuvées par un admin du tailnet |
| `advertisedRoutes` | array of string | `["10.0.0.0/16", "192.168.1.0/24"]` | Subnets que le device demande à exposer |
| `clientConnectivity` | object | (voir ci-dessous) | Rapport sur les conditions réseau physiques actuelles du device |
| `clientConnectivity.endpoints` | array of string | `["199.9.14.201:59128", "192.68.0.21:59128"]` | Endpoints UDP IP:port de magicsock |
| `clientConnectivity.mappingVariesByDestIP` | boolean | `false` | `true` si les mappings NAT varient selon l'IP de destination |
| `clientConnectivity.latency` | object (map) | `{"Dallas": {"latencyMs": 60.46}, "New York City": {"preferred": true, "latencyMs": 31.32}}` | Map des serveurs DERP et leur latence. Chaque entrée a `preferred` (boolean) et `latencyMs` (float) |
| `clientConnectivity.clientSupports` | object | (voir ci-dessous) | Fonctionnalités supportées par le client |
| `clientConnectivity.clientSupports.hairPinning` | boolean \| null | `null` | Plus tracé, toujours null |
| `clientConnectivity.clientSupports.ipv6` | boolean \| null | `false` | `true` si l'OS supporte IPv6 |
| `clientConnectivity.clientSupports.pcp` | boolean \| null | `false` | `true` si un service PCP existe sur le routeur |
| `clientConnectivity.clientSupports.pmp` | boolean \| null | `false` | `true` si un service NAT-PMP existe sur le routeur |
| `clientConnectivity.clientSupports.udp` | boolean \| null | `false` | `true` si le trafic UDP est activé |
| `clientConnectivity.clientSupports.upnp` | boolean \| null | `false` | `true` si un service UPnP existe sur le routeur |
| `tags` | array of string | `["tag:golink"]` | Tags assignés au device (vide pour devices externes) |
| `tailnetLockError` | string | `""` | Indique un problème avec la signature tailnet lock (peuplé uniquement si tailnet lock activé) |
| `tailnetLockKey` | string | `""` | Clé tailnet lock du node (toujours présente même si tailnet lock désactivé) |
| `sshEnabled` | boolean | `false` | `true` si Tailscale SSH est activé sur ce device |
| `postureIdentity` | object | `{"serialNumbers": ["CP74LFQJXM"]}` | Identifiants supplémentaires du device pour la posture. Contient `serialNumbers` (array of string) et/ou `disabled` (boolean) |
| `isEphemeral` | boolean | `false` | `true` si le device est éphémère |
| `distro` | object | `{"name": "ubuntu", "version": "25.04", "codeName": "Plucky Puffin"}` | Détails de la distribution OS |
| `distro.name` | string | `"ubuntu"` | Nom de la distribution |
| `distro.version` | string | `"25.04"` | Version de la distribution |
| `distro.codeName` | string | `"Plucky Puffin"` | Nom de code de la distribution |

### Schema `DeviceRoutes`

| Champ | Type | Exemple | Description |
|---|---|---|---|
| `advertisedRoutes` | array of string | `["10.0.0.0/16", "192.168.1.0/24"]` | Subnets que le device demande à exposer |
| `enabledRoutes` | array of string | `["10.0.0.0/16", "192.168.1.0/24"]` | Routes subnet approuvées par un admin du tailnet |

### Schema `DevicePostureAttributes`

| Champ | Type | Exemple | Description |
|---|---|---|---|
| `attributes` | object (map: string/number/boolean) | `{"custom:myScore": 80, "custom:diskEncryption": true, "node:os": "linux", "node:tsVersion": "1.40.0", ...}` | Tous les attributs posture assignés au node. Les valeurs peuvent être strings, numbers ou booleans |
| `expiries` | object (map: string date-time) | `{"custom:myScore": "2024-04-23T18:25:43.511Z"}` | Temps d'expiration pour chaque attribut posture, si défini |

---

## Endpoint 1 : List tailnet devices

- **Path :** `/tailnet/{tailnet}/devices`
- **Méthode HTTP :** `GET`
- **Operation ID :** `listTailnetDevices`
- **Tags :** `Devices`
- **Summary :** List tailnet devices
- **Description :** Lists the devices in a tailnet. OAuth Scope: `devices:core:read`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `tailnet` | path | string | oui | L'identifiant du tailnet (voir définition ci-dessus) |
| `fields` | query | string (enum: `all`, `default`) | non | Contrôle les champs retournés (voir définition ci-dessus) |
| `<field>=<value> filters` | query | string | non | Filtrage côté serveur des devices sous la forme `<field>=<value>`. Les champs doivent être des propriétés de premier niveau du device (ex: `isEphemeral`, `tags`, `hostname`). Matching exact. Types simples (strings, numbers, dates) et listes supportés. Objets complexes (ex: `clientConnectivity`) non supportés. Plusieurs paramètres = AND logique. Exemple : `isEphemeral=true&tags=tag:prod&tags=tag:subnetrouter` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | `{ "devices": [ Device, ... ] }` — tableau d'objets `Device` (voir schema ci-dessus) |
| `404` | Tailnet not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Request took too long to process, please try again later | Error object |

---

## Endpoint 2 : Batch update custom device posture attributes

- **Path :** `/tailnet/{tailnet}/device-attributes`
- **Méthode HTTP :** `PATCH`
- **Operation ID :** `batchUpdateCustomDevicePostureAttributes`
- **Tags :** `Devices`
- **Summary :** Batch update custom device posture attributes
- **Description :** Batch updates posture attributes across devices in a tailnet. Utilise la sémantique JSON Merge Patch (RFC 7396). Spécifier `null` pour un attribut le supprime. Les attributs doivent être dans le namespace `custom:`. OAuth Scope: `devices:posture_attributes`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `tailnet` | path | string | oui | L'identifiant du tailnet |

### Corps de la requête (JSON, requis)

```json
{
  "nodes": {
    "<deviceId>": {
      "<attributeName>": {
        "value": <string | number | boolean>,   // requis
        "expiry": "<date-time>"                   // optionnel
      }
      // ... ou null pour supprimer l'attribut
    }
    // ... autres devices
  },
  "comment": "<string, max 200 chars>"  // optionnel, ajouté aux audit logs
}
```

**Détails des champs du body :**

| Champ | Type | Requis | Description |
|---|---|---|---|
| `nodes` | object (map de deviceId → map d'attributs) | non | Mapping deviceId → attributs posture |
| `nodes.<deviceId>` | object (map d'attributeName → config ou null) | — | Mapping nom d'attribut → configuration |
| `nodes.<deviceId>.<attr>.value` | string \| number \| boolean | oui (dans chaque config) | Valeur de l'attribut |
| `nodes.<deviceId>.<attr>.expiry` | string (date-time) | non | Temps d'expiration optionnel. Si défini, Tailscale supprime automatiquement l'attribut quelques minutes après |
| `nodes.<deviceId>.<attr>` = `null` | null | — | Supprimer l'attribut |
| `comment` | string (max 200 chars) | non | Commentaire optionnel pour l'audit log |

**Exemple de body :**
```json
{
  "nodes": {
    "nPM2KNuedB21DEVEL": {
      "custom:myattr": {
        "value": "my_value"
      }
    },
    "nPpz3VEKzX11DEVEL": {
      "custom:flag": {
        "value": true,
        "expiry": "2025-09-19T15:00:00Z"
      }
    }
  },
  "comment": "bulk posture attribute update"
}
```

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | `null` |
| `400` | An invalid request payload was sent | Error object |
| `404` | Tailnet not found | Error object |
| `500` | Internal server error | Error object |

---

## Endpoint 3 : Get a device

- **Path :** `/device/{deviceId}`
- **Méthode HTTP :** `GET`
- **Operation ID :** `getDevice`
- **Tags :** `Devices`
- **Summary :** Get a device
- **Description :** Retrieve the details for the specified device. OAuth Scope: `devices:core:read`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device (`nodeId` préféré, `id` numérique accepté) |
| `fields` | query | string (enum: `all`, `default`) | non | Contrôle les champs retournés |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | Objet `Device` (voir schema ci-dessus) |
| `400` | Invalid ID supplied | Error object |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 4 : Delete a device

- **Path :** `/device/{deviceId}`
- **Méthode HTTP :** `DELETE`
- **Operation ID :** `deleteDevice`
- **Tags :** `Devices`
- **Summary :** Delete a device
- **Description :** Deletes the device from its tailnet. The device must belong to the requesting user's tailnet. Deleting devices shared with the tailnet is not supported. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête

Aucun.

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `400` | Invalid device value | Error object |
| `500` | Internal server error | Error object |
| `501` | Device not owned by tailnet | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 5 : Expire a device's key

- **Path :** `/device/{deviceId}/expire`
- **Méthode HTTP :** `POST`
- **Operation ID :** `expireDeviceKey`
- **Tags :** `Devices`
- **Summary :** Expire a device's key
- **Description :** Mark a device's node key as expired. This will require the device to re-authenticate in order to connect to the tailnet. The device must belong to the requesting user's tailnet. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête

Aucun.

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 6 : List device routes

- **Path :** `/device/{deviceId}/routes`
- **Méthode HTTP :** `GET`
- **Operation ID :** `listDeviceRoutes`
- **Tags :** `Devices`
- **Summary :** List device routes
- **Description :** Retrieve the list of subnet routes that a device is advertising, as well as those that are enabled for it. Routes must be both advertised and enabled for a device to act as a subnet router or exit node. If a device has advertised routes, they are not exposed to traffic until they are enabled. Conversely, if routes are enabled before they are advertised, they are not available for routing until the device in question has advertised them. OAuth Scope: `devices:routes:read`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête

Aucun.

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | Objet `DeviceRoutes` : `{ "advertisedRoutes": ["10.0.0.0/16", ...], "enabledRoutes": ["10.0.0.0/16", ...] }` |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 7 : Set device routes

- **Path :** `/device/{deviceId}/routes`
- **Méthode HTTP :** `POST`
- **Operation ID :** `setDeviceRoutes`
- **Tags :** `Devices`
- **Summary :** Set device routes
- **Description :** Set a device's enabled subnet routes by replacing the existing list of subnet routes with the supplied parameters. Advertised routes cannot be set through the API, since they must be set directly on the device. Routes must be both advertised and enabled for a device to act as a subnet router or exit node. OAuth Scope: `devices:routes`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON, requis)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `routes` | array of string | non | La nouvelle liste de routes subnet activées | `["10.0.0.0/16", "192.168.1.0/24"]` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | Objet `DeviceRoutes` : `{ "advertisedRoutes": [...], "enabledRoutes": [...] }` |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 8 : Authorize device

- **Path :** `/device/{deviceId}/authorized`
- **Méthode HTTP :** `POST`
- **Operation ID :** `authorizeDevice`
- **Tags :** `Devices`
- **Summary :** Authorize device
- **Description :** This call marks a device as authorized or revokes its authorization for tailnets where device authorization is required, according to the `authorized` field in the payload. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON)

| Champ | Type | Requis | Description |
|---|---|---|---|
| `authorized` | boolean | **oui** | `true` pour autoriser un nouveau device ou re-autoriser un device précédemment désautorisé. `false` pour désautoriser un device autorisé. |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 9 : Set device name

- **Path :** `/device/{deviceId}/name`
- **Méthode HTTP :** `POST`
- **Operation ID :** `setDeviceName`
- **Tags :** `Devices`
- **Summary :** Set device name
- **Description :** When a device is added to a tailnet, its Tailscale device name (also sometimes referred to as machine name) is generated from its OS hostname. The device name is the canonical name for the device on your tailnet. Device name changes immediately get propagated through your tailnet, so be aware that any existing Magic DNS URLs using the old name will no longer work. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `name` | string | **oui** | Le nouveau nom pour le device. Peut être fourni en FQDN (ex: `"nodename.your-domain.ts.net"`) ou juste le nom de base (ex: `"nodename"`). Si non défini ou vide, le nom est réinitialisé à partir du hostname OS. | `"dev-server"` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 10 : Set device tags

- **Path :** `/device/{deviceId}/tags`
- **Méthode HTTP :** `POST`
- **Operation ID :** `setDeviceTags`
- **Tags :** `Devices`
- **Summary :** Set device tags
- **Description :** Tags let you assign an identity to a device that is separate from human users, and use that identity as part of an ACL to restrict access. Tags are similar to role accounts, but more flexible. Tags are created in the tailnet policy file by defining the tag and an owner of the tag. Once a device is tagged, the tag is the owner of that device. A single node can have multiple tags assigned. Consult the policy file for your tailnet in the admin console for the list of tags that have been created for your tailnet. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `tags` | array of string | non | La nouvelle liste de tags pour le device | `["tag:foo", "tag:bar"]` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `400` | Bad request | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 11 : Update device key

- **Path :** `/device/{deviceId}/key`
- **Méthode HTTP :** `POST`
- **Operation ID :** `updateDeviceKey`
- **Tags :** `Devices`
- **Summary :** Update device key
- **Description :** When a device is added to a tailnet, its key expiry is set according to the tailnet's key expiry setting. If the key is not refreshed and expires, the device can no longer communicate with other devices in the tailnet. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `keyExpiryDisabled` | boolean | **oui** | `true` : désactiver l'expiration de la clé du device (le temps d'expiration original est conservé ; à la réactivation, la clé expirera à ce moment original). `false` : activer l'expiration de la clé du device (la clé peut déjà avoir expiré, dans ce cas le device doit être ré-authentifié). | `true` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 12 : Set device IPv4 address

- **Path :** `/device/{deviceId}/ip`
- **Méthode HTTP :** `POST`
- **Operation ID :** `setDeviceIp`
- **Tags :** `Devices`
- **Summary :** Set device IPv4 address
- **Description :** When a device is added to a tailnet, its Tailscale IPv4 address is set at random either from the CGNAT range, or a subset of the CGNAT range specified by an ip pool. This endpoint can be used to replace the existing IPv4 address with a specific value. This action will break any existing connections to this machine. You will need to reconnect to this machine using the new IP address. You may also need to flush your DNS cache. OAuth Scope: `devices:core`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête (JSON)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `ipv4` | string | **oui** | La nouvelle adresse IPv4 pour le device | `"100.80.0.1"` |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 13 : Get device posture attributes

- **Path :** `/device/{deviceId}/attributes`
- **Méthode HTTP :** `GET`
- **Operation ID :** `getDevicePostureAttributes`
- **Tags :** `Devices`
- **Summary :** Get device posture attributes
- **Description :** Retrieve all posture attributes for the specified device. This returns a JSON object of all the key-value pairs of posture attributes for the device. OAuth Scope: `devices:posture_attributes:read`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |

### Corps de la requête

Aucun.

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | Objet `DevicePostureAttributes` (voir schema). Exemple : `{ "attributes": { "custom:myScore": 80, "custom:diskEncryption": true, "custom:myAttribute": "my_value", "node:os": "linux", "node:osVersion": "5.19.0-42-generic", "node:tsReleaseTrack": "stable", "node:tsVersion": "1.40.0", "node:tsAutoUpdate": false, "node:tsStateEncrypted": false }, "expiries": { "custom:myScore": "2024-04-23T18:25:43.511Z" } }` |
| `404` | Device not found | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 14 : Set custom device posture attributes

- **Path :** `/device/{deviceId}/attributes/{attributeKey}`
- **Méthode HTTP :** `POST`
- **Operation ID :** `setCustomDevicePostureAttributes`
- **Tags :** `Devices`
- **Summary :** Set custom device posture attributes
- **Description :** Create or update a custom posture attribute on the specified device. User-managed attributes must be in the `custom` namespace, which is indicated by prefixing the attribute key with `custom:`. OAuth Scope: `devices:posture_attributes`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |
| `attributeKey` | path | string | oui | Nom de l'attribut posture (doit commencer par `custom:`). Max 128 chars, lettres/chiffres/underscores/deux-points uniquement. Sensible à la casse, unicité vérifiée de manière insensible. |

### Corps de la requête (JSON, requis)

| Champ | Type | Requis | Description | Exemple |
|---|---|---|---|---|
| `value` | string \| number \| boolean | non | Valeur de l'attribut. String : max 50 chars, lettres/chiffres/underscores/points uniquement. Number : entier, JSON safe (jusqu'à 2^53 - 1). | `"my_value"` |
| `expiry` | string (date-time) | non | Temps d'expiration optionnel. Si défini, Tailscale supprime automatiquement l'attribut quelques minutes après le temps spécifié. | `"2022-12-01T05:23:30Z"` |
| `comment` | string (max 200 chars) | non | Commentaire optionnel indiquant pourquoi l'attribut est défini, ajouté à l'audit log. | — |

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | Objet `DevicePostureAttributes` (voir schema) |
| `404` | Device not found | Error object |
| `429` | Too many requests (rate limited) | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Endpoint 15 : Delete custom device posture attributes

- **Path :** `/device/{deviceId}/attributes/{attributeKey}`
- **Méthode HTTP :** `DELETE`
- **Operation ID :** `deleteCustomDevicePostureAttributes`
- **Tags :** `Devices`
- **Summary :** Delete custom device posture attributes
- **Description :** Delete a posture attribute from the specified device. This is only applicable to user-managed posture attributes in the `custom` namespace, which is indicated by prefixing the attribute key with `custom:`. OAuth Scope: `devices:posture_attributes`.

### Paramètres

| Nom | In | Type | Requis | Description |
|---|---|---|---|---|
| `deviceId` | path | string | oui | ID du device |
| `attributeKey` | path | string | oui | Nom de l'attribut posture à supprimer (doit commencer par `custom:`) |

### Corps de la requête

Aucun.

### Réponses

| Code | Description | Body |
|---|---|---|
| `200` | Successful operation | (vide) |
| `404` | Device not found | Error object |
| `429` | Too many requests (rate limited) | Error object |
| `500` | Internal server error | Error object |
| `504` | Timeout | Error object |

---

## Récapitulatif des endpoints

| # | Méthode | Path | Operation ID | Summary |
|---|---|---|---|---|
| 1 | `GET` | `/tailnet/{tailnet}/devices` | `listTailnetDevices` | List tailnet devices |
| 2 | `PATCH` | `/tailnet/{tailnet}/device-attributes` | `batchUpdateCustomDevicePostureAttributes` | Batch update custom device posture attributes |
| 3 | `GET` | `/device/{deviceId}` | `getDevice` | Get a device |
| 4 | `DELETE` | `/device/{deviceId}` | `deleteDevice` | Delete a device |
| 5 | `POST` | `/device/{deviceId}/expire` | `expireDeviceKey` | Expire a device's key |
| 6 | `GET` | `/device/{deviceId}/routes` | `listDeviceRoutes` | List device routes |
| 7 | `POST` | `/device/{deviceId}/routes` | `setDeviceRoutes` | Set device routes |
| 8 | `POST` | `/device/{deviceId}/authorized` | `authorizeDevice` | Authorize device |
| 9 | `POST` | `/device/{deviceId}/name` | `setDeviceName` | Set device name |
| 10 | `POST` | `/device/{deviceId}/tags` | `setDeviceTags` | Set device tags |
| 11 | `POST` | `/device/{deviceId}/key` | `updateDeviceKey` | Update device key |
| 12 | `POST` | `/device/{deviceId}/ip` | `setDeviceIp` | Set device IPv4 address |
| 13 | `GET` | `/device/{deviceId}/attributes` | `getDevicePostureAttributes` | Get device posture attributes |
| 14 | `POST` | `/device/{deviceId}/attributes/{attributeKey}` | `setCustomDevicePostureAttributes` | Set custom device posture attributes |
| 15 | `DELETE` | `/device/{deviceId}/attributes/{attributeKey}` | `deleteCustomDevicePostureAttributes` | Delete custom device posture attributes |
