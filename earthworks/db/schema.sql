SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry, geography, and raster spatial types and functions';


--
-- Name: lab_test_code; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.lab_test_code AS ENUM (
    'moisture_content',
    'grain_size_analysis',
    'hydrometer',
    'proctor',
    'atterberg'
);


--
-- Name: project_file_code; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.project_file_code AS ENUM (
    'report',
    'lab_report',
    'calculation',
    'proposal',
    'budget',
    'field_data',
    'other'
);


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: borehole; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.borehole (
    id integer NOT NULL,
    project integer NOT NULL,
    datapoint integer NOT NULL,
    program integer,
    name text NOT NULL,
    start_date date NOT NULL,
    end_date date,
    field_eng text NOT NULL,
    type text DEFAULT 'borehole'::text,
    drilling_method text,
    CONSTRAINT borehole_drilling_method_check CHECK ((char_length(drilling_method) < 30)),
    CONSTRAINT borehole_field_eng_check CHECK ((char_length(field_eng) < 80)),
    CONSTRAINT borehole_name_check CHECK ((char_length(name) < 40)),
    CONSTRAINT borehole_type_check CHECK ((char_length(type) < 30))
);


--
-- Name: borehole_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.borehole_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: borehole_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.borehole_id_seq OWNED BY public.borehole.id;


--
-- Name: borehole_type; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.borehole_type (
    code text NOT NULL,
    description text,
    CONSTRAINT borehole_type_code_check CHECK ((char_length(code) < 30)),
    CONSTRAINT borehole_type_description_check CHECK ((char_length(description) < 100))
);


--
-- Name: datapoint; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.datapoint (
    id integer NOT NULL,
    location public.geography(Point,4326) NOT NULL
);


--
-- Name: datapoint_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.datapoint_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: datapoint_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.datapoint_id_seq OWNED BY public.datapoint.id;


--
-- Name: drilling_method; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.drilling_method (
    code text NOT NULL,
    description text,
    CONSTRAINT drilling_method_code_check CHECK ((char_length(code) < 30)),
    CONSTRAINT drilling_method_description_check CHECK ((char_length(description) < 100))
);


--
-- Name: field_program; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.field_program (
    id integer NOT NULL,
    project integer NOT NULL,
    start_date date NOT NULL,
    end_date date
);


--
-- Name: field_program_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.field_program_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: field_program_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.field_program_id_seq OWNED BY public.field_program.id;


--
-- Name: gsa_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gsa_data (
    id integer NOT NULL,
    test integer NOT NULL,
    pan boolean NOT NULL,
    size double precision NOT NULL,
    mass_retained double precision
);


--
-- Name: gsa_data_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.gsa_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gsa_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.gsa_data_id_seq OWNED BY public.gsa_data.id;


--
-- Name: gsa_test; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gsa_test (
    id integer NOT NULL,
    tare_mass double precision,
    dry_plus_tare double precision,
    washed_plus_tare double precision
);


--
-- Name: instrument; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.instrument (
    id integer NOT NULL,
    project integer NOT NULL,
    datapoint integer NOT NULL,
    program integer,
    name text NOT NULL,
    device_id uuid,
    install_date date NOT NULL,
    field_eng text NOT NULL,
    CONSTRAINT instrument_field_eng_check CHECK ((char_length(field_eng) < 80)),
    CONSTRAINT instrument_name_check CHECK ((char_length(name) < 40))
);


--
-- Name: instrument_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.instrument_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: instrument_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.instrument_id_seq OWNED BY public.instrument.id;


--
-- Name: lab_test; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lab_test (
    id integer NOT NULL,
    name text,
    type public.lab_test_code NOT NULL,
    start_date date,
    end_date date,
    performed_by text,
    sample integer NOT NULL,
    checked_date date,
    checked_by text,
    CONSTRAINT lab_test_checked_by_check CHECK ((char_length(checked_by) < 200)),
    CONSTRAINT lab_test_name_check CHECK ((char_length(name) < 100)),
    CONSTRAINT lab_test_performed_by_check CHECK ((char_length(performed_by) < 200))
);


--
-- Name: lab_test_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.lab_test_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: lab_test_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.lab_test_id_seq OWNED BY public.lab_test.id;


--
-- Name: moisture_test; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.moisture_test (
    id integer NOT NULL,
    tare_mass double precision,
    sample_plus_tare double precision,
    dry_plus_tare double precision
);


--
-- Name: project; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project (
    id integer NOT NULL,
    name text NOT NULL,
    number text NOT NULL,
    client text NOT NULL,
    pm text NOT NULL,
    location text NOT NULL,
    organization text,
    default_coords public.geography(Point,4326) NOT NULL,
    CONSTRAINT project_client_check CHECK ((char_length(client) < 255)),
    CONSTRAINT project_location_check CHECK ((char_length(location) < 255)),
    CONSTRAINT project_name_check CHECK ((char_length(name) < 255)),
    CONSTRAINT project_number_check CHECK ((char_length(number) < 255)),
    CONSTRAINT project_organization_check CHECK ((char_length(organization) < 255)),
    CONSTRAINT project_pm_check CHECK ((char_length(pm) < 255))
);


--
-- Name: project_file; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project_file (
    id integer NOT NULL,
    project integer NOT NULL,
    category public.project_file_code NOT NULL,
    file bytea NOT NULL,
    filename text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    created_by text NOT NULL,
    expired_at timestamp without time zone,
    CONSTRAINT project_file_created_by_check CHECK ((char_length(created_by) < 250)),
    CONSTRAINT project_file_filename_check CHECK ((char_length(filename) < 250))
);


--
-- Name: project_file_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.project_file_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: project_file_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.project_file_id_seq OWNED BY public.project_file.id;


--
-- Name: project_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.project_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: project_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.project_id_seq OWNED BY public.project.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: soil_sample; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.soil_sample (
    id integer NOT NULL,
    name text NOT NULL,
    borehole integer NOT NULL,
    start_depth double precision NOT NULL,
    end_depth double precision NOT NULL,
    description text NOT NULL,
    uscs text NOT NULL,
    CONSTRAINT soil_sample_description_check CHECK ((char_length(description) < 800)),
    CONSTRAINT soil_sample_name_check CHECK ((char_length(name) < 100)),
    CONSTRAINT soil_sample_uscs_check CHECK ((char_length(uscs) < 20))
);


--
-- Name: soil_sample_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.soil_sample_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: soil_sample_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.soil_sample_id_seq OWNED BY public.soil_sample.id;


--
-- Name: strata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.strata (
    id integer NOT NULL,
    borehole integer NOT NULL,
    start_depth double precision NOT NULL,
    end_depth double precision NOT NULL,
    description text NOT NULL,
    soils text NOT NULL,
    moisture text,
    consistency text,
    CONSTRAINT strata_consistency_check CHECK ((char_length(consistency) < 50)),
    CONSTRAINT strata_description_check CHECK ((char_length(description) < 800)),
    CONSTRAINT strata_moisture_check CHECK ((char_length(moisture) < 50)),
    CONSTRAINT strata_soils_check CHECK ((char_length(soils) < 200))
);


--
-- Name: strata_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.strata_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: strata_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.strata_id_seq OWNED BY public.strata.id;


--
-- Name: time_series_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.time_series_data (
    id integer NOT NULL,
    device_id uuid NOT NULL,
    instrument integer,
    series integer DEFAULT 0,
    "time" timestamp with time zone,
    value double precision
);


--
-- Name: time_series_data_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.time_series_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: time_series_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.time_series_data_id_seq OWNED BY public.time_series_data.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username text NOT NULL,
    CONSTRAINT users_username_check CHECK ((char_length(username) < 40))
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: borehole id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole ALTER COLUMN id SET DEFAULT nextval('public.borehole_id_seq'::regclass);


--
-- Name: datapoint id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.datapoint ALTER COLUMN id SET DEFAULT nextval('public.datapoint_id_seq'::regclass);


--
-- Name: field_program id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.field_program ALTER COLUMN id SET DEFAULT nextval('public.field_program_id_seq'::regclass);


--
-- Name: gsa_data id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gsa_data ALTER COLUMN id SET DEFAULT nextval('public.gsa_data_id_seq'::regclass);


--
-- Name: instrument id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument ALTER COLUMN id SET DEFAULT nextval('public.instrument_id_seq'::regclass);


--
-- Name: lab_test id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lab_test ALTER COLUMN id SET DEFAULT nextval('public.lab_test_id_seq'::regclass);


--
-- Name: project id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project ALTER COLUMN id SET DEFAULT nextval('public.project_id_seq'::regclass);


--
-- Name: project_file id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project_file ALTER COLUMN id SET DEFAULT nextval('public.project_file_id_seq'::regclass);


--
-- Name: soil_sample id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.soil_sample ALTER COLUMN id SET DEFAULT nextval('public.soil_sample_id_seq'::regclass);


--
-- Name: strata id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.strata ALTER COLUMN id SET DEFAULT nextval('public.strata_id_seq'::regclass);


--
-- Name: time_series_data id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.time_series_data ALTER COLUMN id SET DEFAULT nextval('public.time_series_data_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: borehole borehole_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_pkey PRIMARY KEY (id);


--
-- Name: borehole borehole_project_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_project_name_key UNIQUE (project, name);


--
-- Name: borehole_type borehole_type_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole_type
    ADD CONSTRAINT borehole_type_pkey PRIMARY KEY (code);


--
-- Name: datapoint datapoint_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.datapoint
    ADD CONSTRAINT datapoint_pkey PRIMARY KEY (id);


--
-- Name: drilling_method drilling_method_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drilling_method
    ADD CONSTRAINT drilling_method_pkey PRIMARY KEY (code);


--
-- Name: field_program field_program_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.field_program
    ADD CONSTRAINT field_program_pkey PRIMARY KEY (id);


--
-- Name: gsa_data gsa_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gsa_data
    ADD CONSTRAINT gsa_data_pkey PRIMARY KEY (id);


--
-- Name: gsa_test gsa_test_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gsa_test
    ADD CONSTRAINT gsa_test_pkey PRIMARY KEY (id);


--
-- Name: instrument instrument_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument
    ADD CONSTRAINT instrument_pkey PRIMARY KEY (id);


--
-- Name: instrument instrument_project_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument
    ADD CONSTRAINT instrument_project_name_key UNIQUE (project, name);


--
-- Name: lab_test lab_test_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lab_test
    ADD CONSTRAINT lab_test_pkey PRIMARY KEY (id);


--
-- Name: moisture_test moisture_test_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.moisture_test
    ADD CONSTRAINT moisture_test_pkey PRIMARY KEY (id);


--
-- Name: project_file project_file_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project_file
    ADD CONSTRAINT project_file_pkey PRIMARY KEY (id);


--
-- Name: project project_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: soil_sample soil_sample_name_borehole_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.soil_sample
    ADD CONSTRAINT soil_sample_name_borehole_key UNIQUE (name, borehole);


--
-- Name: soil_sample soil_sample_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.soil_sample
    ADD CONSTRAINT soil_sample_pkey PRIMARY KEY (id);


--
-- Name: strata strata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.strata
    ADD CONSTRAINT strata_pkey PRIMARY KEY (id);


--
-- Name: time_series_data time_series_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.time_series_data
    ADD CONSTRAINT time_series_data_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: pan_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX pan_idx ON public.gsa_data USING btree (test) WHERE pan;


--
-- Name: borehole borehole_datapoint_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_datapoint_fkey FOREIGN KEY (datapoint) REFERENCES public.datapoint(id);


--
-- Name: borehole borehole_drilling_method_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_drilling_method_fkey FOREIGN KEY (drilling_method) REFERENCES public.drilling_method(code) ON DELETE SET NULL;


--
-- Name: borehole borehole_program_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_program_fkey FOREIGN KEY (program) REFERENCES public.field_program(id);


--
-- Name: borehole borehole_project_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_project_fkey FOREIGN KEY (project) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: borehole borehole_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borehole
    ADD CONSTRAINT borehole_type_fkey FOREIGN KEY (type) REFERENCES public.borehole_type(code) ON DELETE RESTRICT;


--
-- Name: field_program field_program_project_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.field_program
    ADD CONSTRAINT field_program_project_fkey FOREIGN KEY (project) REFERENCES public.project(id);


--
-- Name: gsa_data gsa_data_test_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gsa_data
    ADD CONSTRAINT gsa_data_test_fkey FOREIGN KEY (test) REFERENCES public.gsa_test(id) ON DELETE CASCADE;


--
-- Name: gsa_test gsa_test_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gsa_test
    ADD CONSTRAINT gsa_test_id_fkey FOREIGN KEY (id) REFERENCES public.lab_test(id) ON DELETE CASCADE;


--
-- Name: instrument instrument_datapoint_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument
    ADD CONSTRAINT instrument_datapoint_fkey FOREIGN KEY (datapoint) REFERENCES public.datapoint(id);


--
-- Name: instrument instrument_program_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument
    ADD CONSTRAINT instrument_program_fkey FOREIGN KEY (program) REFERENCES public.field_program(id);


--
-- Name: instrument instrument_project_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.instrument
    ADD CONSTRAINT instrument_project_fkey FOREIGN KEY (project) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: lab_test lab_test_sample_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lab_test
    ADD CONSTRAINT lab_test_sample_fkey FOREIGN KEY (sample) REFERENCES public.soil_sample(id) ON DELETE CASCADE;


--
-- Name: moisture_test moisture_test_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.moisture_test
    ADD CONSTRAINT moisture_test_id_fkey FOREIGN KEY (id) REFERENCES public.lab_test(id) ON DELETE CASCADE;


--
-- Name: project_file project_file_project_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project_file
    ADD CONSTRAINT project_file_project_fkey FOREIGN KEY (project) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: soil_sample soil_sample_borehole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.soil_sample
    ADD CONSTRAINT soil_sample_borehole_fkey FOREIGN KEY (borehole) REFERENCES public.borehole(id) ON DELETE CASCADE;


--
-- Name: strata strata_borehole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.strata
    ADD CONSTRAINT strata_borehole_fkey FOREIGN KEY (borehole) REFERENCES public.borehole(id) ON DELETE CASCADE;


--
-- Name: time_series_data time_series_data_instrument_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.time_series_data
    ADD CONSTRAINT time_series_data_instrument_fkey FOREIGN KEY (instrument) REFERENCES public.instrument(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20190403054412'),
    ('20190403054428'),
    ('20190403054441'),
    ('20190403054445'),
    ('20190403054506'),
    ('20190403054511'),
    ('20190419065059'),
    ('20190422034654'),
    ('20190607042011');
