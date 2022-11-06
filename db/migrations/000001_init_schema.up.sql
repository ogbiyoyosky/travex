CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table "users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "firstName" varchar(255),
    "lastName" varchar(255),
    "email" varchar(255)  unique NOT NULL,
    "password" varchar(255) NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now()),
    "updatedAt" timestamptz NOT NULL DEFAULT (now()),
    "role" varchar(255)  NOT NULL DEFAULT 'customer'
);

create table "locations" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "address" varchar(255) NOT NULL,
    "postalCode" varchar(255) NOT NULL,
    "locationTypeId" int  NOT NULL,
    "userId" int NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now()),
    "updatedAt" timestamptz NOT NULL DEFAULT (now())
);

create table "location_types" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now()),
    "updatedAt" timestamptz NOT NULL DEFAULT (now())
);

create table "reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "locationId" int NOT NULL,
    "authorId" int NOT NULL,
    "rating" float NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now()),
    "updatedAt" timestamptz NOT NULL DEFAULT (now())
);

create table "comments" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "locationId" int NOT NULL,
    "authorId" int NOT NULL,
    "text" text NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now()),
    "updatedAt" timestamptz NOT NULL DEFAULT (now())
);

