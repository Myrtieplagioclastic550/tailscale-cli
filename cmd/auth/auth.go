package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/config"
	"github.com/dimer47/tailscale-cli/internal/keychain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCmdAuth returns the auth command group.
func NewCmdAuth() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Gestion de l'authentification et des contextes",
		Long: `Commandes pour gérer l'authentification auprès de l'API Tailscale.

Le token API est stocké de façon sécurisée dans le Keychain macOS.
Il n'est jamais écrit en clair sur le disque.

Priorité de résolution du token :
  1. Flag --api-token
  2. Variable d'environnement TSCLI_API_TOKEN
  3. Keychain macOS (via 'tailscale-cli auth login')`,
	}

	authCmd.AddCommand(newLoginCmd())
	authCmd.AddCommand(newStatusCmd())
	authCmd.AddCommand(newSwitchCmd())
	authCmd.AddCommand(newListCmd())
	authCmd.AddCommand(newRemoveCmd())

	return authCmd
}

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authentification interactive — stocke le token dans le Keychain macOS",
		Long: `Demande interactivement le token API et le tailnet, puis :
  - Stocke le token dans le Keychain macOS (chiffré)
  - Sauvegarde le tailnet et le contexte dans ~/.tailscale-cli/config.json

Si un token existe déjà pour ce contexte, il est remplacé.
Les tokens Tailscale expirent après 1 à 90 jours.
Relancez 'tailscale-cli auth login' pour mettre à jour un token expiré.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !keychain.IsAvailable() {
				return fmt.Errorf("le Keychain macOS n'est pas disponible sur ce système")
			}

			scanner := bufio.NewScanner(os.Stdin)

			// Demander le nom du contexte
			fmt.Print("Nom du contexte (default) : ")
			scanner.Scan()
			contextName := strings.TrimSpace(scanner.Text())
			if contextName == "" {
				contextName = "default"
			}

			// Demander le token API
			fmt.Print("Token API Tailscale : ")
			scanner.Scan()
			token := strings.TrimSpace(scanner.Text())
			if token == "" {
				return fmt.Errorf("le token API ne peut pas être vide")
			}

			// Demander le tailnet
			fmt.Print("Tailnet (- pour le tailnet par défaut) : ")
			scanner.Scan()
			tailnet := strings.TrimSpace(scanner.Text())
			if tailnet == "" {
				tailnet = "-"
			}

			// Stocker le token dans le Keychain (remplace l'ancien si existant)
			if err := keychain.Set(contextName, token); err != nil {
				return fmt.Errorf("erreur Keychain : %w", err)
			}

			// Charger ou créer la configuration (sans token)
			cfgPath := viper.GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				cfg = &config.Config{
					Contexts: make(map[string]config.Context),
				}
			}

			// Sauvegarder le contexte (sans token — il est dans le Keychain)
			cfg.Contexts[contextName] = config.Context{
				Tailnet: tailnet,
			}
			cfg.DefaultContext = contextName

			if err := cfg.Save(cfgPath); err != nil {
				return fmt.Errorf("erreur lors de la sauvegarde de la configuration : %w", err)
			}

			fmt.Printf("\nContexte %q configuré :\n", contextName)
			fmt.Printf("  Tailnet : %s\n", tailnet)
			fmt.Printf("  Token   : stocké dans le Keychain macOS\n")
			fmt.Printf("\nPour mettre à jour un token expiré, relancez 'tailscale-cli auth login'.\n")
			return nil
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Affiche le statut d'authentification du contexte actif",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := viper.GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("impossible de charger la configuration : %w", err)
			}

			contextName := viper.GetString("context")
			ctx, resolvedName, err := config.GetActiveContext(cfg, contextName)
			if err != nil {
				return fmt.Errorf("impossible de résoudre le contexte actif : %w", err)
			}

			fmt.Printf("Contexte actif : %s\n", resolvedName)
			fmt.Printf("Tailnet        : %s\n", ctx.Tailnet)

			// Résoudre le token
			token := ""
			tokenSource := ""

			// 1. Flag / env
			if t := viper.GetString("api-token"); t != "" {
				token = t
				tokenSource = "flag/env"
			}

			// 2. Keychain
			if token == "" && keychain.IsAvailable() {
				if t, err := keychain.Get(resolvedName); err == nil && t != "" {
					token = t
					tokenSource = "Keychain macOS"
				}
			}

			// 3. Config file (legacy)
			if token == "" && ctx.APIToken != "" {
				token = ctx.APIToken
				tokenSource = "config file (legacy, non sécurisé)"
			}

			if token == "" {
				fmt.Printf("Token          : non configuré\n")
				return nil
			}

			// Masquer le token sauf les 8 derniers caractères
			masked := maskToken(token)
			fmt.Printf("Token          : %s (source: %s)\n", masked, tokenSource)

			// Vérifier le token en faisant un appel API
			client := api.NewClient(token)
			_, err = client.Get(fmt.Sprintf("/tailnet/%s/devices", ctx.Tailnet))
			if err != nil {
				fmt.Printf("Statut         : invalide ou expiré (%v)\n", err)
				fmt.Printf("\nPour mettre à jour le token, lancez 'tailscale-cli auth login'.\n")
			} else {
				fmt.Printf("Statut         : valide\n")
			}

			return nil
		},
	}
}

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <context>",
		Short: "Change le contexte actif",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			cfgPath := viper.GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("impossible de charger la configuration : %w", err)
			}

			if _, ok := cfg.Contexts[contextName]; !ok {
				return fmt.Errorf("le contexte %q n'existe pas dans la configuration", contextName)
			}

			cfg.DefaultContext = contextName

			if err := cfg.Save(cfgPath); err != nil {
				return fmt.Errorf("erreur lors de la sauvegarde de la configuration : %w", err)
			}

			fmt.Printf("Contexte actif changé pour %q.\n", contextName)
			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Affiche tous les contextes configurés",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := viper.GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("impossible de charger la configuration : %w", err)
			}

			if len(cfg.Contexts) == 0 {
				fmt.Println("Aucun contexte configuré.")
				return nil
			}

			for name, ctx := range cfg.Contexts {
				prefix := "  "
				if name == cfg.DefaultContext {
					prefix = "* "
				}

				// Vérifier si un token existe dans le Keychain
				tokenStatus := "pas de token"
				if keychain.IsAvailable() {
					if t, err := keychain.Get(name); err == nil && t != "" {
						tokenStatus = "token dans Keychain"
					}
				}
				if ctx.APIToken != "" {
					tokenStatus = "token dans config (legacy)"
				}

				fmt.Printf("%s%s (tailnet: %s, %s)\n", prefix, name, ctx.Tailnet, tokenStatus)
			}

			return nil
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <context>",
		Short: "Supprime un contexte et son token du Keychain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			cfgPath := viper.GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("impossible de charger la configuration : %w", err)
			}

			if _, ok := cfg.Contexts[contextName]; !ok {
				return fmt.Errorf("le contexte %q n'existe pas dans la configuration", contextName)
			}

			// Supprimer le token du Keychain
			if keychain.IsAvailable() {
				if err := keychain.Delete(contextName); err != nil {
					fmt.Fprintf(os.Stderr, "Avertissement : impossible de supprimer le token du Keychain : %v\n", err)
				}
			}

			// Supprimer le contexte de la config
			delete(cfg.Contexts, contextName)

			if cfg.DefaultContext == contextName {
				cfg.DefaultContext = ""
			}

			if err := cfg.Save(cfgPath); err != nil {
				return fmt.Errorf("erreur lors de la sauvegarde de la configuration : %w", err)
			}

			fmt.Printf("Contexte %q supprimé (token retiré du Keychain).\n", contextName)
			return nil
		},
	}
}

// maskToken masque un token en ne gardant que les 8 derniers caractères.
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return "****" + token[len(token)-8:]
}
