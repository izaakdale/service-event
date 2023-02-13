CREATE TABLE "events" (
  "event_id" bigserial PRIMARY KEY,
  "event_name" varchar NOT NULL,
  "tickets_remaining" int NOT NULL,
  "event_timestamp" timestamp NOT NULL
);

CREATE INDEX ON "events" ("event_name");