CREATE TABLE "transactions" (
    "id" uuid PRIMARY KEY NOT NULL,
    "trans_date" DATE NOT NULL,
    "type" int NOT NULL,
    "sum" decimal NOT NULL,
    "comment" text,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "income" (
    "trans_id" uuid NOT NULL,
    "id_category" uuid NOT NULL,
    "id_bill_to" uuid NOT NULL,
    "sender" text,
    "repeat" bool
);

CREATE TABLE "expense" (
    "trans_id" uuid NOT NULL,
    "id_category" uuid NOT NULL,
    "id_bill_from" uuid NOT NULL,
    "recipient" text,
    "repeat" bool
);

CREATE TABLE "debt" (
    "trans_id" uuid NOT NULL,
    "type" bool,
    "id_bill_from" uuid,
    "id_bill_to" uuid,
    "creditor_debtor" text NOT NULL
);

CREATE TABLE "transfer" (
    "trans_id" uuid NOT NULL,
    "id_bill_from" uuid NOT NULL,
    "id_bill_to" uuid NOT NULL
);

CREATE TABLE "categories" (
    "id" uuid PRIMARY KEY,
    "type" bool,
    "name" text NOT NULL,
    "mandatory" bool
);

CREATE TABLE "bills" (
   "id" uuid PRIMARY KEY,
   "status" bool
);

INSERT INTO categories (id, type, name, mandatory) VALUES ('00000000-0000-0000-0000-000000000000', true, 'no category', false);
INSERT INTO categories (id, type, name, mandatory) VALUES ('00000000-0000-0000-0000-000000000001', false, 'no category', false);

CREATE INDEX ON "transactions" ("trans_date");

CREATE INDEX ON "transactions" ("type");

CREATE INDEX ON "income" ("id_category");

CREATE INDEX ON "income" ("id_bill_to");

CREATE INDEX ON "income" ("repeat");

CREATE INDEX ON "expense" ("id_category");

CREATE INDEX ON "expense" ("id_bill_from");

CREATE INDEX ON "expense" ("repeat");

CREATE INDEX ON "debt" ("type");

CREATE INDEX ON "debt" ("creditor_debtor");

CREATE INDEX ON "categories" ("type");

CREATE INDEX ON "categories" ("mandatory");

ALTER TABLE "income" ADD FOREIGN KEY ("trans_id") REFERENCES "transactions" ("id");

ALTER TABLE "income" ADD FOREIGN KEY ("id_category") REFERENCES "categories" ("id");

ALTER TABLE "income" ADD FOREIGN KEY ("id_bill_to") REFERENCES "bills" ("id");

ALTER TABLE "expense" ADD FOREIGN KEY ("trans_id") REFERENCES "transactions" ("id");

ALTER TABLE "expense" ADD FOREIGN KEY ("id_category") REFERENCES "categories" ("id");

ALTER TABLE "expense" ADD FOREIGN KEY ("id_bill_from") REFERENCES "bills" ("id");

ALTER TABLE "debt" ADD FOREIGN KEY ("trans_id") REFERENCES "transactions" ("id");

ALTER TABLE "debt" ADD FOREIGN KEY ("id_bill_from") REFERENCES "bills" ("id");

ALTER TABLE "debt" ADD FOREIGN KEY ("id_bill_to") REFERENCES "bills" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("trans_id") REFERENCES "transactions" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("id_bill_from") REFERENCES "bills" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("id_bill_to") REFERENCES "bills" ("id");