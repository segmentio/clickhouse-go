package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/lib/data"
	"github.com/ClickHouse/clickhouse-go/lib/protocol"
)

func (ch *clickhouse) sendQuery(query string, externalTables []ExternalTable) error {
	ch.logf("[send query] %s", query)
	if err := ch.encoder.Uvarint(protocol.ClientQuery); err != nil {
		return err
	}
	if err := ch.encoder.String(""); err != nil {
		return err
	}
	{ // client info
		// TODO: check errors?
		_ = ch.encoder.Uvarint(1)
		_ = ch.encoder.String("")
		_ = ch.encoder.String("") //initial_query_id
		_ = ch.encoder.String("[::ffff:127.0.0.1]:0")
		_ = ch.encoder.Uvarint(1) // iface type TCP
		_ = ch.encoder.String(hostname)
		_ = ch.encoder.String(hostname)
	}
	if err := ch.ClientInfo.Write(ch.encoder); err != nil {
		return err
	}
	if ch.ServerInfo.Revision >= protocol.DBMS_MIN_REVISION_WITH_QUOTA_KEY_IN_CLIENT_INFO {
		// TODO: check errors?
		_ = ch.encoder.String("")
	}

	// the settings are written as list of contiguous name-value pairs, finished with empty name
	if !ch.settings.IsEmpty() {
		ch.logf("[query settings] %s", ch.settings.settingsStr)
		if err := ch.settings.Serialize(ch.encoder); err != nil {
			return err
		}
	}
	// empty string is a marker of the end of the settings
	if err := ch.encoder.String(""); err != nil {
		return err
	}
	if err := ch.encoder.Uvarint(protocol.StateComplete); err != nil {
		return err
	}
	compress := protocol.CompressDisable
	if ch.compress {
		compress = protocol.CompressEnable
	}
	if err := ch.encoder.Uvarint(compress); err != nil {
		return err
	}
	if err := ch.encoder.String(query); err != nil {
		return err
	}
	if err := ch.sendExternalTables(externalTables); err != nil {
		return err
	}
	if err := ch.writeBlock(&data.Block{}, ""); err != nil {
		return err
	}
	return ch.encoder.Flush()
}
