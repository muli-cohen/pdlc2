package subnet

import (
	"fmt"
	"strconv"
	"strings"
)

type Network struct {
	Address        string `json:"Address"`
	NetworkAddress string `json:"Network"`
	Broadcast      string `json:"Broadcast,omitempty"`
	Mask           string `json:"Mask"`
	FirstHost      string `json:"FirstHost,omitempty"`
	LastHost       string `json:"LastHost,omitempty"`
	PrefixLength   int    `json:"PrefixLength"`
	UsableHosts    uint64 `json:"UsableHosts"`
}

func fmtIP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, (ip>>16)&0xFF, (ip>>8)&0xFF, ip&0xFF)
}

func Parse(cidr string) (Network, error) {
	parts := strings.SplitN(cidr, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return Network{}, fmt.Errorf("invalid CIDR %q: missing prefix length", cidr)
	}

	addrPart := parts[0]
	prefixStr := parts[1]

	prefix, err := strconv.Atoi(prefixStr)
	if err != nil {
		return Network{}, fmt.Errorf("invalid CIDR %q: prefix length %q is not a number", cidr, prefixStr)
	}
	if prefix < 0 || prefix > 32 {
		return Network{}, fmt.Errorf("invalid CIDR %q: prefix length %d is out of range (0-32)", cidr, prefix)
	}

	octets := strings.Split(addrPart, ".")
	if len(octets) != 4 {
		return Network{}, fmt.Errorf("invalid CIDR %q: address must have 4 octets, got %d", cidr, len(octets))
	}

	var ipUint32 uint32
	for i, octet := range octets {
		val, err := strconv.Atoi(octet)
		if err != nil {
			return Network{}, fmt.Errorf("invalid CIDR %q: octet %q is not a number", cidr, octet)
		}
		if val < 0 || val > 255 {
			return Network{}, fmt.Errorf("invalid CIDR %q: octet %d is out of range (0-255)", cidr, val)
		}
		ipUint32 |= uint32(val) << uint(24-8*i)
	}

	var mask uint32
	if prefix == 0 {
		mask = 0
	} else {
		mask = ^uint32(0) << uint(32-prefix)
	}

	networkAddr := ipUint32 & mask

	n := Network{
		Address:        addrPart,
		NetworkAddress: fmtIP(networkAddr),
		Mask:           fmtIP(mask),
		PrefixLength:   prefix,
	}

	if prefix <= 30 {
		broadcast := networkAddr | ^mask
		n.Broadcast = fmtIP(broadcast)
		n.FirstHost = fmtIP(networkAddr + 1)
		n.LastHost = fmtIP(broadcast - 1)
		n.UsableHosts = (uint64(1) << uint(32-prefix)) - 2
	}

	return n, nil
}
