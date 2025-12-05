package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/yourusername/jiceng-sanji-exam/internal/db"
	"github.com/yourusername/jiceng-sanji-exam/internal/models"
)

// QuestionService 题库服务
type QuestionService struct{}

// NewQuestionService 创建题库服务
func NewQuestionService() *QuestionService {
	return &QuestionService{}
}

// CreateQuestionBank 创建题库
func (s *QuestionService) CreateQuestionBank(req *models.QuestionBankCreateRequest, createdBy int) (*models.QuestionBank, error) {
	// 插入题库记录
	result, err := db.DB.Exec(`
		INSERT INTO question_banks (name, description, subject, created_by)
		VALUES (?, ?, ?, ?)
	`, req.Name, req.Description, req.Subject, createdBy)
	if err != nil {
		return nil, err
	}

	// 获取插入的题库ID
	bankID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询插入的题库信息
	var bank models.QuestionBank
	err = db.DB.QueryRow(`
		SELECT id, name, description, subject, created_by, created_at, updated_at
		FROM question_banks WHERE id = ?
	`, bankID).Scan(
		&bank.ID, &bank.Name, &bank.Description, &bank.Subject, &bank.CreatedBy,
		&bank.CreatedAt, &bank.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

// GetQuestionBankByID 根据ID获取题库
func (s *QuestionService) GetQuestionBankByID(bankID int) (*models.QuestionBank, error) {
	var bank models.QuestionBank
	err := db.DB.QueryRow(`
		SELECT id, name, description, subject, created_by, created_at, updated_at
		FROM question_banks WHERE id = ?
	`, bankID).Scan(
		&bank.ID, &bank.Name, &bank.Description, &bank.Subject, &bank.CreatedBy,
		&bank.CreatedAt, &bank.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

// ListQuestionBanks 获取题库列表
func (s *QuestionService) ListQuestionBanks(subject string, page, pageSize int) ([]models.QuestionBank, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	var banks []models.QuestionBank
	var total int

	// 获取总记录数
	query := "SELECT COUNT(*) FROM question_banks"
	args := []interface{}{}

	if subject != "" {
		query += " WHERE subject = ?"
		args = append(args, subject)
	}

	err := db.DB.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取题库列表
	listQuery := `
		SELECT id, name, description, subject, created_by, created_at, updated_at
		FROM question_banks
	`
	argsList := []interface{}{}

	if subject != "" {
		listQuery += " WHERE subject = ?"
		argsList = append(argsList, subject)
	}

	listQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	argsList = append(argsList, pageSize, offset)

	rows, err := db.DB.Query(listQuery, argsList...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var bank models.QuestionBank
		err := rows.Scan(
			&bank.ID, &bank.Name, &bank.Description, &bank.Subject, &bank.CreatedBy,
			&bank.CreatedAt, &bank.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		banks = append(banks, bank)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return banks, total, nil
}

// CreateQuestion 创建题目
func (s *QuestionService) CreateQuestion(req *models.QuestionCreateRequest, createdBy int) (*models.Question, error) {
	// 检查题库是否存在
	var bankCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM question_banks WHERE id = ?", req.BankID).Scan(&bankCount)
	if err != nil {
		return nil, err
	}
	if bankCount == 0 {
		return nil, errors.New("题库不存在")
	}

	// 插入题目记录
	result, err := db.DB.Exec(`
		INSERT INTO questions (bank_id, type, content, options, answer, score, difficulty, analysis, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.BankID, req.Type, req.Content, req.Options, req.Answer, req.Score, req.Difficulty, req.Analysis, createdBy)
	if err != nil {
		return nil, err
	}

	// 获取插入的题目ID
	questionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询插入的题目信息
	var question models.Question
	err = db.DB.QueryRow(`
		SELECT id, bank_id, type, content, options, answer, score, difficulty, analysis, created_by, created_at, updated_at
		FROM questions WHERE id = ?
	`, questionID).Scan(
		&question.ID, &question.BankID, &question.Type, &question.Content, &question.Options,
		&question.Answer, &question.Score, &question.Difficulty, &question.Analysis, &question.CreatedBy,
		&question.CreatedAt, &question.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

// GetQuestionByID 根据ID获取题目
func (s *QuestionService) GetQuestionByID(questionID int) (*models.Question, error) {
	var question models.Question
	err := db.DB.QueryRow(`
		SELECT id, bank_id, type, content, options, answer, score, difficulty, analysis, created_by, created_at, updated_at
		FROM questions WHERE id = ?
	`, questionID).Scan(
		&question.ID, &question.BankID, &question.Type, &question.Content, &question.Options,
		&question.Answer, &question.Score, &question.Difficulty, &question.Analysis, &question.CreatedBy,
		&question.CreatedAt, &question.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

// ListQuestionsByBank 获取题库下的题目列表
func (s *QuestionService) ListQuestionsByBank(bankID int, page, pageSize int) ([]models.Question, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	var questions []models.Question
	var total int

	// 获取总记录数
	err := db.DB.QueryRow("SELECT COUNT(*) FROM questions WHERE bank_id = ?", bankID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取题目列表
	rows, err := db.DB.Query(`
		SELECT id, bank_id, type, content, options, answer, score, difficulty, analysis, created_by, created_at, updated_at
		FROM questions WHERE bank_id = ?
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`, bankID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID, &question.BankID, &question.Type, &question.Content, &question.Options,
			&question.Answer, &question.Score, &question.Difficulty, &question.Analysis, &question.CreatedBy,
			&question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		questions = append(questions, question)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// GetRandomQuestions 获取随机题目
func (s *QuestionService) GetRandomQuestions(subject string, difficulty string, count int) ([]models.Question, error) {
	var questions []models.Question
	var err error
	var rows *sql.Rows

	// 构建查询条件
	query := "SELECT id, bank_id, type, content, options, answer, score, difficulty, analysis, created_by, created_at, updated_at FROM questions WHERE 1=1"
	args := []interface{}{}

	// 添加科目条件
	if subject != "" {
		query += " AND bank_id IN (SELECT id FROM question_banks WHERE subject = ?)"
		args = append(args, subject)
	}

	// 添加难度条件
	if difficulty != "" {
		query += " AND difficulty = ?"
		args = append(args, difficulty)
	}

	// 随机排序并限制数量
	query += fmt.Sprintf(" ORDER BY RAND() LIMIT %d", count)

	// 执行查询
	rows, err = db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID, &question.BankID, &question.Type, &question.Content, &question.Options,
			&question.Answer, &question.Score, &question.Difficulty, &question.Analysis, &question.CreatedBy,
			&question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
