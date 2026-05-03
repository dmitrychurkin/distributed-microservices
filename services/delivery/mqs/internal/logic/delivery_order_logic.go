package logic

import (
	"context"
	"delivery/mqs/internal/svc"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderCreated struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderCreated(ctx context.Context, svcCtx *svc.ServiceContext) *OrderCreated {
	return &OrderCreated{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderCreated) Consume(ctx context.Context, key, val string) error {
	logx.Infof("OrderCreated key :%s , val :%s", key, val)

	if val == "" || len(val) == 0 {
		logx.Infof("Received TOMBSTONE for key: %s. Deleting from local cache...", key)
		return nil
	}

	var orderMessage orderMessage
	if err := json.Unmarshal([]byte(val), &orderMessage); err != nil {
		return err
	}

	var orderPayload orderPayload
	if err := json.Unmarshal([]byte(orderMessage.Payload), &orderPayload); err != nil {
		return err
	}

	if _, err := l.svcCtx.SqlConn.ExecCtx(
		ctx,
		`INSERT INTO deliveries (order_id, product_id)
		VALUES ($1, $2)
		ON CONFLICT (order_id)
		DO UPDATE SET
			product_id = EXCLUDED.product_id,
			updated_at = CURRENT_TIMESTAMP`,
		orderPayload.OrderId,
		orderPayload.ProductId,
	); err != nil {
		return err
	}

	return nil
}
