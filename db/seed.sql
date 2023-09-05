CREATE TABLE IF NOT EXISTS client (
			id uuid NOT NULL,
			"name" varchar(100) NOT NULL,
      balance int NOT NULL,
      PRIMARY KEY(id)
);
insert into client (id, name, balance) values ('c7fbc542-4c2f-11ee-be56-0242ac120002', 'User teste', 100);