BEGIN;

CREATE TABLE IF NOT EXISTS ES_AGGREGATE (
  ID              VARCHAR(50)     PRIMARY KEY,
  VERSION         INTEGER  NOT NULL,
  TYPE            TEXT     NOT NULL
);

CREATE TABLE IF NOT EXISTS ES_EVENT (
  ID              BIGSERIAL  PRIMARY KEY,
  TRANSACTION_ID  XID8       NOT NULL,
  AGGREGATE_ID    VARCHAR(50)       NOT NULL REFERENCES ES_AGGREGATE (ID),
  VERSION         INTEGER    NOT NULL,
  TYPE		      TEXT       NOT NULL,
  DATA	          JSON       NOT NULL,
  UNIQUE (AGGREGATE_ID, VERSION)
);

CREATE TABLE IF NOT EXISTS ES_AGGREGATE_SNAPSHOT (
  AGGREGATE_ID	VARCHAR(50) NOT NULL REFERENCES ES_AGGREGATE (ID),
  VERSION       INTEGER  NOT NULL,
  DATA			JSON     NOT NULL,
  PRIMARY KEY (AGGREGATE_ID, VERSION)
);

CREATE TABLE IF NOT EXISTS OMS_ORDER (
  ID				BIGSERIAL PRIMARY KEY,
  VERSION			INTEGER         NOT NULL,
  STATUS          	TEXT            NOT NULL,
  CUSTOMER_ID		BIGINTEGER      NOT NULL,
  PRICE           	DECIMAL(19, 2)  NOT NULL,
  PLACED_DATE     	TIMESTAMPZ       NOT NULL,
  ACCEPTED_DATE   	TIMESTAMPZ,
  CANCELLED_DATE  	TIMESTAMPZ,
  COMPLETED_DATE  	TIMESTAMPZ
);

END;