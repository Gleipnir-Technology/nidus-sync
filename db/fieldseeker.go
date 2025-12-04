package db

import (
	"context"
	"fmt"

	fslayer "github.com/Gleipnir-Technology/arcgis-go/fieldseeker/layer"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	googleuuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/scan"
)

func SaveOrUpdateAerialSpraySession(fs []*fslayer.AerialSpraySession) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring AerialSpraySession rows")
	return 0, 0, nil
}
func SaveOrUpdateAerialSprayLine(fs []*fslayer.AerialSprayLine) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring AerialSprayLine rows")
	return 0, 0, nil
}
func SaveOrUpdateBarrierSpray(fs []*fslayer.BarrierSpray) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring BarrierSpray rows")
	return 0, 0, nil
}
func SaveOrUpdateBarrierSprayRoute(fs []*fslayer.BarrierSprayRoute) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring BarrierSprayRoute rows")
	return 0, 0, nil
}
func SaveOrUpdateContainerRelate(fs []*fslayer.ContainerRelate) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring ContainerRelate rows")
	return 0, 0, nil
}
func SaveOrUpdateFieldScoutingLog(fs []*fslayer.FieldScoutingLog) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring FieldScoutingLog rows")
	return 0, 0, nil
}
func SaveOrUpdateHabitatRelate(fs []*fslayer.HabitatRelate) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring HabitatRelate rows")
	return 0, 0, nil
}
func SaveOrUpdateInspectionSample(fs []*fslayer.InspectionSample) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring InspectionSample rows")
	return 0, 0, nil
}
func SaveOrUpdateInspectionSampleDetail(fs []*fslayer.InspectionSampleDetail) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring InspectionSampleDetail rows")
	return 0, 0, nil
}
func SaveOrUpdateLandingCount(fs []*fslayer.LandingCount) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring LandingCount rows")
	return 0, 0, nil
}
func SaveOrUpdateLandingCountLocation(fs []*fslayer.LandingCountLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring LandingCountLocation rows")
	return 0, 0, nil
}
func SaveOrUpdateLineLocation(fs []*fslayer.LineLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring LineLocation rows")
	return 0, 0, nil
}
func SaveOrUpdateLocationTracking(fs []*fslayer.LocationTracking) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring LocationTracking rows")
	return 0, 0, nil
}
func SaveOrUpdateMosquitoInspection(fs []*fslayer.MosquitoInspection) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring MosquitoInspection rows")
	return 0, 0, nil
}
func SaveOrUpdateOfflineMapAreas(fs []*fslayer.OfflineMapAreas) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring OfflineMapAreas rows")
	return 0, 0, nil
}
func SaveOrUpdateProposedTreatmentArea(fs []*fslayer.ProposedTreatmentArea) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring ProposedTreatmentArea rows")
	return 0, 0, nil
}
func SaveOrUpdatePointLocation(fs []*fslayer.PointLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring PointLocation rows")
	return 0, 0, nil
}
func SaveOrUpdatePolygonLocation(fs []*fslayer.PolygonLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring PolygonLocation rows")
	return 0, 0, nil
}
func SaveOrUpdatePoolDetail(fs []*fslayer.PoolDetail) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring PoolDetail rows")
	return 0, 0, nil
}
func SaveOrUpdatePool(fs []*fslayer.Pool) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring Pool rows")
	return 0, 0, nil
}
func SaveOrUpdatePoolBuffer(fs []*fslayer.PoolBuffer) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring PoolBuffer rows")
	return 0, 0, nil
}
func SaveOrUpdateQALarvCount(fs []*fslayer.QALarvCount) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring QALarvCount rows")
	return 0, 0, nil
}
func SaveOrUpdateQAMosquitoInspection(fs []*fslayer.QAMosquitoInspection) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring QAMosquitoInspection rows")
	return 0, 0, nil
}
func SaveOrUpdateQAProductObservation(fs []*fslayer.QAProductObservation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring QAProductObservation rows")
	return 0, 0, nil
}
func SaveOrUpdateRestrictedArea(fs []*fslayer.RestrictedArea) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring RestrictedArea rows")
	return 0, 0, nil
}
func SaveOrUpdateRodentInspection(fs []*fslayer.RodentInspection) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring RodentInspection rows")
	return 0, 0, nil
}
func toUUID(u googleuuid.UUID) omitnull.Val[uuid.UUID] {
	bytes := u[:]
	converted, err := uuid.FromBytes(bytes)
	if err != nil {
		log.Warn().Str("uuid", u.String()).Msg("Failed to convert uuid")
		return omitnull.FromPtr[uuid.UUID](nil)
	}
	return omitnull.From(converted)
}
func toObjectID(o uint) omit.Val[int64] {
	return omit.From[int64](int64(o))
}

type InsertResultRow struct {
	Inserted bool `db:"row_inserted"`
	Version  int  `db:"version_num"`
}

type rowConverter[T any] func(*T) []SqlParam

func doUpdatesViaFunction[T any](ctx context.Context, org *models.Organization, fs []*T, table string, procedure string, converter rowConverter[T]) (inserts uint, updates uint, err error) {
	//log.Info().Int("rows", len(fs)).Msg("Processing RodentLocation")
	for _, row := range fs {
		params := converter(row)
		q := queryStoredProcedure(procedure, params...)
		query := psql.RawQuery(q)
		log.Info().Str("query", q).Msg("querying")
		result, err := bob.One[InsertResultRow](ctx, PGInstance.BobDB, query, scan.StructMapper[InsertResultRow]())
		if err != nil {
			return inserts, updates, fmt.Errorf("Failed to execute %s: %w", procedure, err)
		}
		if result.Inserted {
			if result.Version == 1 {
				inserts += 1
			} else {
				updates += 1
			}
		}
		log.Info().Bool("inserted", result.Inserted).Int("version", result.Version).Msg("querying")
	}
	return inserts, updates, err

}
func SaveOrUpdateRodentLocation(ctx context.Context, org *models.Organization, fs []*fslayer.RodentLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "RodentLocation", "fieldseeker.insert_rodentlocation", func(row *fslayer.RodentLocation) []SqlParam {
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			String("p_locationname", row.LocationName),
			String("p_zone", row.Zone),
			String("p_zone2", row.Zone2),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.Usetype),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.Accessdesc),
			String("p_comments", row.Comments),
			String("p_symbology", row.Symbology),
			String("p_externalid", row.ExternalID),
			Timestamp("p_nextactiondatescheduled", row.Nextactiondatescheduled),
			Int32("p_locationnumber", row.Locationnumber),
			Timestamp("p_lastinspectdate", row.LastInspectionDate),
			String("p_lastinspectspecies", row.LastInspectionSpecies),
			String("p_lastinspectaction", row.LastInspectionAction),
			String("p_lastinspectconditions", row.LastInspectionConditions),
			String("p_lastinspectrodentevidence", row.LastInspectionRodentEvidence),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_jurisdiction", row.Jurisdiction),
		}
	})
}

func SaveOrUpdateSampleCollection(fs []*fslayer.SampleCollection) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring SampleCollection rows")
	return 0, 0, nil
}
func SaveOrUpdateSampleLocation(fs []*fslayer.SampleLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring SampleLocation rows")
	return 0, 0, nil
}
func SaveOrUpdateServiceRequest(fs []*fslayer.ServiceRequest) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring ServiceRequest rows")
	return 0, 0, nil
}
func SaveOrUpdateSpeciesAbundance(fs []*fslayer.SpeciesAbundance) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring SpeciesAbundance rows")
	return 0, 0, nil
}
func SaveOrUpdateStormDrain(fs []*fslayer.StormDrain) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring StormDrain rows")
	return 0, 0, nil
}
func SaveOrUpdateTracklog(fs []*fslayer.Tracklog) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring Tracklog rows")
	return 0, 0, nil
}
func SaveOrUpdateTrapLocation(fs []*fslayer.TrapLocation) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring TrapLocation rows")
	return 0, 0, nil
}
func SaveOrUpdateTrapData(fs []*fslayer.TrapData) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring TrapData rows")
	return 0, 0, nil
}
func SaveOrUpdateTimeCard(fs []*fslayer.TimeCard) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring TimeCard rows")
	return 0, 0, nil
}
func SaveOrUpdateTreatment(fs []*fslayer.Treatment) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring Treatment rows")
	return 0, 0, nil
}
func SaveOrUpdateTreatmentArea(fs []*fslayer.TreatmentArea) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring TreatmentArea rows")
	return 0, 0, nil
}
func SaveOrUpdateULVSprayRoute(fs []*fslayer.ULVSprayRoute) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring ULVSprayRoute rows")
	return 0, 0, nil
}
func SaveOrUpdateZones(fs []*fslayer.Zones) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring Zones rows")
	return 0, 0, nil
}
func SaveOrUpdateZones2(fs []*fslayer.Zones2) (inserts uint, updates uint, err error) {
	//log.Warn().Int("len", len(fs)).Msg("Ignoring Zones2 rows")
	return 0, 0, nil
}
