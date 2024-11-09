INSERT INTO public.event (
    id, creator_id, sport_type, address, date_start, start_time,
    price, game_level, description, raw_message, capacity, creation_type,
    url_preview, url_photos
) VALUES
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7430', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'football',
       'Москва Госпитальный пер. 4/6', '2024-10-24 00:00:00',
       '2026-10-24 20:00:00', 0, ARRAY['mid_minus']::game_level_enum[],
       'Приходите все! Чисто игровая тренировка',
       '', NULL, 'site',
       'https://127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}'),
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7431', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'football',
       'Москва Госпитальный пер. 4/6', '2024-10-25 00:00:00',
       '2026-10-25 18:00:00', 700, ARRAY['mid_plus']::game_level_enum[],
       'Половину тренировки отрабатываем схему 4-4-2, вторая половина игровая',
       '', 22, 'site',
       'https://127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}'),
      ('9bc13767-82d4-46e2-8dd8-51a5df1c7432', '0bc13767-82d4-46e2-8dd8-51a5df1c7426', 'football',
       'Москва Госпитальный пер. 4/6', '2024-10-25 00:00:00',
       '2026-10-25 20:00:00', 1000, ARRAY['mid']::game_level_enum[],
       'Сегодня чисто игровая тренировка. Вход с улицы напротив школы. На проходной скажите, что на игру',
       '', 0, 'site',
       'https://127.0.0.1/api/v1/img/default_football.jpeg',
       '{"https://127.0.0.1/api/v1/img/default_football.jpeg"}');
