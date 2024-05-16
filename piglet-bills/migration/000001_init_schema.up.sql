CREATE TABLE "bills" (
    "id" uuid PRIMARY KEY,
    "bill_name" VARCHAR(255) UNIQUE NOT NULL,
    "current_sum" decimal,
    "bill_type" bool,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "goals" (
    "bill_id" uuid,
    "goal_sum" decimal NOT NULL,
    "date" DATE NOT NULL,
    "monthly_payment" decimal
);

CREATE TABLE "accounts" (
    "bill_id" uuid,
    "bill_status" bool
);

INSERT INTO bills (id, bill_name, current_sum, bill_type)
    VALUES ('00000000-0000-0000-0000-000000000001', 'default bill', 0, true);
INSERT INTO accounts (bill_id, bill_status) VALUES ('00000000-0000-0000-0000-000000000001', true);

CREATE INDEX ON "bills" ("bill_name");

CREATE INDEX ON "bills" ("bill_type");

CREATE INDEX ON "goals" ("bill_id");

CREATE INDEX ON "accounts" ("bill_id");

CREATE INDEX ON "accounts" ("bill_status");

ALTER TABLE "goals" ADD FOREIGN KEY ("bill_id") REFERENCES "bills" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("bill_id") REFERENCES "bills" ("id");
