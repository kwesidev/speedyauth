--
-- PostgreSQL database dump
--

-- Dumped from database version 12.13 (Ubuntu 12.13-1.pgdg20.04+1)
-- Dumped by pg_dump version 15.1 (Ubuntu 15.1-1.pgdg20.04+1)

-- Started on 2023-09-15 08:58:22 SAST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 7 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 2 (class 3079 OID 107598)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 3061 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 115831)
-- Name: reset_password_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reset_password_requests (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    code character varying NOT NULL,
    created timestamp without time zone NOT NULL,
    expiry_time timestamp without time zone NOT NULL
);


ALTER TABLE public.reset_password_requests OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 115829)
-- Name: reset_password_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reset_password_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reset_password_requests_id_seq OWNER TO postgres;

--
-- TOC entry 3062 (class 0 OID 0)
-- Dependencies: 214
-- Name: reset_password_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reset_password_requests_id_seq OWNED BY public.reset_password_requests.id;


--
-- TOC entry 211 (class 1259 OID 115802)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    type character varying NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 210 (class 1259 OID 115800)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO postgres;

--
-- TOC entry 3063 (class 0 OID 0)
-- Dependencies: 210
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 203 (class 1259 OID 107591)
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- TOC entry 207 (class 1259 OID 115766)
-- Name: two_factor_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.two_factor_requests (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    request_id character varying NOT NULL,
    ip_address character varying,
    user_agent character varying,
    code character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    expiry_time timestamp without time zone NOT NULL
);


ALTER TABLE public.two_factor_requests OWNER TO postgres;

--
-- TOC entry 206 (class 1259 OID 115764)
-- Name: two_factor_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.two_factor_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.two_factor_requests_id_seq OWNER TO postgres;

--
-- TOC entry 3064 (class 0 OID 0)
-- Dependencies: 206
-- Name: two_factor_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.two_factor_requests_id_seq OWNED BY public.two_factor_requests.id;


--
-- TOC entry 209 (class 1259 OID 115786)
-- Name: user_refresh_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_refresh_tokens (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    token character varying NOT NULL,
    created timestamp without time zone NOT NULL,
    ip_address character varying NOT NULL,
    user_agent character varying NOT NULL,
    expiry_time timestamp without time zone NOT NULL
);


ALTER TABLE public.user_refresh_tokens OWNER TO postgres;

--
-- TOC entry 208 (class 1259 OID 115784)
-- Name: user_refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_refresh_tokens_id_seq OWNER TO postgres;

--
-- TOC entry 3065 (class 0 OID 0)
-- Dependencies: 208
-- Name: user_refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_refresh_tokens_id_seq OWNED BY public.user_refresh_tokens.id;


--
-- TOC entry 213 (class 1259 OID 115813)
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_roles (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL
);


ALTER TABLE public.user_roles OWNER TO postgres;

--
-- TOC entry 212 (class 1259 OID 115811)
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_roles_id_seq OWNER TO postgres;

--
-- TOC entry 3066 (class 0 OID 0)
-- Dependencies: 212
-- Name: user_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_roles_id_seq OWNED BY public.user_roles.id;


--
-- TOC entry 205 (class 1259 OID 115750)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    uu_id uuid DEFAULT public.uuid_generate_v4(),
    username character varying NOT NULL,
    password character varying NOT NULL,
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    email_address character varying NOT NULL,
    phone_number character varying NOT NULL,
    active boolean NOT NULL,
    meta_data character varying,
    two_factor_enabled boolean
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 115748)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3067 (class 0 OID 0)
-- Dependencies: 204
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 2886 (class 2604 OID 115834)
-- Name: reset_password_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests ALTER COLUMN id SET DEFAULT nextval('public.reset_password_requests_id_seq'::regclass);


--
-- TOC entry 2884 (class 2604 OID 115805)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 2882 (class 2604 OID 115769)
-- Name: two_factor_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests ALTER COLUMN id SET DEFAULT nextval('public.two_factor_requests_id_seq'::regclass);


--
-- TOC entry 2883 (class 2604 OID 115789)
-- Name: user_refresh_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.user_refresh_tokens_id_seq'::regclass);


--
-- TOC entry 2885 (class 2604 OID 115816)
-- Name: user_roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles ALTER COLUMN id SET DEFAULT nextval('public.user_roles_id_seq'::regclass);


--
-- TOC entry 2880 (class 2604 OID 115753)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3054 (class 0 OID 115831)
-- Dependencies: 215
-- Data for Name: reset_password_requests; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reset_password_requests (id, user_id, code, created, expiry_time) FROM stdin;
\.


--
-- TOC entry 3050 (class 0 OID 115802)
-- Dependencies: 211
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, type) FROM stdin;
1	ADMIN
2	USER
\.


--
-- TOC entry 3042 (class 0 OID 107591)
-- Dependencies: 203
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
1	f
\.


--
-- TOC entry 3046 (class 0 OID 115766)
-- Dependencies: 207
-- Data for Name: two_factor_requests; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.two_factor_requests (id, user_id, request_id, ip_address, user_agent, code, created_at, expiry_time) FROM stdin;
\.


--
-- TOC entry 3048 (class 0 OID 115786)
-- Dependencies: 209
-- Data for Name: user_refresh_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_refresh_tokens (id, user_id, token, created, ip_address, user_agent, expiry_time) FROM stdin;
\.


--
-- TOC entry 3052 (class 0 OID 115813)
-- Dependencies: 213
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_roles (id, user_id, role_id) FROM stdin;
\.


--
-- TOC entry 3044 (class 0 OID 115750)
-- Dependencies: 205
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, uu_id, username, password, first_name, last_name, email_address, phone_number, active, meta_data, two_factor_enabled) FROM stdin;
\.


--
-- TOC entry 3068 (class 0 OID 0)
-- Dependencies: 214
-- Name: reset_password_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reset_password_requests_id_seq', 1, false);


--
-- TOC entry 3069 (class 0 OID 0)
-- Dependencies: 210
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 2, true);


--
-- TOC entry 3070 (class 0 OID 0)
-- Dependencies: 206
-- Name: two_factor_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.two_factor_requests_id_seq', 1, false);


--
-- TOC entry 3071 (class 0 OID 0)
-- Dependencies: 208
-- Name: user_refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_refresh_tokens_id_seq', 1, false);


--
-- TOC entry 3072 (class 0 OID 0)
-- Dependencies: 212
-- Name: user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_roles_id_seq', 1, false);


--
-- TOC entry 3073 (class 0 OID 0)
-- Dependencies: 204
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- TOC entry 2908 (class 2606 OID 115841)
-- Name: reset_password_requests reset_password_requests_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_code_key UNIQUE (code);


--
-- TOC entry 2910 (class 2606 OID 115839)
-- Name: reset_password_requests reset_password_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_pkey PRIMARY KEY (id);


--
-- TOC entry 2904 (class 2606 OID 115810)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 2888 (class 2606 OID 107595)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 2896 (class 2606 OID 115778)
-- Name: two_factor_requests two_factor_requests_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_code_key UNIQUE (code);


--
-- TOC entry 2898 (class 2606 OID 115774)
-- Name: two_factor_requests two_factor_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_pkey PRIMARY KEY (id);


--
-- TOC entry 2900 (class 2606 OID 115776)
-- Name: two_factor_requests two_factor_requests_request_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_request_id_key UNIQUE (request_id);


--
-- TOC entry 2902 (class 2606 OID 115794)
-- Name: user_refresh_tokens user_refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens
    ADD CONSTRAINT user_refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 2906 (class 2606 OID 115818)
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- TOC entry 2890 (class 2606 OID 115763)
-- Name: users users_email_address_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_address_key UNIQUE (email_address);


--
-- TOC entry 2892 (class 2606 OID 115759)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 2894 (class 2606 OID 115761)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- TOC entry 2915 (class 2606 OID 115842)
-- Name: reset_password_requests reset_password_requests_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2911 (class 2606 OID 115779)
-- Name: two_factor_requests two_factor_requests_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2912 (class 2606 OID 115795)
-- Name: user_refresh_tokens user_refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens
    ADD CONSTRAINT user_refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2913 (class 2606 OID 115824)
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 2914 (class 2606 OID 115819)
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3060 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2023-09-15 08:58:22 SAST

--
-- PostgreSQL database dump complete
--

