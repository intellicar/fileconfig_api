
-- DROP TRIGGER IF EXISTS trigger_update_release_config ON "FileRelease"."ReleaseConfig";
-- DROP TRIGGER IF EXISTS trigger_soft_delete_release_config ON "FileRelease"."ReleaseConfig";
-- DROP TRIGGER IF EXISTS trigger_insert_release_config ON "FileRelease"."ReleaseConfig";
-- DROP TABLE IF EXISTS "FileRelease"."ReleaseConfigHistory";


CREATE TABLE IF NOT EXISTS
 "FileRelease"."ReleaseConfig" (
    configid SERIAL PRIMARY KEY, 
    group TEXT,
    grouptype TEXT,
    boardversion VARCHAR(255),
    releasedate BIGINT NOT NULL, 
    filename TEXT UNIQUE NOT NULL , 
    sim VARCHAR NOT NULL,
    nrfbootloader TEXT, 
    releasenote TEXT, 
    islatest BOOLEAN DEFAULT false, 
    isvalid BOOLEAN DEFAULT false, 
    createdby VARCHAR NOT NULL,
    updatedby VARCHAR,
    updatedat bigint DEFAULT CAST(
            extract(
                epoch
                FROM
                    NOW()
            ) * 1000 AS bigint
        ) NOT NULL,, 
    isdelete BOOLEAN DEFAULT false
);


-- Table for ReleaseConfigHistory
CREATE TABLE IF NOT EXISTS
 "FileRelease"."ReleaseConfigHistory" (
    historyid SERIAL PRIMARY KEY,
    configid INTEGER,
    updatedby VARCHAR,
    updatedat BIGINT,
    operation VARCHAR,
    fieldchanged TEXT,
    olddata JSONB,
    newdata JSONB
);

-- Function for UPDATE
CREATE OR REPLACE FUNCTION log_release_config_update()
RETURNS TRIGGER AS $$
DECLARE
    old_row JSONB;
    new_row JSONB;
    field_changes TEXT := '';
BEGIN
    old_row := to_jsonb(OLD) - 'historyid';  
    new_row := to_jsonb(NEW) - 'historyid'; 

    IF (NEW.group IS DISTINCT FROM OLD.group) THEN
        field_changes := field_changes || 'group, ';
    END IF;

    IF (NEW.grouptype IS DISTINCT FROM OLD.grouptype) THEN
        field_changes := field_changes || 'grouptype, ';
    END IF;

    IF (NEW.boardversion IS DISTINCT FROM OLD.boardversion) THEN
        field_changes := field_changes || 'boardversion, ';
    END IF;

    IF (NEW.releasedate IS DISTINCT FROM OLD.releasedate) THEN
        field_changes := field_changes || 'releaseDate, ';
    END IF;

    IF (NEW.filename IS DISTINCT FROM OLD.filename) THEN
        field_changes := field_changes || 'filename, ';
    END IF;

    IF (NEW.sim IS DISTINCT FROM OLD.sim) THEN
        field_changes := field_changes || 'sim, ';
    END IF;

    IF (NEW.nrfbootloader IS DISTINCT FROM OLD.nrfbootloader) THEN
        field_changes := field_changes || 'nrfbootloader, ';
    END IF;

    IF (NEW.releasenote IS DISTINCT FROM OLD.releasenote) THEN
        field_changes := field_changes || 'releasenote, ';
    END IF;

    IF (NEW.islatest IS DISTINCT FROM OLD.islatest) THEN
        field_changes := field_changes || 'islatest, ';
    END IF;

    IF (NEW.isvalid IS DISTINCT FROM OLD.isvalid) THEN
        field_changes := field_changes || 'isvalid, ';
    END IF;

    IF (NEW.updatedby IS DISTINCT FROM OLD.updatedby) THEN
        field_changes := field_changes || 'updatedby, ';
    END IF;

    IF field_changes <> '' THEN
        field_changes := trim(trailing ', ' from field_changes);

        INSERT INTO "FileRelease"."ReleaseConfigHistory" 
        (configid, updatedby, updatedat, operation, fieldchanged, olddata, newdata)
        VALUES 
        (NEW.configid, NEW.updatedby, EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT, 'UPDATE', field_changes, old_row, new_row);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for UPDATE
CREATE TRIGGER trigger_update_release_config
AFTER UPDATE ON "FileRelease"."ReleaseConfig"
FOR EACH ROW
EXECUTE FUNCTION log_release_config_update();

-- Function for INSERT
CREATE OR REPLACE FUNCTION log_release_config_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "FileRelease"."ReleaseConfigHistory" (
        configid, 
        updatedby, 
        updatedat, 
        operation, 
        fieldchanged, 
        olddata, 
        newdata
    )
    VALUES (
        NEW.configid, 
        NEW.createdby, 
        EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT, 
        'INSERT', 
        NULL, 
        NULL, 
        row_to_json(NEW)::jsonb
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for INSERT
CREATE TRIGGER trigger_insert_release_config
AFTER INSERT ON "FileRelease"."ReleaseConfig"
FOR EACH ROW
EXECUTE FUNCTION log_release_config_insert();

-- Function for SOFT DELETE
CREATE OR REPLACE FUNCTION log_release_config_soft_delete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.isdelete = true AND OLD.isdelete = false THEN
        INSERT INTO "FileRelease"."ReleaseConfigHistory" (
            configid, 
            updatedby, 
            updatedat, 
            operation, 
            fieldchanged, 
            olddata, 
            newdata
        )
        VALUES (
            NEW.configid, 
            NEW.updatedby, 
            EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT, 
            'DELETE', 
            'isdelete', 
            row_to_json(OLD)::jsonb, 
            row_to_json(NEW)::jsonb  
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for SOFT DELETE
CREATE TRIGGER trigger_soft_delete_release_config
AFTER UPDATE ON "FileRelease"."ReleaseConfig"
FOR EACH ROW
EXECUTE FUNCTION log_release_config_soft_delete();
