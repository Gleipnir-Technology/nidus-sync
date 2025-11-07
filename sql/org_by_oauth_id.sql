-- OrgByOauthId
SELECT o.id AS organization_id, o.arcgis_id AS arcgis_id, o.fieldseeker_url AS fieldseeker_url
FROM oauth_token ot
JOIN user_ u ON ot.user_id = u.id
JOIN organization o ON u.organization_id = o.id
WHERE ot.id = $1
