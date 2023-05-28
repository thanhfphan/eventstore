BEGIN;

CREATE TABLE IF NOT EXISTS es_aggregate(
  "id"              VARCHAR(50) PRIMARY KEY,
  "version"         INTEGER NOT NULL,
  "aggregate_type"  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS es_event(
  "id"              BIGSERIAL PRIMARY KEY,
  "transaction_id"  XID8 NOT NULL,
  "aggregate_id"    VARCHAR(50) NOT NULL REFERENCES es_aggregate("id"),
  "version"         INTEGER NOT NULL,
  "event_type"		  TEXT NOT NULL,
  "data"	          JSONB NOT NULL,
  "metadata"	      JSONB,

  UNIQUE ("aggregate_id", "version")
);

CREATE TABLE IF NOT EXISTS es_aggregate_snapshot(
  "aggregate_id"	VARCHAR(50) NOT NULL REFERENCES es_aggregate("id"),
  "version"       INTEGER NOT NULL,
  "data"			    JSONB NOT NULL,

  PRIMARY KEY (aggregate_id, VERSION)
);

CREATE TABLE IF NOT EXISTS oms_order(
  "id"                BIGSERIAL PRIMARY KEY,
  "version"			      INTEGER NOT NULL,
  "status"          	TEXT NOT NULL,
  "customer_id"		    BIGINT NOT NULL,
  "price"           	DECIMAL(19, 2) NOT NULL,
  "placed_date"     	TIMESTAMPTZ NOT NULL,
  "accepted_date"   	TIMESTAMPTZ,
  "cancelled_date"  	TIMESTAMPTZ,
  "completed_date"  	TIMESTAMPTZ
);

END;