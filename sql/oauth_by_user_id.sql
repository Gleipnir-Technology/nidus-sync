-- OauthTokenByUserId
SELECT * FROM oauth_token WHERE
	user_id = $1;
