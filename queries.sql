--   clear all tables except FACULTIES AND PROGRAMS
DELETE FROM assignation_models;
DELETE FROM alert_models;
DELETE FROM semester_availability_models;

-- get total assigned resources in a semester
SELECT 
    SUM(am.classrooms) AS assigned_classrooms, 
    SUM(am.labs) AS assigned_labs, 
    SUM(am.mobile_labs) AS assigned_mobile_labs,
    sem.classrooms AS classrooms,
    sem.labs AS labs,
    sem.mobile_labs AS mobile_labs
FROM assignation_models am
INNER JOIN semester_availability_models sem ON sem.id = am.semester_id 
WHERE sem.semester = '2025-10'
GROUP BY sem.classrooms, sem.labs, sem.mobile_labs;


-- get assignation history of specified semester ordered by creation time
SELECT 
    am.id,
    am.classrooms,
    am.labs,
    am.mobile_labs,
    sem.semester,
    pm.name,
    am.created_at
FROM assignation_models am 
INNER JOIN semester_availability_models sem ON sem.id = am.semester_id 
INNER JOIN program_models pm ON pm.id = am.program_id 
WHERE sem.semester = '2025-10'
ORDER BY am.created_at ASC;


-- get alert history log
SELECT 
    am.id,
    am.message,
    am.requested_classrooms,
    am.requested_labs,
    am.available_classrooms,
    am.available_labs,
    am.available_mobile_labs,
    sem.semester,
    pm.name,
    am.created_at
FROM alert_models am 
INNER JOIN semester_availability_models sem ON sem.id = am.semester_id 
INNER JOIN program_models pm ON pm.id = am.program_id 
WHERE sem.semester = '2025-10'
ORDER BY am.created_at ASC;


--check number of logs either asignation or alert on a given semester

WITH counts AS (
    SELECT COUNT(*) AS count, semester_id
    FROM assignation_models am
    GROUP BY semester_id 
    UNION 
    SELECT COUNT(*) AS count, semester_id
    FROM alert_models al
    GROUP BY semester_id 
)
SELECT SUM(c.count) AS total_logs
FROM counts c 
INNER JOIN semester_availability_models sem ON sem.id = c.semester_id 
WHERE sem.semester = '2025-10';


--check all programs with an assignemnt or alert on a given semester
WITH all_logs AS (
    SELECT 
        program_id,
        semester_id 
    FROM assignation_models am
    UNION 
    SELECT 
        program_id,
        semester_id 
    FROM alert_models al
)
SELECT pm.name AS program, fac.name AS faculty
FROM all_logs
INNER JOIN program_models pm ON pm.id = all_logs.program_id 
INNER JOIN faculty_models fac ON fac.id = pm.faculty_id 
INNER JOIN semester_availability_models sem ON sem.id = all_logs.semester_id 
WHERE sem.semester = '2025-10'
ORDER BY fac.name, pm.name;