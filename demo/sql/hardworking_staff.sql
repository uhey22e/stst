SELECT
    staff.staff_id,
    staff.first_name,
    staff.last_name,
    date_trunc('month', payment.payment_date) AS payment_month,
    sum(payment.amount) AS amount
FROM
    staff,
    payment
WHERE
    staff.staff_id = payment.staff_id
GROUP BY
    staff.staff_id,
    staff.first_name,
    staff.last_name,
    date_trunc('month', payment.payment_date)
LIMIT 10;

