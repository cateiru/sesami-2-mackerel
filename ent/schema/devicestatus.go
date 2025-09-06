package schema

import (
	"time"
	
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DeviceStatus holds the schema definition for the DeviceStatus entity.
type DeviceStatus struct {
	ent.Schema
}

// Fields of the DeviceStatus.
func (DeviceStatus) Fields() []ent.Field {
	return []ent.Field{
		field.Int("battery_percentage").
			Comment("バッテリー残量（%）"),
		field.Float("battery_voltage").
			Comment("バッテリー電圧"),
		field.Int("position").
			Comment("デバイスの位置情報"),
		field.String("status").
			Comment("デバイスステータス"),
		field.Int64("timestamp").
			Comment("SESAMI APIからのタイムスタンプ"),
		field.Time("created_at").
			Default(time.Now).
			Comment("レコード作成日時"),
	}
}

// Edges of the DeviceStatus.
func (DeviceStatus) Edges() []ent.Edge {
	return nil
}
