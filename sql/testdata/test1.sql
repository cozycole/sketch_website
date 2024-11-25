--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)

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
-- Data for Name: person; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.person (id, slug, first, last, birthdate, profile_img, description) FROM stdin;
1	kyle-mooney-1	Kyle	Mooney	1984-09-03	kyle-mooney-1.jpg	\N
4	tim-gilbert-4	Tim	Gilbert	1983-05-13	tim-gilbert-4.jpg	this is the description
5	james-hartnett-5	James	Hartnett	\N	james-hartnett-5.jpg	\N
\.


--
-- Data for Name: character; Type: TABLE DATA; Schema: public; Owner: colet
--
COPY public."character" (id, name, description, img_name, person_id, slug) FROM stdin;
1	David S. Pumpkins	\N	\N	\N	david-s-pumpkins-1
2	Dave	\N	\N	\N	dave-2
\.


--
-- Data for Name: creator; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.creator (id, name, profile_img, page_url, date_established, slug, description) FROM stdin;
1	nathanfielder	nathanfielder-1.jpg	https://www.youtube.com/@nathanfielder	2006-10-16	nathanfielder-1	\N
\.



--
-- Data for Name: video; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video (id, title, video_url, thumbnail_name, upload_date, pg_rating, search_vector, insert_timestamp, slug, description) FROM stdin;
1	Good Pals	https://www.youtube.com/watch?v=6aTqXkZHnQE	good-pals-1.jpg	2008-09-08	PG	'good':1 'pal':2	\N	good-pals-1	
\.


--
-- Data for Name: video_creator_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_creator_rel (creator_id, video_id) FROM stdin;
1	1
\.



-- Data for Name: video_person_rel; Type: TABLE DATA; Schema: public; Owner: colet
--

COPY public.video_person_rel (person_id, video_id, character_id, id) FROM stdin;
4	1	\N	4
5	1	\N	5
\.


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.person_id_seq', 5, true);


--
-- Name: character_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.character_id_seq', 1, false);


--
-- Name: creator_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.creator_id_seq', 13, true);


--
-- Name: video_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_id_seq', 31, true);


--
-- Name: video_person_rel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: colet
--

SELECT pg_catalog.setval('public.video_person_rel_id_seq', 5, true);
