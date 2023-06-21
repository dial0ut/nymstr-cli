package imap

import (
    "github.com/emersion/go-imap/client"
    "github.com/dial0ut/nymstr-cli/pkg/socks5"
)

// Replace the dial function in the client to use our SOCKS5 proxy
client.Dial = socks5.DialSocks5

