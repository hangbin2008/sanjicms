-- 基层三基考试系统数据库表结构
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    gender VARCHAR(10) DEFAULT '男',
    email VARCHAR(100) UNIQUE,
    role VARCHAR(20) NOT NULL DEFAULT 'employee',
    phone VARCHAR(20) UNIQUE,
    id_card VARCHAR(18) UNIQUE,
    department VARCHAR(100),
    job_title VARCHAR(50),
    avatar VARCHAR(255),
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 注释：
-- username: 用户名，手机号或身份证号
-- name: 姓名，只能是中文
-- role: 角色：admin（站长）、manager（管理员）、employee（员工/考生）
-- phone: 手机号
-- id_card: 身份证号
-- department: 部门
-- job_title: 职称
-- avatar: 头像URL
-- status: 状态：1-启用，0-禁用

-- 创建题库表
CREATE TABLE IF NOT EXISTS question_banks (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    subject VARCHAR(50) NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建题目表
CREATE TABLE IF NOT EXISTS questions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    bank_id INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    options TEXT,
    answer TEXT NOT NULL,
    score FLOAT NOT NULL,
    difficulty VARCHAR(20) DEFAULT 'medium',
    analysis TEXT,
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (bank_id) REFERENCES question_banks(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建试卷表
CREATE TABLE IF NOT EXISTS exams (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    subject VARCHAR(50) NOT NULL,
    total_score FLOAT NOT NULL,
    duration INT NOT NULL,
    start_time DATETIME,
    end_time DATETIME,
    status VARCHAR(20) DEFAULT 'draft',
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建试卷题目关联表
CREATE TABLE IF NOT EXISTS exam_questions (
    exam_id INT NOT NULL,
    question_id INT NOT NULL,
    sequence INT NOT NULL,
    PRIMARY KEY (exam_id, question_id),
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建考生答卷表
CREATE TABLE IF NOT EXISTS exam_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    exam_id INT NOT NULL,
    user_id INT NOT NULL,
    start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    end_time DATETIME,
    duration INT,
    total_score FLOAT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'ongoing',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建考生答题表
CREATE TABLE IF NOT EXISTS exam_answers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    record_id INT NOT NULL,
    question_id INT NOT NULL,
    user_answer TEXT NOT NULL,
    score FLOAT DEFAULT 0,
    is_correct TINYINT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (record_id) REFERENCES exam_records(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    key_name VARCHAR(50) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description VARCHAR(200),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入初始数据

-- 1. 先插入初始用户（站长）- 确保用户数据在其他依赖表之前插入
-- 密码：Admin@123
INSERT IGNORE INTO users (username, password_hash, name, role, status) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'admin', 1);

-- 2. 插入系统配置
INSERT IGNORE INTO system_configs (key_name, value, description) VALUES
('jwt_secret', 'FP0bnaIYdagkRFNlPEdjQHb7RfNaGMuDF3DLjKtI4zE=', 'JWT签名密钥'),
('password_min_length', '8', '密码最小长度'),
('password_require_letter', '1', '密码必须包含字母'),
('password_require_digit', '1', '密码必须包含数字'),
('password_require_special', '1', '密码必须包含特殊字符'),
('exam_auto_grade', '1', '是否自动评分');

-- 3. 插入额外的系统配置，用于支持动态配置
INSERT IGNORE INTO system_configs (key_name, value, description) VALUES
('server_port', '8080', '服务器端口'),
('db_host', 'db', '数据库主机'),
('db_port', '3306', '数据库端口');
