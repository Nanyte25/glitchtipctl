package login

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Step 1: Define the model
type model struct {
	loading bool
	message string
}

func initialModel() model {
	return model{
		loading: true,
		message: "",
	}
}

// Step 2: Define messages
type errMsg struct{ err error }
type successMsg struct{}

// Step 3: Implement the Bubble Tea Model interface
func (m model) Init() bubbletea.Cmd {
	return performLogin
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		return model{loading: false, message: "Login failed: " + msg.err.Error()}, nil

	case successMsg:
		return model{loading: false, message: "Login successful!"}, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.loading {
		return "Attempting to log in...\n"
	}
	return m.message + "\n"
}

// Step 4: Define command execution
func performLogin() bubbletea.Msg {
	// Simulate network request
	time.Sleep(2 * time.Second) // This is where you'd make the actual network request

	// Here, replace this with actual logic to perform login
	// For demonstration, we'll assume the login is successful
	return successMsg{}
}

// Step 5: Cobra command
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your GlitchTip account",
	Long:  `This command allows you to log in to your GlitchTip account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return err
		}

		// Start the Bubble Tea program
		p := bubbletea.NewProgram(initialModel())
		if err := p.Start(); err != nil {
			fmt.Printf("Could not start the login process: %v\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	// Initialize your command here if needed
}
