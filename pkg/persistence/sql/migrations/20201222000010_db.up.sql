SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
-- SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

SET default_tablespace = '';
SET default_with_oids = false;



CREATE TABLE public.customers (
	id bigserial NOT NULL,
	first_name text NULL,
	last_name text NULL,
	birth_date timestamp NULL,
	gender text NULL,
	email text NULL,
	address text NULL,
	CONSTRAINT customers_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.customers OWNER TO customers;
GRANT ALL ON TABLE public.customers TO customers;
