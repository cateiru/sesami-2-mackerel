package schema

import (
	"time"
	
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DeviceHistory holds the schema definition for the DeviceHistory entity.
type DeviceHistory struct {
	ent.Schema
}

// Fields of the DeviceHistory.
func (DeviceHistory) Fields() []ent.Field {
	return []ent.Field{
		field.String("device_uuid").
			Comment("デバイスUUID"),
		field.String("event_type").
			Comment("イベントタイプ（lock/unlock/manual_lock/manual_unlock等）"),
		field.Int64("timestamp").
			Comment("SESAMI APIからのタイムスタンプ"),
		field.String("user_id").
			Optional().
			Comment("操作を行ったユーザーID"),
		field.String("tag").
			Optional().
			Comment("操作に関連するタグ"),
		field.Time("created_at").
			Default(time.Now).
			Comment("レコード作成日時"),
	}
}

// Edges of the DeviceHistory.
func (DeviceHistory) Edges() []ent.Edge {
	return nil
}