package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"time"
)

// Email represents an email message.
type Email struct {
	from           string
	sender         string
	replyTo        string
	returnPath     string
	recipients     []string
	headers        textproto.MIMEHeader
	parts          []part
	attachments    []*file
	inlines        []*file
	Charset        string
	Encoding       encoding
	Encryption     encryption
	Username       string
	Password       string
	TLSConfig      *tls.Config
	ConnectTimeout int
	Error          error
}

// part represents the different content parts of an email body.
type part struct {
	contentType string
	body        *bytes.Buffer
}

// file represents the files that can be added to the email message.
type file struct {
	filename string
	mimeType string
	data     []byte
}

type encryption int

const (
	// EncryptionTLS sets encryption type to TLS when sending email
	EncryptionTLS encryption = iota
	// EncryptionSSL sets encryption type to SSL when sending email
	EncryptionSSL
	// EncryptionNone uses no encryption when sending email
	EncryptionNone
)

var encryptionTypes = [...]string{"TLS", "SSL", "None"}

func (encryption encryption) String() string {
	return encryptionTypes[encryption]
}

type encoding int

const (
	// EncodingQuotedPrintable sets the message body encoding to quoted-printable
	EncodingQuotedPrintable encoding = iota
	// EncodingBase64 sets the message body encoding to base64
	EncodingBase64
	// EncodingNone turns off encoding on the message body
	EncodingNone
)

var encodingTypes = [...]string{"quoted-printable", "base64", "binary"}

func (encoding encoding) String() string {
	return encodingTypes[encoding]
}

// New creates a new email. It uses UTF-8 by default.
func New() *Email {
	email := &Email{
		headers:    make(textproto.MIMEHeader),
		Charset:    "UTF-8",
		Encoding:   EncodingQuotedPrintable,
		Encryption: EncryptionNone,
		TLSConfig:  new(tls.Config),
	}

	email.AddHeader("MIME-Version", "1.0")

	return email
}

// GetError returns the first email error encountered
func (email *Email) GetError() error {
	return email.Error
}

// SetFrom sets the From address.
func (email *Email) SetFrom(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("From", address)

	return email
}

// SetSender sets the Sender address.
func (email *Email) SetSender(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Sender", address)

	return email
}

// SetReplyTo sets the Reply-To address.
func (email *Email) SetReplyTo(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Reply-To", address)

	return email
}

// SetReturnPath sets the Return-Path address. This is most often used
// to send bounced emails to a different email address.
func (email *Email) SetReturnPath(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Return-Path", address)

	return email
}

// AddTo adds a To address. You can provide multiple
// addresses at the same time.
func (email *Email) AddTo(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("To", addresses...)

	return email
}

// AddCc adds a Cc address. You can provide multiple
// addresses at the same time.
func (email *Email) AddCc(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Cc", addresses...)

	return email
}

// AddBcc adds a Bcc address. You can provide multiple
// addresses at the same time.
func (email *Email) AddBcc(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Bcc", addresses...)

	return email
}

// AddAddresses allows you to add addresses to the specified address header.
func (email *Email) AddAddresses(header string, addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	found := false

	// check for a valid address header
	for _, h := range []string{"To", "Cc", "Bcc", "From", "Sender", "Reply-To", "Return-Path"} {
		if header == h {
			found = true
		}
	}

	if !found {
		email.Error = errors.New("Mail Error: Invalid address header; Header: [" + header + "]")
		return email
	}

	// check to see if the addresses are valid
	for i := range addresses {
		address, err := mail.ParseAddress(addresses[i])
		if err != nil {
			email.Error = errors.New("Mail Error: " + err.Error() + "; Header: [" + header + "] Address: [" + addresses[i] + "]")
			return email
		}

		// check for more than one address
		switch {
		case header == "From" && len(email.from) > 0:
			fallthrough
		case header == "Sender" && len(email.sender) > 0:
			fallthrough
		case header == "Reply-To" && len(email.replyTo) > 0:
			fallthrough
		case header == "Return-Path" && len(email.returnPath) > 0:
			email.Error = errors.New("Mail Error: There can only be one \"" + header + "\" address; Header: [" + header + "] Address: [" + addresses[i] + "]")
			return email
		default:
			// other address types can have more than one address
		}

		// save the address
		switch header {
		case "From":
			email.from = address.Address
		case "Sender":
			email.sender = address.Address
		case "Reply-To":
			email.replyTo = address.Address
		case "Return-Path":
			email.returnPath = address.Address
		default:
			// check that the address was added to the recipients list
			email.recipients, err = addAddress(email.recipients, address.Address)
			if err != nil {
				email.Error = errors.New("Mail Error: " + err.Error() + "; Header: [" + header + "] Address: [" + addresses[i] + "]")
				return email
			}
		}

		// make sure the from and sender addresses are different
		if email.from != "" && email.sender != "" && email.from == email.sender {
			email.sender = ""
			email.headers.Del("Sender")
			email.Error = errors.New("Mail Error: From and Sender should not be set to the same address.")
			return email
		}

		// add all addresses to the headers except for Bcc and Return-Path
		if header != "Bcc" && header != "Return-Path" {
			// add the address to the headers
			email.headers.Add(header, address.String())
		}
	}

	return email
}

// addAddress adds an address to the address list if it hasn't already been added
func addAddress(addressList []string, address string) ([]string, error) {
	// loop through the address list to check for dups
	for _, a := range addressList {
		if address == a {
			return addressList, errors.New("Mail Error: Address: [" + address + "] has already been added")
		}
	}

	return append(addressList, address), nil
}

type priority int

const (
	// PriorityHigh sets the email priority to High
	PriorityHigh priority = iota
	// PriorityLow sets the email priority to Low
	PriorityLow
)

var priorities = [...]string{"High", "Low"}

func (priority priority) String() string {
	return priorities[priority]
}

// SetPriority sets the email message priority. Use with
// either "High" or "Low".
func (email *Email) SetPriority(priority priority) *Email {
	if email.Error != nil {
		return email
	}

	switch priority {
	case PriorityHigh:
		email.AddHeaders(textproto.MIMEHeader{
			"X-Priority":        {"1 (Highest)"},
			"X-MSMail-Priority": {"High"},
			"Importance":        {"High"},
		})
	case PriorityLow:
		email.AddHeaders(textproto.MIMEHeader{
			"X-Priority":        {"5 (Lowest)"},
			"X-MSMail-Priority": {"Low"},
			"Importance":        {"Low"},
		})
	default:
	}

	return email
}

// SetDate sets the date header to the provided date/time.
// The format of the string should be YYYY-MM-DD HH:MM:SS Time Zone.
//
// Example: SetDate("2015-04-28 10:32:00 CDT")
func (email *Email) SetDate(dateTime string) *Email {
	if email.Error != nil {
		return email
	}

	const dateFormat = "2006-01-02 15:04:05 MST"

	// Try to parse the provided date/time
	dt, err := time.Parse(dateFormat, dateTime)
	if err != nil {
		email.Error = errors.New("Mail Error: Setting date failed with: " + err.Error())
		return email
	}

	email.headers.Set("Date", dt.Format(time.RFC1123Z))

	return email
}

// SetSubject sets the subject of the email message.
func (email *Email) SetSubject(subject string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddHeader("Subject", subject)

	return email
}

// SetBody sets the body of the email message.
func (email *Email) SetBody(contentType, body string) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = []part{
		part{
			contentType: contentType,
			body:        bytes.NewBufferString(body),
		},
	}

	return email
}

// AddHeader adds the given "header" with the passed "value".
func (email *Email) AddHeader(header string, values ...string) *Email {
	if email.Error != nil {
		return email
	}

	// check that there is actually a value
	if len(values) < 1 {
		email.Error = errors.New("Mail Error: no value provided; Header: [" + header + "]")
		return email
	}

	switch header {
	case "Sender":
		fallthrough
	case "From":
		fallthrough
	case "To":
		fallthrough
	case "Bcc":
		fallthrough
	case "Cc":
		fallthrough
	case "Reply-To":
		fallthrough
	case "Return-Path":
		email.AddAddresses(header, values...)
	case "Date":
		if len(values) > 1 {
			email.Error = errors.New("Mail Error: To many dates provided")
			return email
		}
		email.SetDate(values[0])
	default:
		email.headers[header] = values
	}

	return email
}

// AddHeaders is used to add mulitple headers at once
func (email *Email) AddHeaders(headers textproto.MIMEHeader) *Email {
	if email.Error != nil {
		return email
	}

	for header, values := range headers {
		email.AddHeader(header, values...)
	}

	return email
}

// AddAlternative allows you to add alternative parts to the body
// of the email message. This is most commonly used to add an
// html version in addition to a plain text version that was
// already added with SetBody.
func (email *Email) AddAlternative(contentType, body string) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = append(email.parts,
		part{
			contentType: contentType,
			body:        bytes.NewBufferString(body),
		},
	)

	return email
}

// AddAttachment allows you to add an attachment to the email message.
// You can optionally provide a different name for the file.
func (email *Email) AddAttachment(file string, name ...string) *Email {
	if email.Error != nil {
		return email
	}

	if len(name) > 1 {
		email.Error = errors.New("Mail Error: Attach can only have a file and an optional name")
		return email
	}

	email.Error = email.attach(file, false, name...)

	return email
}

// AddInline allows you to add an inline attachment to the email message.
// You can optionally provide a different name for the file.
func (email *Email) AddInline(file string, name ...string) *Email {
	if email.Error != nil {
		return email
	}

	if len(name) > 1 {
		email.Error = errors.New("Mail Error: Inline can only have a file and an optional name")
		return email
	}

	email.Error = email.attach(file, true, name...)

	return email
}

// attach does the low level attaching of the files
func (email *Email) attach(f string, inline bool, name ...string) error {
	// Get the file data
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.New("Mail Error: Failed to add file with following error: " + err.Error())
	}

	// get the file mime type
	mimeType := mime.TypeByExtension(filepath.Ext(f))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// get the filename
	_, filename := filepath.Split(f)

	// if an alternative filename was provided, use that instead
	if len(name) == 1 {
		filename = name[0]
	}

	if inline {
		email.inlines = append(email.inlines, &file{
			filename: filename,
			mimeType: mimeType,
			data:     data,
		})
	} else {
		email.attachments = append(email.attachments, &file{
			filename: filename,
			mimeType: mimeType,
			data:     data,
		})
	}

	return nil
}

// getFrom returns the sender of the email, if any
func (email *Email) getFrom() string {
	from := email.returnPath
	if from == "" {
		from = email.sender
		if from == "" {
			from = email.from
			if from == "" {
				from = email.replyTo
			}
		}
	}

	return from
}

func (email *Email) hasMixedPart() bool {
	return (len(email.parts) > 0 && len(email.attachments) > 0) || len(email.attachments) > 1
}

func (email *Email) hasRelatedPart() bool {
	return (len(email.parts) > 0 && len(email.inlines) > 0) || len(email.inlines) > 1
}

func (email *Email) hasAlternativePart() bool {
	return len(email.parts) > 1
}

// GetMessage builds and returns the email message
func (email *Email) GetMessage() string {
	msg := newMessage(email)

	if email.hasMixedPart() {
		msg.openMultipart("mixed")
	}

	if email.hasRelatedPart() {
		msg.openMultipart("related")
	}

	if email.hasAlternativePart() {
		msg.openMultipart("alternative")
	}

	for _, part := range email.parts {
		msg.addBody(part.contentType, part.body.Bytes())
	}

	if email.hasAlternativePart() {
		msg.closeMultipart()
	}

	msg.addFiles(email.inlines, true)
	if email.hasRelatedPart() {
		msg.closeMultipart()
	}

	msg.addFiles(email.attachments, false)
	if email.hasMixedPart() {
		msg.closeMultipart()
	}

	return msg.getHeaders() + msg.body.String()
}

// Send sends the composed email
func (email *Email) Send(address string) error {
	if email.Error != nil {
		return email.Error
	}

	var auth smtp.Auth

	from := email.getFrom()
	if from == "" {
		return errors.New(`Mail Error: No "From" address specified.`)
	}

	if len(email.recipients) < 1 {
		return errors.New("Mail Error: No recipient specified.")
	}

	msg := email.GetMessage()

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return errors.New("Mail Error: " + err.Error())
	}

	if email.Username != "" || email.Password != "" {
		auth = smtp.PlainAuth("", email.Username, email.Password, host)
	}

	return send(host, port, from, email.recipients, msg, auth, email.Encryption, email.TLSConfig, email.ConnectTimeout)
}

// dial connects to the smtp server with the request encryption type
func dial(host string, port string, encryption encryption, config *tls.Config) (*smtp.Client, error) {
	var conn net.Conn
	var err error

	address := host + ":" + port

	// do the actual dial
	switch encryption {
	case EncryptionSSL:
		conn, err = tls.Dial("tcp", address, config)
	default:
		conn, err = net.Dial("tcp", address)
	}

	if err != nil {
		return nil, errors.New("Mail Error on dailing with encryption type " + encryption.String() + ": " + err.Error())
	}

	c, err := smtp.NewClient(conn, host)

	if err != nil {
		return nil, errors.New("Mail Error on smtp dial: " + err.Error())
	}

	return c, err
}

// smtpConnect connects to the smtp server and starts TLS and passes auth
// if necessary
func smtpConnect(host string, port string, from string, to []string, msg string, auth smtp.Auth, encryption encryption, config *tls.Config) (*smtp.Client, error) {
	// connect to the mail server
	c, err := dial(host, port, encryption, config)

	if err != nil {
		return nil, err
	}

	// send Hello
	if err = c.Hello("localhost"); err != nil {
		c.Close()
		return nil, errors.New("Mail Error on Hello: " + err.Error())
	}

	// start TLS if necessary
	if encryption == EncryptionTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if config.ServerName == "" {
				config = &tls.Config{ServerName: host}
			}

			if err = c.StartTLS(config); err != nil {
				c.Close()
				return nil, errors.New("Mail Error on Start TLS: " + err.Error())
			}
		}
	}

	// pass the authentication if necessary
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				c.Close()
				return nil, errors.New("Mail Error on Auth: " + err.Error())
			}
		}
	}

	return c, nil
}

type smtpConnectErrorChannel struct {
	client *smtp.Client
	err    error
}

// send does the low level sending of the email
func send(host string, port string, from string, to []string, msg string, auth smtp.Auth, encryption encryption, config *tls.Config, connectTimeout int) error {
	var smtpConnectChannel chan smtpConnectErrorChannel
	var c *smtp.Client
	var err error

	// set the timeout value
	timeout := time.Duration(connectTimeout) * time.Second

	// if there is a timeout, setup the channel and do the connect under a goroutine
	if timeout != 0 {
		smtpConnectChannel = make(chan smtpConnectErrorChannel, 2)
		go func() {
			c, err = smtpConnect(host, port, from, to, msg, auth, encryption, config)
			// send the result
			smtpConnectChannel <- smtpConnectErrorChannel{
				client: c,
				err:    err,
			}
		}()
	}

	if timeout == 0 {
		// no timeout, just fire the connect
		c, err = smtpConnect(host, port, from, to, msg, auth, encryption, config)
	} else {
		// get the connect result or timeout result, which ever happens first
		select {
		case result := <-smtpConnectChannel:
			c = result.client
			err = result.err
		case <-time.After(timeout):
			return errors.New("Mail Error: SMTP Connection timed out")
		}
	}

	// check for connect error
	if err != nil {
		return err
	}

	defer c.Close()

	// Set the sender
	if err := c.Mail(from); err != nil {
		return err
	}

	// Set the recipients
	for _, address := range to {
		if err := c.Rcpt(address); err != nil {
			return err
		}
	}

	// Send the data command
	w, err := c.Data()
	if err != nil {
		return err
	}

	// write the message
	_, err = fmt.Fprint(w, msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
