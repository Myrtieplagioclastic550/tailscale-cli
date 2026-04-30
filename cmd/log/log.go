package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// LogOptions contains the dependencies for the log commands.
type LogOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdLog returns the log command group.
func NewCmdLog(opts LogOptions) *cobra.Command {
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "Consultation des logs",
		Long:  "Commandes pour consulter les logs de configuration et de réseau du tailnet.",
	}

	logCmd.AddCommand(newAuditCmd(opts))
	logCmd.AddCommand(newNetworkCmd(opts))
	logCmd.AddCommand(newStreamCmd(opts))
	logCmd.AddCommand(newAWSIDCmd(opts))

	return logCmd
}

// --- audit ---

func newAuditCmd(opts LogOptions) *cobra.Command {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Logs d'audit (configuration)",
	}

	auditCmd.AddCommand(newAuditListCmd(opts))

	return auditCmd
}

func newAuditListCmd(opts LogOptions) *cobra.Command {
	var start string
	var end string
	var actors []string
	var targets []string
	var events []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Liste les logs d'audit",
		RunE: func(cmd *cobra.Command, args []string) error {
			if start == "" || end == "" {
				return fmt.Errorf("les flags --start et --end sont requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/logging/configuration?start=%s&end=%s", opts.GetTailnet(), start, end)
			for _, actor := range actors {
				path += fmt.Sprintf("&actor=%s", actor)
			}
			for _, target := range targets {
				path += fmt.Sprintf("&target=%s", target)
			}
			for _, event := range events {
				path += fmt.Sprintf("&event=%s", event)
			}

			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}

	cmd.Flags().StringVar(&start, "start", "", "Date de début (requis)")
	cmd.Flags().StringVar(&end, "end", "", "Date de fin (requis)")
	cmd.Flags().StringSliceVar(&actors, "actor", nil, "Filtrer par acteur")
	cmd.Flags().StringSliceVar(&targets, "target", nil, "Filtrer par cible")
	cmd.Flags().StringSliceVar(&events, "event", nil, "Filtrer par type d'événement")

	return cmd
}

// --- network ---

func newNetworkCmd(opts LogOptions) *cobra.Command {
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "Logs réseau",
	}

	networkCmd.AddCommand(newNetworkListCmd(opts))

	return networkCmd
}

func newNetworkListCmd(opts LogOptions) *cobra.Command {
	var start string
	var end string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Liste les logs réseau",
		RunE: func(cmd *cobra.Command, args []string) error {
			if start == "" || end == "" {
				return fmt.Errorf("les flags --start et --end sont requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/logging/network?start=%s&end=%s", opts.GetTailnet(), start, end)

			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}

	cmd.Flags().StringVar(&start, "start", "", "Date de début (requis)")
	cmd.Flags().StringVar(&end, "end", "", "Date de fin (requis)")

	return cmd
}

// --- stream ---

func newStreamCmd(opts LogOptions) *cobra.Command {
	streamCmd := &cobra.Command{
		Use:   "stream",
		Short: "Gestion du streaming de logs",
	}

	streamCmd.AddCommand(newStreamStatusCmd(opts))
	streamCmd.AddCommand(newStreamGetCmd(opts))
	streamCmd.AddCommand(newStreamSetCmd(opts))
	streamCmd.AddCommand(newStreamDisableCmd(opts))

	return streamCmd
}

func newStreamStatusCmd(opts LogOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "status <logType>",
		Short: "Affiche le statut du streaming de logs",
		Long:  "Affiche le statut du streaming. logType doit être configuration ou network.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/logging/%s/stream/status", opts.GetTailnet(), args[0])
			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newStreamGetCmd(opts LogOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <logType>",
		Short: "Affiche la configuration du streaming de logs",
		Long:  "Affiche la configuration du streaming. logType doit être configuration ou network.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/logging/%s/stream", opts.GetTailnet(), args[0])
			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newStreamSetCmd(opts LogOptions) *cobra.Command {
	var destinationType string
	var url string
	var user string
	var token string
	var compression string
	var s3Bucket string
	var s3Region string
	var s3RoleARN string
	var file string
	var stdin bool

	cmd := &cobra.Command{
		Use:   "set <logType>",
		Short: "Configure le streaming de logs",
		Long:  "Configure le streaming. logType doit être configuration ou network. Utilisez --file ou --stdin pour fournir le body complet, ou les flags individuels.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			var reader io.Reader

			if file != "" {
				f, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("opening file %s: %w", file, err)
				}
				defer f.Close()
				reader = f
			} else if stdin {
				reader = os.Stdin
			} else {
				payload := map[string]interface{}{}
				if destinationType != "" {
					payload["destinationType"] = destinationType
				}
				if url != "" {
					payload["url"] = url
				}
				if user != "" {
					payload["user"] = user
				}
				if token != "" {
					payload["token"] = token
				}
				if compression != "" {
					payload["compression"] = compression
				}
				if s3Bucket != "" {
					payload["s3Bucket"] = s3Bucket
				}
				if s3Region != "" {
					payload["s3Region"] = s3Region
				}
				if s3RoleARN != "" {
					payload["s3RoleArn"] = s3RoleARN
				}

				data, err := json.Marshal(payload)
				if err != nil {
					return fmt.Errorf("marshaling request body: %w", err)
				}
				reader = bytes.NewReader(data)
			}

			path := fmt.Sprintf("/tailnet/%s/logging/%s/stream", opts.GetTailnet(), args[0])
			body, err := client.Put(path, reader)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Configuration du streaming de logs %s mise à jour.\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&destinationType, "destination-type", "", "Type de destination")
	cmd.Flags().StringVar(&url, "url", "", "URL de destination")
	cmd.Flags().StringVar(&user, "user", "", "Nom d'utilisateur")
	cmd.Flags().StringVar(&token, "token", "", "Token d'authentification")
	cmd.Flags().StringVar(&compression, "compression", "", "Type de compression")
	cmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "Nom du bucket S3")
	cmd.Flags().StringVar(&s3Region, "s3-region", "", "Région AWS du bucket S3")
	cmd.Flags().StringVar(&s3RoleARN, "s3-role-arn", "", "ARN du rôle IAM pour S3")
	cmd.Flags().StringVar(&file, "file", "", "Fichier contenant le body JSON complet")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire le body JSON depuis stdin")

	return cmd
}

func newStreamDisableCmd(opts LogOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "disable <logType>",
		Short: "Désactive le streaming de logs",
		Long:  "Désactive le streaming. logType doit être configuration ou network.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/logging/%s/stream", opts.GetTailnet(), args[0])
			body, err := client.Delete(path)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Streaming de logs %s désactivé.\n", args[0])
			return nil
		},
	}
}

// --- aws-id ---

func newAWSIDCmd(opts LogOptions) *cobra.Command {
	awsIDCmd := &cobra.Command{
		Use:   "aws-id",
		Short: "Gestion des identifiants AWS externes",
	}

	awsIDCmd.AddCommand(newAWSIDCreateCmd(opts))
	awsIDCmd.AddCommand(newAWSIDValidateCmd(opts))

	return awsIDCmd
}

func newAWSIDCreateCmd(opts LogOptions) *cobra.Command {
	var reusable bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crée un identifiant AWS externe",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]interface{}{}
			if cmd.Flags().Changed("reusable") {
				payload["reusable"] = reusable
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/aws-external-id", opts.GetTailnet())
			body, err := client.Post(path, bytes.NewReader(data))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), result, nil)
		},
	}

	cmd.Flags().BoolVar(&reusable, "reusable", false, "Créer un identifiant réutilisable")

	return cmd
}

func newAWSIDValidateCmd(opts LogOptions) *cobra.Command {
	var roleARN string

	cmd := &cobra.Command{
		Use:   "validate <id>",
		Short: "Valide une trust policy AWS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if roleARN == "" {
				return fmt.Errorf("le flag --role-arn est requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{"roleArn": roleARN}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/aws-external-id/%s/validate-aws-trust-policy", opts.GetTailnet(), args[0])
			body, err := client.Post(path, bytes.NewReader(data))
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Trust policy AWS validée pour l'identifiant %s.\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&roleARN, "role-arn", "", "ARN du rôle IAM à valider (requis)")

	return cmd
}
