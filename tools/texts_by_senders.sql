-- TextsBySenders
SELECT 
    id,
    content,
    created,
    source,
    destination,
    is_welcome,
    origin
FROM 
    comms.text_log
WHERE 
    (source = :src AND destination = :dst)
    OR 
    (source = :dst AND destination = :src)
ORDER BY 
    created ASC;
