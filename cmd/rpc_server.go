package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/ollama/ollama/llm"
	"github.com/spf13/cobra"
)

func rpcServerRun(cmd *cobra.Command, args []string) error {
	rpcHost, _ := cmd.Flags().GetString("host")
	rpcPort, _ := cmd.Flags().GetInt("port")
	device, _ := cmd.Flags().GetString("device")

	exe, err := llm.FindLlamaCppBinary("llama-rpc-server")
	if err != nil {
		return fmt.Errorf("llama-rpc-server not found: %w", err)
	}

	// llama-rpc-server flags: -H/--host, -p/--port, -d/--device <dev1,dev2,...>
	cliArgs := []string{
		"--host", rpcHost,
		"--port", fmt.Sprintf("%d", rpcPort),
	}
	if device != "" {
		cliArgs = append(cliArgs, "--device", device)
	}

	slog.Info("starting llama-rpc-server", "host", rpcHost, "port", rpcPort)
	proc := exec.Command(exe, cliArgs...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	if err := proc.Start(); err != nil {
		return fmt.Errorf("failed to start RPC server: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		proc.Process.Signal(syscall.SIGTERM)
	}()

	return proc.Wait()
}
