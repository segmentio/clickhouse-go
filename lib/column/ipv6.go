package column

import (
	"net"

	"github.com/ClickHouse/clickhouse-go/lib/binary"
)

type IPv6 struct {
	base
}

func (*IPv6) Read(decoder *binary.Decoder, isNull bool) (interface{}, error) {
	v, err := decoder.Fixed(16)
	if err != nil {
		return nil, err
	}
	return net.IP(v), nil
}

func (ip *IPv6) Write(encoder *binary.Encoder, v interface{}) error {
	var netIP net.IP
	switch val := v.(type) {
	case string:
		netIP = net.ParseIP(val)
	case net.IP:
		netIP = val
	case *net.IP:
		netIP = *(val)
	default:
		return &ErrUnexpectedType{
			T:      v,
			Column: ip,
		}
	}

	if netIP == nil {
		return &ErrUnexpectedType{
			T:      v,
			Column: ip,
		}
	}
	if _, err := encoder.Write([]byte(netIP.To16())); err != nil {
		return err
	}
	return nil
}
