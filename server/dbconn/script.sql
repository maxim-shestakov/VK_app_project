CREATE TABLE filmsactors.films (
	id int GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar(50) NOT NULL,
	"description" varchar(1000) NOT NULL,
	"date" char(8) NOT NULL,
	rating int4 DEFAULT 0 NOT NULL,
	CONSTRAINT films_pk PRIMARY KEY (id)
);

CREATE TABLE filmsactors.actors (
	id int GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar(50) NOT NULL,
	surname varchar(50) NOT NULL,
	fathername varchar(50) NOT NULL,
	birthdate char(8) NOT NULL,
	sex char(1) NO NULL,
    CONSTRAINT actor_pk PRIMARY KEY (id)
);

CREATE TABLE filmsactors.actorsfilms (
	actor_id int4 not NULL,
	film_id int4 not NULL,
	CONSTRAINT actorsfilms_pk PRIMARY KEY (actor_id, film_id),
	CONSTRAINT actors_films_actor_fk FOREIGN KEY (actor_id) REFERENCES filmsactors.actors(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT actors_films_film_fk FOREIGN KEY (film_id) REFERENCES filmsactors.films(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE filmsactors.users (
	"login" varchar(50) NOT NULL,
	"password" varchar(50) NOT NULL,
	"role" int DEFAULT 0 NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY ("login")
);

INSERT INTO filmsactors.actors ("name",surname,fathername,birthdate) VALUES
	 ('Роберт','Дауни-младший',NULL,'04041965', 'm'),
	 ('Киану','Ривз',NULL,'02091964', 'm'),
	 ('Бурунов','Сергей','Александрович','06031977', 'm'),;

INSERT INTO filmsactors.films ("name", "description", "date", rating) VALUES
	 ('Мстители: Финал','Оставшиеся в живых члены команды Мстителей и их союзники должны разработать новый план, который поможет противостоять разрушительным действиям могущественного титана Таноса. После наиболее масштабной и трагической битвы в истории они не могут допустить ошибку.','29042019',7.9),
	 ('Джон Уик 3','Киллер-изгой бежит от байкеров-самураев и других неприятностей. Мощное продолжение остросюжетной франшизы','16052019',7.0),
	 ('Затмение','Самозванец участвует в шоу экстрасенсов. Очаровательный Александр Петров и вихрь мистических происшествий','25112016',5.8);

INSERT INTO filmsactors.actorsfilms (actor_id,film_id) VALUES
	 (2,1),
	 (3,2),
	 (4,3);

INSERT INTO filmsactors.users ("login","password","role") VALUES
	 ('alice_smith','$2a$12$EEWfSU1DD4NYqF9V0sOX7.jxky5YGC.4yTi2CSjsAkGhW9ohDKNdm',0),
	 ('john_doe','$2a$12$RFOEwd0Z8Fw.ZYeTwzpFpeSjPka2nlhoZSjebYqR4V.ZVENwFtCo.',1);

