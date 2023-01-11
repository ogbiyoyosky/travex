CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table "users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "first_name" varchar(255),
    "last_name" varchar(255),
    "email" varchar(255)  unique NOT NULL,
    "password" varchar(255) NOT NULL,
    "role" varchar(255)  NOT NULL DEFAULT 'customer',
    "created_at" timestamp  DEFAULT NOW(),
    "updated_at" timestamp DEFAULT NOW(),
    "deleted_at" timestamp
);

create table "locations" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "address" varchar(255) NOT NULL,
    "image" varchar(255) NOT NULL,
    "description" text NOT NULL,
    "location_type_id" uuid  NOT NULL,
    "user_id" uuid NOT NULL,
    "is_approved" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp  DEFAULT NOW(),
    "is_approved_at" timestamp,
    "deleted_at" timestamp
);

create table "location_types" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "created_at" timestamp  DEFAULT NOW(),
    "updated_at" timestamp  DEFAULT NOW(),
    "deleted_at" timestamp
);

create table "reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "location_id" uuid NOT NULL,
    "author_id" uuid NOT NULL,
    "rating" DECIMAL NOT NULL,
    "text" text DEFAULT NULL,
    "created_at" timestamp  DEFAULT NOW(),
    "updated_at" timestamp DEFAULT NOW(),
    "deleted_at "timestamp,
);

create table "comments" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "location_id" uuid NOT NULL,
    "author_id" uuid NOT NULL,
    "review_id" uuid NOT NULL,
    "text" text NOT NULL,
    "is_approved" boolean NOT NULL DEFAULT FALSE,
    "is_approved_by" uuid DEFAULT  NULL,
    "is_approved_at" timestamp,
    "created_at" timestamp  DEFAULT NOW(),
    "updated_at" timestamp  DEFAULT NOW(),
    "deleted_at" timestamp
);

