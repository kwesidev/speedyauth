--
-- PostgreSQL database dump
--

-- Dumped from database version 12.13 (Ubuntu 12.13-1.pgdg20.04+1)
-- Dumped by pg_dump version 15.1 (Ubuntu 15.1-1.pgdg20.04+1)

-- Started on 2023-09-15 08:35:49 SAST

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
-- TOC entry 2 (class 3079 OID 58366)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 3059 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 212 (class 1259 OID 41964)
-- Name: reset_password_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reset_password_requests (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    code character varying NOT NULL,
    expiry_time timestamp without time zone NOT NULL,
    created timestamp without time zone NOT NULL
);


ALTER TABLE public.reset_password_requests OWNER TO postgres;

--
-- TOC entry 211 (class 1259 OID 41962)
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
-- TOC entry 3060 (class 0 OID 0)
-- Dependencies: 211
-- Name: reset_password_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reset_password_requests_id_seq OWNED BY public.reset_password_requests.id;


--
-- TOC entry 206 (class 1259 OID 33719)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    type character varying NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 205 (class 1259 OID 33717)
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
-- TOC entry 3061 (class 0 OID 0)
-- Dependencies: 205
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 215 (class 1259 OID 107581)
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 58426)
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
-- TOC entry 213 (class 1259 OID 58424)
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
-- TOC entry 3062 (class 0 OID 0)
-- Dependencies: 213
-- Name: two_factor_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.two_factor_requests_id_seq OWNED BY public.two_factor_requests.id;


--
-- TOC entry 210 (class 1259 OID 41947)
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
-- TOC entry 209 (class 1259 OID 41945)
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
-- TOC entry 3063 (class 0 OID 0)
-- Dependencies: 209
-- Name: user_refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_refresh_tokens_id_seq OWNED BY public.user_refresh_tokens.id;


--
-- TOC entry 208 (class 1259 OID 33730)
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_roles (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL
);


ALTER TABLE public.user_roles OWNER TO postgres;

--
-- TOC entry 207 (class 1259 OID 33728)
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
-- TOC entry 3064 (class 0 OID 0)
-- Dependencies: 207
-- Name: user_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_roles_id_seq OWNED BY public.user_roles.id;


--
-- TOC entry 204 (class 1259 OID 33692)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying NOT NULL,
    password character varying NOT NULL,
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    email_address character varying NOT NULL,
    phone_number character varying NOT NULL,
    active boolean NOT NULL,
    uu_id uuid DEFAULT public.uuid_generate_v4(),
    meta_data character varying,
    two_factor_enabled boolean
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 33690)
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
-- TOC entry 3065 (class 0 OID 0)
-- Dependencies: 203
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 2885 (class 2604 OID 41967)
-- Name: reset_password_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests ALTER COLUMN id SET DEFAULT nextval('public.reset_password_requests_id_seq'::regclass);


--
-- TOC entry 2882 (class 2604 OID 33722)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 2886 (class 2604 OID 58429)
-- Name: two_factor_requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests ALTER COLUMN id SET DEFAULT nextval('public.two_factor_requests_id_seq'::regclass);


--
-- TOC entry 2884 (class 2604 OID 41950)
-- Name: user_refresh_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.user_refresh_tokens_id_seq'::regclass);


--
-- TOC entry 2883 (class 2604 OID 33733)
-- Name: user_roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles ALTER COLUMN id SET DEFAULT nextval('public.user_roles_id_seq'::regclass);


--
-- TOC entry 2880 (class 2604 OID 33695)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3049 (class 0 OID 41964)
-- Dependencies: 212
-- Data for Name: reset_password_requests; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reset_password_requests (id, user_id, code, expiry_time, created) FROM stdin;
\.


--
-- TOC entry 3043 (class 0 OID 33719)
-- Dependencies: 206
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, type) FROM stdin;
1	ADMIN
2	USER
\.


--
-- TOC entry 3052 (class 0 OID 107581)
-- Dependencies: 215
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
1	f
\.


--
-- TOC entry 3051 (class 0 OID 58426)
-- Dependencies: 214
-- Data for Name: two_factor_requests; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.two_factor_requests (id, user_id, request_id, ip_address, user_agent, code, created_at, expiry_time) FROM stdin;
61	5	585ce7b90c6d470bd277714751628710f2ee8b4c			813447	2023-08-09 07:46:35.423047	2023-08-09 07:51:35.42318
64	7	b28a01aa47e1049552b9349c77e319d819307e58			781817	2023-08-09 18:36:41.162146	2023-08-09 18:41:41.162312
65	5	f8b894a87c6dcaadaeaf7fa2880adbec298eda02			026573	2023-08-09 21:24:31.812427	2023-08-09 21:29:31.813812
68	5	1752b60b18c7a87e1405bf1ab814b84a4be98352			422666	2023-08-10 14:32:41.100978	2023-08-10 14:37:41.102013
74	5	5f2500ff7caced204fb644e6727dc7eda44b7ab5			871214	2023-08-14 14:08:17.416056	2023-08-14 14:13:17.416512
76	5	d1081537de4331a3d637c0508b260cd356bec23d			304181	2023-08-30 10:51:49.234786	2023-08-30 10:56:49.234982
77	5	74a9be17175ef220fb49a5cef45866cb36b5ffbf			033357	2023-08-31 08:56:28.952844	2023-08-31 09:01:28.95328
95	5	dd38b5b0c2ee47c05b52e169f1a4e063175141d8			358710	2023-08-31 10:26:47.713834	2023-08-31 10:31:47.714007
\.


--
-- TOC entry 3047 (class 0 OID 41947)
-- Dependencies: 210
-- Data for Name: user_refresh_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_refresh_tokens (id, user_id, token, created, ip_address, user_agent, expiry_time) FROM stdin;
1946	5	b16295909889c41468ae69a0c35cb479d57ffff0	2023-08-10 14:07:46.506524	127.0.0.1:40586		2023-08-10 14:22:46.506849
2708	5	b9abad2eb1f375b6eae63f89f2df11b7290b4be6	2023-09-01 08:46:33.172018			2023-09-01 09:01:33.171864
2709	5	b990bb586cd930aa146561d3776ec1440ce70948	2023-09-10 10:29:18.324171			2023-09-10 10:44:18.32396
2710	5	a516c542465e5b02ce1575186eb74824b89e81ab	2023-09-10 14:06:57.273727			2023-09-10 14:06:57.273556
2711	5	736bbdb9a1bbc1b6387df43446add88aa0a45fb8	2023-09-10 14:07:06.342452			2023-09-10 14:07:06.342242
2712	5	85c5e3e9621d7ea67b50de2ceb98bd75514afd97	2023-09-10 14:09:20.165465			2023-09-10 14:09:20.165291
2713	5	ee0094965e51a416494f0477536fa434f52e245d	2023-09-10 19:59:23.92016			2023-09-10 19:59:23.913406
2024	5	484e0f1d34ff6ff75786dfbc16f2a91649b714e6	2023-08-14 10:44:18.551227	127.0.0.1:52864		2023-08-14 10:59:18.55147
1868	5	3a5eac29bc089e48f1e10de04a2a4c18b0c354d3	2023-08-09 21:43:50.299511	127.0.0.1:39714		2023-08-09 21:58:50.299719
1029	5	3032633173af150f2a589787ed066cbe035e98bf	2023-08-06 08:28:24.986244	127.0.0.1:40992		2023-08-06 08:43:24.986483
2507	5	c7623dc53b7df88fadacb45a790ef7db856123c8	2023-08-15 13:07:05.844344	127.0.0.1:35170		2023-08-15 13:22:05.844766
\.


--
-- TOC entry 3045 (class 0 OID 33730)
-- Dependencies: 208
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_roles (id, user_id, role_id) FROM stdin;
2	3	2
3	5	2
1	2	1
4	7	2
5	20	2
6	24	2
\.


--
-- TOC entry 3041 (class 0 OID 33692)
-- Dependencies: 204
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, first_name, last_name, email_address, phone_number, active, uu_id, meta_data, two_factor_enabled) FROM stdin;
7	prince	$2a$10$qlK90Ys6WWHOMkTgKJCzQ.yafSy31n6pHao8vsIGGHkgvA/l9I6qC	Benjamin 	Owusu Koti botor	oprinnce61@yahoo.com		t	aeba0411-98b6-4858-9f8a-2a421a309c64	\N	t
5	jackie	$2a$10$nBtWC.KOJyDqizRu9Sv14OXn84rPjZtlPCRxB/ahTxcM0YwISg8IW	Kwesi	Jones	willzako@aol.com	0738455979	t	0fbdc180-6fc1-42bc-a7d7-0c6b2a065bbe	\N	f
20	tupac	$2a$10$IS0BMDyzJYfPiKbN4cMyKeKANJAIV5lVdzYc4PaVUYYg4t98SLx3e	Tupac	Amaru Shakur	tupac@gmail.com	0660517444	t	28bfa7c7-2c46-4dd4-8f52-2d1bbc525ccd	\N	f
24	jackson	$2a$10$mqjPVBm3opIi6kEp2zg/Ve6c6Iar4.yLDNT0RPPuYW8JCuJEtjrMa	jackie	Opoku	jackson@gmail.com	0731482947	t	56c7a965-ef6b-44b0-9aae-f64da3132a4f	\N	f
3	apalo	$2b$10$4KAofVbTqHGRrE8YQfWLieFnizEcqNDVyEvkbjiGJmGyGfiNVhWq.	Jane	Doe	viljoend0@gmail.com	010101010101	t	9f724883-ed0a-4716-9ca9-2666f0da9cef	\N	f
2	kwesidev	$2b$10$Nt0zoFFWzEIbWYhn6baF6OF.gmhj9Ew2bsMemIg6Qz165RrORQxnK	William	Akomaning	william@lexpro.co.za	0660517444	t	ab54de7b-930b-4306-aa78-11c4b1a18bff	\N	t
\.


--
-- TOC entry 3066 (class 0 OID 0)
-- Dependencies: 211
-- Name: reset_password_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reset_password_requests_id_seq', 86, true);


--
-- TOC entry 3067 (class 0 OID 0)
-- Dependencies: 205
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 2, true);


--
-- TOC entry 3068 (class 0 OID 0)
-- Dependencies: 213
-- Name: two_factor_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.two_factor_requests_id_seq', 105, true);


--
-- TOC entry 3069 (class 0 OID 0)
-- Dependencies: 209
-- Name: user_refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_refresh_tokens_id_seq', 2713, true);


--
-- TOC entry 3070 (class 0 OID 0)
-- Dependencies: 207
-- Name: user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_roles_id_seq', 6, true);


--
-- TOC entry 3071 (class 0 OID 0)
-- Dependencies: 203
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 31, true);


--
-- TOC entry 2898 (class 2606 OID 41974)
-- Name: reset_password_requests reset_password_requests_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_code_key UNIQUE (code);


--
-- TOC entry 2900 (class 2606 OID 41972)
-- Name: reset_password_requests reset_password_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_pkey PRIMARY KEY (id);


--
-- TOC entry 2892 (class 2606 OID 33727)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 2908 (class 2606 OID 107585)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 2902 (class 2606 OID 58438)
-- Name: two_factor_requests two_factor_requests_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_code_key UNIQUE (code);


--
-- TOC entry 2904 (class 2606 OID 58434)
-- Name: two_factor_requests two_factor_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_pkey PRIMARY KEY (id);


--
-- TOC entry 2906 (class 2606 OID 58436)
-- Name: two_factor_requests two_factor_requests_request_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_request_id_key UNIQUE (request_id);


--
-- TOC entry 2888 (class 2606 OID 33747)
-- Name: users unique_value; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_value UNIQUE (username, email_address);


--
-- TOC entry 2896 (class 2606 OID 41955)
-- Name: user_refresh_tokens user_refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens
    ADD CONSTRAINT user_refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 2894 (class 2606 OID 33735)
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- TOC entry 2890 (class 2606 OID 33700)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 2912 (class 2606 OID 41975)
-- Name: reset_password_requests reset_password_requests_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reset_password_requests
    ADD CONSTRAINT reset_password_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2913 (class 2606 OID 58439)
-- Name: two_factor_requests two_factor_requests_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.two_factor_requests
    ADD CONSTRAINT two_factor_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2911 (class 2606 OID 41956)
-- Name: user_refresh_tokens user_refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_refresh_tokens
    ADD CONSTRAINT user_refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2909 (class 2606 OID 33741)
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 2910 (class 2606 OID 33736)
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3058 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2023-09-15 08:35:49 SAST

--
-- PostgreSQL database dump complete
--

