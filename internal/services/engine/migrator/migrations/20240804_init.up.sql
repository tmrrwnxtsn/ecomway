CREATE TABLE IF NOT EXISTS operation
(
    id              SERIAL PRIMARY KEY,
    user_id         BIGINT                   NOT NULL,
    type            VARCHAR(255)             NOT NULL,
    currency        VARCHAR(255)             NOT NULL,
    amount          NUMERIC(20, 8)           NOT NULL,
    status          VARCHAR(255)             NOT NULL,
    external_id     VARCHAR(255),
    external_system VARCHAR(255)             NOT NULL,
    external_method VARCHAR(255)             NOT NULL,
    external_status VARCHAR(255),
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (external_method, external_id)
);

CREATE INDEX ix_operation_user_id ON operation (user_id);
CREATE INDEX ix_operation_external_id ON operation (external_id);
CREATE INDEX ix_operation_status ON operation (status);
CREATE INDEX ix_operation_created_at ON operation (created_at);

CREATE TABLE IF NOT EXISTS operation_metadata
(
    operation_id      BIGINT NOT NULL REFERENCES operation ON DELETE CASCADE,
    tool_id           VARCHAR(255),
    additional        JSONB,
    fail_reason       VARCHAR(255),
    confirmation_code VARCHAR(255),
    processed_at      TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS tool
(
    id              VARCHAR(255),
    user_id         BIGINT,
    external_method VARCHAR(255),
    type            VARCHAR(255),
    details         JSONB,
    displayed       VARCHAR(255)             NOT NULL,
    name            VARCHAR(255)             NOT NULL,
    status          VARCHAR(255)             NOT NULL DEFAULT 'ACTIVE',
    fake            BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, user_id, external_method)
);
