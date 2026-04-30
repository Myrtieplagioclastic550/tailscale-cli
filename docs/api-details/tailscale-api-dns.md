# Tailscale API - Section DNS - Documentation exhaustive des endpoints

---

## Parametre commun a tous les endpoints DNS

### Path parameter: `tailnet`

| Champ       | Valeur                                                                                                                                                                                                                                                       |
|-------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Nom**     | `tailnet`                                                                                                                                                                                                                                                    |
| **In**      | path                                                                                                                                                                                                                                                         |
| **Type**    | `string`                                                                                                                                                                                                                                                     |
| **Requis**  | Oui                                                                                                                                                                                                                                                          |
| **Exemple** | `example.com`                                                                                                                                                                                                                                                |
| **Description** | Le tailnet ID. Les tailnets crees avant octobre 2025 peuvent encore utiliser l'ID legacy. On peut fournir un tiret (`-`) pour referencer le tailnet par defaut du token d'acces, ou fournir le **tailnet ID** visible dans la page General Settings de la console admin Tailscale. |

---

## Reponses d'erreur communes

Toutes les reponses d'erreur suivent le schema `Error` :

```json
{
  "message": "string"   // (requis) Message d'erreur
}
```

| Code | Description              | Exemple de message         |
|------|--------------------------|----------------------------|
| 404  | Tailnet not found.       | `"not found"`              |
| 500  | Internal server error.   | `"internal server error"`  |

---

## Endpoint 1 : List DNS nameservers

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/nameservers`                     |
| **Methode HTTP**| `GET`                                                    |
| **operationId** | `listDnsNameservers`                                     |
| **Tags**        | `DNS`                                                    |
| **Summary**     | List DNS nameservers                                     |
| **Description** | Lists the global DNS nameservers for a tailnet.          |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

Aucun query parameter. Aucun body.

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:**

```json
{
  "dns": ["string"]   // Tableau de strings - DNS nameservers
}
```

| Champ | Type              | Description       | Exemple                  |
|-------|-------------------|--------------------|--------------------------|
| `dns` | `array` of `string` | DNS nameservers. | `["8.8.8.8", "1.2.3.4"]` |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 2 : Set DNS nameservers

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/nameservers`                     |
| **Methode HTTP**| `POST`                                                   |
| **operationId** | `setDnsNameservers`                                      |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Set DNS nameservers                                      |
| **Description** | Replaces the list of global DNS nameservers for the given tailnet with the list supplied in the request. Note that changing the list of DNS nameservers may also affect the status of MagicDNS (if MagicDNS is on; learn about MagicDNS). If all nameservers have been removed, MagicDNS will be automatically disabled (until explicitly turned back on by the user). |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:**

```json
{
  "dns": ["string"]   // Tableau de strings - DNS nameservers
}
```

| Champ | Type              | Description       | Exemple                  |
|-------|-------------------|--------------------|--------------------------|
| `dns` | `array` of `string` | DNS nameservers. | `["8.8.8.8", "1.2.3.4"]` |

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:**

```json
{
  "dns": ["string"],
  "magicDNS": boolean
}
```

| Champ      | Type              | Description                                       | Exemple                  |
|------------|-------------------|---------------------------------------------------|--------------------------|
| `dns`      | `array` of `string` | DNS nameservers.                                 | `["8.8.8.8", "1.2.3.4"]` |
| `magicDNS` | `boolean`         | Whether MagicDNS is active for this tailnet.      | `true`                   |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 3 : Get DNS preferences

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/preferences`                     |
| **Methode HTTP**| `GET`                                                    |
| **operationId** | `getDnsPreferences`                                      |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Get DNS preferences                                      |
| **Description** | Retrieves the DNS preferences that are currently set for the given tailnet. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

Aucun query parameter. Aucun body.

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsPreferences`

```json
{
  "magicDNS": boolean
}
```

| Champ      | Type      | Requis | Description                                  | Exemple |
|------------|-----------|--------|----------------------------------------------|---------|
| `magicDNS` | `boolean` | Oui    | Whether MagicDNS is active for this tailnet. | `true`  |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 4 : Set DNS preferences

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/preferences`                     |
| **Methode HTTP**| `POST`                                                   |
| **operationId** | `setDnsPreferences`                                      |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Set DNS preferences                                      |
| **Description** | Set the DNS preferences for a tailnet; specifically, the MagicDNS setting. Note that MagicDNS is dependent on DNS servers. Learn about MagicDNS. If there is at least one DNS server, then MagicDNS can be enabled. Otherwise, it returns an error. Note that removing all nameservers will turn off MagicDNS. To reenable it, nameservers must be added back, and MagicDNS must be explicitly turned on. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:** `DnsPreferences`

```json
{
  "magicDNS": boolean
}
```

| Champ      | Type      | Requis | Description                                  | Exemple |
|------------|-----------|--------|----------------------------------------------|---------|
| `magicDNS` | `boolean` | Oui    | Whether MagicDNS is active for this tailnet. | `true`  |

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsPreferences`

```json
{
  "magicDNS": boolean
}
```

| Champ      | Type      | Requis | Description                                  | Exemple |
|------------|-----------|--------|----------------------------------------------|---------|
| `magicDNS` | `boolean` | Oui    | Whether MagicDNS is active for this tailnet. | `true`  |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 5 : List DNS search paths

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/searchpaths`                     |
| **Methode HTTP**| `GET`                                                    |
| **operationId** | `listDnsSearchPaths`                                     |
| **Tags**        | `DNS`                                                    |
| **Summary**     | List DNS search paths                                    |
| **Description** | Retrieves the list of search paths, also referred to as *search domains*, that is currently set for the given tailnet. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

Aucun query parameter. Aucun body.

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsSearchPaths`

```json
{
  "searchPaths": ["string"]
}
```

| Champ         | Type              | Requis | Description                               | Exemple                                    |
|---------------|-------------------|--------|-------------------------------------------|--------------------------------------------|
| `searchPaths` | `array` of `string` | Oui    | The search domains for the given tailnet. | `["user1.example.com", "user2.example.com"]` |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 6 : Set DNS search paths

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/searchpaths`                     |
| **Methode HTTP**| `POST`                                                   |
| **operationId** | `setDnsSearchPaths`                                      |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Set DNS search paths                                     |
| **Description** | Replaces the list of search paths for the given tailnet. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:** `DnsSearchPaths`

```json
{
  "searchPaths": ["string"]
}
```

| Champ         | Type              | Requis | Description                               | Exemple                                    |
|---------------|-------------------|--------|-------------------------------------------|--------------------------------------------|
| `searchPaths` | `array` of `string` | Oui    | The search domains for the given tailnet. | `["user1.example.com", "user2.example.com"]` |

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsSearchPaths`

```json
{
  "searchPaths": ["string"]
}
```

| Champ         | Type              | Requis | Description                               | Exemple                                    |
|---------------|-------------------|--------|-------------------------------------------|--------------------------------------------|
| `searchPaths` | `array` of `string` | Oui    | The search domains for the given tailnet. | `["user1.example.com", "user2.example.com"]` |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 7 : Get split DNS

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/split-dns`                       |
| **Methode HTTP**| `GET`                                                    |
| **operationId** | `getSplitDns`                                            |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Get split DNS                                            |
| **Description** | Retrieves the split DNS settings, which is a map from domains to lists of nameservers, that is currently set for the given tailnet. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

Aucun query parameter. Aucun body.

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `SplitDns`

Le schema est un objet avec des proprietes additionnelles (map dynamique). Chaque cle est un nom de domaine, et la valeur est un tableau de strings (adresses de nameservers) ou `null`.

```json
{
  "<domain_name>": ["string"] | null
}
```

**Exemple:**

```json
{
  "example.com": ["1.1.1.1", "1.2.3.4"],
  "other.com": ["2.2.2.2"]
}
```

| Champ                  | Type                           | Description                                           |
|------------------------|--------------------------------|-------------------------------------------------------|
| `<domain_name>` (cle) | `array` of `string` ou `null`  | Map of domain names to lists of nameservers or to `null`. |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 8 : Update split DNS

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/split-dns`                       |
| **Methode HTTP**| `PATCH`                                                  |
| **operationId** | `updateSplitDns`                                         |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Update split DNS                                         |
| **Description** | Performs partial updates of the split DNS settings for a given tailnet. Only domains specified in the request map will be modified. Setting the value of a mapping to `null` clears the nameservers for that domain. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:** `SplitDns`

```json
{
  "<domain_name>": ["string"] | null
}
```

| Champ                  | Type                           | Description                                           |
|------------------------|--------------------------------|-------------------------------------------------------|
| `<domain_name>` (cle) | `array` of `string` ou `null`  | Map of domain names to lists of nameservers or to `null`. Mettre `null` pour effacer les nameservers d'un domaine. |

**Exemple de body:**

```json
{
  "example.com": ["1.1.1.1", "1.2.3.4"],
  "other.com": null
}
```

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `SplitDns` (meme format que le body - map domaines vers nameservers)

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 9 : Set split DNS

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/split-dns`                       |
| **Methode HTTP**| `PUT`                                                    |
| **operationId** | `setSplitDns`                                            |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Set split DNS                                            |
| **Description** | Replaces the split DNS settings for a given tailnet. Setting the value of a mapping to `null` clears the nameservers for that domain. Sending an empty object clears nameservers for all domains. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:** `SplitDns`

```json
{
  "<domain_name>": ["string"] | null
}
```

| Champ                  | Type                           | Description                                           |
|------------------------|--------------------------------|-------------------------------------------------------|
| `<domain_name>` (cle) | `array` of `string` ou `null`  | Map of domain names to lists of nameservers or to `null`. Envoyer `{}` pour effacer tous les domaines. |

**Exemple de body:**

```json
{
  "example.com": ["1.1.1.1", "1.2.3.4"],
  "other.com": ["2.2.2.2"]
}
```

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `SplitDns` (meme format que le body - map domaines vers nameservers)

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 10 : Get DNS configuration

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/configuration`                   |
| **Methode HTTP**| `GET`                                                    |
| **operationId** | `getDnsConfiguration`                                    |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Get DNS configuration                                    |
| **Description** | Retrieves the full DNS configuration for a tailnet, including global nameservers, split DNS routes, search paths, and MagicDNS configuration. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

Aucun query parameter. Aucun body.

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsConfiguration`

```json
{
  "nameservers": [
    {
      "address": "string",
      "useWithExitNode": boolean
    }
  ],
  "splitDNS": {
    "<domain_name>": [
      {
        "address": "string",
        "useWithExitNode": boolean
      }
    ] | null
  },
  "searchPaths": ["string"],
  "preferences": {
    "overrideLocalDNS": boolean,
    "magicDNS": boolean
  }
}
```

#### Champ `nameservers` (array of `DnsConfigurationResolver`)

Global DNS resolvers to use. If `preferences.overrideLocalDNS` is true, these override the local OS configuration; otherwise they are used as fallback resolvers.

| Champ              | Type      | Description                                                                                                              | Exemple    |
|--------------------|-----------|--------------------------------------------------------------------------------------------------------------------------|------------|
| `address`          | `string`  | IPv4 or IPv6 address of the DNS resolver.                                                                                | `"1.1.1.1"` |
| `useWithExitNode`  | `boolean` | If true, this resolver should still be used when a device is configured to use a Tailscale exit node. Requires Tailscale v1.88.1 or later. | `true`     |

**Exemple:**

```json
[
  { "address": "8.8.8.8", "useWithExitNode": true },
  { "address": "1.1.1.1", "useWithExitNode": false }
]
```

#### Champ `splitDNS` (object, additionalProperties)

Map of DNS name suffixes (domains) to lists of resolvers for Split DNS and advanced routing overlays. Chaque valeur est un tableau de `DnsConfigurationResolver` ou `null`.

**Exemple:**

```json
{
  "corp.example.com": [
    { "address": "10.0.0.53", "useWithExitNode": true },
    { "address": "10.0.1.53", "useWithExitNode": true }
  ],
  "other.internal": [
    { "address": "10.0.2.53", "useWithExitNode": false }
  ]
}
```

#### Champ `searchPaths` (array of string)

| Champ         | Type              | Description                  | Exemple                                    |
|---------------|-------------------|------------------------------|--------------------------------------------|
| `searchPaths` | `array` of `string` | Search domain paths to apply. | `["user1.example.com", "user2.example.com"]` |

#### Champ `preferences` (schema `DnsConfigurationPreferences`)

| Champ              | Type      | Description                                                                                                     | Exemple |
|--------------------|-----------|-----------------------------------------------------------------------------------------------------------------|---------|
| `overrideLocalDNS` | `boolean` | If true, resolvers in `nameservers` override the local OS DNS configuration; if false, local resolvers are used. | `true`  |
| `magicDNS`         | `boolean` | Whether MagicDNS is enabled for this tailnet.                                                                    | `true`  |

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Endpoint 11 : Set DNS configuration

| Champ           | Valeur                                                   |
|-----------------|----------------------------------------------------------|
| **Path**        | `/tailnet/{tailnet}/dns/configuration`                   |
| **Methode HTTP**| `POST`                                                   |
| **operationId** | `setDnsConfiguration`                                    |
| **Tags**        | `DNS`                                                    |
| **Summary**     | Set DNS configuration                                    |
| **Description** | Replaces the DNS configuration for the given tailnet. `nameservers` defines the global resolvers to use when `preferences.overrideLocalDNS` is true. `splitDNS` maps DNS name suffixes (domains) to lists of resolvers for Split DNS. `searchPaths` sets custom DNS search domain paths. `preferences.overrideLocalDNS` controls whether resolvers in `nameservers` override the local OS configuration (true) or local resolvers are used (false). Defaults to false. `preferences.magicDNS` enables MagicDNS. Defaults to false. |

### Parametres

| Nom       | In   | Type     | Requis | Description       |
|-----------|------|----------|--------|-------------------|
| `tailnet` | path | `string` | Oui    | Le tailnet ID.    |

### Request Body

**Content-Type:** `application/json`

**Schema:** `DnsConfiguration`

```json
{
  "nameservers": [
    {
      "address": "string",
      "useWithExitNode": boolean
    }
  ],
  "splitDNS": {
    "<domain_name>": [
      {
        "address": "string",
        "useWithExitNode": boolean
      }
    ] | null
  },
  "searchPaths": ["string"],
  "preferences": {
    "overrideLocalDNS": boolean,
    "magicDNS": boolean
  }
}
```

#### Champ `nameservers` (array of `DnsConfigurationResolver`)

Global DNS resolvers to use. If `preferences.overrideLocalDNS` is true, these override the local OS configuration; otherwise they are used as fallback resolvers.

| Champ              | Type      | Description                                                                                                              | Exemple    |
|--------------------|-----------|--------------------------------------------------------------------------------------------------------------------------|------------|
| `address`          | `string`  | IPv4 or IPv6 address of the DNS resolver.                                                                                | `"1.1.1.1"` |
| `useWithExitNode`  | `boolean` | If true, this resolver should still be used when a device is configured to use a Tailscale exit node. Requires Tailscale v1.88.1 or later. | `true`     |

#### Champ `splitDNS` (object, additionalProperties)

Map of DNS name suffixes (domains) to lists of resolvers for Split DNS and advanced routing overlays. Chaque valeur est un tableau de `DnsConfigurationResolver` ou `null`.

#### Champ `searchPaths` (array of string)

| Champ         | Type              | Description                  | Exemple                                    |
|---------------|-------------------|------------------------------|--------------------------------------------|
| `searchPaths` | `array` of `string` | Search domain paths to apply. | `["user1.example.com", "user2.example.com"]` |

#### Champ `preferences` (schema `DnsConfigurationPreferences`)

| Champ              | Type      | Description                                                                                                     | Exemple |
|--------------------|-----------|-----------------------------------------------------------------------------------------------------------------|---------|
| `overrideLocalDNS` | `boolean` | If true, resolvers in `nameservers` override the local OS DNS configuration; if false, local resolvers are used. | `true`  |
| `magicDNS`         | `boolean` | Whether MagicDNS is enabled for this tailnet.                                                                    | `true`  |

### Reponse 200 - Successful operation

**Content-Type:** `application/json`

**Schema:** `DnsConfiguration` (meme structure que le request body ci-dessus)

### Reponses d'erreur

| Code | Description        |
|------|--------------------|
| 404  | Tailnet not found. |
| 500  | Internal server error. |

---

## Recapitulatif des 11 endpoints DNS

| #  | Methode | Path                                          | operationId            | Summary                  |
|----|---------|-----------------------------------------------|------------------------|--------------------------|
| 1  | GET     | `/tailnet/{tailnet}/dns/nameservers`          | `listDnsNameservers`   | List DNS nameservers     |
| 2  | POST    | `/tailnet/{tailnet}/dns/nameservers`          | `setDnsNameservers`    | Set DNS nameservers      |
| 3  | GET     | `/tailnet/{tailnet}/dns/preferences`          | `getDnsPreferences`    | Get DNS preferences      |
| 4  | POST    | `/tailnet/{tailnet}/dns/preferences`          | `setDnsPreferences`    | Set DNS preferences      |
| 5  | GET     | `/tailnet/{tailnet}/dns/searchpaths`          | `listDnsSearchPaths`   | List DNS search paths    |
| 6  | POST    | `/tailnet/{tailnet}/dns/searchpaths`          | `setDnsSearchPaths`    | Set DNS search paths     |
| 7  | GET     | `/tailnet/{tailnet}/dns/split-dns`            | `getSplitDns`          | Get split DNS            |
| 8  | PATCH   | `/tailnet/{tailnet}/dns/split-dns`            | `updateSplitDns`       | Update split DNS         |
| 9  | PUT     | `/tailnet/{tailnet}/dns/split-dns`            | `setSplitDns`          | Set split DNS            |
| 10 | GET     | `/tailnet/{tailnet}/dns/configuration`        | `getDnsConfiguration`  | Get DNS configuration    |
| 11 | POST    | `/tailnet/{tailnet}/dns/configuration`        | `setDnsConfiguration`  | Set DNS configuration    |

---

## Schemas references

### `DnsPreferences`

```yaml
type: object
properties:
  magicDNS:
    type: boolean
    example: true
    description: Whether MagicDNS is active for this tailnet.
required:
  - magicDNS
```

### `DnsSearchPaths`

```yaml
type: object
properties:
  searchPaths:
    type: array
    items:
      type: string
    example: ["user1.example.com", "user2.example.com"]
    description: The search domains for the given tailnet.
required:
  - searchPaths
```

### `SplitDns`

```yaml
type: object
additionalProperties:
  x-additionalPropertiesName: Domain names to DNS
  type: [array, 'null']
  items:
    type: string
example:
  example.com: ["1.1.1.1", "1.2.3.4"]
  other.com: ["2.2.2.2"]
description: Map of domain names to lists of nameservers or to null.
```

### `DnsConfigurationResolver`

```yaml
type: object
properties:
  address:
    type: string
    description: IPv4 or IPv6 address of the DNS resolver.
    example: "1.1.1.1"
  useWithExitNode:
    type: boolean
    description: If true, this resolver should still be used when a device is configured to use a Tailscale exit node. Requires Tailscale v1.88.1 or later.
    example: true
```

### `DnsConfigurationPreferences`

```yaml
type: object
properties:
  overrideLocalDNS:
    type: boolean
    description: If true, resolvers in nameservers override the local OS DNS configuration; if false, local resolvers are used.
    example: true
  magicDNS:
    type: boolean
    description: Whether MagicDNS is enabled for this tailnet.
    example: true
```

### `DnsConfiguration`

```yaml
type: object
properties:
  nameservers:
    type: array
    description: Global DNS resolvers to use.
    items:
      $ref: DnsConfigurationResolver
    example:
      - address: "8.8.8.8"
        useWithExitNode: true
      - address: "1.1.1.1"
        useWithExitNode: false
  splitDNS:
    type: object
    description: Map of DNS name suffixes (domains) to lists of resolvers for Split DNS.
    additionalProperties:
      type: [array, 'null']
      items:
        $ref: DnsConfigurationResolver
    example:
      corp.example.com:
        - address: "10.0.0.53"
          useWithExitNode: true
        - address: "10.0.1.53"
          useWithExitNode: true
      other.internal:
        - address: "10.0.2.53"
          useWithExitNode: false
  searchPaths:
    type: array
    items:
      type: string
    description: Search domain paths to apply.
    example: ["user1.example.com", "user2.example.com"]
  preferences:
    $ref: DnsConfigurationPreferences
```

### `Error`

```yaml
type: object
properties:
  message:
    type: string
required:
  - message
```
