CREATE TABLE IF NOT EXISTS task_assignees (
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, user_id)
);

CREATE INDEX idx_task_assignees_user_id ON task_assignees(user_id);