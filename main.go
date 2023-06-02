package main

import (
    "fmt"
    "github.com/abiosoft/ishell"
    "github.com/gorilla/websocket"
    "encoding/json"
    "os"
    "os/exec"
    "io/ioutil"
    "log"
    "strings"
)

const contactFile = "$HOME/.nym/clients/nymstr-chat/config/contacts.lst"

type NymCLI struct {
    *ishell.Shell
    conn *websocket.Conn
    contacts map[string]string // Contact list
}

func New(uri string) *NymCLI {
    conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
    if err != nil {
        panic(err)
    }

    cli := &NymCLI{
        Shell: ishell.New(),
        conn: conn,
        contacts: make(map[string]string), // Initialize contact list
    }

    sendCmd := &ishell.Cmd{
        Name: "send",
        Help: "send a message to a recipient",
        Func: cli.sendCommand,
    }
    cli.AddCmd(sendCmd)

    contactsCmd := &ishell.Cmd{
        Name: "contacts",
        Help: "save a new contact or list all the contacts",
        Func: cli.contactsCommand,
    }
    cli.AddCmd(contactsCmd)

    testCmd := &ishell.Cmd{
        Name: "test",
        Help: "send a message to self",
        Func: cli.testCommand,
    }
    cli.AddCmd(testCmd)

    return cli
}


func (cli *NymCLI) sendCommand(c *ishell.Context) {
    if len(cli.contacts) == 0 {
        fmt.Println("No contacts in the list. First add an address with `contacts add <address>`")
        return
    }

    // Display list of contacts
    contacts := make([]string, 0, len(cli.contacts))
    for name := range cli.contacts {
        contacts = append(contacts, name)
    }

    choice := c.MultiChoice(contacts, "Choose a contact")
    recipient := cli.contacts[contacts[choice]]
    c.Print("Message to ", recipient, ": ")
    message := c.ReadLine()

    // Send message
    sendRequest, err := json.Marshal(map[string]interface{}{
        "type": "send",
        "recipient": recipient,
        "message": message,
        "withReplySurb": true,
    })
    if err != nil {
        panic(err)
    }
    if err = cli.conn.WriteMessage(websocket.TextMessage, []byte(sendRequest)); err != nil {
        panic(err)
    }
}
func (cli *NymCLI) contactsCommand(c *ishell.Context) {
    if len(c.Args) == 0 {
        for name, address := range cli.contacts {
            fmt.Printf("%s: %s\n", name, address)
        }
        return
    }

    if c.Args[0] == "add" && len(c.Args) == 3 {
        name := c.Args[1]
        address := c.Args[2]
        cli.contacts[name] = address
        fmt.Printf("Added %s to contacts\n", name)
    } else {
        fmt.Println("usage: contacts add <name> <address>")
    }
}
func (cli *NymCLI) getSelfAddress() string {
    selfAddressRequest, err := json.Marshal(map[string]string{"type": "selfAddress"})
    if err != nil {
        panic(err)
    }

    if err = cli.conn.WriteMessage(websocket.TextMessage, []byte(selfAddressRequest)); err != nil {
        panic(err)
    }

    responseJSON := make(map[string]interface{})
    err = cli.conn.ReadJSON(&responseJSON)
    if err != nil {
        panic(err)
    }

    return responseJSON["address"].(string)
}

func (cli *NymCLI) testCommand(c *ishell.Context) {
    selfAddress := cli.getSelfAddress()
    fmt.Printf("our address is: %v\n", selfAddress)

    c.Print("Message to self: ")
    message := c.ReadLine()

    // Send message to self
    sendRequest, err := json.Marshal(map[string]interface{}{
        "type": "send",
        "recipient": selfAddress,
        "message": message,
        "withReplySurb": true    })
    if err != nil {
        panic(err)
    }
    if err = cli.conn.WriteMessage(websocket.TextMessage, []byte(sendRequest)); err != nil {
        panic(err)
    }
}
// ...

// Launch nym-client if it is not running and exit if connection is not possible
func (cli *NymCLI) ConnectOrStartNymClient(uri string) error {
    conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
    if err == nil {
        // we managed to connect, so we close the connection and return
        conn.Close()
        return nil
    }
    // there was an error connecting, so we try to start nym-client

    logFile, err := os.Create("/path/to/your/project/error.log")
    if err != nil {
        return fmt.Errorf("failed to create log file: %w", err)
    }

    cmd := exec.Command("./nym-client", "run", "--id", "nymstr-cli")
    cmd.Dir = "." // this is the directory in which the command will be run
    cmd.Stdout = logFile
    cmd.Stderr = logFile

    err = cmd.Start()
    if err != nil {
        return fmt.Errorf("could not start nym-client: %w", err)
    }
    // wait for the command to finish
    err = cmd.Wait()
    if err != nil {
        // the command exited with an error
        return fmt.Errorf("nym-client exited with error: %w", err)
    }
    return nil
}

func saveContacts(contacts []string) {
    // Convert the slice of strings into a single string with line breaks
    content := strings.Join(contacts, "\n")

    // Write the content to the file
    err := ioutil.WriteFile(os.ExpandEnv(contactFile), []byte(content), 0644)
    if err != nil {
        log.Fatalf("Error writing to file: %v", err)
    }
}

func readContacts() []string {
    // Check if the file exists
    if _, err := os.Stat(os.ExpandEnv(contactFile)); os.IsNotExist(err) {
        // File does not exist, return an empty slice
        return []string{}
    }

    // Read the file
    content, err := ioutil.ReadFile(os.ExpandEnv(contactFile))
    if err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    // Split the content by line break to get the contacts
    contacts := strings.Split(string(content), "\n")

    return contacts
}

func main() {
    uri := "ws://localhost:1977"
    cli := New(uri)
    cli.Run()
}

