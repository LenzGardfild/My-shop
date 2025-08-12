package db

import (
	"context"

	"my-shop/internal/model"
)

func (db *DB) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	order := &model.Order{}
	// Получаем основной заказ
	row := db.Pool.QueryRow(ctx, `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid=$1`, orderUID)
	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, err
	}
	// Доставка
	row = db.Pool.QueryRow(ctx, `SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid=$1`, orderUID)
	err = row.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, err
	}
	// Оплата
	row = db.Pool.QueryRow(ctx, `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE order_uid=$1`, orderUID)
	err = row.Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, err
	}
	// Товары
	rows, err := db.Pool.Query(ctx, `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid=$1`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}

func (db *DB) SaveOrder(ctx context.Context, order *model.Order) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (order_uid) DO NOTHING`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (order_uid) DO NOTHING`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (order_uid) DO NOTHING`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (db *DB) LoadAllOrders(ctx context.Context) ([]*model.Order, error) {
	orders := []*model.Order{}
	rows, err := db.Pool.Query(ctx, `SELECT order_uid FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		order, err := db.GetOrder(ctx, id)
		if err == nil && order != nil {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
