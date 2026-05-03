package logic

import (
	"context"
	"encoding/json"

	"ordering/rpc/internal/svc"
	"ordering/rpc/ordering"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PlaceOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlaceOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlaceOrderLogic {
	return &PlaceOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PlaceOrderLogic) PlaceOrder(in *ordering.Request) (*ordering.Response, error) {
	orderId := uuid.New().String()

	orderPayload := &orderPayload{
		OrderId:   orderId,
		ProductId: in.ProductId,
	}

	payload, err := json.Marshal(&orderPayload)

	if err != nil {
		return nil, err
	}

	err = l.svcCtx.SqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(
			ctx,
			`INSERT INTO orders (uuid, product_id) VALUES ($1, $2)`,
			orderId,
			in.ProductId,
		); err != nil {
			return err
		}

		if _, err := session.ExecCtx(
			ctx,
			`INSERT INTO outbox (id, topic, payload) VALUES ($1, $2, $3)`,
			orderId,
			"order_created",
			payload,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ordering.Response{
		OrderId: orderId,
	}, nil
}
