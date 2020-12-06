CREATE TABLE IF NOT EXISTS users (
	key SERIAL PRIMARY KEY NOT NULL,
	id BIGINT NOT NULL,
	name VARCHAR(25) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX users_id_idx ON users(id);
CREATE INDEX users_name_idx ON users(name);
CREATE INDEX users_desc_idx ON users(created_at DESC NULLS LAST);

CREATE TABLE IF NOT EXISTS permissions (
	id SERIAL PRIMARY KEY NOT NULL,
	user_type VARCHAR(20) NOT NULL,
	user_id BIGINT NOT NULL,
	flags BIGINT DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_type, user_id)
);
CREATE INDEX permissions_user_type_user_id_idx ON permissions(user_type, user_id);
CREATE INDEX permissions_user_type_flags_idx ON permissions(user_type, flags);

CREATE TABLE IF NOT EXISTS messages (
	key SERIAL PRIMARY KEY NOT NULL,
	room_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL DEFAULT -1,
	username VARCHAR(25) NOT NULL,
	message VARCHAR(500) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(room_id, user_id, message, created_at),
	UNIQUE(room_id, username, message, created_at)
);
CREATE INDEX messages_room_user_idx ON messages(room_id, user_id);
CREATE INDEX messages_room_username_idx ON messages(room_id, username);
CREATE INDEX messages_desc_idx ON messages(created_at DESC NULLS LAST);

CREATE TABLE IF NOT EXISTS count (
	key SERIAL PRIMARY KEY NOT NULL,
	room_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	total BIGINT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(room_id, user_id)
);
CREATE INDEX count_room_user_idx ON count(room_id, user_id);
CREATE INDEX count_total_desc_idx ON count(total DESC NULLS LAST);
CREATE INDEX count_created_desc_idx ON count(created_at DESC NULLS LAST);
CREATE INDEX count_updated_desc_idx ON count(updated_at DESC NULLS LAST);

-- Update username history avoiding duplicates
CREATE OR REPLACE FUNCTION on_name_change() RETURNS TRIGGER AS $$
	DECLARE LAST users%ROWTYPE;
	BEGIN
		SELECT * INTO LAST FROM users WHERE id=NEW.id ORDER BY created_at DESC LIMIT 1;
		IF NEW.name = LAST.name OR NEW.created_at < LAST.created_at THEN
			RETURN NULL;
		END IF;
		RETURN NEW;
	END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER avoid_dupes BEFORE INSERT ON users FOR EACH ROW
	EXECUTE PROCEDURE on_name_change();

-- Increment message count for users
CREATE OR REPLACE FUNCTION on_message_insert() RETURNS TRIGGER AS $$
	BEGIN
		IF(TG_OP = 'INSERT') THEN
            IF NEW.user_id >= 0 THEN
                INSERT INTO users(id, name, created_at) VALUES(NEW.user_id, NEW.username, NEW.created_at);
                INSERT INTO count(room_id, user_id, total) VALUES(NEW.room_id, NEW.user_id, 1)
                    ON CONFLICT(room_id, user_id) DO UPDATE SET total=count.total+1, updated_at=now() WHERE count.room_id=NEW.room_id AND count.user_id=NEW.user_id;
            END IF;
            RETURN NEW;
		ELSEIF(TG_OP = 'DELETE') THEN
            IF OLD.user_id >= 0 THEN
                INSERT INTO count(room_id, user_id, total) VALUES(OLD.room_id, OLD.user_id, 0)
                    ON CONFLICT(room_id, user_id) DO UPDATE SET total=count.total-1, updated_at=now() WHERE count.room_id=OLD.room_id AND count.user_id=OLD.user_id;
			END IF;
            RETURN OLD;
		END IF;
		RETURN NEW;
	END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increment AFTER INSERT OR DELETE ON messages FOR EACH ROW
	EXECUTE PROCEDURE on_message_insert();
