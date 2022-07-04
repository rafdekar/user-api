CREATE TABLE "users" (
                         "id" uuid DEFAULT MD5(RANDOM()::TEXT || CLOCK_TIMESTAMP()::TEXT)::UUID PRIMARY KEY,
                         "first_name" varchar NOT NULL,
                         "last_name" varchar NOT NULL,
                         "nickname" varchar UNIQUE NOT NULL,
                         "password" varchar NOT NULL,
                         "email" varchar NOT NULL,
                         "country" varchar NOT NULL,
                         "modified_at" timestamp NOT NULL DEFAULT (now()),
                         "created_at" timestamp NOT NULL DEFAULT (now())
);
