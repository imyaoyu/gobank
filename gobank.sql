CREATE TABLE IF NOT EXISTS t_prod_info (
	id INTEGER PRIMARY KEY,
	prod_id TEXT NOT NULL UNIQUE,
	prod_type TEXT NOT NULL,
	prod_name TEXT NOT NULL,
	prod_rate INTEGER NOT NULL,
	prod_rate_str TEXT NOT NULL,
	prod_term INTEGER NOT NULL,
	prod_note TEXT NOT NULL
);
INSERT into t_prod_info 
(prod_id,prod_type,prod_name,prod_rate,prod_rate_str,prod_term,prod_note) 
values 
('101','1','1年定存',150,'1.50%',12,'当日起息，保本保息'),
('102','1','2年定存',180,'1.80%',24,'当日起息，保本保息'),
('103','1','3年定存',200,'2.00%',36,'当日起息，保本保息');

CREATE TABLE IF NOT EXISTS t_acct_info (
	id INTEGER PRIMARY KEY,
	user_id TEXT NOT NULL,
	prod_id TEXT NOT NULL,
	prod_type TEXT NOT NULL,
	prod_name TEXT NOT NULL,
	rate INTEGER NOT NULL,
	balance INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	open_date TEXT NOT NULL,
	end_date TEXT NOT NULL,
	close_date TEXT,
	status TEXT NOT NULL,
	interest INTEGER
);
CREATE INDEX IF NOT EXISTS idx_t_acct_info ON t_acct_info(user_id);