package mapper

type Orderpiece struct {
	ID    string `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Price int32  `bson:"price" json:"price"`
}

type CreateOrderRequest struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
}
