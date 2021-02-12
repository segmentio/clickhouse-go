package data

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/lib/binary"
)

const ClientName = "Golang SQLDriver"

const (
	ClickHouseRevision         = 54213
	ClickHouseDBMSVersionMajor = 1
	ClickHouseDBMSVersionMinor = 1
)

type ClientInfo struct{}

func (ClientInfo) Write(encoder *binary.Encoder) error {
	// TODO: check errors?
	_ = encoder.String(ClientName)
	_ = encoder.Uvarint(ClickHouseDBMSVersionMajor)
	_ = encoder.Uvarint(ClickHouseDBMSVersionMinor)
	_ = encoder.Uvarint(ClickHouseRevision)
	return nil
}

func (ClientInfo) String() string {
	return fmt.Sprintf("%s %d.%d.%d", ClientName, ClickHouseDBMSVersionMajor, ClickHouseDBMSVersionMinor, ClickHouseRevision)
}
