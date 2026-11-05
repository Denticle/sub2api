package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIKeyMultiGroupMigration(t *testing.T) {
	content, err := FS.ReadFile("250_api_key_multi_group.sql")
	require.NoError(t, err)

	sql := strings.Join(strings.Fields(string(content)), " ")
	// 可重复执行
	require.Contains(t, sql, "CREATE TABLE IF NOT EXISTS api_key_groups")
	require.Contains(t, sql, "CREATE INDEX IF NOT EXISTS idx_api_key_groups_group")
	require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS group_mode VARCHAR(16) NOT NULL DEFAULT 'manual'")
	require.Contains(t, sql, "UNIQUE (api_key_id, position)")
	// 回填 position 0，冲突时跳过
	require.Contains(t, sql, "SELECT k.id, k.group_id, 0 FROM api_keys k WHERE k.group_id IS NOT NULL")
	require.Contains(t, sql, "ON CONFLICT DO NOTHING")
	// 工作台密钥改为自动模式
	require.Contains(t, sql, "SET group_mode = 'auto' WHERE id IN (SELECT api_key_id FROM oauth_client_api_keys)")
}
