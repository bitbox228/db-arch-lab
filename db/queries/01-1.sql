WITH anime_ranking_by_genre AS (SELECT a.genre,
                                       a.id                                                         as anime_id,
                                       a.title                                                      as anime_title,
                                       COUNT(DISTINCT ual.user_id)                                  AS views_count,
                                       ROW_NUMBER()
                                       OVER (PARTITION BY a.genre ORDER BY COUNT(ual.user_id) DESC) AS genre_rank
                                FROM anime a
                                         JOIN user_anime_status ual
                                              ON a.id = ual.anime_id AND ual.list <> 'WANT_TO_WATCH'
                                GROUP BY a.genre, a.id
                                ORDER BY a.genre, genre_rank),
     anime_reviews_stats AS (SELECT a.id                                        AS id,
                                    COUNT(DISTINCT r.id)                        AS reviews_count,
                                    SUM(CASE WHEN rr.is_like THEN 1 ELSE 0 END) AS reviews_likes_count
                             FROM anime a
                                      LEFT JOIN
                                  reviews r ON a.id = r.anime_id
                                      LEFT JOIN reactions rr
                                                ON r.id = rr.review_id
                             GROUP BY a.id),
     reviews_stats AS (SELECT r.id                                        AS id,
                              r.anime_id                                  AS anime_id,
                              SUM(CASE WHEN rr.is_like THEN 1 ELSE 0 END) AS likes_count
                       FROM reviews r
                                LEFT JOIN reactions rr ON r.id = rr.review_id
                       GROUP BY r.id)

SELECT rg.genre,
       rg.genre_rank,
       rg.anime_id,
       rg.anime_title,
       rg.views_count,
       ars.reviews_count,
       ars.reviews_likes_count,
       rs.id          AS review_id,
       rs.likes_count AS likes_count
FROM anime_ranking_by_genre rg
         JOIN anime_reviews_stats ars
              ON rg.anime_id = ars.id
         LEFT JOIN reviews_stats rs
                   ON rg.anime_id = rs.anime_id
WHERE rg.genre_rank <= 10
ORDER BY rg.genre, rg.genre_rank;