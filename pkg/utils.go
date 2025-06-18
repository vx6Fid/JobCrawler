package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateHash generates a SHA256 hash of input strings concatenated
func GenerateHash(parts ...string) string {
	joined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(hash[:])
}

// Returns PostedOn + 30 days for TTL
func CalculateExpireAt(postedOn time.Time) time.Time {
	return postedOn.AddDate(0, 0, 30)
}

func FormatDate(dateStr string) (time.Time, error) {
	return time.Parse("Jan 2, 2006", dateStr)
}

// parseRelativeTime parses strings like "X days ago", "X hours ago", etc.
func ParseRelativeTime(relativeStr string) (time.Time, error) {
	relativeStr = strings.ToLower(strings.TrimSpace(relativeStr))

	// Regex to capture the number and the unit
	re := regexp.MustCompile(`(\d+)\s+(day|days|hour|hours|minute|minutes|week|weeks|month|months|year|years)\s+ago`)
	matches := re.FindStringSubmatch(relativeStr)

	if len(matches) < 3 {
		return time.Time{}, fmt.Errorf("could not parse relative time: %s", relativeStr)
	}

	valueStr := matches[1]
	unit := matches[2]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number in relative time: %s", valueStr)
	}

	var duration time.Duration
	switch {
	case strings.HasPrefix(unit, "minute"):
		duration = time.Duration(value) * time.Minute
	case strings.HasPrefix(unit, "hour"):
		duration = time.Duration(value) * time.Hour
	case strings.HasPrefix(unit, "day"):
		duration = time.Duration(value) * 24 * time.Hour // A day is 24 hours
	case strings.HasPrefix(unit, "week"):
		duration = time.Duration(value) * 7 * 24 * time.Hour // A week is 7 days
	case strings.HasPrefix(unit, "month"):
		// Months and years are tricky without a specific reference date
		// For simplicity, we'll approximate a month as 30 days.
		// For precision, you'd need to consider the start date.
		duration = time.Duration(value) * 30 * 24 * time.Hour
	case strings.HasPrefix(unit, "year"):
		// Approximate a year as 365 days
		duration = time.Duration(value) * 365 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}

	// Calculate the time 'duration' ago from now
	return time.Now().Add(-duration), nil
}
