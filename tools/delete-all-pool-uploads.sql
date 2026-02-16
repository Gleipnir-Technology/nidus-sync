BEGIN TRANSACTION;
	DELETE FROM fileupload.pool;
	DELETE FROM fileupload.error_csv;
	DELETE FROM fileupload.csv;
	DELETE FROM fileupload.error_file;
	DELETE FROM fileupload.file;
COMMIT;

