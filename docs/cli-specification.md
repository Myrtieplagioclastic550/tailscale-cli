# Telscale CLI - Specification complete

## Vue d'ensemble

**Nom** : `telscale`
**Objectif** : CLI pour automatiser toutes les operations de l'API Tailscale v2, utilisable en standalone ou comme outil MCP dans Claude Code.
**Pattern** : `telscale [FLAGS_GLOBAUX] RESSOURCE ACTION [FLAGS_LOCAUX] [ARGS]`
**Inspirations** : `gh`, `doctl`, `kubectl`

---

## Flags globaux

| Flag | Short | Env var | Description | Default |
|------|-------|---------|-------------|---------|
| `--api-token` | `-t` | `TELSCALE_API_TOKEN` | Token API (`tskey-api-xxx`) ou OAuth (`tskey-client-xxx`) | config file |
| `--tailnet` | `-n` | `TELSCALE_TAILNET` | Tailnet ID (ou `-` pour le defaut du token) | `-` |
| `--output` | `-o` | `TELSCALE_OUTPUT` | Format de sortie : `table`, `json`, `yaml`, `csv` | `table` |
| `--json` | | | Raccourci pour `--output json` | |
| `--quiet` | `-q` | | Supprime la sortie (utile en scripts, retourne uniquement le code de sortie) | `false` |
| `--debug` | | `TELSCALE_DEBUG` | Affiche les requetes HTTP et logs de debug | `false` |
| `--no-color` | | `NO_COLOR` | Desactive la coloration | `false` |
| `--config` | | `TELSCALE_CONFIG` | Chemin du fichier de configuration | `~/.telscale/config.json` |
| `--context` | `-c` | `TELSCALE_CONTEXT` | Contexte nomme (multi-comptes) | `default` |

---

## Configuration (`~/.telscale/config.json`)

```json
{
  "default_context": "work",
  "contexts": {
    "work": {
      "api_token": "tskey-api-xxxxx",
      "tailnet": "mycompany.com"
    },
    "personal": {
      "api_token": "tskey-api-yyyyy",
      "tailnet": "-"
    }
  }
}
```

---

## Commandes d'administration

### `telscale auth`

| Commande | Description |
|----------|-------------|
| `telscale auth login` | Configuration interactive (token + tailnet) |
| `telscale auth status` | Affiche le contexte actif et verifie le token |
| `telscale auth switch <context>` | Change le contexte actif |
| `telscale auth list` | Liste les contextes configures |
| `telscale auth remove <context>` | Supprime un contexte |

### `telscale completion`

| Commande | Description |
|----------|-------------|
| `telscale completion bash` | Genere l'autocompletion bash |
| `telscale completion zsh` | Genere l'autocompletion zsh |
| `telscale completion fish` | Genere l'autocompletion fish |
| `telscale completion powershell` | Genere l'autocompletion PowerShell |

### `telscale version`

Affiche la version de la CLI.

---

## Commandes par ressource

---

### 1. `telscale device` — Gestion des devices (15 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale device list` | `GET /tailnet/{tailnet}/devices` | Lister les devices |
| `telscale device get <deviceId>` | `GET /device/{deviceId}` | Obtenir un device |
| `telscale device delete <deviceId>` | `DELETE /device/{deviceId}` | Supprimer un device |
| `telscale device expire <deviceId>` | `POST /device/{deviceId}/expire` | Expirer la cle d'un device |
| `telscale device authorize <deviceId>` | `POST /device/{deviceId}/authorized` | Autoriser un device |
| `telscale device deauthorize <deviceId>` | `POST /device/{deviceId}/authorized` | Desautoriser un device |
| `telscale device set-name <deviceId> <name>` | `POST /device/{deviceId}/name` | Renommer un device |
| `telscale device set-tags <deviceId>` | `POST /device/{deviceId}/tags` | Definir les tags |
| `telscale device set-key <deviceId>` | `POST /device/{deviceId}/key` | Activer/desactiver l'expiration de cle |
| `telscale device set-ip <deviceId> <ipv4>` | `POST /device/{deviceId}/ip` | Definir l'adresse IPv4 |

#### `telscale device routes` — Sous-commande routes

| Commande | API | Description |
|----------|-----|-------------|
| `telscale device routes list <deviceId>` | `GET /device/{deviceId}/routes` | Lister les routes |
| `telscale device routes set <deviceId>` | `POST /device/{deviceId}/routes` | Definir les routes activees |

#### `telscale device posture` — Sous-commande posture

| Commande | API | Description |
|----------|-----|-------------|
| `telscale device posture get <deviceId>` | `GET /device/{deviceId}/attributes` | Obtenir les attributs posture |
| `telscale device posture set <deviceId> <key>` | `POST /device/{deviceId}/attributes/{key}` | Definir un attribut posture |
| `telscale device posture delete <deviceId> <key>` | `DELETE /device/{deviceId}/attributes/{key}` | Supprimer un attribut posture |
| `telscale device posture batch-update` | `PATCH /tailnet/{tailnet}/device-attributes` | Mise a jour batch |

#### Flags specifiques a `device`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--fields` | `all` ou `default` | `list`, `get` |
| `--filter` | Filtre serveur `key=value` (ex: `--filter isEphemeral=true`) | `list` |
| `--tags` | Tags a definir (ex: `--tags tag:prod,tag:server`) | `set-tags` |
| `--routes` | Routes a activer (ex: `--routes 10.0.0.0/16,192.168.1.0/24`) | `routes set` |
| `--key-expiry-disabled` | `true`/`false` | `set-key` |
| `--value` | Valeur de l'attribut posture | `posture set` |
| `--expiry` | Date d'expiration (RFC 3339) | `posture set` |
| `--comment` | Commentaire pour l'audit log | `posture set`, `posture batch-update` |
| `--confirm` | Bypass la confirmation pour les actions destructives | `delete`, `expire` |

---

### 2. `telscale acl` — Policy File / ACL (4 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale acl get` | `GET /tailnet/{tailnet}/acl` | Recuperer le policy file |
| `telscale acl set` | `POST /tailnet/{tailnet}/acl` | Definir le policy file |
| `telscale acl preview` | `POST /tailnet/{tailnet}/acl/preview` | Previsualiser les regles |
| `telscale acl validate` | `POST /tailnet/{tailnet}/acl/validate` | Valider et tester |

#### Flags specifiques a `acl`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--file` | Chemin du fichier ACL (JSON ou HuJSON) | `set`, `preview`, `validate` |
| `--format` | `json` ou `hujson` (en-tete Accept) | `get`, `set` |
| `--details` | Retourner les details (warnings/errors) | `get` |
| `--if-match` | ETag pour mise a jour optimiste | `set` |
| `--type` | `user` ou `ipport` | `preview` |
| `--preview-for` | Utilisateur ou IP:port a tester | `preview` |
| `--stdin` | Lire le policy file depuis stdin | `set`, `validate` |

---

### 3. `telscale key` — Auth keys et trust credentials (5 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale key list` | `GET /tailnet/{tailnet}/keys` | Lister les cles |
| `telscale key create` | `POST /tailnet/{tailnet}/keys` | Creer une cle |
| `telscale key get <keyId>` | `GET /tailnet/{tailnet}/keys/{keyId}` | Obtenir une cle |
| `telscale key delete <keyId>` | `DELETE /tailnet/{tailnet}/keys/{keyId}` | Supprimer une cle |
| `telscale key update <keyId>` | `PUT /tailnet/{tailnet}/keys/{keyId}` | Configurer un OAuth/federated |

#### Flags specifiques a `key`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--all` | Retourner toutes les cles (pas seulement les siennes) | `list` |
| `--type` | Type de cle : `auth`, `client`, `federated` | `create` |
| `--description` | Description (max 50 chars) | `create`, `update` |
| `--expiry` | Duree d'expiration en secondes | `create` |
| `--reusable` | Cle reutilisable | `create` |
| `--ephemeral` | Cle ephemere | `create` |
| `--preauthorized` | Cle pre-autorisee | `create` |
| `--tags` | Tags associes | `create`, `update` |
| `--scopes` | Scopes OAuth | `create`, `update` |
| `--issuer` | OIDC issuer (federated) | `create`, `update` |
| `--subject` | OIDC subject pattern (federated) | `create`, `update` |
| `--audience` | OIDC audience (federated) | `create`, `update` |
| `--confirm` | Bypass confirmation | `delete` |

---

### 4. `telscale dns` — Configuration DNS (11 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale dns nameservers list` | `GET .../dns/nameservers` | Lister les nameservers |
| `telscale dns nameservers set` | `POST .../dns/nameservers` | Definir les nameservers |
| `telscale dns preferences get` | `GET .../dns/preferences` | Obtenir les preferences |
| `telscale dns preferences set` | `POST .../dns/preferences` | Definir les preferences |
| `telscale dns searchpaths list` | `GET .../dns/searchpaths` | Lister les search paths |
| `telscale dns searchpaths set` | `POST .../dns/searchpaths` | Definir les search paths |
| `telscale dns split get` | `GET .../dns/split-dns` | Obtenir le split DNS |
| `telscale dns split update` | `PATCH .../dns/split-dns` | Mise a jour partielle split DNS |
| `telscale dns split set` | `PUT .../dns/split-dns` | Remplacer le split DNS |
| `telscale dns config get` | `GET .../dns/configuration` | Obtenir la config DNS complete |
| `telscale dns config set` | `POST .../dns/configuration` | Definir la config DNS complete |

#### Flags specifiques a `dns`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--nameservers` | Liste de nameservers (ex: `--nameservers 8.8.8.8,1.1.1.1`) | `nameservers set` |
| `--magic-dns` | Activer/desactiver MagicDNS (`true`/`false`) | `preferences set` |
| `--search-paths` | Domaines de recherche | `searchpaths set` |
| `--domain` | Nom de domaine pour split DNS | `split update`, `split set` |
| `--servers` | Serveurs DNS pour un domaine | `split update`, `split set` |
| `--file` | Fichier JSON avec la config complete | `config set`, `split set` |
| `--stdin` | Lire depuis stdin | `config set`, `split set` |

---

### 5. `telscale log` — Logging et streaming (8 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale log audit list` | `GET .../logging/configuration` | Lister les logs d'audit |
| `telscale log network list` | `GET .../logging/network` | Lister les logs reseau |
| `telscale log stream status <logType>` | `GET .../logging/{logType}/stream/status` | Statut du streaming |
| `telscale log stream get <logType>` | `GET .../logging/{logType}/stream` | Config du streaming |
| `telscale log stream set <logType>` | `PUT .../logging/{logType}/stream` | Configurer le streaming |
| `telscale log stream disable <logType>` | `DELETE .../logging/{logType}/stream` | Desactiver le streaming |
| `telscale log aws-id create` | `POST .../aws-external-id` | Creer/obtenir un external ID AWS |
| `telscale log aws-id validate <id>` | `POST .../aws-external-id/{id}/validate-aws-trust-policy` | Valider l'integration AWS |

#### Flags specifiques a `log`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--start` | Debut de la fenetre temporelle (RFC 3339) | `audit list`, `network list` |
| `--end` | Fin de la fenetre temporelle (RFC 3339) | `audit list`, `network list` |
| `--actor` | Filtrer par acteur (ID ou `~search`) | `audit list` |
| `--target` | Filtrer par cible | `audit list` |
| `--event` | Filtrer par type d'evenement | `audit list` |
| `--destination-type` | Type de destination streaming (splunk/elastic/s3/...) | `stream set` |
| `--url` | URL du endpoint de streaming | `stream set` |
| `--user` | Username pour l'auth streaming | `stream set` |
| `--token` | Token pour l'auth streaming | `stream set` |
| `--s3-bucket` | Bucket S3 | `stream set` |
| `--s3-region` | Region S3 | `stream set` |
| `--s3-role-arn` | ARN du role IAM | `stream set` |
| `--compression` | Format de compression (zstd/gzip/none) | `stream set` |
| `--role-arn` | ARN du role AWS a valider | `aws-id validate` |
| `--reusable` | External ID reutilisable | `aws-id create` |
| `--file` | Fichier JSON avec la config streaming | `stream set` |

---

### 6. `telscale user` — Gestion des utilisateurs (7 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale user list` | `GET /tailnet/{tailnet}/users` | Lister les utilisateurs |
| `telscale user get <userId>` | `GET /users/{userId}` | Obtenir un utilisateur |
| `telscale user set-role <userId> <role>` | `POST /users/{userId}/role` | Modifier le role |
| `telscale user approve <userId>` | `POST /users/{userId}/approve` | Approuver un utilisateur |
| `telscale user suspend <userId>` | `POST /users/{userId}/suspend` | Suspendre un utilisateur |
| `telscale user restore <userId>` | `POST /users/{userId}/restore` | Restaurer un utilisateur |
| `telscale user delete <userId>` | `POST /users/{userId}/delete` | Supprimer un utilisateur |

#### Flags specifiques a `user`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--type` | Filtre : `member`, `shared`, `all` | `list` |
| `--role` | Filtre par role (ou nouveau role pour `set-role`) | `list`, `set-role` |
| `--confirm` | Bypass confirmation | `delete`, `suspend` |

---

### 7. `telscale invite` — Invitations (11 endpoints)

#### `telscale invite user` — Invitations utilisateur

| Commande | API | Description |
|----------|-----|-------------|
| `telscale invite user list` | `GET /tailnet/{tailnet}/user-invites` | Lister les invitations |
| `telscale invite user create` | `POST /tailnet/{tailnet}/user-invites` | Creer des invitations |
| `telscale invite user get <id>` | `GET /user-invites/{id}` | Obtenir une invitation |
| `telscale invite user delete <id>` | `DELETE /user-invites/{id}` | Supprimer une invitation |
| `telscale invite user resend <id>` | `POST /user-invites/{id}/resend` | Renvoyer l'invitation |

#### `telscale invite device` — Invitations device

| Commande | API | Description |
|----------|-----|-------------|
| `telscale invite device list <deviceId>` | `GET /device/{deviceId}/device-invites` | Lister les invites |
| `telscale invite device create <deviceId>` | `POST /device/{deviceId}/device-invites` | Creer des invites |
| `telscale invite device get <id>` | `GET /device-invites/{id}` | Obtenir une invite |
| `telscale invite device delete <id>` | `DELETE /device-invites/{id}` | Supprimer une invite |
| `telscale invite device resend <id>` | `POST /device-invites/{id}/resend` | Renvoyer l'invite |
| `telscale invite device accept` | `POST /device-invites/-/accept` | Accepter une invite |

#### Flags specifiques a `invite`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--email` | Email du destinataire | `create` |
| `--role` | Role a assigner (user invites) | `user create` |
| `--multi-use` | Invite multi-usage | `device create` |
| `--allow-exit-node` | Permettre l'utilisation comme exit node | `device create` |
| `--invite-url` | URL ou code de l'invite | `device accept` |
| `--confirm` | Bypass confirmation | `delete` |

---

### 8. `telscale posture` — Device Posture integrations (5 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale posture list` | `GET /tailnet/{tailnet}/posture/integrations` | Lister les integrations |
| `telscale posture create` | `POST /tailnet/{tailnet}/posture/integrations` | Creer une integration |
| `telscale posture get <id>` | `GET /posture/integrations/{id}` | Obtenir une integration |
| `telscale posture update <id>` | `PATCH /posture/integrations/{id}` | Modifier une integration |
| `telscale posture delete <id>` | `DELETE /posture/integrations/{id}` | Supprimer une integration |

#### Flags specifiques a `posture`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--provider` | Fournisseur : falcon/intune/jamfpro/kandji/kolide/sentinelone | `create` |
| `--cloud-id` | Identifiant cloud du provider | `create`, `update` |
| `--client-id` | Identifiant client | `create`, `update` |
| `--tenant-id` | Tenant ID (Intune) | `create`, `update` |
| `--client-secret` | Secret d'auth (interactif si omis) | `create`, `update` |
| `--confirm` | Bypass confirmation | `delete` |

---

### 9. `telscale contact` — Contacts (3 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale contact get` | `GET /tailnet/{tailnet}/contacts` | Obtenir les contacts |
| `telscale contact update <type> <email>` | `PATCH .../contacts/{type}` | Modifier un contact |
| `telscale contact resend-verification <type>` | `POST .../contacts/{type}/resend-verification-email` | Renvoyer l'email de verification |

`<type>` : `account`, `support`, `security`

---

### 10. `telscale webhook` — Webhooks (7 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale webhook list` | `GET /tailnet/{tailnet}/webhooks` | Lister les webhooks |
| `telscale webhook create` | `POST /tailnet/{tailnet}/webhooks` | Creer un webhook |
| `telscale webhook get <id>` | `GET /webhooks/{id}` | Obtenir un webhook |
| `telscale webhook update <id>` | `PATCH /webhooks/{id}` | Modifier un webhook |
| `telscale webhook delete <id>` | `DELETE /webhooks/{id}` | Supprimer un webhook |
| `telscale webhook test <id>` | `POST /webhooks/{id}/test` | Tester un webhook |
| `telscale webhook rotate-secret <id>` | `POST /webhooks/{id}/rotate` | Rotation du secret |

#### Flags specifiques a `webhook`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--url` | URL du endpoint | `create`, `update` |
| `--provider` | Type de provider : slack/mattermost/googlechat/discord | `create` |
| `--events` | Evenements (ex: `--events nodeCreated,userDeleted`) | `create`, `update` |
| `--confirm` | Bypass confirmation | `delete` |

---

### 11. `telscale settings` — Tailnet Settings (2 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale settings get` | `GET /tailnet/{tailnet}/settings` | Obtenir les parametres |
| `telscale settings update` | `PATCH /tailnet/{tailnet}/settings` | Modifier les parametres |

#### Flags specifiques a `settings`

| Flag | Description | Type |
|------|-------------|------|
| `--devices-approval` | Device approval | bool |
| `--devices-auto-updates` | Auto-updates | bool |
| `--devices-key-duration` | Duree expiration cles (jours, 1-180) | int |
| `--users-approval` | User approval | bool |
| `--users-external-tailnets` | Qui peut rejoindre des tailnets externes : none/admin/member | string |
| `--network-flow-logging` | Logs de flux reseau | bool |
| `--regional-routing` | Routage regional | bool |
| `--posture-identity-collection` | Collecte d'identite posture | bool |
| `--https` | Certificats HTTPS | bool |
| `--acls-external` | ACLs geres en externe | bool |
| `--acls-external-link` | Lien vers la gestion externe des ACLs | string |

---

### 12. `telscale service` — Services VIP (7 endpoints)

| Commande | API | Description |
|----------|-----|-------------|
| `telscale service list` | `GET /tailnet/{tailnet}/services` | Lister les Services |
| `telscale service get <name>` | `GET .../services/{name}` | Obtenir un Service |
| `telscale service create <name>` | `PUT .../services/{name}` | Creer un Service |
| `telscale service update <name>` | `PUT .../services/{name}` | Modifier un Service |
| `telscale service delete <name>` | `DELETE .../services/{name}` | Supprimer un Service |
| `telscale service hosts <name>` | `GET .../services/{name}/devices` | Lister les devices d'un Service |
| `telscale service approve <name> <deviceId>` | `POST .../device/{deviceId}/approved` | Approuver/revoquer un device |

`<name>` : Prefixe `svc:` requis (ex: `svc:my-service`)

#### Flags specifiques a `service`

| Flag | Description | Commandes |
|------|-------------|-----------|
| `--comment` | Commentaire | `create`, `update` |
| `--ports` | Ports (ex: `--ports tcp:80,tcp:443`) | `create`, `update` |
| `--tags` | Tags | `create`, `update` |
| `--ipv4` | Adresse IPv4 a assigner | `create`, `update` |
| `--new-name` | Nouveau nom (pour renommer) | `update` |
| `--approved` | `true`/`false` | `approve` |
| `--confirm` | Bypass confirmation | `delete` |

---

## Arborescence complete des commandes

```
telscale
├── auth
│   ├── login
│   ├── status
│   ├── switch <context>
│   ├── list
│   └── remove <context>
├── completion {bash|zsh|fish|powershell}
├── version
│
├── device
│   ├── list
│   ├── get <deviceId>
│   ├── delete <deviceId>
│   ├── expire <deviceId>
│   ├── authorize <deviceId>
│   ├── deauthorize <deviceId>
│   ├── set-name <deviceId> <name>
│   ├── set-tags <deviceId>
│   ├── set-key <deviceId>
│   ├── set-ip <deviceId> <ipv4>
│   ├── routes
│   │   ├── list <deviceId>
│   │   └── set <deviceId>
│   └── posture
│       ├── get <deviceId>
│       ├── set <deviceId> <key>
│       ├── delete <deviceId> <key>
│       └── batch-update
│
├── acl
│   ├── get
│   ├── set
│   ├── preview
│   └── validate
│
├── key
│   ├── list
│   ├── create
│   ├── get <keyId>
│   ├── delete <keyId>
│   └── update <keyId>
│
├── dns
│   ├── nameservers
│   │   ├── list
│   │   └── set
│   ├── preferences
│   │   ├── get
│   │   └── set
│   ├── searchpaths
│   │   ├── list
│   │   └── set
│   ├── split
│   │   ├── get
│   │   ├── update
│   │   └── set
│   └── config
│       ├── get
│       └── set
│
├── log
│   ├── audit
│   │   └── list
│   ├── network
│   │   └── list
│   ├── stream
│   │   ├── status <logType>
│   │   ├── get <logType>
│   │   ├── set <logType>
│   │   └── disable <logType>
│   └── aws-id
│       ├── create
│       └── validate <id>
│
├── user
│   ├── list
│   ├── get <userId>
│   ├── set-role <userId> <role>
│   ├── approve <userId>
│   ├── suspend <userId>
│   ├── restore <userId>
│   └── delete <userId>
│
├── invite
│   ├── user
│   │   ├── list
│   │   ├── create
│   │   ├── get <id>
│   │   ├── delete <id>
│   │   └── resend <id>
│   └── device
│       ├── list <deviceId>
│       ├── create <deviceId>
│       ├── get <id>
│       ├── delete <id>
│       ├── resend <id>
│       └── accept
│
├── posture
│   ├── list
│   ├── create
│   ├── get <id>
│   ├── update <id>
│   └── delete <id>
│
├── contact
│   ├── get
│   ├── update <type> <email>
│   └── resend-verification <type>
│
├── webhook
│   ├── list
│   ├── create
│   ├── get <id>
│   ├── update <id>
│   ├── delete <id>
│   ├── test <id>
│   └── rotate-secret <id>
│
├── settings
│   ├── get
│   └── update
│
└── service
    ├── list
    ├── get <name>
    ├── create <name>
    ├── update <name>
    ├── delete <name>
    ├── hosts <name>
    └── approve <name> <deviceId>
```

---

## Exemples d'utilisation

### Configuration initiale
```bash
# Configuration interactive
telscale auth login

# Ou via env var
export TELSCALE_API_TOKEN="tskey-api-xxxxx"
export TELSCALE_TAILNET="mycompany.com"
```

### Devices
```bash
# Lister tous les devices
telscale device list

# Lister en JSON
telscale device list --json

# Lister avec tous les champs
telscale device list --fields all

# Filtrer les devices ephemeres
telscale device list --filter isEphemeral=true

# Obtenir un device specifique
telscale device get n292kg92CNTRL

# Autoriser un device
telscale device authorize n292kg92CNTRL

# Definir les tags
telscale device set-tags n292kg92CNTRL --tags tag:prod,tag:server

# Gerer les routes
telscale device routes list n292kg92CNTRL
telscale device routes set n292kg92CNTRL --routes 10.0.0.0/16,192.168.1.0/24
```

### ACL / Policy File
```bash
# Recuperer l'ACL actuel
telscale acl get

# Recuperer en JSON avec details
telscale acl get --format json --details

# Appliquer un nouveau policy file
telscale acl set --file policy.hujson

# Valider avant d'appliquer
telscale acl validate --file policy.hujson

# Preview des regles pour un utilisateur
telscale acl preview --type user --preview-for "admin@company.com" --file policy.hujson
```

### Keys
```bash
# Creer une auth key reutilisable, pre-autorisee
telscale key create --type auth --reusable --preauthorized --tags tag:ci --expiry 86400

# Lister toutes les cles
telscale key list --all

# Supprimer une cle
telscale key delete k123456CNTRL --confirm
```

### DNS
```bash
# Voir la config DNS complete
telscale dns config get

# Activer MagicDNS
telscale dns preferences set --magic-dns true

# Configurer split DNS
telscale dns split update --domain corp.internal --servers 10.0.0.53,10.0.1.53

# Definir les nameservers
telscale dns nameservers set --nameservers 8.8.8.8,1.1.1.1
```

### Users
```bash
# Lister les utilisateurs actifs
telscale user list

# Promouvoir un admin
telscale user set-role u123456 admin

# Suspendre un utilisateur
telscale user suspend u789012

# Restaurer un utilisateur
telscale user restore u789012
```

### Webhooks
```bash
# Creer un webhook Slack
telscale webhook create --url https://hooks.slack.com/xxx --provider slack --events nodeCreated,nodeDeleted

# Tester un webhook
telscale webhook test 12345

# Rotation du secret
telscale webhook rotate-secret 12345
```

### Services
```bash
# Lister les Services
telscale service list

# Creer un Service
telscale service create svc:web --ports tcp:80,tcp:443 --tags tag:prod

# Voir les devices d'un Service
telscale service hosts svc:web

# Approuver un device pour un Service
telscale service approve svc:web n292kg92CNTRL --approved true
```

### Logging
```bash
# Logs d'audit des 24 dernieres heures
telscale log audit list --start 2024-06-05T00:00:00Z --end 2024-06-06T00:00:00Z

# Filtrer par type d'evenement
telscale log audit list --start ... --end ... --event NODE.CREATE,USER.CREATE

# Configurer le streaming vers Elastic
telscale log stream set configuration --destination-type elastic --url http://elastic:9200/logs --token xxx
```

### Settings
```bash
# Voir les parametres
telscale settings get

# Activer le device approval
telscale settings update --devices-approval true

# Modifier la duree des cles
telscale settings update --devices-key-duration 90
```

---

## Integration MCP (Claude Code)

La CLI peut etre utilisee comme outil MCP dans Claude Code via un wrapper qui expose chaque commande comme un tool :

```json
{
  "mcpServers": {
    "telscale": {
      "command": "telscale",
      "args": ["--output", "json", "mcp-serve"],
      "env": {
        "TELSCALE_API_TOKEN": "tskey-api-xxxxx"
      }
    }
  }
}
```

La commande `telscale mcp-serve` lance un serveur MCP (stdio) qui expose automatiquement toutes les commandes comme tools.

---

## Stack technique recommandee

| Composant | Choix | Justification |
|-----------|-------|---------------|
| **Langage** | Go | Binaire unique, cross-compile, performant, standard CLI |
| **Framework CLI** | `cobra` + `viper` | Standard de l'industrie (kubectl, gh, docker, hugo...) |
| **Client HTTP** | `net/http` standard | Pas de dependance externe necessaire |
| **Sortie table** | `tablewriter` ou `lipgloss` | Formatage console elegant |
| **Sortie JSON** | `encoding/json` | Standard library |
| **Config** | `viper` | Multi-format, env vars, fichiers, flags |
| **Completion** | `cobra` built-in | Support bash/zsh/fish/powershell |
| **MCP** | `mcp-go` | SDK Go pour Model Context Protocol |
