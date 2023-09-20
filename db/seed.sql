CREATE TABLE IF NOT EXISTS client (
			id uuid NOT NULL,
			"name" varchar(100) NOT NULL,
      balance int NOT NULL,
      PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS transfers (
			id uuid NOT NULL,
			sender_id uuid not null,
			receiver_id uuid not null,
      total  int NOT NULL,
      operation_date timestamp  not null, 
      PRIMARY KEY(id),
      FOREIGN KEY (sender_id) REFERENCES client(id),
      FOREIGN KEY (receiver_id) REFERENCES client(id)
);
insert into client (id, name, balance) values ('c7fbc542-4c2f-11ee-be56-0242ac120002', 'User teste', 100);
insert into client (id, name, balance) values ('88ee33d6-4c33-11ee-be56-0242ac120002', 'User teste 2', 100);