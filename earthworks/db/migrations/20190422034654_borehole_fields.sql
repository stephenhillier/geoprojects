-- migrate:up
CREATE TABLE borehole_type (
    code TEXT PRIMARY KEY CHECK (char_length(code) < 30),
    description TEXT CHECK (char_length(description) < 100)
);

INSERT INTO borehole_type (code, description) VALUES
    ('borehole', 'Borehole'),
    ('test_pit', 'Test pit'),
    ('other', 'Other');

CREATE TABLE drilling_method (
    code TEXT PRIMARY KEY CHECK (char_length(code) < 30),
    description TEXT CHECK (char_length(description) < 100)
);

INSERT INTO drilling_method (code, description) VALUES
    ('air_rotary', 'Air rotary'),
    ('solid_stem', 'Solid stem auger'),
    ('hollow_stem', 'Hollow stem auger'),
    ('mud_rotary', 'Mud rotary'),
    ('excavator', 'Backhoe/excavator'),
    ('other', 'Other')
    ;

ALTER TABLE borehole
     ADD COLUMN type TEXT REFERENCES borehole_type(code) ON DELETE RESTRICT DEFAULT 'borehole' CHECK (char_length(type) < 30),
     ADD COLUMN drilling_method TEXT NULL REFERENCES drilling_method(code) ON DELETE SET NULL CHECK (char_length(drilling_method) < 30)
    ;

-- migrate:down

