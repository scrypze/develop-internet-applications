BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO users (uuid, login, role, pass)
VALUES (uuid_generate_v4(), 'creator', 2, 'hash')
ON CONFLICT (login) DO NOTHING;

INSERT INTO stars (
  id, title, description, image_path, spectral_type, temperature,
  radius, mass, luminosity, metallicity, age, distance
) VALUES
  (1,'Gliese 667','Красный карлик спектрального класса M1.5V','http://localhost:9000/exocalc/gliese-667.png','M1.5V','~3700 K','~0.42 R☉','~0.33 M☉','~0.014 L☉','-0.55','~2-10 млрд лет','~6.8 pc (~22 световых лет)'),
  (2,'TRAPPIST-1','Ультрахолодный красный карлик спектрального класса M8V','http://localhost:9000/exocalc/trappist-1.png','M8V','~2550 K','~0.12 R☉','~0.09 M☉','~0.0005 L☉','0.04','~7.6 млрд лет','~12.1 pc (~39 световых лет)'),
  (3,'Ross 128','Красный карлик спектрального класса M4V','http://localhost:9000/exocalc/ross-128.png','M4V','~3192 K','~0.21 R☉','~0.16 M☉','~0.0036 L☉','0.00','~9.45 млрд лет','~3.37 pc (~11 световых лет)'),
  (4,'Proxima Centauri','Красный карлик спектрального класса M5.5V','http://localhost:9000/exocalc/proxima.png','M5.5V','~3042 K','~0.15 R☉','~0.12 M☉','~0.0017 L☉','0.21','~4.85 млрд лет','~1.30 pc (~4.24 световых лет)')
ON CONFLICT (id) DO NOTHING;

WITH creator AS (
  SELECT uuid FROM users WHERE login = 'creator' LIMIT 1
)
INSERT INTO selected_stars (id, status, scientist, creator_id, date)
SELECT 1, 'draft', 'Михаил Ломоносов', creator.uuid, DATE '2003-10-12'
FROM creator
ON CONFLICT (id) DO NOTHING;

INSERT INTO calculate_exoplanets (
  selected_stars_id, star_id, comment, probable_number_of_planets, habitable_zone
) VALUES
  (1,1,'',2.5,'0.12-0.24 a.e.'),
  (1,2,'',2.5,'0.12-0.24 a.e.'),
  (1,3,'',2.5,'0.12-0.24 a.e.')
ON CONFLICT DO NOTHING;

-- Обновление sequence для таблицы stars после явной вставки ID
SELECT setval('stars_id_seq', (SELECT COALESCE(MAX(id), 1) FROM stars));

COMMIT;