package storage

import (
	"errors"
	"fmt"
	"sync"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"` 
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	OrderID     int    `json:"order_id"` 
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	ID 			      int	   `json:"id"`
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"` 
	Payment           Payment  `json:"payment"`  
	Items             []Item   `json:"items"`    
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Storage interface {
	Insert(o *Order) int
	Get(id int) (Order, error)
	// Update(id int, o Order)
	// Delete(id int)
}

type MemoryStorage struct {
	counter int
	data    map[int]Order
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[int]Order),
		counter: 1,
	}
}

func (s *MemoryStorage) Insert(o *Order) int { 
	s.Lock()
	o.ID = s.counter
	s.data[s.counter] = *o
	s.counter++
	fmt.Printf("Order added to cache\n")
	s.Unlock()
	return s.counter - 1

}

func (s *MemoryStorage) Get(id int) (Order, error) {
	s.Lock()
	defer s.Unlock()

	order, ok := s.data[id]
	if !ok {
		return order, errors.New("order not found")
	}

	return order, nil
}

func (s *MemoryStorage) ShowCache() {
	fmt.Printf("\n\nShowCache")
	for key, order := range s.data {
		fmt.Printf("\n%d - %+v\n", key, order)
		for k, item := range order.Items {
			fmt.Printf("\n       - item %d - %+v\n", k, item)
		}
	}
}

func (s *MemoryStorage) Update(id int, o Order) {
	s.Lock()
	s.data[id] = o
	s.Unlock()
}

func (s *MemoryStorage) Delete(id int) {
	s.Lock()
	delete(s.data, id)
	s.Unlock()
}
