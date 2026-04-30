# tailscale-cli

CLI pour l'API Tailscale v2 — gérez votre tailnet depuis le terminal.

## Fonctionnalités

- **85 endpoints** couverts : devices, ACL, DNS, clés, utilisateurs, webhooks, services, etc.
- **Token sécurisé** : stocké dans le Keychain macOS (jamais en clair sur le disque)
- **Multi-contextes** : gérez plusieurs comptes Tailscale
- **Sortie flexible** : table, JSON, YAML, CSV
- **Intégration MCP** : utilisable comme outil dans Claude Code

## Installation

### Depuis les releases GitHub

```bash
# macOS (Apple Silicon)
curl -sL https://github.com/dimer47/tailscale-cli/releases/latest/download/tailscale-cli_darwin_arm64.tar.gz | tar xz
sudo mv tailscale-cli /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/dimer47/tailscale-cli/releases/latest/download/tailscale-cli_darwin_amd64.tar.gz | tar xz
sudo mv tailscale-cli /usr/local/bin/

# Linux (amd64)
curl -sL https://github.com/dimer47/tailscale-cli/releases/latest/download/tailscale-cli_linux_amd64.tar.gz | tar xz
sudo mv tailscale-cli /usr/local/bin/
```

### Depuis les sources

```bash
go install github.com/dimer47/tailscale-cli@latest
```

### Homebrew (macOS)

```bash
brew tap dimer47/tap
brew install tailscale-cli
```

## Configuration

```bash
# Configuration interactive (token stocké dans le Keychain macOS)
tailscale-cli auth login

# Ou via variable d'environnement
export TSCLI_API_TOKEN="tskey-api-xxxxx"
```

### Priorité de résolution du token

1. Flag `--api-token`
2. Variable d'environnement `TSCLI_API_TOKEN`
3. Keychain macOS
4. Fichier de config (legacy)

### Multi-contextes

```bash
tailscale-cli auth login          # Configure le contexte "default"
tailscale-cli auth login          # Entrez "work" comme nom de contexte
tailscale-cli auth switch work    # Change le contexte actif
tailscale-cli auth list           # Liste tous les contextes
```

## Utilisation

```bash
# Devices
tailscale-cli device list
tailscale-cli device list --json
tailscale-cli device get <nodeId>
tailscale-cli device authorize <nodeId>
tailscale-cli device set-tags <nodeId> --tags tag:prod,tag:server
tailscale-cli device routes list <nodeId>

# ACL / Policy File
tailscale-cli acl get
tailscale-cli acl set --file policy.hujson
tailscale-cli acl validate --file policy.hujson

# DNS
tailscale-cli dns config get
tailscale-cli dns preferences set --magic-dns true
tailscale-cli dns nameservers set --nameservers 8.8.8.8,1.1.1.1

# Keys
tailscale-cli key list --all
tailscale-cli key create --type auth --reusable --preauthorized --tags tag:ci

# Users
tailscale-cli user list
tailscale-cli user set-role <userId> admin
tailscale-cli user suspend <userId>

# Webhooks
tailscale-cli webhook create --url https://hooks.slack.com/xxx --provider slack --events nodeCreated,nodeDeleted
tailscale-cli webhook test <id>

# Settings
tailscale-cli settings get
tailscale-cli settings update --devices-approval true

# Services
tailscale-cli service list
tailscale-cli service create svc:web --ports tcp:80,tcp:443
```

## Variables d'environnement

| Variable | Description |
|----------|-------------|
| `TSCLI_API_TOKEN` | Token API Tailscale |
| `TSCLI_TAILNET` | Tailnet par défaut |
| `TSCLI_OUTPUT` | Format de sortie (table/json/yaml) |
| `TSCLI_DEBUG` | Mode debug |
| `TSCLI_CONFIG` | Chemin du fichier de config |
| `TSCLI_CONTEXT` | Contexte actif |

## Autocomplétion

```bash
# Bash
tailscale-cli completion bash > /etc/bash_completion.d/tailscale-cli

# Zsh
tailscale-cli completion zsh > "${fpath[1]}/_tailscale-cli"

# Fish
tailscale-cli completion fish > ~/.config/fish/completions/tailscale-cli.fish
```

## Licence

MIT
