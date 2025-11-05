-- UserByUsername
SELECT * FROM user_ WHERE
	username = $1 AND
	password_hash_type = 'bcrypt-14';
