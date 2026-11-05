-- OAuth 2.0 授权服务器（只服务一个预配置的客户端：视频工作台）。
-- 授权码与访问令牌都只存 SHA-256，明文只在签发瞬间出现在响应里。

CREATE TABLE IF NOT EXISTS oauth_authorization_codes (
    id                    BIGSERIAL PRIMARY KEY,
    code_hash             VARCHAR(64)  NOT NULL UNIQUE,
    client_id             VARCHAR(128) NOT NULL,
    user_id               BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    redirect_uri          TEXT         NOT NULL,
    code_challenge        VARCHAR(128) NOT NULL,
    expires_at            TIMESTAMPTZ  NOT NULL,
    consumed_at           TIMESTAMPTZ,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_oauth_codes_expires_at ON oauth_authorization_codes (expires_at);

CREATE TABLE IF NOT EXISTS oauth_access_tokens (
    id                    BIGSERIAL PRIMARY KEY,
    token_hash            VARCHAR(64)  NOT NULL UNIQUE,
    client_id             VARCHAR(128) NOT NULL,
    user_id               BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- 由哪个授权码签发：授权码被重放时据此作废令牌（RFC 6749 §4.1.2）
    authorization_code_id BIGINT       REFERENCES oauth_authorization_codes(id) ON DELETE CASCADE,
    expires_at            TIMESTAMPTZ  NOT NULL,
    revoked_at            TIMESTAMPTZ,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_oauth_tokens_expires_at ON oauth_access_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_oauth_tokens_code ON oauth_access_tokens (authorization_code_id);

-- 用户对客户端的授权记录：授权过一次之后不再弹确认页
CREATE TABLE IF NOT EXISTS oauth_consents (
    user_id               BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id             VARCHAR(128) NOT NULL,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, client_id)
);

-- 外部工作台托管的专用密钥：哪个用户、哪个客户端、对应哪把 API Key。
-- API Key 本体仍在 api_keys 里，计费、限速、用量记录完全走现有链路。
CREATE TABLE IF NOT EXISTS oauth_client_api_keys (
    user_id               BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id             VARCHAR(128) NOT NULL,
    api_key_id            BIGINT       NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, client_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_oauth_client_api_keys_key ON oauth_client_api_keys (api_key_id);

COMMENT ON TABLE oauth_client_api_keys IS 'OAuth 客户端（如视频工作台）为用户托管的专用 API Key';
