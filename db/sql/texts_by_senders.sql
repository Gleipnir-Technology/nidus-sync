-- TextsBySenders
SELECT 
    id,
    content,
    created,
    source,
    destination,
    is_visible_to_llm,
    is_welcome,
    origin
FROM 
    comms.text_log
WHERE 
    (source = $1 AND destination = $2)
    OR 
    (source = $2 AND destination = $1)
ORDER BY 
    created ASC;
