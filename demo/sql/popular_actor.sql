SELECT
    actor.actor_id,
    actor.first_name,
    actor.last_name,
    count(*) AS films_count
FROM
    actor,
    film_actor
WHERE
    actor.actor_id = film_actor.actor_id
GROUP BY
    actor.actor_id,
    actor.first_name,
    actor.last_name
ORDER BY
    actor.actor_id
LIMIT 10;

