package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dimer47/tailscale-cli/internal/update"
	"github.com/spf13/cobra"
)

const (
	selfUpdateOwner  = "dimer47"
	selfUpdateRepo   = "tailscale-cli"
	selfUpdateBinary = "tailscale-cli"
)

func newSelfUpdateCmd() *cobra.Command {
	var check bool

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Met à jour la CLI vers la dernière version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if Version == "dev" {
				return fmt.Errorf("self-update n'est pas disponible pour les builds dev — utilisez 'go build' à la place")
			}

			result := update.Check(Version)

			if check {
				if result != nil && result.UpdateAvailable {
					fmt.Printf("Nouvelle version disponible : %s → %s\n", result.CurrentVersion, result.LatestVersion)
				} else {
					fmt.Printf("Déjà à jour (v%s)\n", Version)
				}
				return nil
			}

			if result == nil || !result.UpdateAvailable {
				fmt.Printf("Déjà à jour (v%s)\n", Version)
				return nil
			}

			fmt.Printf("Mise à jour de tailscale-cli : %s → %s\n", result.CurrentVersion, result.LatestVersion)

			// Determine download URL
			goos := runtime.GOOS
			goarch := runtime.GOARCH
			ext := "tar.gz"
			if goos == "windows" {
				ext = "zip"
			}
			dlURL := fmt.Sprintf(
				"https://github.com/%s/%s/releases/latest/download/%s_%s_%s.%s",
				selfUpdateOwner, selfUpdateRepo, selfUpdateRepo, goos, goarch, ext,
			)

			// Download to temp file
			fmt.Printf("Téléchargement depuis %s...\n", dlURL)
			client := &http.Client{Timeout: 120 * time.Second}
			resp, err := client.Get(dlURL)
			if err != nil {
				return fmt.Errorf("échec du téléchargement : %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("échec du téléchargement : HTTP %d", resp.StatusCode)
			}

			tmpDir, err := os.MkdirTemp("", "tailscale-cli-update-*")
			if err != nil {
				return fmt.Errorf("création du répertoire temporaire : %w", err)
			}
			defer os.RemoveAll(tmpDir)

			archivePath := filepath.Join(tmpDir, selfUpdateRepo+"."+ext)
			f, err := os.Create(archivePath)
			if err != nil {
				return fmt.Errorf("création du fichier temporaire : %w", err)
			}

			written, err := io.Copy(f, resp.Body)
			f.Close()
			if err != nil {
				return fmt.Errorf("téléchargement : %w", err)
			}
			fmt.Printf("Téléchargé %.1f Mo\n", float64(written)/1024/1024)

			// Extract
			binaryName := selfUpdateBinary
			if goos == "windows" {
				binaryName += ".exe"
			}

			if ext == "tar.gz" {
				extractCmd := exec.Command("tar", "xzf", archivePath, "-C", tmpDir)
				if out, err := extractCmd.CombinedOutput(); err != nil {
					return fmt.Errorf("extraction de l'archive : %s : %w", string(out), err)
				}
			} else {
				extractCmd := exec.Command("unzip", "-o", archivePath, "-d", tmpDir)
				if out, err := extractCmd.CombinedOutput(); err != nil {
					return fmt.Errorf("extraction de l'archive : %s : %w", string(out), err)
				}
			}

			newBinary := filepath.Join(tmpDir, binaryName)
			if _, err := os.Stat(newBinary); os.IsNotExist(err) {
				return fmt.Errorf("binaire introuvable dans l'archive")
			}

			// Find current binary location
			currentBinary, err := os.Executable()
			if err != nil {
				return fmt.Errorf("localisation du binaire actuel : %w", err)
			}
			currentBinary, err = filepath.EvalSymlinks(currentBinary)
			if err != nil {
				return fmt.Errorf("résolution des liens symboliques : %w", err)
			}

			// Try direct replace first
			fmt.Printf("Installation dans %s...\n", currentBinary)
			if err := replaceBinary(newBinary, currentBinary); err != nil {
				// Need sudo
				fmt.Println("Permission refusée — tentative avec sudo...")
				sudoCmd := exec.Command("sudo", "cp", newBinary, currentBinary)
				sudoCmd.Stdin = os.Stdin
				sudoCmd.Stdout = os.Stdout
				sudoCmd.Stderr = os.Stderr
				if err := sudoCmd.Run(); err != nil {
					return fmt.Errorf("installation échouée : %w\nVous pouvez exécuter manuellement : sudo cp %s %s", err, newBinary, currentBinary)
				}
			}

			// Clear update cache
			update.ClearCache()

			fmt.Printf("Mise à jour réussie vers v%s\n", result.LatestVersion)
			return nil
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "vérifie les mises à jour sans installer")

	return cmd
}

func replaceBinary(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Chmod(0755)
}
