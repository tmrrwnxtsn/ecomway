DROP INDEX IF EXISTS ix_operation_user_id;
DROP INDEX IF EXISTS ix_operation_external_id;
DROP INDEX IF EXISTS ix_operation_status;
DROP INDEX IF EXISTS ix_operation_created_at;

DROP TABLE IF EXISTS tool;
DROP TABLE IF EXISTS operation;
DROP TABLE IF EXISTS operation_metadata;
DROP TABLE IF EXISTS "user";
