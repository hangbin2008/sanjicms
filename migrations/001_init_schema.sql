-- 基层三基考试系统数据库表结构
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL, -- 用户名，手机号或身份证号
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL, -- 姓名，只能是中文
    role VARCHAR(20) NOT NULL DEFAULT 'employee', -- 角色：admin（站长）、manager（管理员）、employee（员工/考生）
    phone VARCHAR(20) UNIQUE, -- 手机号
    id_card VARCHAR(18) UNIQUE, -- 身份证号
    department VARCHAR(100), -- 部门
    job_title VARCHAR(50), -- 职称
    status TINYINT DEFAULT 1, -- 状态：1-启用，0-禁用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建题库表
CREATE TABLE IF NOT EXISTS question_banks (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL, -- 题库名称
    description TEXT, -- 题库描述
    subject VARCHAR(50) NOT NULL, -- 科目：基础知识、基本技能、专业知识、专业实践技能
    created_by INT NOT NULL, -- 创建人ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建题目表
CREATE TABLE IF NOT EXISTS questions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    bank_id INT NOT NULL, -- 所属题库ID
    type VARCHAR(20) NOT NULL, -- 题目类型：single（单选题）、multiple（多选题）、judgment（判断题）、essay（简答题）
    content TEXT NOT NULL, -- 题目内容
    options TEXT, -- 选项，JSON格式存储
    answer TEXT NOT NULL, -- 答案
    score FLOAT NOT NULL, -- 分值
    difficulty VARCHAR(20) DEFAULT 'medium', -- 难度：easy（简单）、medium（中等）、hard（困难）
    analysis TEXT, -- 解析
    created_by INT NOT NULL, -- 创建人ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (bank_id) REFERENCES question_banks(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建试卷表
CREATE TABLE IF NOT EXISTS exams (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL, -- 试卷标题
    description TEXT, -- 试卷描述
    subject VARCHAR(50) NOT NULL, -- 科目
    total_score FLOAT NOT NULL, -- 总分
    duration INT NOT NULL, -- 考试时长（分钟）
    start_time DATETIME, -- 开始时间
    end_time DATETIME, -- 结束时间
    status VARCHAR(20) DEFAULT 'draft', -- 状态：draft（草稿）、published（已发布）、completed（已结束）
    created_by INT NOT NULL, -- 创建人ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建试卷题目关联表
CREATE TABLE IF NOT EXISTS exam_questions (
    exam_id INT NOT NULL, -- 试卷ID
    question_id INT NOT NULL, -- 题目ID
    sequence INT NOT NULL, -- 题目顺序
    PRIMARY KEY (exam_id, question_id),
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

-- 创建考生答卷表
CREATE TABLE IF NOT EXISTS exam_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    exam_id INT NOT NULL, -- 试卷ID
    user_id INT NOT NULL, -- 考生ID
    start_time DATETIME DEFAULT CURRENT_TIMESTAMP, -- 开始答题时间
    end_time DATETIME, -- 结束答题时间
    duration INT, -- 实际答题时长（秒）
    total_score FLOAT DEFAULT 0, -- 总分
    status VARCHAR(20) DEFAULT 'ongoing', -- 状态：ongoing（进行中）、submitted（已提交）、graded（已评分）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建考生答题表
CREATE TABLE IF NOT EXISTS exam_answers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    record_id INT NOT NULL, -- 答卷ID
    question_id INT NOT NULL, -- 题目ID
    user_answer TEXT NOT NULL, -- 考生答案
    score FLOAT DEFAULT 0, -- 得分
    is_correct TINYINT DEFAULT 0, -- 是否正确：1-正确，0-错误
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (record_id) REFERENCES exam_records(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    key_name VARCHAR(50) UNIQUE NOT NULL, -- 配置键
    value TEXT NOT NULL, -- 配置值
    description VARCHAR(200), -- 配置描述
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 插入初始数据
-- 插入系统配置
INSERT IGNORE INTO system_configs (key_name, value, description) VALUES
('jwt_secret', 'your-secret-key', 'JWT签名密钥'),
('password_min_length', '8', '密码最小长度'),
('password_require_letter', '1', '密码必须包含字母'),
('password_require_digit', '1', '密码必须包含数字'),
('password_require_special', '1', '密码必须包含特殊字符'),
('exam_auto_grade', '1', '是否自动评分');

-- 插入初始用户（站长）
INSERT IGNORE INTO users (username, password_hash, name, role, status) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'admin', 1);
-- 密码：Admin@123
