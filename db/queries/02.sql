WITH RECURSIVE friends_of_friends AS (SELECT f.user_id2 as user_id, 0 AS level
                                      FROM friends f
                                      WHERE user_id1 = 3
                                      UNION ALL
                                      SELECT f.user_id2 AS user_id, fof.level + 1 AS level
                                      FROM friends f
                                               JOIN friends_of_friends fof ON fof.user_id = f.user_id1
                                      WHERE fof.level < 1)

SELECT DISTINCT fof.user_id,
                fof.level
FROM friends_of_friends fof
ORDER BY fof.level, fof.user_id;