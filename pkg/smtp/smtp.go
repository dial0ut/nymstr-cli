package smtp

import (
    "github.com/emersion/go-smtp"
    "github.com/dial0ut/nymstr-cli/pkg/socks5"

// Replace the dial function in the client to use our SOCKS5 proxy
smtp.Dial = socks5.DialSocks5

