package repository

import (
	"context"
	"errors"
	"fmt"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"

	"entgo.io/ent/dialect"
)

// api_key_groups 关联表（手动模式的有序分组列表）的读写。
//
// 用原生 SQL 管理，不改 ent schema（见 add-api-key-multi-group-routing design D1）。
// 约束：position 0 恒等于 api_keys.group_id（主路由），只读 group_id 的旧代码继续成立。
// 读取时若行缺失或首位与 group_id 不一致，退回只绑 group_id，兼容迁移前的数据。

// supportsAPIKeyGroups 报告当前驱动是否有 api_key_groups 表。
// 该表只由 Postgres 迁移创建；其他方言（SQLite 单测）退回只认 group_id 的旧行为。
func (r *apiKeyRepository) supportsAPIKeyGroups() bool {
	return r.client != nil && r.client.Driver().Dialect() == dialect.Postgres
}

// withTx 在已有事务内直接复用，否则开启新事务。
func (r *apiKeyRepository) withTx(ctx context.Context, fn func(txCtx context.Context, client *dbent.Client) error) error {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return fn(ctx, tx.Client())
	}
	tx, err := r.client.Tx(ctx)
	if errors.Is(err, dbent.ErrTxStarted) {
		// client 本身已是事务客户端（如测试夹具或外层事务传入），直接复用，不嵌套。
		return fn(ctx, r.client)
	}
	if err != nil {
		return fmt.Errorf("begin api key transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := fn(dbent.NewTxContext(ctx, tx), tx.Client()); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit api key transaction: %w", err)
	}
	return nil
}

// writeAPIKeyGroups 把密钥的分组列表整体重写为 groupIDs 的顺序（调用方负责去重与校验）。
func writeAPIKeyGroups(ctx context.Context, exec *dbent.Client, apiKeyID int64, groupIDs []int64) error {
	if _, err := exec.ExecContext(ctx, `DELETE FROM api_key_groups WHERE api_key_id = $1`, apiKeyID); err != nil {
		return fmt.Errorf("delete api key groups: %w", err)
	}
	if len(groupIDs) == 0 {
		return nil
	}
	if _, err := exec.ExecContext(ctx, `
		INSERT INTO api_key_groups (api_key_id, group_id, position)
		SELECT $1, t.group_id, (t.ord - 1)::int
		FROM unnest($2::bigint[]) WITH ORDINALITY AS t(group_id, ord)`,
		apiKeyID, pq.Array(groupIDs)); err != nil {
		return fmt.Errorf("insert api key groups: %w", err)
	}
	return nil
}

// renumberAPIKeyGroups 把指定密钥的 position 压紧为 0..n-1，并把 api_keys.group_id
// 同步为新的首位。没有剩余行的密钥不动 group_id（与删除单分组时的既有行为一致）。
func renumberAPIKeyGroups(ctx context.Context, exec *dbent.Client, apiKeyIDs []int64) error {
	if len(apiKeyIDs) == 0 {
		return nil
	}
	// (api_key_id, position) 唯一约束不可延迟，逐行移动会和尚未移动的行撞车，
	// 所以先整体挪到负数区，再写回目标值。
	if _, err := exec.ExecContext(ctx, `
		UPDATE api_key_groups g
		SET position = -(r.rn + 1)
		FROM (
			SELECT api_key_id, group_id,
				ROW_NUMBER() OVER (PARTITION BY api_key_id ORDER BY position) - 1 AS rn
			FROM api_key_groups
			WHERE api_key_id = ANY($1)
		) r
		WHERE g.api_key_id = r.api_key_id AND g.group_id = r.group_id`,
		pq.Array(apiKeyIDs)); err != nil {
		return fmt.Errorf("renumber api key groups: %w", err)
	}
	if _, err := exec.ExecContext(ctx, `
		UPDATE api_key_groups SET position = -position - 1
		WHERE api_key_id = ANY($1) AND position < 0`,
		pq.Array(apiKeyIDs)); err != nil {
		return fmt.Errorf("renumber api key groups: %w", err)
	}
	if _, err := exec.ExecContext(ctx, `
		UPDATE api_keys k
		SET group_id = g.group_id, updated_at = NOW()
		FROM api_key_groups g
		WHERE g.api_key_id = k.id AND g.position = 0
			AND k.id = ANY($1) AND k.deleted_at IS NULL
			AND k.group_id IS DISTINCT FROM g.group_id`,
		pq.Array(apiKeyIDs)); err != nil {
		return fmt.Errorf("sync api key primary group: %w", err)
	}
	return nil
}

// removeGroupFromAPIKeys 删除分组时清理所有密钥里的该分组，并重排剩余分组。
//
// 分组是软删除，外键级联不会触发，必须由 deleteCascade 显式调用。
func removeGroupFromAPIKeys(ctx context.Context, exec *dbent.Client, groupID int64) error {
	if exec.Driver().Dialect() != dialect.Postgres {
		return nil
	}
	rows, err := exec.QueryContext(ctx,
		`DELETE FROM api_key_groups WHERE group_id = $1 RETURNING api_key_id`, groupID)
	if err != nil {
		return fmt.Errorf("delete api key groups by group: %w", err)
	}
	affected, err := scanInt64Rows(rows)
	if err != nil {
		return err
	}
	return renumberAPIKeyGroups(ctx, exec, affected)
}

// replaceGroupInAPIKeys 把用户密钥里的 oldGroupID 换成 newGroupID：
// 已含 newGroupID 的密钥直接去掉 oldGroupID，其余原位替换。
func replaceGroupInAPIKeys(ctx context.Context, exec *dbent.Client, userID, oldGroupID, newGroupID int64) error {
	rows, err := exec.QueryContext(ctx, `
		SELECT g.api_key_id FROM api_key_groups g
		JOIN api_keys k ON k.id = g.api_key_id
		WHERE k.user_id = $1 AND k.deleted_at IS NULL AND g.group_id = $2`,
		userID, oldGroupID)
	if err != nil {
		return fmt.Errorf("query api key groups by user: %w", err)
	}
	affected, err := scanInt64Rows(rows)
	if err != nil {
		return err
	}
	if len(affected) == 0 {
		return nil
	}
	if _, err := exec.ExecContext(ctx, `
		DELETE FROM api_key_groups g
		WHERE g.api_key_id = ANY($1) AND g.group_id = $2
			AND EXISTS (SELECT 1 FROM api_key_groups n WHERE n.api_key_id = g.api_key_id AND n.group_id = $3)`,
		pq.Array(affected), oldGroupID, newGroupID); err != nil {
		return fmt.Errorf("drop replaced api key group: %w", err)
	}
	if _, err := exec.ExecContext(ctx, `
		UPDATE api_key_groups SET group_id = $3
		WHERE api_key_id = ANY($1) AND group_id = $2`,
		pq.Array(affected), oldGroupID, newGroupID); err != nil {
		return fmt.Errorf("replace api key group: %w", err)
	}
	return renumberAPIKeyGroups(ctx, exec, affected)
}

func scanInt64Rows(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}) ([]int64, error) {
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return nil, err
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	return out, rows.Close()
}

// listGroupIDsByAPIKeyIDs 一次查出多把密钥的有序分组 ID，避免列表页 N+1。
func (r *apiKeyRepository) listGroupIDsByAPIKeyIDs(ctx context.Context, apiKeyIDs []int64) (map[int64][]int64, error) {
	out := make(map[int64][]int64, len(apiKeyIDs))
	if len(apiKeyIDs) == 0 || !r.supportsAPIKeyGroups() {
		return out, nil
	}
	exec := clientFromContext(ctx, r.client)
	rows, err := exec.QueryContext(ctx, `
		SELECT api_key_id, group_id FROM api_key_groups
		WHERE api_key_id = ANY($1)
		ORDER BY api_key_id, position`, pq.Array(apiKeyIDs))
	if err != nil {
		return nil, fmt.Errorf("query api key groups: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var apiKeyID, groupID int64
		if err := rows.Scan(&apiKeyID, &groupID); err != nil {
			return nil, fmt.Errorf("scan api key group: %w", err)
		}
		out[apiKeyID] = append(out[apiKeyID], groupID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate api key groups: %w", err)
	}
	return out, nil
}

// attachAPIKeyGroups 给密钥填上 GroupMode 之外的 GroupIDs / Groups。
//
// 只有绑了多个分组的密钥才填 Groups；Groups[0] 复用已预加载的 Group。
// 已软删除的分组被 ent 拦截器过滤掉，不会出现在 Groups 里。
func (r *apiKeyRepository) attachAPIKeyGroups(ctx context.Context, keys []*service.APIKey) error {
	ids := make([]int64, 0, len(keys))
	for _, k := range keys {
		if k != nil && k.GroupID != nil {
			ids = append(ids, k.ID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	byKey, err := r.listGroupIDsByAPIKeyIDs(ctx, ids)
	if err != nil {
		return err
	}

	extra := make(map[int64]struct{})
	for _, k := range keys {
		if k == nil || k.GroupID == nil {
			continue
		}
		groupIDs := byKey[k.ID]
		if len(groupIDs) == 0 || groupIDs[0] != *k.GroupID {
			// 旧数据或不一致：按只绑 group_id 处理。
			continue
		}
		k.GroupIDs = groupIDs
		if len(groupIDs) > 1 {
			for _, id := range groupIDs[1:] {
				extra[id] = struct{}{}
			}
		}
	}
	if len(extra) == 0 {
		return nil
	}

	extraIDs := make([]int64, 0, len(extra))
	for id := range extra {
		extraIDs = append(extraIDs, id)
	}
	ents, err := clientFromContext(ctx, r.client).Group.Query().Where(group.IDIn(extraIDs...)).All(ctx)
	if err != nil {
		return fmt.Errorf("load api key groups: %w", err)
	}
	groups := make(map[int64]*service.Group, len(ents))
	for _, g := range ents {
		groups[g.ID] = groupEntityToService(g)
	}

	for _, k := range keys {
		if k == nil || len(k.GroupIDs) < 2 {
			continue
		}
		k.Groups = make([]*service.Group, 0, len(k.GroupIDs))
		for i, id := range k.GroupIDs {
			if i == 0 && k.Group != nil && k.Group.ID == id {
				k.Groups = append(k.Groups, k.Group)
				continue
			}
			if g, ok := groups[id]; ok {
				k.Groups = append(k.Groups, g)
			}
		}
	}
	return nil
}

func (r *apiKeyRepository) attachAPIKeyGroupsToSlice(ctx context.Context, keys []service.APIKey) error {
	ptrs := make([]*service.APIKey, len(keys))
	for i := range keys {
		ptrs[i] = &keys[i]
	}
	return r.attachAPIKeyGroups(ctx, ptrs)
}
