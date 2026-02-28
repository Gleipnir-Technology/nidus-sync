-- +goose Up
-- Hat tip to https://jesseamundsen.github.io/2025-11-02-geojsontogeom/
-- +goose StatementBegin
create function public.geojsontogeom(geojson jsonb)
returns table (geometrytype text, properties jsonb, geom geometry)
language plpgsql
as $function$
begin
    return query
    select f.features -> 'geometry' ->> 'type' geometrytype
        ,f.features -> 'properties' properties
        ,st_setsrid(st_geomfromgeojson(f.features ->> 'geometry'),4326) geometry
    from (
        select jsonb_array_elements(case
            when lower(geojson ->> 'type')='featurecollection' then geojson -> 'features'
            when lower(geojson ->> 'type')='feature' then jsonb_build_array(geojson)
            else jsonb_build_array(jsonb_build_object('type','Feature','geometry',geojson)) end) features
    ) f;
end
$function$;
-- +goose StatementEnd
-- +goose Down
DROP FUNCTION public.geojsontogeom;
