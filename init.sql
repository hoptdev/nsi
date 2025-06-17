CREATE TABLE dashboards (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    parentId int NULL REFERENCES dashboards
);

CREATE TYPE widgetType AS ENUM ('square');

CREATE TABLE widgets (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    dashboardId int NOT NULL REFERENCES dashboards ON DELETE CASCADE,
    type widgetType NOT NULL,
    config jsonb NOT NULL
);

CREATE TYPE grantType AS ENUM ('read', 'update', 'admin');

CREATE TABLE accessRights (
    id SERIAL PRIMARY KEY,
    userId int NULL,
    userGroupId int NULL,
    accessToken varchar(512) NULL,
    type grantType NOT NULL
);

CREATE TABLE dashboardOnAccessRights (
    accessRightId INT REFERENCES accessRights(id) ON DELETE CASCADE,
    dashboardId INT REFERENCES dashboards(id) ON DELETE CASCADE,
    PRIMARY KEY (accessRightId, dashboardId)
);

CREATE TABLE widgetOnAccessRights (
    accessRightId INT REFERENCES accessRights(id) ON DELETE CASCADE,
    widgetId INT REFERENCES widgets(id) ON DELETE CASCADE,
    PRIMARY KEY (accessRightId, widgetId)
);