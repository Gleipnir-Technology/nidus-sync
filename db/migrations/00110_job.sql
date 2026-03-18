-- +goose Up
CREATE TYPE JobType AS ENUM (
	'audio-transcode',
	'csv-commit',
	'csv-import',
	'label-studio-audio-create',
	'email-send',
	'text-respond',
	'text-send'
);
ALTER TYPE comms.TextJobType ADD VALUE 'report-message' AFTER 'report-confirmation';
CREATE TABLE job (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	type_ JobType NOT NULL,
	row_id INTEGER NOT NULL,
	PRIMARY KEY(id)
);
COMMENT ON TABLE job IS 'A temporary holding place for jobs that are pushed to backend workers. Once work is completed the job should be deleted';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_new_job()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_job', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER job_insert_trigger
    AFTER INSERT ON job
    FOR EACH ROW
    EXECUTE FUNCTION notify_new_job();

-- +goose Down
DROP TRIGGER job_insert_trigger ON job;
DROP TABLE job;
-- ALTER TYPE comms.TextJobType DROP VALUE 'report-message';
DROP TYPE JobType;

