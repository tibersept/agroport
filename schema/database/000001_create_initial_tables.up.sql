CREATE TABLE agroport_user (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL
);

CREATE TABLE agroport_region (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL
);

CREATE TABLE agroport_field (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL,
    region_id INTEGER NOT NULL,

    CONSTRAINT fk_agroport_region FOREIGN KEY (region_id) REFERENCES agroport_region(id) ON DELETE CASCADE
);

CREATE TABLE agroport_operation (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL
);

CREATE TABLE agroport_op_on_day (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    op_day DATE NOT NULL,
    note VARCHAR,
    field_id INTEGER NOT NULL,
    op_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    CONSTRAINT fk_agroport_field FOREIGN KEY (field_id) REFERENCES agroport_field(id) ON DELETE CASCADE,
    CONSTRAINT fk_agroport_op FOREIGN KEY (op_id) REFERENCES agroport_operation(id) ON DELETE CASCADE,
    CONSTRAINT fk_agroport_user FOREIGN KEY (user_id) REFERENCES agroport_user(id) ON DELETE CASCADE
);
