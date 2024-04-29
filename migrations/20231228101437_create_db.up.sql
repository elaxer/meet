CREATE SEQUENCE assessments_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."assessments" (
    "id" integer DEFAULT nextval('assessments_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "from_user_id" integer NOT NULL,
    "to_user_id" integer NOT NULL,
    "message" character varying,
    "decision" character varying NOT NULL,
    "is_mutual" boolean DEFAULT false NOT NULL,
    CONSTRAINT "assessments_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE "public"."cities" (
    "id" integer NOT NULL,
    "country_id" smallint NOT NULL,
    "name" character varying NOT NULL,
    "latitude" numeric NOT NULL,
    "longitude" numeric NOT NULL,
    CONSTRAINT "cities_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE "public"."countries" (
    "id" smallint NOT NULL,
    "region_id" smallint,
    "name" character varying NOT NULL,
    "name_native" character varying NOT NULL,
    "emoji" character(2) NOT NULL,
    CONSTRAINT "countries_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE couples_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 3 CACHE 1;

CREATE TABLE "public"."couples" (
    "id" integer DEFAULT nextval('couples_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "from_user_id" integer NOT NULL,
    "to_user_id" integer NOT NULL,
    CONSTRAINT "couples_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE messages_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."messages" (
    "id" integer DEFAULT nextval('messages_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "from_user_id" integer NOT NULL,
    "to_user_id" integer NOT NULL,
    "text" text NOT NULL,
    "is_readed" boolean DEFAULT false NOT NULL,
    CONSTRAINT "messages_pkey" PRIMARY KEY ("id")
) WITH (oids = false);



CREATE SEQUENCE questionnaire_photos_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."questionnaire_photos" (
    "id" integer DEFAULT nextval('questionnaire_photos_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "questionnaire_id" integer NOT NULL,
    "path" character varying NOT NULL,
    CONSTRAINT "questionnaire_photos_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE questionnaires_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 69 CACHE 1;

CREATE TABLE "public"."questionnaires" (
    "id" integer DEFAULT nextval('questionnaires_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "user_id" integer NOT NULL,
    "name" character varying NOT NULL,
    "birth_date" date,
    "gender" boolean,
    "orientation" character varying,
    "meeting_purpose" character varying,
    "age_range_min" smallint,
    "age_range_max" smallint,
    "city_id" integer,
    "about" text,
    "is_active" boolean DEFAULT false NOT NULL,
    "state" character varying,
    CONSTRAINT "questionnaires_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE "public"."regions" (
    "id" smallint NOT NULL,
    "name" character varying NOT NULL,
    CONSTRAINT "regions_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 79 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "created_at" timestamp DEFAULT now() NOT NULL,
    "updated_at" timestamp,
    "login" character varying NOT NULL,
    "password_hash" character varying,
    "is_blocked" boolean DEFAULT false NOT NULL,
    "tg_id" integer,
    CONSTRAINT "users_login" UNIQUE ("login"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

ALTER TABLE ONLY "public"."assessments" ADD CONSTRAINT "assessments_from_user_id_fkey" FOREIGN KEY (from_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."assessments" ADD CONSTRAINT "assessments_to_user_id_fkey" FOREIGN KEY (to_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."cities" ADD CONSTRAINT "cities_country_id_fkey" FOREIGN KEY (country_id) REFERENCES countries(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."countries" ADD CONSTRAINT "countries_region_id_fkey" FOREIGN KEY (region_id) REFERENCES regions(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."couples" ADD CONSTRAINT "couples_from_user_id_fkey" FOREIGN KEY (from_user_id) REFERENCES users(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."couples" ADD CONSTRAINT "couples_to_user_id_fkey" FOREIGN KEY (to_user_id) REFERENCES users(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."messages" ADD CONSTRAINT "messages_from_user_id_fkey" FOREIGN KEY (from_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."messages" ADD CONSTRAINT "messages_to_user_id_fkey" FOREIGN KEY (to_user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."questionnaire_photos" ADD CONSTRAINT "questionnaire_photos_questionnaire_id_fkey" FOREIGN KEY (questionnaire_id) REFERENCES questionnaires(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."questionnaires" ADD CONSTRAINT "questionnaires_city_id_fkey" FOREIGN KEY (city_id) REFERENCES cities(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."questionnaires" ADD CONSTRAINT "questionnaires_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;