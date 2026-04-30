# Tailscale API - Sections Keys et Policy File (ACL)

Extraction exhaustive des endpoints des lignes 1793 a 2645 du fichier OpenAPI.

---

## Schemas references

### Schema `Key`

Objet representant une cle API, un auth key, un OAuth client ou une federated identity.

| Champ | Type | Format | Description | Exemple |
|-------|------|--------|-------------|---------|
| `id` | string | - | Identifiant de la cle | `k123456CNTRL` |
| `key` | string | - | Materiau secret de la cle (uniquement renseigne a la creation) | `tskey-auth-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` |
| `keyType` | string | enum: `auth`, `client`, `api`, `federated` | Type de cle. `auth` = machine auth keys, `client` = OAuth clients, `federated` = federated identities, `api` = personal API access tokens ou tokens generes via OAuth client/federated identity | `auth` |
| `expirySeconds` | integer | int64 | Duree en secondes avant expiration. Uniquement pour les auth keys. | `7776000` |
| `created` | string | date-time | Date de creation | `2021-12-09T23:22:39Z` |
| `updated` | string | date-time | Date de derniere mise a jour | `2021-12-09T23:22:39Z` |
| `expires` | string | date-time | Date d'expiration | `2022-03-09T23:22:39Z` |
| `revoked` | string | date-time | Date de revocation | `2022-03-12T23:22:39Z` |
| `capabilities` | objet `KeyCapabilities` | - | Permissions de la cle (voir ci-dessous). Uniquement pour les auth keys. | - |
| `scopes` | array of string | - | Liste des scopes accordes. Uniquement pour OAuth clients, API access tokens et federated identities. | `["all:read"]` |
| `tags` | array of string | - | Liste des tags associes au trust credential. Les auth keys crees avec ce client doivent avoir exactement ces tags ou des tags possedes par les tags du client. Obligatoire si les scopes incluent "devices:core" ou "auth_keys". Uniquement pour OAuth clients et federated identities. | `["tag:example"]` |
| `description` | string | - | Description courte | `dev access` |
| `invalid` | boolean | - | `true` si la cle est revoquee (supprimee) ou expiree | `false` |
| `userId` | string | - | ID de l'utilisateur qui a cree cette cle, vide pour les cles creees par des trust credentials | `uscwcTtzzo11DEVEL` |
| `audience` | string | - | Valeur utilisee pour matcher le claim `aud` d'un token OIDC. Tailscale genere un audience securise par defaut si non specifie. Uniquement pour federated identities. | `api.tailscale.com/Tz8TefihCR11DEVEL-kqc11MVpwu11DEVEL` |
| `issuer` | string | uri | Issuer du token OIDC. Doit etre une URL https:// valide et publiquement accessible. Uniquement pour federated identities. | `https://example.com` |
| `subject` | string | - | Pattern utilise pour matcher le claim `sub` d'un token OIDC. Peut inclure des `*` comme wildcards. Uniquement pour federated identities. | `my-example-subject-*` |
| `customClaimRules` | object (additionalProperties: string) | - | Map de noms de claims vers des patterns pour matcher des claims arbitraires dans le token OIDC. Les patterns peuvent inclure `*`. Uniquement pour federated identities. | `{"exampleAdditionalClaim": "valueToMatch", "otherAdditionalClaim": "valueWithWildcard*"}` |

### Schema `KeyCapabilities`

Mapping de ressources vers des actions autorisees.

| Champ | Type | Description |
|-------|------|-------------|
| `devices` | object | Permissions de la cle sur les devices. Uniquement pour les auth keys. |
| `devices.create` | object | Permissions lors de la creation de devices. |
| `devices.create.reusable` | boolean | Auth keys reutilisables (peuvent etre utilises plusieurs fois pour enregistrer differents devices). Exemple: `true` |
| `devices.create.ephemeral` | boolean | Cles ephemeres, utilisees pour connecter et nettoyer des devices de courte duree. Exemple: `false` |
| `devices.create.preauthorized` | boolean | Cles pre-autorisees (pre-approved). `true` signifie que les devices enregistres avec cette cle n'auront pas besoin d'approbation supplementaire d'un admin. Exemple: `true` |
| `devices.create.tags` | array of string | Tags qui seront definis sur les devices enregistres avec cette cle. Pour les auth keys possedes par le tailnet (via OAuth), les tags sont obligatoires et doivent correspondre exactement aux tags du OAuth client. Pour les auth keys possedes par un utilisateur, les tags sont optionnels. Exemple: `["tag:example"]` |

### Parametres communs references

#### Parametre `tailnet`

| Propriete | Valeur |
|-----------|--------|
| In | path |
| Name | `tailnet` |
| Required | true |
| Type | string |
| Exemple | `example.com` |
| Description | L'ID du tailnet. Les tailnets crees avant oct 2025 peuvent encore utiliser l'ID legacy, mais le Tailnet ID est l'identifiant prefere. On peut fournir un tiret (`-`) pour referencer le tailnet par defaut du token d'acces utilise, ou fournir le **tailnet ID** depuis la page General Settings de la console admin. |

#### Parametre `keyId`

| Propriete | Valeur |
|-----------|--------|
| In | path |
| Name | `keyId` |
| Required | true |
| Type | string |
| Exemple | `k123456CNTRL` |
| Description | L'ID de la cle. L'ID de la cle peut etre trouve dans la console admin (https://login.tailscale.com/admin/settings/keys). |

#### Parametre `all`

| Propriete | Valeur |
|-----------|--------|
| In | query |
| Name | `all` |
| Required | true |
| Type | boolean |
| Exemple | `true` |
| Description | Si defini a `true`, retourne tous les auth keys, API access tokens et OAuth clients pour le tailnet. |

#### Parametre `AcceptHeaderParam`

| Propriete | Valeur |
|-----------|--------|
| In | header |
| Name | `Accept` |
| Required | false |
| Type | string |
| Description | La reponse est encodee en JSON si `application/json` est demande, sinon HuJSON sera retourne. |

---

## ENDPOINT 1 : List tailnet keys

- **Chemin** : `/tailnet/{tailnet}/keys`
- **Methode HTTP** : `GET`
- **Operation ID** : `listTailnetKeys`
- **Tags** : `Keys`

### Summary / Description

**Summary** : List tailnet keys

**Description** :
Retourne une liste des auth keys actifs, API access tokens et trust credentials.

Si le parametre `{all}` n'est pas specifie, l'ensemble des cles retournees depend du token d'acces utilise :
- Si l'appel est fait avec un user-owned API access token : retourne uniquement les cles possedees par cet utilisateur.
- Si l'appel est fait avec un access token derive d'un OAuth client : retourne tous les OAuth clients pour le tailnet.
- Si l'appel est fait avec un access token derive d'une federated identity : retourne toutes les federated identities pour le tailnet.

OAuth Scopes :
- `api_access_tokens:read` : acces aux personal API access tokens
- `auth_keys:read` : acces aux machine auth keys
- `oauth_keys:read` : acces aux OAuth clients et OAuth tokens
- `federated_keys:read` : acces aux federated identities

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `all` | query | boolean | oui | Si `true`, retourne tous les auth keys, API access tokens et OAuth clients pour le tailnet |

### Reponses

#### 200 - Successful operation

```json
{
  "keys": [
    {
      "id": "XXXX14CNTRL",
      "keyType": "client",
      "created": "2021-12-09T23:22:39Z",
      "scopes": ["all"],
      "description": "test key",
      "userId": "uscwcTtzzo11DEVEL"
    },
    {
      "id": "XXXXZ3CNTRL",
      "keyType": "api",
      "expirySeconds": 7776000,
      "created": "2021-12-09T23:22:39Z",
      "expires": "2022-03-09T23:22:39Z",
      "scopes": ["all"],
      "description": "production key",
      "userId": "uscwcTtzzo11DEVEL"
    },
    {
      "id": "XXXX43CNTRL",
      "keyType": "auth",
      "expirySeconds": 7776000,
      "created": "2021-12-09T23:22:39Z",
      "expires": "2022-03-09T23:22:39Z",
      "capabilities": {
        "devices": {
          "create": {
            "reusable": true,
            "ephemeral": false,
            "preauthorized": true,
            "tags": ["tag:example"]
          }
        }
      },
      "description": "dev access",
      "userId": "uscwcTtzzo11DEVEL"
    }
  ]
}
```

Structure de la reponse :

| Champ | Type | Description |
|-------|------|-------------|
| `keys` | array of `Key` | Liste des cles actives. Chaque element est un objet `Key` (voir schema ci-dessus). |

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 2 : Create an auth key or trust credential

- **Chemin** : `/tailnet/{tailnet}/keys`
- **Methode HTTP** : `POST`
- **Operation ID** : `createKey`
- **Tags** : `Keys`

### Summary / Description

**Summary** : Create an auth key or trust credential

**Description** :
Cree un nouvel auth key ou trust credential dans le tailnet specifie.
Les trust credentials incluent les OAuth clients et les federated identities.
La cle sera associee a l'utilisateur proprietaire du API access token utilise pour faire l'appel, ou, si l'appel est fait avec un access token derive d'un OAuth client ou d'une federated identity, la cle sera possedee par le tailnet.

Retourne un objet JSON avec la cle generee. La cle doit etre enregistree et gardee en securite car elle porte les capabilities ou scopes specifies dans la requete. L'identite de la cle est incorporee dans la cle elle-meme et peut etre utilisee pour effectuer des operations sur la cle (ex: la revoquer ou recuperer des informations). La cle complete ne peut plus etre recuperee apres la reponse initiale.

OAuth Scopes :
- `auth_keys` : acces pour creer des machine auth keys
- `oauth_keys` : acces pour creer des OAuth clients
- `federated_keys` : acces pour creer des federated identities

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |

### Corps de la requete (Request Body)

Content-Type : `application/json`

Les champs supportes varient selon la valeur du champ `keyType`.

Pour les auth keys : au minimum, le body doit avoir un objet `capabilities` avec un objet `devices` (peut etre un objet JSON vide). Sans rien d'autre, cela genere une cle a usage unique sans tags.

Pour les OAuth clients : au minimum, le body doit avoir au moins un scope.

Pour les federated identities : au minimum, le body doit avoir au moins un scope, un issuer valide et un subject.

| Champ | Type | Format | Requis | Description | Exemple |
|-------|------|--------|--------|-------------|---------|
| `keyType` | string | enum: `auth`, `client`, `federated` | non | Type de cle a creer. Defaut a "auth" si omis. | `"auth"` |
| `description` | string | - | non | Courte description du but de la cle. Maximum 50 caracteres alphanumeriques. Les tirets et espaces sont aussi autorises. | `"dev access"` |
| `capabilities` | objet `KeyCapabilities` | - | non (mais minimum requis pour auth keys) | Mapping de ressources vers actions autorisees (voir schema KeyCapabilities). | voir schema |
| `expirySeconds` | integer | int64 | non | Duree en secondes avant expiration de la cle. Defaut a 90 jours si non fourni. Uniquement pour les auth keys. | `86400` |
| `scopes` | array of string | - | non (mais requis pour OAuth clients et federated identities) | Liste des scopes a accorder a la cle. Au moins un scope requis pour OAuth clients et federated identities. Voir trust credentials scopes. Uniquement pour OAuth clients et federated identities. | `["all:read"]` |
| `tags` | array of string | - | non (mais obligatoire si scopes incluent "devices:core" ou "auth_keys") | Tags associes au trust credential. Les auth keys crees avec ce credential doivent avoir exactement ces tags ou des tags possedes par les tags du credential. Uniquement pour OAuth clients et federated identities. | `["tag:example"]` |
| `issuer` | string | uri | non (mais requis pour federated identities) | Issuer du token OIDC. Doit etre une URL https:// valide et publiquement accessible. Uniquement pour federated identities. | `"https://example.com"` |
| `subject` | string | - | non (mais requis pour federated identities) | Pattern utilise pour matcher le claim `sub` du token OIDC. Peut inclure `*` comme wildcard. Uniquement pour federated identities. | `"my-example-subject-*"` |
| `audience` | string | - | non | Valeur utilisee pour matcher le claim `aud` du token OIDC. Optionnel car Tailscale genere un audience securise par defaut a la creation. Recommande de laisser Tailscale generer sauf si l'IdP necessite un format specifique. Uniquement pour federated identities. | `"api.tailscale.com/Tz8TefihCR11DEVEL-kqc11MVpwu11DEVEL"` |
| `customClaimRules` | object (additionalProperties: string) | - | non | Map de noms de claims vers des patterns pour matcher des claims arbitraires dans le token OIDC. Peut inclure `*`. Uniquement pour federated identities. | `{"exampleAdditionalClaim": "valueToMatch", "otherAdditionalClaim": "valueWithWildcard*"}` |

### Reponses

#### 200 - Successful operation

Retourne un objet `Key` (voir schema complet ci-dessus) avec tous les champs, y compris le champ `key` contenant le materiau secret (uniquement disponible a la creation).

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 3 : Get key

- **Chemin** : `/tailnet/{tailnet}/keys/{keyId}`
- **Methode HTTP** : `GET`
- **Operation ID** : `getKey`
- **Tags** : `Keys`

### Summary / Description

**Summary** : Get key

**Description** :
Retourne un objet JSON avec des informations sur un API access token, un OAuth client, une federated identity ou un auth key specifique, comme ses dates de creation et d'expiration et ses capabilities.

OAuth Scopes :
- `api_access_tokens:read` : acces aux personal API access tokens
- `auth_keys:read` : acces aux machine auth keys
- `oauth_keys:read` : acces aux OAuth clients et OAuth tokens
- `federated_keys:read` : acces aux federated identities

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `keyId` | path | string | oui | L'ID de la cle. Peut etre trouve dans la console admin. Exemple: `k123456CNTRL` |

### Reponses

#### 200 - Successful operation

La reponse pour une cle revoquee (supprimee) ou expiree aura un champ `invalid` defini a `true`.

Retourne un objet `Key` (voir schema complet ci-dessus).

#### 404 - Tailnet or key not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 4 : Delete key

- **Chemin** : `/tailnet/{tailnet}/keys/{keyId}`
- **Methode HTTP** : `DELETE`
- **Operation ID** : `deleteKey`
- **Tags** : `Keys`

### Summary / Description

**Summary** : Delete key

**Description** :
Supprime un API access token ou auth key specifique.

OAuth Scopes :
- `api_access_tokens` : acces aux personal API access tokens
- `auth_keys` : acces aux machine auth keys
- `oauth_keys` : acces aux OAuth clients et OAuth tokens
- `federated_keys` : acces aux federated identities

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `keyId` | path | string | oui | L'ID de la cle. Exemple: `k123456CNTRL` |

### Reponses

#### 200 - Successful operation

Pas de corps de reponse.

#### 403 - User does not have sufficient access to delete this key

Reponse d'erreur standard.

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 5 : Set key

- **Chemin** : `/tailnet/{tailnet}/keys/{keyId}`
- **Methode HTTP** : `PUT`
- **Operation ID** : `setKey`
- **Tags** : `Keys`

### Summary / Description

**Summary** : Set key

**Description** :
Definit la configuration pour un OAuth client ou une federated identity existant(e).

OAuth Scopes :
- `oauth_keys` : acces aux OAuth clients
- `federated_keys` : acces aux federated identities

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `keyId` | path | string | oui | L'ID de la cle. Exemple: `k123456CNTRL` |

### Corps de la requete (Request Body)

Content-Type : `application/json`

Les champs supportes varient selon la valeur du champ `keyType`.

| Champ | Type | Format | Requis | Description | Exemple |
|-------|------|--------|--------|-------------|---------|
| `keyType` | string | enum: `client`, `federated` | non | Le type de la cle mise a jour. | `"client"` |
| `description` | string | - | non | Courte description du but de la cle. Maximum 50 caracteres alphanumeriques. Tirets et espaces autorises. | `"dev access"` |
| `scopes` | array of string | - | non | Liste des scopes a accorder a la cle. Au moins un scope requis. Voir trust credentials scopes. | `["all:read"]` |
| `tags` | array of string | - | non | Tags associes au trust credential. Obligatoire si les scopes incluent "devices:core" ou "auth_keys". Les auth keys crees avec ce credential doivent avoir exactement ces tags ou des tags possedes par les tags du credential. | `["tag:example"]` |
| `issuer` | string | uri | non | Issuer du token OIDC. Doit etre une URL https:// valide. Uniquement pour federated identities. | `"https://example.com"` |
| `subject` | string | - | non | Pattern utilise pour matcher le claim `sub` du token OIDC. Peut inclure `*`. Uniquement pour federated identities. | `"my-example-subject-*"` |
| `audience` | string | - | non | Valeur utilisee pour matcher le claim `aud` du token OIDC. Uniquement pour federated identities. | `"api.tailscale.com/Tz8TefihCR11DEVEL-kqc11MVpwu11DEVEL"` |
| `customClaimRules` | object (additionalProperties: string) | - | non | Map de noms de claims vers des patterns pour matcher des claims arbitraires dans le token OIDC. Peut inclure `*`. Uniquement pour federated identities. | `{"exampleAdditionalClaim": "valueToMatch", "otherAdditionalClaim": "valueWithWildcard*"}` |

### Reponses

#### 200 - Successful operation

Retourne un objet `Key` (voir schema complet ci-dessus).

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 6 : Get policy file

- **Chemin** : `/tailnet/{tailnet}/acl`
- **Methode HTTP** : `GET`
- **Operation ID** : `getPolicyFile`
- **Tags** : `PolicyFile`

### Summary / Description

**Summary** : Get policy file

**Description** :
Recupere le fichier de politique (policy file) actuel pour le tailnet donne ; cela inclut l'ACL ainsi que les regles et les tests qui ont ete definis.

Cette methode peut retourner le policy file en JSON ou HuJSON, selon l'en-tete Accept.
La reponse inclut aussi un en-tete `ETag`, qui peut etre optionnellement inclus lors du setting du policy file pour eviter les mises a jour manquees.

En savoir plus sur la syntaxe ACL du policy file : https://tailscale.com/kb/1337/acl-syntax

OAuth Scope : `policy_file:read`

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `Accept` | header | string | non | `application/json` pour JSON, sinon HuJSON sera retourne |
| `details` | query | boolean | non | Demander une description detaillee du policy file en fournissant `details=true`. Toute autre valeur ou absence est traitee comme `details=false`. Si utilise, ne pas fournir de parametre `Accept` dans le header. |

### Reponses

#### 200 - Successful operation

**Format JSON (details=false)** :

Le policy file complet en JSON :

```json
{
  "acls": [
    {
      "action": "accept",
      "ports": ["*:*"],
      "users": ["*"]
    }
  ],
  "groups": {
    "group:example": ["user1@example.com", "user2@example.com"]
  },
  "hosts": {
    "example-host-1": "100.100.100.100"
  }
}
```

**Format JSON (details=true)** :

| Champ | Type | Description |
|-------|------|-------------|
| `acl` | string | Representation en base64 du format huJSON |
| `warnings` | array of string | Entrees syntaxiquement valides mais nonsensiques |
| `errors` | array of string \| null | Echecs de parsing |

Exemple :
```json
{
  "acl": "Ly8gUG9raW5nIGFyb3VuZC...",
  "warnings": [
    "\"group:example\": user not found: \"user1@example.com\"",
    "\"group:example\": user not found: \"user2@example.com\""
  ],
  "errors": null
}
```

**Format HuJSON** (`Accept: application/hujson`) :

Retourne le policy file en texte brut au format HuJSON (JSON avec commentaires).

#### 400 - Bad Request

Reponse d'erreur standard.

#### 403 - Forbidden

Reponse d'erreur standard.

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 7 : Set policy file

- **Chemin** : `/tailnet/{tailnet}/acl`
- **Methode HTTP** : `POST`
- **Operation ID** : `setPolicyFile`
- **Tags** : `PolicyFile`

### Summary / Description

**Summary** : Set policy file

**Description** :
Definit l'ACL pour le tailnet donne. Les formats HuJSON et JSON sont tous deux acceptes en entree.
Un en-tete `If-Match` peut etre defini pour eviter les mises a jour manquees.

En cas de succes, retourne l'ACL mis a jour en JSON ou HuJSON selon l'en-tete `Accept`.
Sinon, des erreurs sont retournees pour les ACLs incorrectement definis, les ACLs avec des tests en echec lors des tentatives de mise a jour, et les discordances entre l'en-tete `If-Match` et l'`ETag`.

En savoir plus sur la syntaxe ACL du policy file : https://tailscale.com/kb/1337/acl-syntax

OAuth Scope : `policy_file`

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `Accept` | header | string | non | `application/json` pour JSON, sinon HuJSON sera retourne |
| `If-Match` | header | string | non | Mecanisme de securite pour eviter d'ecraser les mises a jour d'autres utilisateurs. Definir la valeur a celle de l'en-tete `ETag` retourne par un GET sur `/api/v2/tailnet/{tailnet}/acl`. Tailscale compare la valeur `ETag` de la requete avec celle du fichier courant et ne remplace le fichier que s'il y a correspondance. Alternativement, definir a `ts-default` pour s'assurer que le policy file est remplace uniquement si le policy file actuel est encore celui par defaut. Exemples : `-H "If-Match: \"e0b2816b418\""` ou `-H "If-Match: \"ts-default\""` |

### Corps de la requete (Request Body)

Content-Types acceptes : `application/json`, `application/hujson`

Le body contient le policy file complet au format JSON ou HuJSON.

**Exemple JSON** :
```json
{
  "acls": [
    {
      "action": "accept",
      "ports": ["*:*"],
      "users": ["*"]
    }
  ],
  "groups": {
    "group:example": ["user1@example.com", "user2@example.com"]
  },
  "hosts": {
    "example-host-1": "100.100.100.100"
  }
}
```

**Exemple HuJSON** :
```
// Example/default ACLs for unrestricted connections.
{
  // Declare static groups of users beyond those in the identity service.
  "groups": {
    "group:example": ["user1@example.com", "user2@example.com"],
  },

  // Declare convenient hostname aliases to use in place of IP addresses.
  "hosts": {
    "example-host-1": "100.100.100.100",
  },

  // Access control lists.
  "acls": [
    // Match absolutely everything.
    {"action": "accept", "src": ["*"], "dst": ["*:*"]},
  ],
}
```

### Reponses

#### 200 - Successful operation

Retourne l'ACL mis a jour en JSON ou HuJSON selon l'en-tete `Accept` (meme structure que le body de la requete).

#### 400 - ACL validation or test error

Reponse d'erreur standard.

#### 403 - Forbidden

Reponse d'erreur standard.

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 412 - If-Match hash mismatch

```json
{
  "message": "precondition failed, invalid old hash"
}
```

Schema : objet `Error` avec champ `message` (string).

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 8 : Preview rule matches

- **Chemin** : `/tailnet/{tailnet}/acl/preview`
- **Methode HTTP** : `POST`
- **Operation ID** : `previewRuleMatches`
- **Tags** : `PolicyFile`

### Summary / Description

**Summary** : Preview rule matches

**Description** :
Lorsqu'on lui donne un utilisateur ou un IP port a matcher, retourne les regles de politique du tailnet qui s'appliquent a cette ressource, sans sauvegarder le policy file sur le serveur.

OAuth Scope : `policy_file:read`

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |
| `type` | query | string, enum: `user`, `ipport` | oui | Specifier pour quel type de ressource (utilisateur ou IP port) les regles correspondantes doivent etre recuperees. `user` : si la valeur `previewFor` est l'email d'un utilisateur (note : `user` reste dans l'API pour la compatibilite mais a ete remplace par `src` dans les policy files). `ipport` : si la valeur `previewFor` est une adresse IP et un port (note : `ipport` reste dans l'API pour la compatibilite mais a ete remplace par `dst` dans les policy files). Exemple : `user` |
| `previewFor` | query | string | oui | Si `type` est `user`, fournir l'email d'un utilisateur valide avec des machines enregistrees. Si `type` est `ipport`, fournir une adresse IP + port : `10.0.0.1:80`. Le policy file fourni est interroge avec ce parametre pour determiner quelles regles correspondent. Exemple : `10.0.0.1:80` |

### Corps de la requete (Request Body)

Content-Types acceptes : `application/json`, `application/hujson`

Le body contient le policy file hypothetique complet au format JSON ou HuJSON (meme format que pour Set policy file).

### Reponses

#### 200 - The list of rules that apply to the resource

```json
{
  "matches": [
    {
      "users": ["*"],
      "ports": ["*.*"],
      "lineNumber": 19
    }
  ],
  "type": "user",
  "previewFor": "user1@example.com"
}
```

Structure de la reponse :

| Champ | Type | Requis | Description | Exemple |
|-------|------|--------|-------------|---------|
| `matches` | array of object | oui | Liste des regles qui s'appliquent a la ressource | voir ci-dessous |
| `matches[].users` | array of string | oui | Entites sources affectees par la regle | `["*"]` |
| `matches[].ports` | array of string | oui | Destinations accessibles | `["*.*"]` |
| `matches[].lineNumber` | integer | oui | Emplacement de la regle dans le policy file | `19` |
| `type` | string | oui | Echo du `type` fourni dans la requete | `"user"` |
| `previewFor` | string | oui | Echo du `previewFor` fourni dans la requete | `"user1@example.com"` |

#### 400 - Bad Request

Reponse d'erreur standard.

#### 403 - Forbidden

Reponse d'erreur standard.

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## ENDPOINT 9 : Validate and test policy file

- **Chemin** : `/tailnet/{tailnet}/acl/validate`
- **Methode HTTP** : `POST`
- **Operation ID** : `validateAndTestPolicyFile`
- **Tags** : `PolicyFile`

### Summary / Description

**Summary** : Validate and test policy file

**Description** :
Cet endpoint fonctionne dans l'un de deux modes, aucun des deux ne modifie le policy file actuel du tailnet :

1. **Executer des tests ACL** : Lorsque le corps de la requete contient des tests ACL sous forme de tableau JSON, Tailscale execute les tests ACL contre le policy file actuel du tailnet. En savoir plus sur les tests ACL : https://tailscale.com/kb/1337/acl-syntax#tests

2. **Valider un nouveau policy file** : Lorsque le corps de la requete est un objet JSON, Tailscale interprete le body comme un policy file hypothetique avec de nouveaux ACLs, incluant toute nouvelle regle et tout nouveau test. Il valide que le policy file est parsable et execute les tests pour valider les regles existantes.

Dans les deux cas, cette methode ne modifie le policy file du tailnet en aucune maniere.

OAuth Scope : `policy_file:read`

### Parametres

| Nom | In | Type | Requis | Description |
|-----|-----|------|--------|-------------|
| `tailnet` | path | string | oui | L'ID du tailnet (ou `-` pour le defaut) |

### Corps de la requete (Request Body)

Content-Types acceptes : `application/json`, `application/hujson`

Le body peut etre soit un **tableau JSON** (mode tests ACL), soit un **objet JSON/string HuJSON** (mode validation de policy file).

#### Mode 1 : Tests ACL (tableau JSON)

Schema : `array` d'objets de test.

Chaque objet de test :

| Champ | Type | Requis | Description | Exemple |
|-------|------|--------|-------------|---------|
| `src` | string | oui | Identite utilisateur a tester : email d'utilisateur, groupe, tag, ou host qui mappe vers une adresse IP. Le test s'execute du point de vue d'un device authentifie avec l'identite fournie. | `"dave@example.com"` |
| `srcPostureAttrs` | object (additionalProperties: string \| number \| boolean) | non | Attributs de posture du device sous forme de paires cle-valeur, utilises pour evaluer les conditions de posture dans les regles d'acces. Necessaire uniquement si les regles d'acces contiennent des conditions de posture de device. | `{"node:os": "windows"}` |
| `proto` | string | non | Protocole IP pour les regles `accept` et `deny`, similaire au champ `proto` dans les regles ACL. Si omis, le test verifie l'acces TCP ou UDP. | `"tcp"` |
| `accept` | array of string | non | Destinations a accepter. Chaque destination est de la forme `host:port` ou `port` est un port numerique unique et `host` est au format decrit dans la documentation de syntaxe ACL. Les sources dans `src` et les destinations ne supportent pas les wildcards `*`. | `["example-host-1:22"]` |
| `deny` | array of string | non | Destinations a refuser. Meme format que `accept`. Les sources dans `src` et les destinations ne supportent pas les wildcards `*`. | `["example-host-2:100"]` |

**Exemple** :
```json
[
  {
    "src": "user1@example.com",
    "accept": ["example-host-1:22"],
    "deny": ["example-host-2:100"]
  }
]
```

#### Mode 2 : Validation de policy file (objet JSON ou HuJSON)

Schema JSON : `string` (representation JSON du policy file)

Schema HuJSON : `string` (representation HuJSON du policy file)

**Exemple JSON** :
```json
{
  "acls": [
    { "action": "accept", "src": ["100.105.106.107"], "dst": ["1.2.3.4:*"] }
  ],
  "tests": [
    {"src": "100.105.106.107", "allow": ["1.2.3.4:80"]}
  ]
}
```

**Exemple HuJSON** :
```
// Example/default ACLs for unrestricted connections.
{
  // Declare static groups of users beyond those in the identity service.
  "groups": {
    "group:example": ["user1@example.com", "user2@example.com"]
  },
  // Declare convenient hostname aliases to use in place of IP addresses.
  "hosts": {
    "example-host-1": "100.100.100.100"
  },
  // Access control lists.
  "acls": [
    // Match absolutely everything. Comment out this section if you want
    // to define specific ACL restrictions.
    { "action": "accept", "users": ["*"], "ports": ["*:*"] }
  ]
}
```

### Reponses

#### 200 - Validation or tests have run

Un corps de reponse vide (`{}`) implique que la validation ou les tests ont reussi.

Structure de la reponse (en cas d'echec ou d'avertissement) :

| Champ | Type | Description | Exemple |
|-------|------|-------------|---------|
| `message` | string | Message decrivant le resultat | `"test(s) failed"` ou `"warning(s) found"` |
| `data` | array of object | Details des erreurs ou avertissements | voir ci-dessous |
| `data[].user` | string | L'utilisateur concerne | `"user1@example.com"` |
| `data[].errors` | array of string | Liste des erreurs | `["address \"2.2.2.2:22\": want: Drop, got: Accept"]` |
| `data[].warnings` | array of string | Liste des avertissements | `["group is not syncing from SCIM and will be ignored by rules in the policy file"]` |

**Exemple - Tests echoues** :
```json
{
  "message": "test(s) failed",
  "data": [
    {
      "user": "user1@example.com",
      "errors": [
        "address \"2.2.2.2:22\": want: Drop, got: Accept"
      ]
    }
  ]
}
```

**Exemple - Avertissements SCIM** :
```json
{
  "message": "warning(s) found",
  "data": [
    {
      "user": "group:unknown@example.com",
      "warnings": [
        "group is not syncing from SCIM and will be ignored by rules in the policy file"
      ]
    }
  ]
}
```

**Exemple - Succes** :
```json
{}
```

#### 400 - Bad Request

Reponse d'erreur standard.

#### 403 - Forbidden

Reponse d'erreur standard.

#### 404 - Tailnet not found

Reponse d'erreur standard.

#### 500 - Internal Server Error

Reponse d'erreur standard.

---

## Resume des endpoints

| # | Methode | Chemin | Operation ID | Tags | Description |
|---|---------|--------|-------------|------|-------------|
| 1 | `GET` | `/tailnet/{tailnet}/keys` | `listTailnetKeys` | Keys | Lister les cles du tailnet |
| 2 | `POST` | `/tailnet/{tailnet}/keys` | `createKey` | Keys | Creer un auth key ou trust credential |
| 3 | `GET` | `/tailnet/{tailnet}/keys/{keyId}` | `getKey` | Keys | Obtenir les details d'une cle |
| 4 | `DELETE` | `/tailnet/{tailnet}/keys/{keyId}` | `deleteKey` | Keys | Supprimer une cle |
| 5 | `PUT` | `/tailnet/{tailnet}/keys/{keyId}` | `setKey` | Keys | Definir la configuration d'un OAuth client ou federated identity |
| 6 | `GET` | `/tailnet/{tailnet}/acl` | `getPolicyFile` | PolicyFile | Recuperer le policy file |
| 7 | `POST` | `/tailnet/{tailnet}/acl` | `setPolicyFile` | PolicyFile | Definir le policy file |
| 8 | `POST` | `/tailnet/{tailnet}/acl/preview` | `previewRuleMatches` | PolicyFile | Previsualiser les regles correspondantes |
| 9 | `POST` | `/tailnet/{tailnet}/acl/validate` | `validateAndTestPolicyFile` | PolicyFile | Valider et tester un policy file |
