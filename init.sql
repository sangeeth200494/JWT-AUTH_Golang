--
-- PostgreSQL database dump
--

\restrict WMyAe1f1arYJLdopy3YIDhabyxw5tyOVSEH5qA4o4rgCcbu7C5IHNkvzmavRVSl

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_username_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS public.users ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.users;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password text NOT NULL,
    full_name character varying(255),
    phone_number character varying(50),
    profile_picture_url text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_login timestamp without time zone,
    status character varying(50) DEFAULT 'active'::character varying,
    role character varying(50) DEFAULT 'user'::character varying,
    two_factor_enabled boolean DEFAULT false,
    failed_login_attempts integer DEFAULT 0,
    lock_status boolean DEFAULT false,
    deleted_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password, full_name, phone_number, profile_picture_url, created_at, updated_at, last_login, status, role, two_factor_enabled, failed_login_attempts, lock_status, deleted_at) FROM stdin;
2	Anagha C	anagham477@gmail.com	$2a$10$VguTdRA.qlPc0yfcm7zUzu15mZtTVVNBN0VtTKpStqoF5AY7BEGeK	Anagha Chadayan	9567702458		2026-07-09 08:15:53.011127	2026-07-09 08:15:53.011127	0001-01-01 00:00:00	true	user	t	0	f	\N
1	sangeeth_jayaraj	sangeethvv554@gmail.com	$2a$10$9TeBtMHHVn94sCbUU3YbH.uBJx2AMJhdSOa/H8vZqpb9ZBCZ.BG0q	Sangeeth Jayaraj	9567386525		2026-07-09 02:28:59.907427	2026-07-09 02:28:59.907427	0001-01-01 00:00:00	true	admin	t	0	f	\N
3	Mahitha k	mahithajay@gmail.com	$2a$10$MTC3sSWZnxMaTDJ1s4FuieB1D6MVC.nrKEYGiWgnTxQ9.Lf98/gDK	Mahitha Kadavathupurayil	9895268988		2026-07-10 03:44:11.034255	2026-07-10 03:44:11.034255	0001-01-01 00:00:00	true	user	t	0	f	\N
11	test_user_6427	test_user_6427@example.com	$2a$10$BTzK76eKZJ1MOlzxtG4DNutFKdcUhG3qWfQb/uKo6bdfVOghG2u4K	Automated Test User			2026-11-22 09:39:00.489108	2026-11-22 09:39:00.489108	0001-01-01 00:00:00	active	user	f	0	f	\N
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- PostgreSQL database dump complete
--

\unrestrict WMyAe1f1arYJLdopy3YIDhabyxw5tyOVSEH5qA4o4rgCcbu7C5IHNkvzmavRVSl

