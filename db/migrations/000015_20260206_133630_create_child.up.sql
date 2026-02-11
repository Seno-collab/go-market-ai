CREATE TABLE IF NOT EXISTS childrens (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name   TEXT NOT NULL,
  dob         DATE,          -- hoặc year_of_birth INT nếu muốn tối giản
  grade       INT CHECK (grade BETWEEN 1 AND 12),

  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_children_grade ON childrens(grade);


CREATE TABLE IF NOT EXISTS user_childrens (
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  child_id    UUID NOT NULL REFERENCES childrens(id) ON DELETE CASCADE,

  relation    TEXT NOT NULL DEFAULT 'parent'
              CHECK (relation IN ('parent','mother','father','guardian')),

  is_primary  BOOLEAN NOT NULL DEFAULT false,

  -- Quyền cơ bản (MVP: dùng cột boolean dễ query hơn JSONB)
  can_edit_plan      BOOLEAN NOT NULL DEFAULT true,
  can_manage_rewards BOOLEAN NOT NULL DEFAULT true,
  can_view_reports   BOOLEAN NOT NULL DEFAULT true,

  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  PRIMARY KEY (user_id, child_id)
);

CREATE INDEX IF NOT EXISTS idx_user_children_user  ON user_childrens(user_id);
CREATE INDEX IF NOT EXISTS idx_user_children_child ON user_childrens(child_id);
