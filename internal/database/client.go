package database

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/cateiru/sesami-2-mackerel/ent"
	"github.com/cateiru/sesami-2-mackerel/ent/devicehistory"
	"github.com/cateiru/sesami-2-mackerel/ent/devicestatus"
	"github.com/cateiru/sesami-2-mackerel/internal/sesami"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	entClient *ent.Client
}

func NewClient(dbPath string) (*Client, error) {
	drv, err := sql.Open(dialect.SQLite, fmt.Sprintf("file:%s?cache=shared&_fk=1", dbPath))
	if err != nil {
		return nil, fmt.Errorf("SQLiteドライバー開封エラー: %w", err)
	}

	entClient := ent.NewClient(ent.Driver(drv))

	if err := entClient.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("スキーマ作成エラー: %w", err)
	}

	return &Client{
		entClient: entClient,
	}, nil
}

func (c *Client) Close() error {
	return c.entClient.Close()
}

func (c *Client) InsertDeviceStatus(status *sesami.DeviceStatus) error {
	_, err := c.entClient.DeviceStatus.
		Create().
		SetBatteryPercentage(status.BatteryPercentage).
		SetBatteryVoltage(status.BatteryVoltage).
		SetPosition(status.Position).
		SetStatus(status.CHSesame2Status).
		SetTimestamp(status.Timestamp).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("デバイスステータス挿入エラー: %w", err)
	}

	return nil
}

func (c *Client) GetDeviceStatusHistory(limit int) ([]*ent.DeviceStatus, error) {
	records, err := c.entClient.DeviceStatus.
		Query().
		Order(ent.Desc(devicestatus.FieldCreatedAt)).
		Limit(limit).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("デバイスステータス取得エラー: %w", err)
	}

	return records, nil
}

func (c *Client) InsertDeviceHistory(deviceUUID string, history []sesami.HistoryEntry) error {
	for _, entry := range history {
		exists, err := c.entClient.DeviceHistory.
			Query().
			Where(
				devicehistory.DeviceUUID(deviceUUID),
				devicehistory.Timestamp(entry.TimeStamp),
				devicehistory.EventType(entry.Type),
			).
			Exist(context.Background())

		if err != nil {
			return fmt.Errorf("履歴重複チェックエラー: %w", err)
		}

		if exists {
			continue
		}

		_, err = c.entClient.DeviceHistory.
			Create().
			SetDeviceUUID(deviceUUID).
			SetEventType(entry.Type).
			SetTimestamp(entry.TimeStamp).
			SetTag(entry.Tag).
			Save(context.Background())

		if err != nil {
			return fmt.Errorf("履歴データ挿入エラー: %w", err)
		}
	}

	return nil
}

func (c *Client) GetDeviceHistory(deviceUUID string, limit int) ([]*ent.DeviceHistory, error) {
	records, err := c.entClient.DeviceHistory.
		Query().
		Where(devicehistory.DeviceUUID(deviceUUID)).
		Order(ent.Desc(devicehistory.FieldTimestamp)).
		Limit(limit).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("履歴データ取得エラー: %w", err)
	}

	return records, nil
}
