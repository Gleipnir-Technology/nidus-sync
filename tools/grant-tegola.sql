CREATE ROLE "tegola";
ALTER ROLE "tegola" WITH LOGIN;
GRANT SELECT ON import.district TO tegola;
GRANT USAGE on SCHEMA import TO "tegola";
