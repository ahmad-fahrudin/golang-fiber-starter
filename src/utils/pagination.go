package utils

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PaginationParams holds pagination and filtering parameters
type PaginationParams struct {
	Page      int     `json:"page" validate:"omitempty,number,max=1000"`
	Limit     int     `json:"limit" validate:"omitempty,number,max=100"`
	Search    string  `json:"search,omitempty" validate:"omitempty,max=100"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

// PaginationResult holds the result of pagination
type PaginationResult[T any] struct {
	Results      []T   `json:"results"`
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalPages   int64 `json:"total_pages"`
	TotalResults int64 `json:"total_results"`
}

// ApplyPagination applies pagination and date filtering to a GORM query
func ApplyPagination[T any](db *gorm.DB, params *PaginationParams, dateField string) (*PaginationResult[T], error) {
	var results []T
	var totalResults int64

	// Set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Clone the query for counting
	countQuery := db.Model(new(T))
	dataQuery := db.Model(new(T))

	// Apply date filtering if provided
	if params.StartDate != nil && *params.StartDate != "" {
		startDate, err := ParseDate(*params.StartDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		}
		countQuery = countQuery.Where(dateField+" >= ?", startDate)
		dataQuery = dataQuery.Where(dateField+" >= ?", startDate)
	}

	if params.EndDate != nil && *params.EndDate != "" {
		endDate, err := ParseDate(*params.EndDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		}
		// Add 24 hours to include the entire end date
		endDate = endDate.Add(24 * time.Hour)
		countQuery = countQuery.Where(dateField+" < ?", endDate)
		dataQuery = dataQuery.Where(dateField+" < ?", endDate)
	}

	// Count total results
	if err := countQuery.Count(&totalResults).Error; err != nil {
		return nil, err
	}

	// Apply pagination and get results
	if err := dataQuery.Limit(params.Limit).Offset(offset).Order(dateField + " DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalResults) / float64(params.Limit)))

	return &PaginationResult[T]{
		Results:      results,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}

// ParseDate parses date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ExtractPaginationParams extracts pagination parameters from Fiber context
func ExtractPaginationParams(c *fiber.Ctx) *PaginationParams {
	params := &PaginationParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	// Extract date filters
	if startDate := c.Query("start_date"); startDate != "" {
		params.StartDate = &startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		params.EndDate = &endDate
	}

	return params
}

// CreatePaginationResponse creates a standardized pagination response
func CreatePaginationResponse[T any](code int, message string, result *PaginationResult[T]) map[string]interface{} {
	return map[string]interface{}{
		"code":          code,
		"status":        "success",
		"message":       message,
		"results":       result.Results,
		"page":          result.Page,
		"limit":         result.Limit,
		"total_pages":   result.TotalPages,
		"total_results": result.TotalResults,
	}
}

// PaginateQuery is a simple helper function to paginate any query
func PaginateQuery[T any](db *gorm.DB, params *PaginationParams, orderBy string) (*PaginationResult[T], error) {
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	var results []T
	var totalResults int64

	// Set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Count total results
	if err := db.Model(new(T)).Count(&totalResults).Error; err != nil {
		return nil, err
	}

	// Get paginated results
	if err := db.Order(orderBy).Limit(params.Limit).Offset(offset).Find(&results).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int64(totalResults / int64(params.Limit))
	if totalResults%int64(params.Limit) != 0 {
		totalPages++
	}

	return &PaginationResult[T]{
		Results:      results,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}

// ApplyPaginationWithSearch applies pagination with custom search function
func ApplyPaginationWithSearch[T any](
	db *gorm.DB,
	params *PaginationParams,
	dateField string,
	searchCallback func(*gorm.DB, string) *gorm.DB,
) (*PaginationResult[T], error) {
	var results []T
	var totalResults int64

	// Set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Create base query
	query := db.Model(new(T))

	// Apply search callback if provided and search term exists
	if searchCallback != nil && params.Search != "" {
		query = searchCallback(query, params.Search)
	}

	// Apply date filtering if provided
	if params.StartDate != nil && *params.StartDate != "" {
		startDate, err := ParseDate(*params.StartDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		}
		query = query.Where(dateField+" >= ?", startDate)
	}

	if params.EndDate != nil && *params.EndDate != "" {
		endDate, err := ParseDate(*params.EndDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		}
		// Add 24 hours to include the entire end date
		endDate = endDate.Add(24 * time.Hour)
		query = query.Where(dateField+" < ?", endDate)
	}

	// Count total results
	if err := query.Count(&totalResults).Error; err != nil {
		return nil, err
	}

	// Apply pagination and get results
	if err := query.Limit(params.Limit).Offset(offset).Order(dateField + " DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalResults) / float64(params.Limit)))

	return &PaginationResult[T]{
		Results:      results,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}
