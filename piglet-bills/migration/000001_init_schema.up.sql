CREATE TABLE "bills" (
    "id" uuid PRIMARY KEY,
    "billName" VARCHAR(255) NOT NULL UNIQUE,
    "currentSum" decimal,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "goals" (
    "bill_id" uuid,
    "date" DATE NOT NULL,
    "monthlyPayment" decimal
);

CREATE TABLE "accounts" (
      "bill_id" uuid,
      "billStatus" bool
);

CREATE INDEX ON "bills" ("billName");

CREATE INDEX ON "goals" ("bill_id");

CREATE INDEX ON "accounts" ("bill_id");

CREATE INDEX ON "accounts" ("billStatus");

ALTER TABLE "goals" ADD FOREIGN KEY ("bill_id") REFERENCES "bills" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("bill_id") REFERENCES "bills" ("id");
