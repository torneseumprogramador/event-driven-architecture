package projections

import (
	"context"
	"time"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductProjectionView representa a projeção de produto no MongoDB
type ProductProjectionView struct {
	ID        int       `bson:"_id"`
	Name      string    `bson:"name"`
	Price     float64   `bson:"price"`
	Stock     int       `bson:"stock"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// ProductProjection gerencia as projeções de produto
type ProductProjection struct {
	collection *mongo.Collection
}

// NewProductProjection cria uma nova projeção de produto
func NewProductProjection(db *mongo.Database) *ProductProjection {
	return &ProductProjection{
		collection: db.Collection("views.products"),
	}
}

// HandleProductCreated processa evento de produto criado
func (p *ProductProjection) HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error {
	productView := ProductProjectionView{
		ID:        int(event.Product.ID),
		Name:      event.Product.Name,
		Price:     event.Product.Price,
		Stock:     event.Product.Stock,
		CreatedAt: event.OccurredAt,
		UpdatedAt: time.Now(),
	}
	
	// Usa ReplaceOne com upsert para evitar erro de chave duplicada
	filter := bson.M{"_id": event.Product.ID}
	_, err := p.collection.ReplaceOne(ctx, filter, productView, options.Replace().SetUpsert(true))
	return err
}

// HandleProductUpdated processa evento de produto atualizado
func (p *ProductProjection) HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error {
	filter := bson.M{"_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       event.Product.Name,
			"price":      event.Product.Price,
			"stock":      event.Product.Stock,
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// HandleStockReserved processa evento de estoque reservado
func (p *ProductProjection) HandleStockReserved(ctx context.Context, event pkgevents.StockReserved) error {
	filter := bson.M{"_id": event.ProductID}
	update := bson.M{
		"$inc": bson.M{
			"stock": -event.Quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// HandleStockReleased processa evento de estoque liberado
func (p *ProductProjection) HandleStockReleased(ctx context.Context, event pkgevents.StockReleased) error {
	filter := bson.M{"_id": event.ProductID}
	update := bson.M{
		"$inc": bson.M{
			"stock": event.Quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAll busca todos os produtos
func (p *ProductProjection) GetAll(ctx context.Context) ([]ProductProjectionView, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	
	cursor, err := p.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var products []ProductProjectionView
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	
	return products, nil
}

// GetByID busca produto por ID
func (p *ProductProjection) GetByID(ctx context.Context, id int) (*ProductProjectionView, error) {
	var product ProductProjectionView
	err := p.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
