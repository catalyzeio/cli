package prompts

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// IPrompts is the interface in which to interact with the user and accept
// input.
type IPrompts interface {
	UsernamePassword() (string, string, error)
	KeyPassphrase(string) string
	Password(msg string) string
	PHI() error
	YesNo(msg string) error
}

// SPrompts is a concrete implementation of IPrompts
type SPrompts struct{}

// New returns a new instance of IPrompts
func New() IPrompts {
	return &SPrompts{}
}

var validAnswers = map[string]bool{
	"y":   true,
	"yes": true,
	"n":   false,
	"no":  false,
}

// UsernamePassword prompts a user to enter their username and password.
func (p *SPrompts) UsernamePassword() (string, string, error) {
	var username string
	fmt.Print("Username or Email: ")
	in := bufio.NewReader(os.Stdin)
	username, err := in.ReadString('\n')
	if err != nil {
		return "", "", errors.New("Invalid username")
	}
	username = strings.TrimRight(username, "\n")
	if runtime.GOOS == "windows" {
		username = strings.TrimRight(username, "\r")
	}
	fmt.Print("Password: ")
	bytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	return username, string(bytes), nil
}

// KeyPassphrase prompts a user to enter a passphrase for a named key.
func (p *SPrompts) KeyPassphrase(filepath string) string {
	fmt.Printf("Enter passphrase for %s: ", filepath)
	bytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	return string(bytes)
}

// PHI prompts a user to accept liability for downloading PHI to their local
// machine.
func (p *SPrompts) PHI() error {
	var answer string
	for {
		fmt.Println("This operation might result in PHI data being downloaded and decrypted to your local machine. By entering \"y\" at the prompt below, you warrant that you have the necessary privileges to view the data, have taken all necessary precautions to secure this data, and absolve Catalyze of any issues that might arise from its loss.")
		fmt.Print("Do you wish to proceed? (y/n) ")
		fmt.Scanln(&answer)
		fmt.Println("")
		if _, contains := validAnswers[strings.ToLower(answer)]; !contains {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if !validAnswers[strings.ToLower(answer)] {
		return fmt.Errorf("Exiting")
	}
	return nil
}

// YesNo outputs a given message and waits for a user to answer `y/n`.
// If yes, flow continues as normal. If no, an error is returned. The given
// message SHOULD contain the string "(y/n)" or some other form of y/n
// indicating that the user needs to type in y or n. This method does not do
// that for you. The message will not have a new line appended to it. If you
// require a newline, add this to the given message.
func (p *SPrompts) YesNo(msg string) error {
	var answer string
	for {
		fmt.Printf(msg)
		fmt.Scanln(&answer)
		fmt.Println("")
		if _, contains := validAnswers[strings.ToLower(answer)]; !contains {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if !validAnswers[strings.ToLower(answer)] {
		return fmt.Errorf("Exiting")
	}
	return nil
}

// Password prompts the user for a password displaying the given message.
// The password will be hidden while typed. A newline is not added to the given
// message. If a newline is required, it should be part of the passed in string.
func (p *SPrompts) Password(msg string) string {
	fmt.Print(msg)
	bytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	return string(bytes)
}
