INSERT INTO compliance_report_request(created, creator, id, public_id, site_id)
VALUES (NOW(), :user_id, DEFAULT, :public_id, :site_id);


-- INSERT INTO compliance_report_request (created, creator, public_id, site_id, site_version)
-- SELECT 
    -- NOW(),
    -- 1,
    -- generate_alphanumeric_code(8),
    -- id,
    -- version
-- FROM site;
