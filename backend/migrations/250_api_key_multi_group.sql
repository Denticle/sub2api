-- API Key 绑定多个有序分组。
-- position 0 是主路由，始终与 api_keys.group_id 一致；只认 group_id 的旧代码继续成立。
-- 分组是软删除，外键级联不会触发，删除分组时由 groupRepository.deleteCascade 显式清理并重排。

CREATE TABLE IF NOT EXISTS api_key_groups (
    api_key_id  BIGINT NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    group_id    BIGINT NOT NULL REFERENCES groups(id)   ON DELETE CASCADE,
    position    INT    NOT NULL,
    PRIMARY KEY (api_key_id, group_id),
    UNIQUE (api_key_id, position)
);

CREATE INDEX IF NOT EXISTS idx_api_key_groups_group ON api_key_groups (group_id);

COMMENT ON TABLE api_key_groups IS 'API Key 的有序分组列表（手动模式）；position 0 = 主路由 = api_keys.group_id';

-- manual：按 api_key_groups 顺序；auto：按用户当前权限实时计算（仅工作台密钥使用）
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS group_mode VARCHAR(16) NOT NULL DEFAULT 'manual';

COMMENT ON COLUMN api_keys.group_mode IS 'manual = api_key_groups 顺序；auto = 按用户权限实时计算';

-- 回填：每把有分组的密钥写一行 position 0
INSERT INTO api_key_groups (api_key_id, group_id, position)
SELECT k.id, k.group_id, 0
FROM api_keys k
WHERE k.group_id IS NOT NULL
  AND k.deleted_at IS NULL
ON CONFLICT DO NOTHING;

-- 工作台托管的密钥改为自动模式
UPDATE api_keys
SET group_mode = 'auto'
WHERE id IN (SELECT api_key_id FROM oauth_client_api_keys)
  AND group_mode <> 'auto';
