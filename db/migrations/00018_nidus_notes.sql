-- +goose Up
CREATE TABLE note_audio (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_ (id) NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	deletor_id INTEGER REFERENCES user_ (id),
	duration REAL NOT NULL,
	organization_id INTEGER REFERENCES organization (id) NOT NULL,
	transcription TEXT,
	transcription_user_edited BOOLEAN NOT NULL,
	version INTEGER NOT NULL,
	uuid UUID NOT NULL,

	PRIMARY KEY(version, uuid)
);

CREATE TABLE note_audio_breadcrumb (
	cell h3index NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	manually_selected BOOLEAN NOT NULL,
	note_audio_version INTEGER NOT NULL,
	note_audio_uuid UUID NOT NULL,
	position INTEGER NOT NULL,

	FOREIGN KEY (note_audio_version, note_audio_uuid) REFERENCES note_audio (version, uuid),
	PRIMARY KEY (note_audio_version, note_audio_uuid, position)
);

CREATE TYPE AudioDataType AS ENUM (
	'raw',
	'raw_normalized',
	'ogg');

CREATE TABLE note_audio_data (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	note_audio_version INTEGER NOT NULL,
	note_audio_uuid UUID NOT NULL,
	type_ AudioDataType NOT NULL,

	FOREIGN KEY (note_audio_version, note_audio_uuid) REFERENCES note_audio (version, uuid),
	PRIMARY KEY (note_audio_version, note_audio_uuid, type_)
);

CREATE TABLE note_image (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_ (id) NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	deletor_id INTEGER REFERENCES user_ (id),
	organization_id INTEGER REFERENCES organization (id) NOT NULL,
	version INTEGER NOT NULL,
	uuid UUID NOT NULL,

	PRIMARY KEY(version, uuid)
);

CREATE TYPE ImageDataType AS ENUM (
	'raw',
	'png');

CREATE TABLE note_image_data (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	note_image_version INTEGER NOT NULL,
	note_image_uuid UUID NOT NULL,
	type_ AudioDataType NOT NULL,

	FOREIGN KEY (note_image_version, note_image_uuid) REFERENCES note_image (version, uuid),
	PRIMARY KEY (note_image_version, note_image_uuid, type_)
);


CREATE TABLE note_image_breadcrumb (
	cell h3index NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	manually_selected BOOLEAN NOT NULL,
	note_image_version INTEGER NOT NULL,
	note_image_uuid UUID NOT NULL,
	position INTEGER NOT NULL,

	FOREIGN KEY (note_image_version, note_image_uuid) REFERENCES note_image (version, uuid),
	PRIMARY KEY (note_image_version, note_image_uuid, position)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION insert_note_audio(
	p_uuid UUID,
	
	p_created TIMESTAMP WITHOUT TIME ZONE,
	p_creator_id INTEGER,
	p_deleted TIMESTAMP WITHOUT TIME ZONE,
	p_deletor_id INTEGER,
	p_duration REAL,
	p_organization_id INTEGER,
	p_transcription TEXT,
	p_transcription_user_edited BOOLEAN
) RETURNS TABLE(row_inserted BOOLEAN, version_num INTEGER) AS $$
DECLARE
	v_next_version INTEGER;
	v_changes_exist BOOLEAN;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM note_audio lv 
		WHERE lv.uuid = p_uuid
		
		AND lv.created IS NOT DISTINCT FROM p_created 
		AND lv.creator_id IS NOT DISTINCT FROM p_creator_id 
		AND lv.deleted IS NOT DISTINCT FROM p_deleted 
		AND lv.deletor_id IS NOT DISTINCT FROM p_deletor_id 
		AND lv.duration IS NOT DISTINCT FROM p_duration 
		AND lv.organization_id IS NOT DISTINCT FROM p_organization_id 
		AND lv.transcription IS NOT DISTINCT FROM p_transcription 
		AND lv.transcription_user_edited IS NOT DISTINCT FROM p_transcription_user_edited 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM note_audio 
				 WHERE uuid = p_uuid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM note_audio
	WHERE uuid = p_uuid;
	
	-- Insert new version
	INSERT INTO note_audio (
		uuid,
		
		created,
		creator_id,
		deleted,
		deletor_id,
		duration,
		organization_id,
		transcription,
		transcription_user_edited
	) VALUES (
		p_uuid,
	
		p_created,
		p_creator_id,
		p_deleted,
		p_deletor_id,
		p_duration,
		p_organization_id,
		p_transcription,
		p_transcription_user_edited
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose Down
DROP TABLE note_image_breadcrumb;
DROP TABLE note_image_data;
DROP TABLE note_image;
DROP TYPE ImageDataType;
DROP TABLE note_audio_breadcrumb;
DROP TABLE note_audio_data;
DROP TABLE note_audio;
DROP TYPE AudioDataType;
