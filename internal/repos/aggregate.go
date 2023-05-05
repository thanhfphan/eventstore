package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/thanhfphan/eventstore/internal/domain/aggregate"
	"github.com/thanhfphan/eventstore/internal/models"
	"github.com/thanhfphan/eventstore/pkg/dtype"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/utils"
)

type AggregateRepo interface {
	ReadAggregateSnapshot(ctx context.Context, aggId string, version int) (aggregate.AggregateRoot, error)
}

type aggregateRepo struct {
	db *sql.DB
}

func NewAggregate(db *sql.DB) AggregateRepo {
	return &aggregateRepo{
		db: db,
	}
}

func (r *aggregateRepo) ReadAggregateSnapshot(ctx context.Context, aggId string, version int) (aggregate.AggregateRoot, error) {

	query := r.db.QueryRowContext(ctx, `
			SELECT  a."type" as aggregate_type, s."data"
			FROM es_aggregate_snapshot s 
			JOIN es_aggregate a ON s.aggregate_id = a.id
			WHERE s.aggregate_id = ?
				AND (? = -1 OR s."version" <= ?)
			ORDER BY s."version" DESC 
			LIMIT 1	`, aggId, version)

	record := &models.AggregateSnapshot{}
	err := query.Scan(&record.AggregateType, &record.Data)
	if err != nil {
		return nil, err
	}

	if record.AggregateType == string(dtype.AggregateTypeOrder) {
		result := &aggregate.OrderAggregate{
			Status:     record.Data["status"].(string),
			CustomerID: record.Data["customer_id"].(int64),
			Version:    record.Data["version"].(int),
		}

		result.Price, _ = utils.GetFloat(record.Data["price"])
		if val, ok := record.Data["placed_date"]; ok {
			result.PlacedDate = val.(*time.Time)
		}
		if val, ok := record.Data["accepted_date"]; ok {
			result.AcceptedDate = val.(*time.Time)
		}
		if val, ok := record.Data["cancelled_date"]; ok {
			result.CancelledDate = val.(*time.Time)
		}
		if val, ok := record.Data["completed_date"]; ok {
			result.CompletedDate = val.(*time.Time)
		}

		return result, nil
	}

	return nil, errors.New("not support aggregateType=%s", record.AggregateType)
}
