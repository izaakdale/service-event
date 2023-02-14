CREATE TABLE "events" (
  "event_id" bigserial PRIMARY KEY,
  "event_name" varchar NOT NULL,
  "tickets_remaining" int NOT NULL,
  "event_timestamp" timestamp NOT NULL,
  CONSTRAINT "n_of_tickets_valid" CHECK (tickets_remaining >= 0)
);

CREATE INDEX ON "events" ("event_name");