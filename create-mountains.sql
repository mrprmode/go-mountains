DROP TABLE IF EXISTS mountain;
CREATE TABLE mountain (
  id          INT AUTO_INCREMENT NOT NULL,
  name        VARCHAR(128) NOT NULL,
  height      INT NOT NULL,
  local_name  VARCHAR(255),
  PRIMARY KEY (`id`)
);

INSERT INTO mountain
  (name, height, local_name)
VALUES
  ('Mt. Everest', 29032, 'Sagarmatha || Qomolangma'),
  ('Annapurna', 26545, ''),
  ('Gasherbrum III', 26089, ''),
  ('Gyachung Kang', 26089, ''),
  ('Fishtail', 22943, 'Machapuchare'),
  ('Mt. McKinley', 20310, 'Denali'),
  ('Mt. Rainier', 14410, 'Tahoma');
