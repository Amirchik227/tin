package database

import (
	"database/sql"
	"fmt"
	"projectZero/storage"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "amir"
	password = "1234"
	dbname   = "orders"
)

var db *sql.DB

func Close_db() {
	db.Close()
}

func GetCache(s *storage.MemoryStorage)  { 
	// orders
	selectOrderQuery := `
	SELECT  o.id, o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, 
        d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
	 	p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders AS o
		INNER JOIN deliveries AS d USING(id)
		INNER JOIN payments AS p USING(id)`

	rows, err := db.Query(selectOrderQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var orders []storage.Order
	for rows.Next() {
		var order storage.Order
		err := rows.Scan(&order.ID, &order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	// items
	selectItemsQuery := `
	SELECT order_id, chrt_id, track_number, price, 
		rid, name, sale, size, total_price, nm_id, brand, status
	FROM items`
	rows, err = db.Query(selectItemsQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var items []storage.Item
	for rows.Next() {
		var item storage.Item

		err := rows.Scan(&item.OrderID, &item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	
	//add orders to cache
	for _, ord := range orders {
		id := ord.ID
		var newItems []storage.Item
		for _, itm := range items {
			if itm.OrderID == id {
				newItems = append(newItems, itm)
			}
		}
		ord.Items = newItems
		s.Insert(&ord)
	}
	fmt.Println("Cache recieved")
}

func ConnectDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dataBase, err := sql.Open("postgres", psqlInfo)
	fmt.Println()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully connected to %s database \n", dbname)
	db = dataBase
}

func InsertOrder(o storage.Order, orderId int) {
	ordersQuery := `
	INSERT INTO orders (id, order_uid, track_number, entry,
		locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12)
	RETURNING id`
	var id int
	err := db.QueryRow(ordersQuery, orderId, o.OrderUID, o.TrackNumber,
		o.Entry, o.Locale,
		o.InternalSignature, o.CustomerID, o.DeliveryService,
		o.Shardkey, o.SmID, o.DateCreated, o.OofShard).Scan(&id)
	if err != nil {
		panic(err)
	}

	d := o.Delivery
	deliveriesQuery := `
	INSERT INTO deliveries (id, name, phone, zip, city, 
		address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

	err = db.QueryRow(deliveriesQuery, orderId, d.Name, d.Phone,
		d.Zip, d.City, d.Address, d.Region, d.Email).Scan(&id)
	if err != nil {
		panic(err)
	}

	p := o.Payment
	paymentsQuery := `
	INSERT INTO payments ( id, transaction, request_id,
		currency, provider, amount, payment_dt, bank,
		delivery_cost,	goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id`
	err = db.QueryRow(paymentsQuery, orderId, p.Transaction,
		p.RequestId, p.Currency, p.Provider, p.Amount,
		p.PaymentDt, p.Bank, p.DeliveryCost, p.GoodsTotal,
		p.CustomFee).Scan(&id)
	if err != nil {
		panic(err)
	}

	i := o.Items
	itemsQuery := `
	INSERT INTO items (order_id, chrt_id, track_number, price, 
		rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id`

	for idx := 0; idx < len(i); idx++ {
		err = db.QueryRow(itemsQuery, orderId,
			i[idx].ChrtID, i[idx].TrackNumber, i[idx].Price,
			i[idx].RID, i[idx].Name, i[idx].Sale, i[idx].Size,
			i[idx].TotalPrice, i[idx].NmID, i[idx].Brand,
			i[idx].Status).Scan(&id)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Order added to database\n")
}
