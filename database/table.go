package database

import (
	"github.com/MihaPecnik/order-maching-system/models"
	"gorm.io/gorm"
)

func (d *Database) GetBottomBuy(ticker string) (models.UpdateOrderBookResponse, error) {
	var response models.UpdateOrderBookResponse
	err := d.db.Table("tables").
		Select("value", "quantity").
		Where("ticker = ? AND buy = ?", ticker, true).
		Order("value desc").
		First(&response).Error
	if err != nil {
		return models.UpdateOrderBookResponse{}, err
	}
	return response, err
}

func (d *Database) GetTopSell(ticker string) (models.UpdateOrderBookResponse, error) {
	var response models.UpdateOrderBookResponse
	err := d.db.Table("tables").
		Select("value", "quantity").
		Where("ticker = ? AND buy = ?", ticker, false).
		Order("value asc").
		First(&response).Error
	if err != nil {
		return models.UpdateOrderBookResponse{}, err
	}
	return response, err
}

func (d *Database) UpdateOrdersBook(request models.UpdateOrderBookRequest) ([]models.UpdateOrderBookResponse, error) {
	query := `
with ordersUsable as (
  select * , sum(quantity) over (order by value) as qu
  from tables
  where value < ? and buy = ? and ticker = ?
)
select *
from (
       select *
       from ordersUsable
       where qu >= ?
       limit 1
     ) as ordersAdditional
union
select *
from ordersUsable
where  qu < ?
order by value asc ;
`
	if !request.Buy {
		query = `
with ordersUsable as (
  select * , sum(quantity) over (order by value) as qu
  from tables
  where value > ? and buy = ? and ticker = ?
)
select *
from (
       select *
       from ordersUsable
       where qu >= ?
       limit 1
     ) as ordersAdditional
union
select *
from ordersUsable
where  qu < ?
order by value desc ;
`
	}
	stmt := d.db.Session(&gorm.Session{
		PrepareStmt: true,
	})
	response := []models.UpdateOrderBookResponse{}

	// If any error is return, transaction will take care of a rollback
	err := stmt.Transaction(func(tx *gorm.DB) error {
		var suitableOrders []models.Table

		// Get all orders that are suitable for our request
		err := tx.
			Raw(query, request.Value, !request.Buy, request.Ticker, request.Quantity, request.Quantity).
			Scan(&suitableOrders).Error
		if err != nil {
			return err
		}

		for _, order := range suitableOrders {
			// If request has bigger quantity, we can delete it (execute the order)
			// If request has smaller quantity, we update it's quantity (partly execute the order)
			if request.Quantity >= order.Quantity {
				err = tx.Where("id = ?", order.ID).Delete(&models.Table{}).Error
				if err != nil {
					return err
				}

				request.Quantity -= order.Quantity
				response = append(response, models.UpdateOrderBookResponse{
					Value:    order.Value,
					Quantity: order.Quantity,
				})
			} else {
				q := order.Quantity - request.Quantity
				err = tx.Model(&models.Table{}).Where("id = ?", order.ID).
					Updates(models.Table{
						Quantity: q,
					}).Error
				if err != nil {
					return err
				}

				response = append(response, models.UpdateOrderBookResponse{
					Value:    order.Value,
					Quantity: request.Quantity,
				})
			}
		}
		// If there is not enough suitable offers, we create new order
		if request.Quantity > 0 {
			err := tx.Create(&models.Table{
				UserId:   request.UserId,
				Value:    request.Value,
				Ticker:   request.Ticker,
				Quantity: request.Quantity,
				Buy:      request.Buy,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return response, err
}