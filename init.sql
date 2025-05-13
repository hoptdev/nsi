заготовка под миграции

CREATE TABLE dashboards (id SERIAL PRIMARY KEY, name varchar(255), parentId int NULL REFERENCES dashboards);

CREATE TYPE widgetType AS ENUM ('square');

CREATE TABLE widgets (id SERIAL PRIMARY KEY, name varchar(255), dashboardId int NOT NULL REFERENCES dashboards, type widgetType NOT NULL, config json NOT NULL);