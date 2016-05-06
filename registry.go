package manet

import (
	"fmt"
	"net"
	"sync"

	ma "github.com/jbenet/go-multiaddr"
)

type AddrParser func(a net.Addr) (ma.Multiaddr, error)
type MaddrParser func(ma ma.Multiaddr) (net.Addr, error)

var maddrParsers map[string]MaddrParser
var addrParsers map[string]AddrParser
var addrParsersLock sync.Mutex

type AddressSpec struct {
	// NetNames is an array of strings that may be returned
	// by net.Addr.Network() calls on addresses belonging to this type
	NetNames []string

	// Key is the string value for Multiaddr address keys
	Key string

	// ParseNetAddr parses a net.Addr belonging to this type into a multiaddr
	ParseNetAddr AddrParser

	// ConvertMultiaddr converts a multiaddr of this type back into a net.Addr
	ConvertMultiaddr MaddrParser

	// Protocol returns the multiaddr protocol struct for this type
	Protocol ma.Protocol
}

func RegisterAddressType(a *AddressSpec) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()
	for _, n := range a.NetNames {
		addrParsers[n] = a.ParseNetAddr
	}

	maddrParsers[a.Key] = a.ConvertMultiaddr
}

func init() {
	addrParsers = make(map[string]AddrParser)
	maddrParsers = make(map[string]MaddrParser)

	RegisterAddressType(tcpAddrSpec)
	RegisterAddressType(udpAddrSpec)
	RegisterAddressType(utpAddrSpec)
	RegisterAddressType(ip4AddrSpec)
	RegisterAddressType(ip6AddrSpec)

	addrParsers["ip+net"] = parseIpPlusNetAddr
}

var tcpAddrSpec = &AddressSpec{
	Key:              "tcp",
	NetNames:         []string{"tcp", "tcp4", "tcp6"},
	ParseNetAddr:     parseTcpNetAddr,
	ConvertMultiaddr: parseBasicNetMaddr,
}

var udpAddrSpec = &AddressSpec{
	Key:              "udp",
	NetNames:         []string{"udp", "udp4", "udp6"},
	ParseNetAddr:     parseUdpNetAddr,
	ConvertMultiaddr: parseBasicNetMaddr,
}

var utpAddrSpec = &AddressSpec{
	Key:              "utp",
	NetNames:         []string{"utp", "utp4", "utp6"},
	ParseNetAddr:     parseUtpNetAddr,
	ConvertMultiaddr: parseBasicNetMaddr,
}

var ip4AddrSpec = &AddressSpec{
	Key:              "ip4",
	NetNames:         []string{"ip4"},
	ParseNetAddr:     parseIpNetAddr,
	ConvertMultiaddr: parseBasicNetMaddr,
}

var ip6AddrSpec = &AddressSpec{
	Key:              "ip6",
	NetNames:         []string{"ip6"},
	ParseNetAddr:     parseIpNetAddr,
	ConvertMultiaddr: parseBasicNetMaddr,
}

func getAddrParser(net string) (AddrParser, error) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()

	parser, ok := addrParsers[net]
	if !ok {
		return nil, fmt.Errorf("unknown network %v", net)
	}
	return parser, nil
}

func getMaddrParser(name string) (MaddrParser, error) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()
	p, ok := maddrParsers[name]
	if !ok {
		return nil, fmt.Errorf("network not supported: %s", name)
	}

	return p, nil
}
