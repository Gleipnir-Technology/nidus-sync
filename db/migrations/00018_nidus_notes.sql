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

-- +goose Down
DROP TABLE note_image_breadcrumb;
DROP TABLE note_image_data;
DROP TABLE note_image;
DROP TYPE ImageDataType
DROP TABLE note_audio_breadcrumb;
DROP TABLE note_audio_data;
DROP TABLE note_audio;
DROP TYPE AudioDataType;
