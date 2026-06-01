package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wzv5/pping/pkg/ping"

	"github.com/spf13/cobra"
)

type tcpFlags struct {
	timeout time.Duration
	source  string
}

var tcpflag tcpFlags

func addTcpCommand() {
	var cmd = &cobra.Command{
		Use:   "tcp <host> <port>",
		Short: "tcp ping",
		Long:  "tcp ping",
		Args:  cobra.ExactArgs(2),
		RunE:  runtcp,
	}

	cmd.Flags().DurationVarP(&tcpflag.timeout, "timeout", "w", time.Second*4, "timeout")
	cmd.Flags().StringVarP(&tcpflag.source, "source", "s", "", "source address")
	rootCmd.AddCommand(cmd)
}

func runtcp(cmd *cobra.Command, args []string) error {
	host := args[0]
	port, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return err
	}
	fmt.Printf("Ping %s (%d):\n", host, port)
	p := ping.NewTcpPing(host, uint16(port), tcpflag.timeout)
	if tcpflag.source != "" {
		if err := p.SetSource(tcpflag.source); err != nil {
			return err
		}
	}
	return RunPing(p)
}
