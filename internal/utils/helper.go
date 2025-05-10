package util

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	errdefs "tracking-service/internal/errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bwmarrin/snowflake"
)

func RandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func GenerateAPIKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ParseTimeDefaultFormat 使用預設格式 "2006-01-02 15:04:05" 轉換
func ParseTimeDefaultFormat(timeStr string) (time.Time, error) {
	const defaultLayout = "2006-01-02 15:04:05"
	return ParseTime(timeStr, defaultLayout)
}

func ParseTime(timeStr string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, errdefs.ErrorInvalidRequest
	}
	return t, nil
}

func ConvertTimeToTimeStamp(t *time.Time) string {
	if t == nil {
		return ""
	}
	return strconv.FormatInt(t.Unix(), 10)
}

func ConvertGormDeletedAtToTimeStamp(t gorm.DeletedAt) string {
	if !t.Valid {
		return ""
	}
	return strconv.FormatInt(t.Time.Unix(), 10)
}

func ParseBearerToken(authHeader string) (string, error) {
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" || splitToken[1] == "" {
		return "", errdefs.ErrorForbidden
	}
	return splitToken[1], nil
}

// 將 Snowflake ID 轉換為 time.Time
func ConvertSnowflakeToTime(snowflakeId string) (*time.Time, error) {
	// 將字串轉換成整數型 Snowflake ID
	idInt, err := strconv.ParseInt(snowflakeId, 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %v", err)
	}

	// 使用 snowflake.NewIDFromString 將 ID 轉換回 Snowflake ID 類型
	id := snowflake.ID(idInt)

	// 反查 ID 中的時間戳
	timestamp := id.Time()

	timestampTime := time.Unix(0, timestamp*int64(time.Millisecond))

	return &timestampTime, nil
}

func Md5(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ChunkArray 將切片分割為指定大小的子切片
func ChunkArray(data []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		return [][]string{}
	}

	result := make([][]string, 0, (len(data)+chunkSize-1)/chunkSize)
	for start := 0; start < len(data); start += chunkSize {
		end := min(start+chunkSize, len(data))
		result = append(result, data[start:end])
	}
	return result
}

// 輔助函數：取得兩數最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func WithRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * time.Duration(i+1)):
		}
	}
	return err
}
