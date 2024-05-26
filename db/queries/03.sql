SELECT DATE_TRUNC('month', ua.time) AS month_year,
       COUNT(*)                     AS achievements_count
FROM users u
         JOIN user_achievements ua ON u.id = ua.user_id
         JOIN achievements a ON a.id = ua.achievement_id
WHERE ua.time >= NOW() - INTERVAL '5 years'
GROUP BY DATE_TRUNC('month', ua.time)
ORDER BY month_year;