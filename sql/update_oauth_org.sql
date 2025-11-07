-- UpdateOauthTokenOrg
UPDATE oauth_token SET arcgis_id = $1, arcgis_license_type_id = $2
	WHERE refresh_token = $3;
