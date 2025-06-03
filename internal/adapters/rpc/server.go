package rpc

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	primaryColor  = lipgloss.Color("#65DEF1") // Electric blue
	accentColor   = lipgloss.Color("#CC5A71") // Blush
	subtleColor   = lipgloss.Color("#34344A") // Purple
	borderColor   = lipgloss.Color("#3772FF")
	successColor  = lipgloss.Color("#98FF98") // Mint
	typeColor     = lipgloss.Color("#A8DCD1")
	highlightText = lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
	headerColor   = lipgloss.Color("#3772FF")

	headerStyle = lipgloss.NewStyle().Foreground(headerColor).Bold(true).Align(lipgloss.Center)
)

type Server struct {
	grpc *grpc.Server
}

func NewServer() *Server {
	return &Server{grpc: grpc.NewServer()}
}

func (s *Server) Run(handler *Handler) error {
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	mediav1.RegisterMediaServiceServer(s.grpc, handler.Media)

	reflection.Register(s.grpc)

	printServicesTable(s.grpc)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		s.grpc.GracefulStop()
	}()

	fmt.Printf("\n%s\n\n",
		lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Render("⇨ gRPC server started on "+conn.Addr().String()))

	return s.grpc.Serve(conn)
}

func printServicesTable(server *grpc.Server) {
	services := server.GetServiceInfo()
	if len(services) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(accentColor).Render("⚠️ No registered services"))
		return
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor)).
		Headers("SERVICE", "METHOD", "TYPE", "ENDPOINT").
		Width(100).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			switch {
			case row == table.HeaderRow: // Header
				return headerStyle
			default:
				return style.Foreground(lipgloss.Color("255"))
			}
		})

	for serviceName, serviceInfo := range services {
		for _, method := range serviceInfo.Methods {
			t.Row(
				serviceName,
				method.Name,
				getMethodType(method),
				fmt.Sprintf("/%s/%s", serviceName, method.Name),
			)
		}
	}

	fmt.Println("\n" + t.Render() + "\n")
}

func getMethodType(m grpc.MethodInfo) string {
	switch {
	case m.IsClientStream && m.IsServerStream:
		return "BIDI STREAM"
	case m.IsClientStream:
		return "CLIENT STREAM"
	case m.IsServerStream:
		return "SERVER STREAM"
	default:
		return "UNARY"
	}
}
