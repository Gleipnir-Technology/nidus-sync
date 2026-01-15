-- UserByUsername
SELECT * FROM public.user_ WHERE
	username = $1 AND
	password_hash_type = 'bcrypt-14';
