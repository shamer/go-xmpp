package xmpp

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

var streamInit string = `<stream:stream from="example.net" id="CAFEBEEF00000000" version="1.0" xmlns:stream="http://etherx.jabber.org/streams" xmlns="jabber:client"><stream:features><mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl"><mechanism>PLAIN</mechanism></mechanisms></stream:features><success xmlns="urn:ietf:params:xml:ns:xmpp-sasl"/><stream:stream from="example.net" id="CAFEBEEF11111111" version="1.0" xmlns:stream="http://etherx.jabber.org/streams" xmlns="jabber:client"><stream:features><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"/><session xmlns="urn:ietf:params:xml:ns:xmpp-session"/></stream:features><iq id="x" type="result"><bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"><jid>romeo@example.net/BEEFCAFE</jid></bind></iq>`


type tlsFake struct {
	io.ReadWriter
}

func (c tlsFake) LocalAddr() net.Addr                  { return nil }
func (c tlsFake) RemoteAddr() net.Addr                 { return nil }
func (c tlsFake) SetDeadline(t time.Time) error        { return nil }
func (c tlsFake) SetReadDeadline(t time.Time) error    { return nil }
func (c tlsFake) SetWriteDeadline(t time.Time) error   { return nil }
func (c tlsFake) Write(b []byte) (int, error)          { return 0, nil }
func (c tlsFake) Read(b []byte) (n int, err error)     { return c.ReadWriter.Read(b) }
func (c tlsFake) Close() error                         { return nil }
func (c tlsFake) Handshake() error                     { return nil }
func (c tlsFake) ConnectionState() tls.ConnectionState { return tls.ConnectionState{} }
func (c tlsFake) OCSPResponse() []byte                 { return nil }
func (c tlsFake) VerifyHostname(host string) error     { return nil }

// Initialize a client with a fake XML source
func createClientFake(xmlStr, jid string) *Client {
	var cmdbuf bytes.Buffer
	bcmdbuf := bufio.NewWriter(&cmdbuf)
	var fake tlsFake
	fake.ReadWriter = bufio.NewReadWriter(bufio.NewReader(strings.NewReader(xmlStr)), bcmdbuf)
	c := &Client{tls: fake, jid: jid}
	return c
}

func TestRecvChatMessage(t *testing.T) {
	// example XML messages from http://xmpp.org/rfcs/rfc3921.html
	chatMsgXML := streamInit + `<message
		to='romeo@example.net'
		from='juliet@example.com/balcony'
		type='chat'
		xml:lang='en'>
	  <body>Wherefore art thou, Romeo?</body>
	</message>`

	c := createClientFake(chatMsgXML, "romeo@example.net")
	if err := c.init("romeo@example.net/BEEFCAFE", "passwd"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	msg, err := c.Recv()
	expectedMsg := Chat{Remote: "juliet@example.com/balcony", Type: "chat", Text: "Wherefore art thou, Romeo?"}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if msg != expectedMsg {
		t.Error(msg, "!=", expectedMsg)
	}
}
