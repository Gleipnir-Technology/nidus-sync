BEGIN TRANSACTION;
	DELETE FROM fileupload.pool;
	DELETE FROM fileupload.error_csv;
	DELETE FROM fileupload.csv;
	DELETE FROM fileupload.error_file;
	DELETE FROM lead WHERE site_id IN (SELECT id FROM SITE WHERE file_id IS NOT NULL);
	DELETE FROM site WHERE file_id IS NOT NULL;
	DELETE FROM review_task_pool;
	DELETE FROM review_task;
	DELETE FROM fileupload.file;
COMMIT;

