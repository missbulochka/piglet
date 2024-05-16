-- INSERT INTO bills (id, bill_name, current_sum, bill_type)
--     VALUES ('00000000-0000-0000-0000-000000000001', 'default_bill', 0, true);
-- INSERT INTO accounts (bill_status) VALUES (true);

INSERT INTO bills (id, bill_name, current_sum, bill_type)
    VALUES ('00000000-0000-0000-0000-000000000001', 'default account', 0.0, true);
INSERT INTO accounts (bill_id, bill_status) VALUES ('00000000-0000-0000-0000-000000000001', true);
