package operator_ent_provider

import (
	"context"
	"fmt"
	biz_operator "via/internal/biz/operator"
	ent_client "via/internal/client/ent"
	"via/internal/ent"
	"via/internal/ent/operator"
	"via/internal/log"
	"via/internal/model"

	"entgo.io/ent/dialect/sql"
)

type OperatorEntProvider struct {
	client *ent.Client
}

func New() OperatorEntProvider {
	return OperatorEntProvider{ent_client.Get()}
}

func fromEntOperator(operator ent.Operator) model.Operator {
	return model.Operator{
		ID:      operator.ID,
		Name:    operator.Name,
		Account: operator.Account,
		Enabled: operator.Enabled,
	}
}

func (o OperatorEntProvider) GetOperators(ctx context.Context) ([]model.Operator, error) {
	operators := []model.Operator{}
	ops, err := o.client.Operator.
		Query().
		Where(operator.IDNEQ(biz_operator.OPERATOR_SYSTEM)).
		Order(operator.ByID(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed getting Operators")
		return operators, fmt.Errorf("failed getting Operators: %w", err)
	}
	for _, op := range ops {
		operators = append(operators, fromEntOperator(*op))
	}
	return operators, err
}

func (o OperatorEntProvider) GetOperatorByAccount(ctx context.Context, account string) (model.Operator, error) {
	operator, err := o.client.Operator.
		Query().
		Where(operator.AccountEQ(account)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Operator{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying operator by account", "account", account)
		return model.Operator{}, fmt.Errorf("failed querying operator by account: %w", err)

	}
	return fromEntOperator(*operator), err
}

func (o OperatorEntProvider) GetOperatorById(ctx context.Context, id int) (model.Operator, error) {
	operator, err := o.client.Operator.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Operator{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying operator by id", "operator_id", id)
		return model.Operator{}, fmt.Errorf("failed querying operator by id: %w", err)

	}
	return fromEntOperator(*operator), err
}
