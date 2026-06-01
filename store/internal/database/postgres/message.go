package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"store/internal/model"
)

func (d *Database) QueryMessagesByCursor(ctx context.Context, request model.MessagesByCursorQueryRequest) (model.MessagesByCursorQueryResponse, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if request.IsFirstOne() {
		query := `
			select id, tenant_id, external_id, message, timestamp
			from store_message where tenant_id = $1
			order by id desc
			limit $2;`
		rows, err = d.db.QueryContext(ctx, query, request.TenantID, request.Limit)
	} else {
		query := `
			select id, tenant_id, external_id, message, timestamp
			from store_message
			where tenant_id = $1 and id < $2
			order by id desc
			limit $3;`
		rows, err = d.db.QueryContext(ctx, query, request.TenantID, request.Continue, request.Limit)
	}
	if err != nil {
		return model.MessagesByCursorQueryResponse{}, err
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}
	messages := make([]model.Message, 0)
	for rows.Next() {
		message := model.Message{}
		err = rows.Scan(
			&message.ID,
			&message.TenantID,
			&message.ExternalID,
			&message.Message,
			&message.Timestamp,
		)
		if err != nil {
			return model.MessagesByCursorQueryResponse{}, err
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return model.MessagesByCursorQueryResponse{}, err
	}
	var nextContinue string
	if len(messages) == request.Limit {
		nextContinue = messages[len(messages)-1].ID
	}
	return model.MessagesByCursorQueryResponse{
		Messages: messages,
		Continue: nextContinue,
	}, nil
}

func (d *Database) QueryMessagesByOffset(ctx context.Context, request model.MessageByOffsetQueryRequest) (model.MessageByOffsetQueryResponse, error) {
	const query = `
		select id, tenant_id, external_id, message, timestamp, count(1) over()
		from store_message where tenant_id = $1
		order by id desc
		limit $2 offset $3;`

	limit, offset := calculateLimitAndOffset(request.Page, request.PageSize)
	rows, err := d.db.QueryContext(ctx, query, request.TenantID, limit, offset)
	if err != nil {
		return model.MessageByOffsetQueryResponse{}, err
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}
	var totalCount int
	messages := make([]model.Message, 0)
	for rows.Next() {
		message := model.Message{}
		err = rows.Scan(
			&message.ID,
			&message.TenantID,
			&message.ExternalID,
			&message.Message,
			&message.Timestamp,
			&totalCount,
		)
		if err != nil {
			return model.MessageByOffsetQueryResponse{}, err
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return model.MessageByOffsetQueryResponse{}, err
	}
	return model.MessageByOffsetQueryResponse{
		Messages:   messages,
		TotalPages: calculateTotalPages(totalCount, request.PageSize),
		TotalCount: totalCount,
	}, nil
}

func calculateLimitAndOffset(page, pageSize int) (limit, offset int) {
	if page < 0 || pageSize < 0 {
		return 0, 0
	}
	return pageSize, (page - 1) * pageSize
}

func calculateTotalPages(totalCount, pageSize int) int {
	if totalCount == 0 || pageSize == 0 {
		return 1
	}
	return (totalCount + pageSize - 1) / pageSize
}

const (
	lockKindUpdateMessages = "update_messages"
)

func (d *Database) UpdateMessages(ctx context.Context, request model.UpdateMessagesRequest) (model.UpdateMessagesResponse, error) {
	const (
		query = `
			insert into store_message (tenant_id, external_id, message, timestamp) 
			values %s 
			on conflict do nothing`
		fieldsCount = 4
	)

	valueStrings := make([]string, 0, len(request.Messages))
	valueArgs := make([]interface{}, 0, len(request.Messages)*fieldsCount)

	var messageInd int

	for _, message := range request.Messages {
		valueString := fmt.Sprintf(
			"($%d, $%d, $%d, $%d)",
			messageInd*fieldsCount+1,
			messageInd*fieldsCount+2,
			messageInd*fieldsCount+3,
			messageInd*fieldsCount+4,
		)
		valueStrings = append(valueStrings, valueString)
		valueArgs = append(valueArgs,
			message.TenantID,
			message.ExternalID,
			message.Message,
			message.Timestamp,
		)
		messageInd++
	}

	var recordsChanged int64

	err := func() error {
		tx, err := d.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback()
		}()

		_, err = tx.ExecContext(
			ctx,
			`select 1 from store_lock where kind = $1 for update`,
			lockKindUpdateMessages,
		)
		if err != nil {
			return err
		}

		stmt := fmt.Sprintf(query, strings.Join(valueStrings, ","))
		stmtRes, err := tx.ExecContext(ctx, stmt, valueArgs...)
		if err != nil {
			return err
		}

		recordsChanged, _ = stmtRes.RowsAffected()

		return tx.Commit()
	}()

	return model.UpdateMessagesResponse{MessagesUpdated: recordsChanged}, err
}
