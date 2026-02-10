package db

import (
	"context"
	"fmt"

	fslayer "github.com/Gleipnir-Technology/arcgis-go/fieldseeker/layer"
	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	googleuuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func SaveOrUpdateAerialSpraySession(ctx context.Context, org *models.Organization, fs []*fslayer.AerialSpraySession) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring AerialSpraySession data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "AerialSpraySession", "fieldseeker.insert_aerialspraysession", func(row *fslayer.AerialSpraySession) ([]SqlParam, error) {
			return []SqlParam{
				//Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateAerialSprayLine(ctx context.Context, org *models.Organization, fs []*fslayer.AerialSprayLine) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring AerialSprayLine data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "AerialSprayLine", "fieldseeker.insert_aerialsprayline", func(row *fslayer.AerialSprayLine) ([]SqlParam, error) {
			return []SqlParam{
			}, nil
		})
	*/
}
func SaveOrUpdateBarrierSpray(ctx context.Context, org *models.Organization, fs []*fslayer.BarrierSpray) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring BarrierSpray data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "BarrierSpray", "fieldseeker.insert_barrierspray", func(row *fslayer.BarrierSpray) ([]SqlParam, error) {
			return []SqlParam{
			}, nil
		})
	*/
}
func SaveOrUpdateBarrierSprayRoute(ctx context.Context, org *models.Organization, fs []*fslayer.BarrierSprayRoute) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring BarrierSprayRoute data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "BarrierSprayRoute", "fieldseeker.insert_barriersprayroute", func(row *fslayer.BarrierSprayRoute) ([]SqlParam, error) {
			return []SqlParam{
			}, nil
		})
	*/
}
func SaveOrUpdateContainerRelate(ctx context.Context, org *models.Organization, fs []*fslayer.ContainerRelate) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "ContainerRelate", "fieldseeker.insert_containerrelate", func(row *fslayer.ContainerRelate) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			UUID("p_inspsampleid", row.InspsampleID),
			UUID("p_mosquitoinspid", row.MosquitoinspID),
			UUID("p_treatmentid", row.TreatmentID),
			String("p_containertype", row.ContainerType),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
}
func SaveOrUpdateFieldScoutingLog(ctx context.Context, org *models.Organization, fs []*fslayer.FieldScoutingLog) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "FieldScoutingLog", "fieldseeker.insert_fieldscoutinglog", func(row *fslayer.FieldScoutingLog) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Int16("p_status", row.Status),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateHabitatRelate(ctx context.Context, org *models.Organization, fs []*fslayer.HabitatRelate) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "HabitatRelate", "fieldseeker.insert_habitatrelate", func(row *fslayer.HabitatRelate) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_foreign_id", row.ForeignID),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			String("p_habitattype", row.HabitatType),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateInspectionSample(ctx context.Context, org *models.Organization, fs []*fslayer.InspectionSample) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "InspectionSample", "fieldseeker.insert_inspectionsample", func(row *fslayer.InspectionSample) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_insp_id", row.InspID),
			String("p_sampleid", row.SampleID),
			Int16("p_processed", row.Processed),
			String("p_idbytech", row.TechIdentifyingSpeciesInLab),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateInspectionSampleDetail(ctx context.Context, org *models.Organization, fs []*fslayer.InspectionSampleDetail) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "InspectionSampleDetail", "fieldseeker.insert_inspectionsampledetail", func(row *fslayer.InspectionSampleDetail) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_inspsample_id", row.InspsampleID),
			String("p_fieldspecies", row.FieldSpecies),
			Int16("p_flarvcount", row.FieldLarvaCount),
			Int16("p_fpupcount", row.FieldPupaCount),
			Int16("p_feggcount", row.FieldEggCount),
			String("p_flstages", row.FieldLarvalStages),
			String("p_fdomstage", row.FieldDominantStage),
			String("p_fadultact", row.FieldAdultActivity),
			String("p_labspecies", row.LabSpecies),
			Int16("p_llarvcount", row.LabLarvaCount),
			Int16("p_lpupcount", row.LabPupaCount),
			Int16("p_leggcount", row.LabEggCount),
			String("p_ldomstage", row.LabDominantStage),
			String("p_comments", row.Comments),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int16("p_processed", row.Processed),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateLandingCount(ctx context.Context, org *models.Organization, fs []*fslayer.LandingCount) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring LandingCount data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "LandingCount", "fieldseeker.insert_landingcount", func(row *fslayer.LandingCount) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateLandingCountLocation(ctx context.Context, org *models.Organization, fs []*fslayer.LandingCountLocation) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring LandingCountLocation data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "LandingCountLocation", "fieldseeker.insert_landingcountlocation", func(row *fslayer.LandingCountLocation) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateLineLocation(ctx context.Context, org *models.Organization, fs []*fslayer.LineLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "LineLocation", "fieldseeker.insert_linelocation", func(row *fslayer.LineLocation) ([]SqlParam, error) {
		gisPoint, err := lineOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			String("p_zone", row.Zone),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.UseType),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.AccessDescription),
			String("p_comments", row.Comments),
			String("p_symbology", row.Symbology),
			String("p_externalid", row.ExternalID),
			Float64("p_acres", row.Acres),
			Timestamp("p_nextactiondatescheduled", row.NextScheduledAction),
			Int16("p_larvinspectinterval", row.LarvalInspectionInterval),
			Float64("p_length_ft", row.Length),
			Float64("p_width_ft", row.Width),
			String("p_zone2", row.Zone2),
			Int32("p_locationnumber", row.Locationnumber),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_lastinspectdate", row.LastInspectionDate),
			String("p_lastinspectbreeding", row.LastInspectionBreeding),
			Float64("p_lastinspectavglarvae", row.LastInspectionAverageLarvae),
			Float64("p_lastinspectavgpupae", row.LastInspectionAveragePupae),
			String("p_lastinspectlstages", row.LastInspectionLarvalStages),
			String("p_lastinspectactiontaken", row.LastInspectionAction),
			String("p_lastinspectfieldspecies", row.LastInspectionFieldSpecies),
			Timestamp("p_lasttreatdate", row.LastTreatmentDate),
			String("p_lasttreatproduct", row.LastTreatmentProduct),
			Float64("p_lasttreatqty", row.LastTreatmentQuantity),
			String("p_lasttreatqtyunit", row.LastTreatmentQuantityUnit),
			Float64("p_hectares", row.Hectares),
			String("p_lastinspectactivity", row.LastInspectionActivity),
			String("p_lasttreatactivity", row.LastTreatmentActivity),
			Float64("p_length_meters", row.LengthMeters),
			Float64("p_width_meters", row.WidthMeters),
			String("p_lastinspectconditions", row.LastInspectionConditions),
			String("p_waterorigin", row.WaterOrigin),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_jurisdiction", row.Jurisdiction),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateLocationTracking(ctx context.Context, org *models.Organization, fs []*fslayer.LocationTracking) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "LocationTracking", "fieldseeker.insert_locationtracking", func(row *fslayer.LocationTracking) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Float64("p_accuracy", row.Accuracym),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			UUID("p_globalid", row.GlobalID),
			String("p_fieldtech", row.FieldTech),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateMosquitoInspection(ctx context.Context, org *models.Organization, fs []*fslayer.MosquitoInspection) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "MosquitoInspection", "fieldseeker.insert_mosquitoinspection", func(row *fslayer.MosquitoInspection) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Int16("p_numdips", row.Dips),
			String("p_activity", row.Activity),
			String("p_breeding", row.Breeding),
			Int16("p_totlarvae", row.TotalLarvae),
			Int16("p_totpupae", row.TotalPupae),
			Int16("p_eggs", row.Eggs),
			Int16("p_posdips", row.PositiveDips),
			String("p_adultact", row.AdultActivity),
			String("p_lstages", row.LarvalStages),
			String("p_domstage", row.DominantStage),
			String("p_actiontaken", row.Action),
			String("p_comments", row.Comments),
			Float64("p_avetemp", row.AverageTemperature),
			Float64("p_windspeed", row.WindSpeed),
			Float64("p_raingauge", row.RainGauge),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			String("p_winddir", row.WindDirection),
			Float64("p_avglarvae", row.AverageLarvae),
			Float64("p_avgpupae", row.AveragePupae),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			String("p_locationname", row.LocationName),
			String("p_zone", row.Zone),
			Int16("p_recordstatus", row.RecordStatus),
			String("p_zone2", row.Zone2),
			Int16("p_personalcontact", row.PersonalContact),
			Int16("p_tirecount", row.TireCount),
			Int16("p_cbcount", row.CatchBasinCount),
			Int16("p_containercount", row.ContainerCount),
			String("p_fieldspecies", row.FieldSpecies),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			UUID("p_linelocid", row.LinelocID),
			UUID("p_pointlocid", row.PointlocID),
			UUID("p_polygonlocid", row.PolygonlocID),
			UUID("p_srid", row.SrID),
			String("p_fieldtech", row.FieldTech),
			Int16("p_larvaepresent", row.LarvaePresent),
			Int16("p_pupaepresent", row.PupaePresent),
			UUID("p_sdid", row.StormDrainID),
			String("p_sitecond", row.Conditions),
			Int16("p_positivecontainercount", row.PositiveContainerCount),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_jurisdiction", row.Jurisdiction),
			Int16("p_visualmonitoring", row.VisualMonitoring),
			String("p_vmcomments", row.VmComments),
			String("p_adminaction", row.AdminAction),
			UUID("p_ptaid", row.PtaID),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
}
func SaveOrUpdateOfflineMapAreas(ctx context.Context, org *models.Organization, fs []*fslayer.OfflineMapAreas) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring OfflineMapAreas data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "OfflineMapAreas", "fieldseeker.insert_offlinemapareas", func(row *fslayer.OfflineMapAreas) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateProposedTreatmentArea(ctx context.Context, org *models.Organization, fs []*fslayer.ProposedTreatmentArea) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "ProposedTreatmentArea", "fieldseeker.insert_proposedtreatmentarea", func(row *fslayer.ProposedTreatmentArea) ([]SqlParam, error) {
		gisPoint, err := polygonOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		// At this point we've got data that's bad and can't actually be inserted in the database
		// so let's just always make the geo null
		gisPoint = NullParam{"p_geospatial"}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_method", row.Method),
			String("p_comments", row.Comments),
			String("p_zone", row.Zone),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			String("p_zone2", row.Zone2),
			Timestamp("p_completeddate", row.CompletedDate),
			String("p_completedby", row.CompletedBy),
			Int16("p_completed", row.Completed),
			Int16("p_issprayroute", row.IsSprayRoute),
			String("p_name", row.Name),
			Float64("p_acres", row.Acres),
			UUID("p_globalid", row.GlobalID),
			Int16("p_exported", row.Exported),
			String("p_targetproduct", row.TargetProduct),
			Float64("p_targetapprate", row.TargetAppRate),
			Float64("p_hectares", row.Hectares),
			String("p_lasttreatactivity", row.LastTreatmentActivity),
			Timestamp("p_lasttreatdate", row.LastTreatmentDate),
			String("p_lasttreatproduct", row.LastTreatmentProduct),
			Float64("p_lasttreatqty", row.LastTreatmentQuantity),
			String("p_lasttreatqtyunit", row.LastTreatmentQuantityUnit),
			String("p_priority", row.Priority),
			Timestamp("p_duedate", row.DueDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_targetspecies", row.TargetSpecies),
			Float64("p_shape__area", row.ShapeArea),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
}
func SaveOrUpdatePointLocation(ctx context.Context, org *models.Organization, fs []*fslayer.PointLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "PointLocation", "fieldseeker.insert_pointlocation", func(row *fslayer.PointLocation) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			String("p_zone", row.Zone),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.UseType),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.AccessDescription),
			String("p_comments", row.Comments),
			String("p_symbology", row.Symbology),
			String("p_externalid", row.ExternalID),
			Timestamp("p_nextactiondatescheduled", row.NextScheduledAction),
			Int16("p_larvinspectinterval", row.LarvalInspectionInterval),
			String("p_zone2", row.Zone2),
			Int32("p_locationnumber", row.Locationnumber),
			UUID("p_globalid", row.GlobalID),
			String("p_stype", row.SourceType),
			Timestamp("p_lastinspectdate", row.LastInspectionDate),
			String("p_lastinspectbreeding", row.LastInspectionBreeding),
			Float64("p_lastinspectavglarvae", row.LastInspectionAverageLarvae),
			Float64("p_lastinspectavgpupae", row.LastInspectionAveragePupae),
			String("p_lastinspectlstages", row.LastInspectionLarvalStages),
			String("p_lastinspectactiontaken", row.LastInspectionAction),
			String("p_lastinspectfieldspecies", row.LastInspectionFieldSpecies),
			Timestamp("p_lasttreatdate", row.LastTreatmentDate),
			String("p_lasttreatproduct", row.LastTreatmentProduct),
			Float64("p_lasttreatqty", row.LastTreatmentQuantity),
			String("p_lasttreatqtyunit", row.LastTreatmentQuantityUnit),
			String("p_lastinspectactivity", row.LastInspectionActivity),
			String("p_lasttreatactivity", row.LastTreatmentActivity),
			String("p_lastinspectconditions", row.LastInspectionConditions),
			String("p_waterorigin", row.WaterOrigin),
			Float64("p_x", row.X),
			Float64("p_y", row.Y),
			String("p_assignedtech", row.AssignedTech),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_jurisdiction", row.Jurisdiction),
			String("p_deactivate_reason", row.ReasonForDeactivation),
			Int32("p_scalarpriority", row.ScalarPriority),
			String("p_sourcestatus", row.SourceStatus),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdatePolygonLocation(ctx context.Context, org *models.Organization, fs []*fslayer.PolygonLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "PolygonLocation", "fieldseeker.insert_polygonlocation", func(row *fslayer.PolygonLocation) ([]SqlParam, error) {
		gisPoint, err := polygonOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		// At this point we've got data that's bad and can't actually be inserted in the database
		// so let's just always make the geo null
		gisPoint = NullParam{"p_geospatial"}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			String("p_zone", row.Zone),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.UseType),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.AccessDescription),
			String("p_comments", row.Comments),
			String("p_symbology", row.Symbology),
			String("p_externalid", row.ExternalID),
			Float64("p_acres", row.Acres),
			Timestamp("p_nextactiondatescheduled", row.NextScheduledAction),
			Int16("p_larvinspectinterval", row.LarvalInspectionInterval),
			String("p_zone2", row.Zone2),
			Int32("p_locationnumber", row.Locationnumber),
			UUID("p_globalid", row.GlobalID),
			Timestamp("p_lastinspectdate", row.LastInspectionDate),
			String("p_lastinspectbreeding", row.LastInspectionBreeding),
			Float64("p_lastinspectavglarvae", row.LastInspectionAverageLarvae),
			Float64("p_lastinspectavgpupae", row.LastInspectionAveragePupae),
			String("p_lastinspectlstages", row.LastInspectionLarvalStages),
			String("p_lastinspectactiontaken", row.LastInspectionAction),
			String("p_lastinspectfieldspecies", row.LastInspectionFieldSpecies),
			Timestamp("p_lasttreatdate", row.LastTreatmentDate),
			String("p_lasttreatproduct", row.LastTreatmentProduct),
			Float64("p_lasttreatqty", row.LastTreatmentQuantity),
			String("p_lasttreatqtyunit", row.LastTreatmentQuantityUnit),
			Float64("p_hectares", row.Hectares),
			String("p_lastinspectactivity", row.LastInspectionActivity),
			String("p_lasttreatactivity", row.LastTreatmentActivity),
			String("p_lastinspectconditions", row.LastInspectionConditions),
			String("p_waterorigin", row.WaterOrigin),
			String("p_filter", row.Filter),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_jurisdiction", row.Jurisdiction),
			Float64("p_shape__area", row.ShapeArea),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdatePoolDetail(ctx context.Context, org *models.Organization, fs []*fslayer.PoolDetail) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "PoolDetail", "fieldseeker.insert_pooldetail", func(row *fslayer.PoolDetail) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_trapdata_id", row.TrapDataID),
			UUID("p_pool_id", row.PoolID),
			String("p_species", row.Species),
			Int16("p_females", row.Females),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdatePool(ctx context.Context, org *models.Organization, fs []*fslayer.Pool) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "Pool", "fieldseeker.insert_pool", func(row *fslayer.Pool) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_trapdata_id", row.TrapDataID),
			Timestamp("p_datesent", row.DateSent),
			String("p_survtech", row.SurveyTech),
			Timestamp("p_datetested", row.DateTested),
			String("p_testtech", row.TestTech),
			String("p_comments", row.Comments),
			String("p_sampleid", row.SampleID),
			Int16("p_processed", row.Processed),
			UUID("p_lab_id", row.LabID),
			String("p_testmethod", row.TestMethods),
			String("p_diseasetested", row.DiseasesTested),
			String("p_diseasepos", row.DiseasesPositive),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			String("p_lab", row.Lab),
			Int16("p_poolyear", row.PoolYear),
			Int16("p_gatewaysync", row.GatewaySync),
			String("p_vectorsurvcollectionid", row.VectorsurvcollectionID),
			String("p_vectorsurvpoolid", row.VectorsurvpoolID),
			String("p_vectorsurvtrapdataid", row.VectorsurvtrapdataID),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdatePoolBuffer(ctx context.Context, org *models.Organization, fs []*fslayer.PoolBuffer) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring PoolBuffer data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "PoolBuffer", "fieldseeker.insert_poolbuffer", func(row *fslayer.PoolBuffer) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateQALarvCount(ctx context.Context, org *models.Organization, fs []*fslayer.QALarvCount) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring QALarvCount data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "QALarvCount", "fieldseeker.insert_qalarvcount", func(row *fslayer.QALarvCount) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateQAMosquitoInspection(ctx context.Context, org *models.Organization, fs []*fslayer.QAMosquitoInspection) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "QAMosquitoInspection", "fieldseeker.insert_qamosquitoinspection", func(row *fslayer.QAMosquitoInspection) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Int16("p_posdips", row.PositiveDips),
			String("p_actiontaken", row.Action),
			String("p_comments", row.Comments),
			Float64("p_avetemp", row.AverageTemperature),
			Float64("p_windspeed", row.WindSpeed),
			Float64("p_raingauge", row.RainGauge),
			UUID("p_globalid", row.GlobalID),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			String("p_winddir", row.WindDirection),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.Reviewedby),
			Timestamp("p_revieweddate", row.Revieweddate),
			String("p_locationname", row.Locationname),
			String("p_zone", row.Zone),
			Int16("p_recordstatus", row.Recordstatus),
			String("p_zone2", row.Zone2),
			Int16("p_lr", row.LandingRate),
			Int16("p_negdips", row.NegativeDips),
			Float64("p_totalacres", row.TotalAcres),
			Float64("p_acresbreeding", row.AcresBreeding),
			Int16("p_fish", row.FishPresent),
			String("p_sitetype", row.SiteType),
			String("p_breedingpotential", row.BreedingPotential),
			Int16("p_movingwater", row.MovingWater),
			Int16("p_nowaterever", row.NoEvidenceOfWaterEver),
			String("p_mosquitohabitat", row.MosquitoHabitatIndicators),
			Int16("p_habvalue1", row.HabitatValue),
			Int16("p_habvalue1percent", row.Habvalue1percent),
			Int16("p_habvalue2", row.HabitatValue2),
			Int16("p_habvalue2percent", row.Habvalue2percent),
			Int16("p_potential", row.Potential),
			Int16("p_larvaepresent", row.LarvaePresent),
			Int16("p_larvaeinsidetreatedarea", row.LarvaeInsideTreatedArea),
			Int16("p_larvaeoutsidetreatedarea", row.LarvaeOutsideTreatedArea),
			String("p_larvaereason", row.ReasonLarvaePresent),
			String("p_aquaticorganisms", row.AquaticOrganisms),
			String("p_vegetation", row.Vegetation),
			String("p_sourcereduction", row.SourceReduction),
			Int16("p_waterpresent", row.WaterPresent),
			String("p_watermovement1", row.WaterMovement),
			Int16("p_watermovement1percent", row.Watermovement1percent),
			String("p_watermovement2", row.WaterMovement2),
			Int16("p_watermovement2percent", row.Watermovement2percent),
			String("p_soilconditions", row.SoilConditions),
			String("p_waterduration", row.HowLongWaterPresent),
			String("p_watersource", row.WaterSource),
			String("p_waterconditions", row.WaterConditions),
			Int16("p_adultactivity", row.AdultActivity),
			UUID("p_linelocid", row.LinelocID),
			UUID("p_pointlocid", row.PointlocID),
			UUID("p_polygonlocid", row.PolygonlocID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			String("p_fieldtech", row.FieldTech),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateQAProductObservation(ctx context.Context, org *models.Organization, fs []*fslayer.QAProductObservation) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring QAProductObservation data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "QAProductObservation", "fieldseeker.insert_qaproductobservation", func(row *fslayer.QAProductObservation) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateRestrictedArea(ctx context.Context, org *models.Organization, fs []*fslayer.RestrictedArea) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring RestrictedArea data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "RestrictedArea", "fieldseeker.insert_restrictedarea", func(row *fslayer.RestrictedArea) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateRodentInspection(ctx context.Context, org *models.Organization, fs []*fslayer.RodentInspection) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring RodentInspection data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "RodentInspection", "fieldseeker.insert_rodentinspection", func(row *fslayer.RodentInspection) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateRodentLocation(ctx context.Context, org *models.Organization, fs []*fslayer.RodentLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "RodentLocation", "fieldseeker.insert_rodentlocation", func(row *fslayer.RodentLocation) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
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
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
}

func SaveOrUpdateSampleCollection(ctx context.Context, org *models.Organization, fs []*fslayer.SampleCollection) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "SampleCollection", "fieldseeker.insert_samplecollection", func(row *fslayer.SampleCollection) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_loc_id", row.LocID),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			String("p_sitecond", row.Conditions),
			String("p_sampleid", row.SampleID),
			String("p_survtech", row.SurveillanceTechnician),
			Timestamp("p_datesent", row.Sent),
			Timestamp("p_datetested", row.Tested),
			String("p_testtech", row.TestTechnician),
			String("p_comments", row.Comments),
			Int16("p_processed", row.Processed),
			String("p_sampletype", row.SampleType),
			String("p_samplecond", row.SampleCondition),
			String("p_species", row.Species),
			String("p_sex", row.Sex),
			Float64("p_avetemp", row.AverageTemperature),
			Float64("p_windspeed", row.WindSpeed),
			String("p_winddir", row.WindDirection),
			Float64("p_raingauge", row.RainGauge),
			String("p_activity", row.Activity),
			String("p_testmethod", row.TestMethod),
			String("p_diseasetested", row.DiseaseTested),
			String("p_diseasepos", row.DiseasePositive),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			String("p_locationname", row.LocationName),
			String("p_zone", row.Zone),
			Int16("p_recordstatus", row.RecordStatus),
			String("p_zone2", row.Zone2),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			String("p_lab", row.Lab),
			String("p_fieldtech", row.FieldTech),
			UUID("p_flockid", row.FlockID),
			Int16("p_samplecount", row.SampleCount),
			UUID("p_chickenid", row.ChickenID),
			Int16("p_gatewaysync", row.GatewaySync),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateSampleLocation(ctx context.Context, org *models.Organization, fs []*fslayer.SampleLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "SampleLocation", "fieldseeker.insert_samplelocation", func(row *fslayer.SampleLocation) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			String("p_zone", row.Zone),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.UseType),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.AccessDescription),
			String("p_comments", row.Comments),
			String("p_externalid", row.ExternalID),
			Timestamp("p_nextactiondatescheduled", row.NextScheduledAction),
			String("p_zone2", row.Zone2),
			Int32("p_locationnumber", row.Locationnumber),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int16("p_gatewaysync", row.GatewaySync),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateServiceRequest(ctx context.Context, org *models.Organization, fs []*fslayer.ServiceRequest) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "ServiceRequest", "fieldseeker.insert_servicerequest", func(row *fslayer.ServiceRequest) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Timestamp("p_recdatetime", row.Received),
			String("p_source", row.Source),
			String("p_entrytech", row.EnteredBy),
			String("p_priority", row.Priority),
			String("p_supervisor", row.Supervisor),
			String("p_assignedtech", row.AssignedTo),
			String("p_status", row.Status),
			Int16("p_clranon", row.AnonymousCaller),
			String("p_clrfname", row.CallerName),
			String("p_clrphone1", row.CallerPhone),
			String("p_clrphone2", row.CallerAlternatePhone),
			String("p_clremail", row.CallerEmail),
			String("p_clrcompany", row.CallerCompany),
			String("p_clraddr1", row.CallerAddress),
			String("p_clraddr2", row.CallerAddress2),
			String("p_clrcity", row.CallerCity),
			String("p_clrstate", row.CallerState),
			String("p_clrzip", row.CallerZip),
			String("p_clrother", row.CallerOther),
			String("p_clrcontpref", row.CallerContactPreference),
			String("p_reqcompany", row.RequestCompany),
			String("p_reqaddr1", row.RequestAddress),
			String("p_reqaddr2", row.RequestAddress2),
			String("p_reqcity", row.RequestCity),
			String("p_reqstate", row.RequestState),
			String("p_reqzip", row.RequestZip),
			String("p_reqcrossst", row.RequestCrossStreet),
			String("p_reqsubdiv", row.RequestSubdivision),
			String("p_reqmapgrid", row.RequestMapGrID),
			Int16("p_reqpermission", row.PermissionToEnter),
			String("p_reqtarget", row.RequestTarget),
			String("p_reqdescr", row.RequestDescription),
			String("p_reqnotesfortech", row.NotesForFieldTechnician),
			String("p_reqnotesforcust", row.NotesForCustomer),
			String("p_reqfldnotes", row.RequestFieldNotes),
			String("p_reqprogramactions", row.RequestProgramActions),
			Timestamp("p_datetimeclosed", row.Closed),
			String("p_techclosed", row.ClosedBy),
			Int32("p_sr_number", row.Sr),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			Int16("p_accepted", row.Accepted),
			Timestamp("p_accepteddate", row.AcceptedDate),
			String("p_rejectedby", row.RejectedBy),
			Timestamp("p_rejecteddate", row.RejectedDate),
			String("p_rejectedreason", row.RejectedReason),
			Timestamp("p_duedate", row.DueDate),
			String("p_acceptedby", row.AcceptedBy),
			String("p_comments", row.Comments),
			Timestamp("p_estcompletedate", row.EstimatedCompletionDate),
			String("p_nextaction", row.NextAction),
			Int16("p_recordstatus", row.RecordStatus),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_firstresponsedate", row.FirstResponseDate),
			Int16("p_responsedaycount", row.ResponseDayCount),
			String("p_allowed", row.VerifyCorrectLocation),
			String("p_xvalue", row.Xvalue),
			String("p_yvalue", row.Yvalue),
			String("p_validx", row.ValidX),
			String("p_validy", row.ValidY),
			String("p_externalid", row.ExternalID),
			String("p_externalerror", row.ExternalError),
			UUID("p_pointlocid", row.PointlocID),
			Int16("p_notified", row.Notified),
			Timestamp("p_notifieddate", row.NotifiedDate),
			Int16("p_scheduled", row.Scheduled),
			Timestamp("p_scheduleddate", row.ScheduledDate),
			Int32("p_dog", row.Dog),
			String("p_schedule_period", row.SchedulePeriod),
			String("p_schedule_notes", row.ScheduleNotes),
			Int32("p_spanish", row.PreferSpeakingSpanish),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_issuesreported", row.IssuesReported),
			String("p_jurisdiction", row.Jurisdiction),
			String("p_notificationtimestamp", row.NotificationTimestamp),
			String("p_zone", row.Zone),
			String("p_zone2", row.Zone2),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateSpeciesAbundance(ctx context.Context, org *models.Organization, fs []*fslayer.SpeciesAbundance) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "SpeciesAbundance", "fieldseeker.insert_speciesabundance", func(row *fslayer.SpeciesAbundance) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_trapdata_id", row.TrapDataID),
			String("p_species", row.Species),
			Int16("p_males", row.Males),
			Int16("p_unknown", row.Unknown),
			Int16("p_bloodedfem", row.BloodedFemales),
			Int16("p_gravidfem", row.GravidFemales),
			Int16("p_larvae", row.Larvae),
			Int16("p_poolstogen", row.PoolsToGenerate),
			Int16("p_processed", row.Processed),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int16("p_pupae", row.Pupae),
			Int16("p_eggs", row.Eggs),
			Int32("p_females", row.Females),
			Int32("p_total", row.TotalAdults),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			Int32("p_yearweek", row.YearWeek),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateStormDrain(ctx context.Context, org *models.Organization, fs []*fslayer.StormDrain) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "StormDrain", "fieldseeker.insert_stormdrain", func(row *fslayer.StormDrain) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			Timestamp("p_nexttreatmentdate", row.NextTreatmentDate),
			Timestamp("p_lasttreatdate", row.LastTreatDate),
			String("p_lastaction", row.LastAction),
			String("p_symbology", row.Symbology),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			String("p_laststatus", row.LastStatus),
			String("p_zone", row.Zone),
			String("p_zone2", row.Zone2),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_type", row.Type),
			String("p_jurisdiction", row.Jurisdiction),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateTracklog(ctx context.Context, org *models.Organization, fs []*fslayer.Tracklog) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring RodentInspection data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "Tracklog", "fieldseeker.insert_tracklog", func(row *fslayer.Tracklog) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateTrapLocation(ctx context.Context, org *models.Organization, fs []*fslayer.TrapLocation) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "TrapLocation", "fieldseeker.insert_traplocation", func(row *fslayer.TrapLocation) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			String("p_zone", row.Zone),
			String("p_habitat", row.Habitat),
			String("p_priority", row.Priority),
			String("p_usetype", row.UseType),
			Int16("p_active", row.Active),
			String("p_description", row.Description),
			String("p_accessdesc", row.AccessDescription),
			String("p_comments", row.Comments),
			String("p_externalid", row.ExternalID),
			Timestamp("p_nextactiondatescheduled", row.NextScheduledAction),
			String("p_zone2", row.Zone2),
			Int32("p_locationnumber", row.Locationnumber),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int16("p_gatewaysync", row.GatewaySync),
			Int32("p_route", row.Route),
			Int32("p_set_dow", row.SetDayOfWeek),
			Int32("p_route_order", row.RouteOrder),
			String("p_vectorsurvsiteid", row.VectorsurvsiteID),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateTrapData(ctx context.Context, org *models.Organization, fs []*fslayer.TrapData) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "TrapData", "fieldseeker.insert_trapdata", func(row *fslayer.TrapData) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_traptype", row.TrapType),
			String("p_trapactivitytype", row.TrapActivityType),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			String("p_comments", row.Comments),
			String("p_idbytech", row.TechIdentifyingSpeciesInLab),
			String("p_sortbytech", row.TechSortingTrapResultsInLab),
			Int16("p_processed", row.Processed),
			String("p_sitecond", row.SiteConditions),
			String("p_locationname", row.LocationName),
			Int16("p_recordstatus", row.RecordStatus),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			String("p_trapcondition", row.TrapCondition),
			Int16("p_trapnights", row.TrapNights),
			String("p_zone", row.Zone),
			String("p_zone2", row.Zone2),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			UUID("p_srid", row.SrID),
			String("p_fieldtech", row.FieldTech),
			Int16("p_gatewaysync", row.GatewaySync),
			UUID("p_loc_id", row.LocID),
			Float64("p_voltage", row.Voltage),
			String("p_winddir", row.Winddir),
			Float64("p_windspeed", row.Windspeed),
			Float64("p_avetemp", row.Avetemp),
			Float64("p_raingauge", row.Raingauge),
			Int16("p_lr", row.LandingRate),
			Int32("p_field", row.Field),
			String("p_vectorsurvtrapdataid", row.VectorsurvtrapdataID),
			String("p_vectorsurvtraplocationid", row.VectorsurvtraplocationID),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_lure", row.Lure),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateTimeCard(ctx context.Context, org *models.Organization, fs []*fslayer.TimeCard) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "TimeCard", "fieldseeker.insert_timecard", func(row *fslayer.TimeCard) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_activity", row.Activity),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			String("p_comments", row.Comments),
			String("p_externalid", row.ExternalID),
			String("p_equiptype", row.EquipmentType),
			String("p_locationname", row.LocationName),
			String("p_zone", row.Zone),
			String("p_zone2", row.Zone2),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			UUID("p_linelocid", row.LinelocID),
			UUID("p_pointlocid", row.PointlocID),
			UUID("p_polygonlocid", row.PolygonlocID),
			UUID("p_lclocid", row.LclocID),
			UUID("p_samplelocid", row.SamplelocID),
			UUID("p_srid", row.SrID),
			UUID("p_traplocid", row.TraplocID),
			String("p_fieldtech", row.FieldTech),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			UUID("p_rodentlocid", row.RodentlocID),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateTreatment(ctx context.Context, org *models.Organization, fs []*fslayer.Treatment) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "Treatment", "fieldseeker.insert_treatment", func(row *fslayer.Treatment) ([]SqlParam, error) {
		gisPoint, err := pointOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_activity", row.Activity),
			Float64("p_treatarea", row.AreaTreated),
			String("p_areaunit", row.AreaUnit),
			String("p_product", row.Product),
			Float64("p_qty", row.Quantity),
			String("p_qtyunit", row.QuantityUnit),
			String("p_method", row.Method),
			String("p_equiptype", row.EquipmentType),
			String("p_comments", row.Comments),
			Float64("p_avetemp", row.AverageTemperature),
			Float64("p_windspeed", row.WindSpeed),
			String("p_winddir", row.WindDirection),
			Float64("p_raingauge", row.RainGauge),
			Timestamp("p_startdatetime", row.Start),
			Timestamp("p_enddatetime", row.Finish),
			UUID("p_insp_id", row.InspID),
			Int16("p_reviewed", row.Reviewed),
			String("p_reviewedby", row.ReviewedBy),
			Timestamp("p_revieweddate", row.ReviewedDate),
			String("p_locationname", row.LocationName),
			String("p_zone", row.Zone),
			Int16("p_warningoverride", row.WarningOverride),
			Int16("p_recordstatus", row.RecordStatus),
			String("p_zone2", row.Zone2),
			Float64("p_treatacres", row.TreatedAcres),
			Int16("p_tirecount", row.TireCount),
			Int16("p_cbcount", row.CatchBasinCount),
			Int16("p_containercount", row.ContainerCount),
			UUID("p_globalid", row.GlobalID),
			Float64("p_treatmentlength", row.TreatmentLength),
			Float64("p_treatmenthours", row.TreatmentHours),
			String("p_treatmentlengthunits", row.TreatmentLengthUnits),
			UUID("p_linelocid", row.LinelocID),
			UUID("p_pointlocid", row.PointlocID),
			UUID("p_polygonlocid", row.PolygonlocID),
			UUID("p_srid", row.SrID),
			UUID("p_sdid", row.SdID),
			UUID("p_barrierrouteid", row.BarrierrouteID),
			UUID("p_ulvrouteid", row.UlvrouteID),
			String("p_fieldtech", row.FieldTech),
			UUID("p_ptaid", row.PtaID),
			Float64("p_flowrate", row.Flowrate),
			String("p_habitat", row.Habitat),
			Float64("p_treathectares", row.TreatHectares),
			String("p_invloc", row.InventoryLocation),
			String("p_temp_sitecond", row.TempConditions),
			String("p_sitecond", row.Conditions),
			Float64("p_totalcostprodcut", row.TotalCostProduct),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			String("p_targetspecies", row.TargetSpecies),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateTreatmentArea(ctx context.Context, org *models.Organization, fs []*fslayer.TreatmentArea) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "TreatmentArea", "fieldseeker.insert_treatmentarea", func(row *fslayer.TreatmentArea) ([]SqlParam, error) {
		gisPoint, err := polygonOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		gisPoint = NullParam{"p_geospatial"}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			UUID("p_treat_id", row.TreatID),
			UUID("p_session_id", row.SessionID),
			Timestamp("p_treatdate", row.TreatmentDate),
			String("p_comments", row.Comments),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int16("p_notified", row.Notified),
			String("p_type", row.Type),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			Float64("p_shape__area", row.ShapeArea),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateULVSprayRoute(ctx context.Context, org *models.Organization, fs []*fslayer.ULVSprayRoute) (inserts uint, updates uint, err error) {
	log.Warn().Msg("Ignoring RodentInspection data")
	return 0, 0, nil
	/*
		return doUpdatesViaFunction(ctx, org, fs, "ULVSprayRoute", "fieldseeker.insert_ulvsprayroute", func(row *fslayer.ULVSprayRoute) ([]SqlParam, error) {
			return []SqlParam{
				Uint("p_objectid", row.ObjectID),
			}, nil
		})
	*/
}
func SaveOrUpdateZones(ctx context.Context, org *models.Organization, fs []*fslayer.Zones) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "Zones", "fieldseeker.insert_zones", func(row *fslayer.Zones) ([]SqlParam, error) {
		gisPoint, err := polygonOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		// At this point we've got data that's bad and can't actually be inserted in the database
		// so let's just always make the geo null
		gisPoint = NullParam{"p_geospatial"}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Int32("p_active", row.Active),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			Float64("p_shape__area", row.ShapeArea),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}
func SaveOrUpdateZones2(ctx context.Context, org *models.Organization, fs []*fslayer.Zones2) (inserts uint, updates uint, err error) {
	return doUpdatesViaFunction(ctx, org, fs, "Zones2", "fieldseeker.insert_zones2", func(row *fslayer.Zones2) ([]SqlParam, error) {
		gisPoint, err := polygonOrNull(row.Geometry)
		if err != nil {
			return []SqlParam{}, err
		}
		// At this point we've got data that's bad and can't actually be inserted in the database
		// so let's just always make the geo null
		gisPoint = NullParam{"p_geospatial"}
		return []SqlParam{
			Uint("p_objectid", row.ObjectID),
			Int32("p_organization_id", org.ID),
			String("p_name", row.Name),
			UUID("p_globalid", row.GlobalID),
			String("p_created_user", row.CreatedUser),
			Timestamp("p_created_date", row.CreatedDate),
			String("p_last_edited_user", row.LastEditedUser),
			Timestamp("p_last_edited_date", row.LastEditedDate),
			Timestamp("p_creationdate", row.CreationDate),
			String("p_creator", row.Creator),
			Timestamp("p_editdate", row.EditDate),
			String("p_editor", row.Editor),
			Float64("p_shape__area", row.ShapeArea),
			Float64("p_shape__length", row.ShapeLength),
			JsonB("p_geometry", row.Geometry),
			gisPoint,
		}, nil
	})
	return 0, 0, nil
}

type InsertResultRow struct {
	Inserted bool `db:"row_inserted"`
	Version  int  `db:"version_num"`
}

type rowConverter[T any] func(*T) ([]SqlParam, error)

func doUpdatesViaFunction[T any](ctx context.Context, org *models.Organization, fs []*T, table string, procedure string, converter rowConverter[T]) (inserts uint, updates uint, err error) {
	//log.Info().Int("rows", len(fs)).Msg("Processing RodentLocation")
	for _, row := range fs {
		params, err := converter(row)
		if err != nil {
			return inserts, updates, fmt.Errorf("Failed to convert row '%s': %w", row, err)
		}
		q := queryStoredProcedure(procedure, params...)
		query := psql.RawQuery(q)
		result, err := bob.One[InsertResultRow](ctx, PGInstance.BobDB, query, scan.StructMapper[InsertResultRow]())
		if err != nil {
			log.Error().Str("query", q).Msg("Query failed")
			return inserts, updates, fmt.Errorf("Failed to execute %s: %w", procedure, err)
		}
		if result.Inserted {
			if result.Version == 1 {
				inserts += 1
			} else {
				updates += 1
			}
		}
	}
	return inserts, updates, err

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
