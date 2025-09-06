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
		field.Uint("event_type").
			Comment("イベントタイプ（lock/unlock/manual_lock/manual_unlock等）"),
		field.Int64("timestamp").
			Comment("SESAMI APIからのタイムスタンプ"),
		field.String("history_tag").
			Comment("鍵に付けられたタグやメモ 0 ~ 21bytes"),
		field.Uint("record_id").
			Comment("非連続（将来的には連続になる予定）、セサミデバイスが再起動するまでの履歴の一意のID、 小→大"),
		field.String("parameter"),
		field.Time("created_at").
			Default(time.Now).
			Comment("レコード作成日時"),
	}
}

// Edges of the DeviceHistory.
func (DeviceHistory) Edges() []ent.Edge {
	return nil
}
