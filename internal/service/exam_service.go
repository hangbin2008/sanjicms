package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/jiceng-sanji-exam/internal/db"
	"github.com/yourusername/jiceng-sanji-exam/internal/models"
)

// ExamService 试卷服务
type ExamService struct {
	questionService *QuestionService
}

// NewExamService 创建试卷服务
func NewExamService(questionService *QuestionService) *ExamService {
	return &ExamService{
		questionService: questionService,
	}
}

// GenerateExam 生成试卷
func (s *ExamService) GenerateExam(req *models.ExamGenerateRequest, createdBy int) (*models.Exam, error) {
	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}

	// 检查结束时间是否大于开始时间
	if endTime.Before(startTime) {
		return nil, errors.New("结束时间必须大于开始时间")
	}

	// 获取随机题目
	questions, err := s.questionService.GetRandomQuestions(req.Subject, req.Difficulty, req.QuestionCount)
	if err != nil {
		return nil, err
	}

	// 计算总分
	var totalScore float64
	for _, q := range questions {
		totalScore += q.Score
	}

	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 插入试卷记录
	result, err := tx.Exec(`
		INSERT INTO exams (title, description, subject, total_score, duration, start_time, end_time, status, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Title, req.Description, req.Subject, totalScore, req.Duration, startTime, endTime, "published", createdBy)
	if err != nil {
		return nil, err
	}

	// 获取插入的试卷ID
	examID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 插入试卷题目关联
	for i, q := range questions {
		_, err := tx.Exec(`
			INSERT INTO exam_questions (exam_id, question_id, sequence)
			VALUES (?, ?, ?)
		`, examID, q.ID, i+1)
		if err != nil {
			return nil, err
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// 查询生成的试卷信息
	var exam models.Exam
	err = db.DB.QueryRow(`
		SELECT id, title, description, subject, total_score, duration, start_time, end_time, status, created_by, created_at, updated_at
		FROM exams WHERE id = ?
	`, examID).Scan(
		&exam.ID, &exam.Title, &exam.Description, &exam.Subject, &exam.TotalScore, &exam.Duration,
		&exam.StartTime, &exam.EndTime, &exam.Status, &exam.CreatedBy, &exam.CreatedAt, &exam.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 添加题目到试卷
	exam.Questions = questions

	return &exam, nil
}

// GetExamByID 根据ID获取试卷
func (s *ExamService) GetExamByID(examID int) (*models.Exam, error) {
	var exam models.Exam
	err := db.DB.QueryRow(`
		SELECT id, title, description, subject, total_score, duration, start_time, end_time, status, created_by, created_at, updated_at
		FROM exams WHERE id = ?
	`, examID).Scan(
		&exam.ID, &exam.Title, &exam.Description, &exam.Subject, &exam.TotalScore, &exam.Duration,
		&exam.StartTime, &exam.EndTime, &exam.Status, &exam.CreatedBy, &exam.CreatedAt, &exam.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 获取试卷题目
	rows, err := db.DB.Query(`
		SELECT q.id, q.bank_id, q.type, q.content, q.options, q.answer, q.score, q.difficulty, q.analysis, q.created_by, q.created_at, q.updated_at
		FROM questions q
		JOIN exam_questions eq ON q.id = eq.question_id
		WHERE eq.exam_id = ?
		ORDER BY eq.sequence
	`, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
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

	exam.Questions = questions
	return &exam, nil
}

// StartExam 开始考试
func (s *ExamService) StartExam(examID, userID int) (*models.ExamRecord, error) {
	// 检查试卷是否存在
	var exam models.Exam
	err := db.DB.QueryRow("SELECT id, status, start_time, end_time FROM exams WHERE id = ?", examID).Scan(
		&exam.ID, &exam.Status, &exam.StartTime, &exam.EndTime,
	)
	if err != nil {
		return nil, errors.New("试卷不存在")
	}

	// 检查试卷状态
	if exam.Status != "published" {
		return nil, errors.New("试卷未发布")
	}

	// 检查考试时间
	now := time.Now()
	if now.Before(exam.StartTime) {
		return nil, errors.New("考试尚未开始")
	}
	if now.After(exam.EndTime) {
		return nil, errors.New("考试已结束")
	}

	// 检查是否已经参加过该考试
	var count int
	err = db.DB.QueryRow(
		"SELECT COUNT(*) FROM exam_records WHERE exam_id = ? AND user_id = ?",
		examID, userID,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("您已经参加过该考试")
	}

	// 插入考试记录
	result, err := db.DB.Exec(`
		INSERT INTO exam_records (exam_id, user_id, start_time, status)
		VALUES (?, ?, NOW(), 'ongoing')
	`, examID, userID)
	if err != nil {
		return nil, err
	}

	// 获取插入的记录ID
	recordID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询插入的记录信息
	var record models.ExamRecord
	err = db.DB.QueryRow(`
		SELECT id, exam_id, user_id, start_time, duration, total_score, status, created_at, updated_at
		FROM exam_records WHERE id = ?
	`, recordID).Scan(
		&record.ID, &record.ExamID, &record.UserID, &record.StartTime, &record.Duration,
		&record.TotalScore, &record.Status, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// SubmitExam 提交试卷
func (s *ExamService) SubmitExam(req *models.ExamSubmitRequest) (*models.ExamRecord, error) {
	// 查询考试记录
	var record models.ExamRecord
	err := db.DB.QueryRow(`
		SELECT id, exam_id, user_id, start_time, status
		FROM exam_records WHERE id = ?
	`, req.RecordID).Scan(
		&record.ID, &record.ExamID, &record.UserID, &record.StartTime, &record.Status,
	)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}

	// 检查考试状态
	if record.Status != "ongoing" {
		return nil, errors.New("考试已提交或已结束")
	}

	// 计算考试时长
	now := time.Now()
	duration := int(now.Sub(record.StartTime).Seconds())

	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 更新考试记录状态
	_, err = tx.Exec(`
		UPDATE exam_records
		SET end_time = NOW(), duration = ?, status = 'submitted'
		WHERE id = ?
	`, duration, req.RecordID)
	if err != nil {
		return nil, err
	}

	// 插入答题记录
	var totalScore float64
	for _, answer := range req.Answers {
		// 查询题目信息
		var question models.Question
		err = tx.QueryRow(
			"SELECT id, answer, score FROM questions WHERE id = ?",
			answer.QuestionID,
		).Scan(&question.ID, &question.Answer, &question.Score)
		if err != nil {
			return nil, err
		}

		// 判断答案是否正确
		isCorrect := 0
		score := 0.0
		if answer.UserAnswer == question.Answer {
			isCorrect = 1
			score = question.Score
			totalScore += score
		}

		// 插入答题记录
		_, err = tx.Exec(`
			INSERT INTO exam_answers (record_id, question_id, user_answer, score, is_correct)
			VALUES (?, ?, ?, ?, ?)
		`, req.RecordID, answer.QuestionID, answer.UserAnswer, score, isCorrect)
		if err != nil {
			return nil, err
		}
	}

	// 更新考试总分
	_, err = tx.Exec(`
		UPDATE exam_records
		SET total_score = ?, status = 'graded'
		WHERE id = ?
	`, totalScore, req.RecordID)
	if err != nil {
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// 查询更新后的考试记录
	var updatedRecord models.ExamRecord
	err = db.DB.QueryRow(`
		SELECT id, exam_id, user_id, start_time, end_time, duration, total_score, status, created_at, updated_at
		FROM exam_records WHERE id = ?
	`, req.RecordID).Scan(
		&updatedRecord.ID, &updatedRecord.ExamID, &updatedRecord.UserID, &updatedRecord.StartTime, &updatedRecord.EndTime,
		&updatedRecord.Duration, &updatedRecord.TotalScore, &updatedRecord.Status, &updatedRecord.CreatedAt, &updatedRecord.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedRecord, nil
}

// GetExamRecord 获取考试记录
func (s *ExamService) GetExamRecord(recordID int) (*models.ExamRecord, error) {
	var record models.ExamRecord
	err := db.DB.QueryRow(`
		SELECT id, exam_id, user_id, start_time, end_time, duration, total_score, status, created_at, updated_at
		FROM exam_records WHERE id = ?
	`, recordID).Scan(
		&record.ID, &record.ExamID, &record.UserID, &record.StartTime, &record.EndTime,
		&record.Duration, &record.TotalScore, &record.Status, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 获取考试答案
	rows, err := db.DB.Query(`
		SELECT id, record_id, question_id, user_answer, score, is_correct, created_at
		FROM exam_answers WHERE record_id = ?
	`, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.ExamAnswer
	for rows.Next() {
		var answer models.ExamAnswer
		err := rows.Scan(
			&answer.ID, &answer.RecordID, &answer.QuestionID, &answer.UserAnswer, &answer.Score,
			&answer.IsCorrect, &answer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	record.Answers = answers
	return &record, nil
}

// ListExamRecords 获取用户考试记录
func (s *ExamService) ListExamRecords(userID, page, pageSize int) ([]models.ExamRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	var records []models.ExamRecord
	var total int

	// 获取总记录数
	err := db.DB.QueryRow("SELECT COUNT(*) FROM exam_records WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取考试记录列表
	rows, err := db.DB.Query(`
		SELECT id, exam_id, user_id, start_time, end_time, duration, total_score, status, created_at, updated_at
		FROM exam_records WHERE user_id = ?
		ORDER BY start_time DESC LIMIT ? OFFSET ?
	`, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExamRecord
		err := rows.Scan(
			&record.ID, &record.ExamID, &record.UserID, &record.StartTime, &record.EndTime,
			&record.Duration, &record.TotalScore, &record.Status, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetExamStats 获取考试统计数据
func (s *ExamService) GetExamStats(userID int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总考试次数
	var totalExams int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM exam_records WHERE user_id = ?", userID).Scan(&totalExams)
	if err != nil {
		return nil, err
	}
	stats["total_exams"] = totalExams

	// 总得分
	var totalScore float64
	err = db.DB.QueryRow("SELECT COALESCE(SUM(total_score), 0) FROM exam_records WHERE user_id = ? AND status = 'graded'", userID).Scan(&totalScore)
	if err != nil {
		return nil, err
	}
	stats["total_score"] = totalScore

	// 平均得分
	avgScore := 0.0
	if totalExams > 0 {
		avgScore = totalScore / float64(totalExams)
	}
	stats["avg_score"] = avgScore

	// 最高得分
	var maxScore float64
	err = db.DB.QueryRow("SELECT COALESCE(MAX(total_score), 0) FROM exam_records WHERE user_id = ? AND status = 'graded'", userID).Scan(&maxScore)
	if err != nil {
		return nil, err
	}
	stats["max_score"] = maxScore

	// 最低得分
	var minScore float64
	err = db.DB.QueryRow("SELECT COALESCE(MIN(total_score), 0) FROM exam_records WHERE user_id = ? AND status = 'graded'", userID).Scan(&minScore)
	if err != nil {
		return nil, err
	}
	stats["min_score"] = minScore

	// 最近5次考试记录
	rows, err := db.DB.Query(`
		SELECT id, exam_id, total_score, start_time, status
		FROM exam_records
		WHERE user_id = ?
		ORDER BY start_time DESC
		LIMIT 5
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recentExams []map[string]interface{}
	for rows.Next() {
		var id, examID int
		var totalScore float64
		var startTime time.Time
		var status string
		
		err := rows.Scan(&id, &examID, &totalScore, &startTime, &status)
		if err != nil {
			return nil, err
		}
		
		recentExams = append(recentExams, map[string]interface{}{
			"id":          id,
			"exam_id":     examID,
			"total_score": totalScore,
			"start_time":  startTime,
			"status":      status,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	stats["recent_exams"] = recentExams

	return stats, nil
}
