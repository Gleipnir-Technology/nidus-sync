INSERT INTO compliance_report_request(created, creator, id, public_id, site_id, site_version)
VALUES (NOW(), :user_id, DEFAULT, :public_id, :site_id, 1);


-- INSERT INTO compliance_report_request (created, creator, public_id, site_id, site_version)
-- SELECT 
    -- NOW(),
    -- 1,
    -- generate_alphanumeric_code(8),
    -- id,
    -- version
-- FROM site;
