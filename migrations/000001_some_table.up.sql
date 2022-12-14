CREATE TABLE if not exists "categories"(
    "id" serial PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE default current_timestamp
);

CREATE TABLE if not exists "users"(
    "id" serial PRIMARY KEY,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255),
    "phone_number" VARCHAR(255) UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "gender" VARCHAR(255) CHECK("gender" IN('male','female','')),
    "password" VARCHAR(255) NOT NULL,
    "username" VARCHAR(255) NOT NULL UNIQUE,
    "profile_image_url" VARCHAR(255),
    "type" VARCHAR(255)CHECK("type" IN('superadmin','user')) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE default current_timestamp
);

CREATE TABLE if not exists "posts"(
    "id" serial PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "image_url" VARCHAR(255),
    "user_id" INTEGER NOT NULL REFERENCES users(id)ON DELETE CASCADE,
    "category_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "views_count" INTEGER not NULL default 0
);

CREATE TABLE if not exists "comments"(
    "id" serial PRIMARY KEY,
    "post_id" INTEGER NOT NULL REFERENCES posts(id)ON DELETE CASCADE,
    "user_id" INTEGER NOT NULL REFERENCES users(id)ON DELETE CASCADE,
    "description" TEXT,
    "created_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITHOUT TIME ZONE
);


CREATE TABLE if not exists "likes"(
    "id" serial PRIMARY KEY,
    "post_id" INTEGER NOT NULL REFERENCES posts(id)ON DELETE CASCADE,
    "user_id" INTEGER NOT NULL REFERENCES users(id)ON DELETE CASCADE,
    "status" BOOLEAN NOT NULL,
    UNIQUE(post_id, user_id)
);
