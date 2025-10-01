-- 1) Создаём пользователя-автора (если нет)
INSERT INTO users (login, password, is_moderator)
VALUES ('creator', 'x', false)
ON CONFLICT (login) DO NOTHING;

-- 2) Вставляем набор выбранных звёзд с заполненным creator_id
WITH creator AS (
  SELECT id FROM users WHERE login = 'creator' LIMIT 1
)
INSERT INTO selected_stars (
  id,
  status,
  scientist,
  creator_id,
  date,          -- опционально, если хочешь заполнить
  created_at,    -- опционально
  formed_at,     -- опционально
  completed_at   -- опционально
)
SELECT
  1,
  'draft',
  'Михаил Ломоносов',
  creator.id,
  DATE '2003-10-12',
  NULL,
  NULL,
  NULL
FROM creator
ON CONFLICT (id) DO NOTHING;

-- 3) Линки набора к звёздам (id 1..3 уже должны быть вставлены в insert_stars.sql)
INSERT INTO calculate_exoplanets (selection_stars_id, star_id, comment) VALUES
(1, 1, ''),
(1, 2, ''),
(1, 3, '')
ON CONFLICT DO NOTHING;
