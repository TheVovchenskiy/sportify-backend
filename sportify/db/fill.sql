INSERT INTO public.event (
    id, creator_id, sport_type, address, date_start, start_time, end_time,
    price, game_level, description, raw_message, capacity, creation_type,
    url_preview, url_photos
) VALUES
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7430', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'basketball',
       'Москва Ленинградский проспект, 39с79', '2026-10-11 00:00:00',
       '2026-10-11 20:00:00', '2026-10-11 23:00:00', 0, ARRAY['low']::game_level_enum[],
       'Приходите все! Чисто игровая тренировка',
       '', NULL, 'site',
       'https://127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}'),
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7431', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'football',
       'Москва Госпитальный пер. д 4/6', '2026-10-11 00:00:00',
       '2026-10-11 18:00:00', '2026-10-11 18:00:00',700, ARRAY['mid_plus']::game_level_enum[],
       'Половину тренировки отрабатываем схему 4-4-2, вторая половина игровая',
       '', 22, 'site',
       'https:/127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}'),
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7432', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'volleyball',
       'Москва Ленинградский проспект 60', '2026-10-11 00:00:00',
       '2026-10-11 18:00:00','2026-10-11 18:00:00', 1000, ARRAY['mid']::game_level_enum[],
       'Сегодня чисто игровая тренировка. Вход с улицы напротив школы. На проходной скажите, что на игру',
       '', 0, 'site',
       'https://127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}');

-- password=user1234
INSERT INTO "user" (id, username, password) VALUES
    ('0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'user', '36189b41f377c669cb9144a7cab5ad11261dbe33ffe066cdd0896ef82360423bf1efe94cade268ca');
