CREATE TABLE films (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(50) NOT NULL,
	description varchar(1000) NOT NULL,
	"date" bpchar(8) NOT NULL,
	rating float8 DEFAULT 0 NOT NULL,
	CONSTRAINT films_pk PRIMARY KEY (id)
);

CREATE TABLE actors (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(50) NOT NULL,
	surname varchar(50) NOT NULL,
	fathername varchar(50) DEFAULT ' ' NULL,
	birthdate bpchar(8) NOT NULL,
	sex bpchar(1) NOT NULL,
	CONSTRAINT actor_pk PRIMARY KEY (id)
);

CREATE TABLE actorsfilms (
	actor_id int4 NOT NULL,
	film_id int4 NOT NULL,
	CONSTRAINT actorsfilms_pk PRIMARY KEY (actor_id, film_id),
	CONSTRAINT actors_films_actor_fk FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT actors_films_film_fk FOREIGN KEY (film_id) REFERENCES films(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE users (
	login varchar(120) NOT NULL,
	"password" varchar(500) NOT NULL,
	"role" int4 DEFAULT 0 NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (login)
);

INSERT INTO actors ("name",surname,fathername,birthdate,sex) VALUES
        ('Бурунов','Сергей','Александрович','19770306', 'm');

INSERT INTO actors ("name",surname,fathername,birthdate,sex) VALUES
        ('Роберт','Дауни-младший','','19650404', 'm'),
        ('Киану','Ривз','','19640902', 'm');

INSERT INTO films ("name", "description", "date", rating) VALUES
        ('Мстители: Финал','Оставшиеся в живых члены команды Мстителей и их союзники должны разработать новый план, который поможет противостоять разрушительным действиям могущественного титана Таноса. После наиболее масштабной и трагической битвы в истории они не могут допустить ошибку.','20190429',7.9),
        ('Джон Уик 3','Киллер-изгой бежит от байкеров-самураев и других неприятностей. Мощное продолжение остросюжетной франшизы','20190516',7.0),
        ('Затмение','Самозванец участвует в шоу экстрасенсов. Очаровательный Александр Петров и вихрь мистических происшествий','20161125',5.8);

INSERT INTO actorsfilms (actor_id,film_id) VALUES
        (1,1),
        (2,2),
        (3,3);

INSERT INTO users ("login","password","role") VALUES
        ('alice_smith','$2a$12$EEWfSU1DD4NYqF9V0sOX7.jxky5YGC.4yTi2CSjsAkGhW9ohDKNdm',0),
        ('john_doe','$2a$12$RFOEwd0Z8Fw.ZYeTwzpFpeSjPka2nlhoZSjebYqR4V.ZVENwFtCo.',1);

