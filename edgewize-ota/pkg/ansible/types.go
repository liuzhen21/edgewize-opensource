package ansible

import (
	"fmt"
	"strconv"
)

// AnsiblePlaybookJSONResults
type AnsiblePlaybookJSONResults struct {
	Playbook          string                                      `json:"-"`
	CustomStats       interface{}                                 `json:"custom_stats"`
	GlobalCustomStats interface{}                                 `json:"global_custom_stats"`
	Plays             []AnsiblePlaybookJSONResultsPlay            `json:"plays"`
	Stats             map[string]*AnsiblePlaybookJSONResultsStats `json:"stats"`
}

// AnsiblePlaybookJSONResultsPlay
type AnsiblePlaybookJSONResultsPlay struct {
	Play  *AnsiblePlaybookJSONResultsPlaysPlay `json:"play"`
	Tasks []AnsiblePlaybookJSONResultsPlayTask `json:"tasks"`
}

// AnsiblePlaybookJSONResultsPlaysPlay
type AnsiblePlaybookJSONResultsPlaysPlay struct {
	Name     string                                  `json:"name"`
	Id       string                                  `json:"id"`
	Duration *AnsiblePlaybookJSONResultsPlayDuration `json:"duration"`
}

type AnsiblePlaybookJSONResultsPlayTask struct {
	Task  *AnsiblePlaybookJSONResultsPlayTaskItem `json:"task"`
	Hosts map[string]interface{}                  `json:"hosts"`
}

type AnsiblePlaybookJSONResultsPlayTaskItem struct {
	Name     string                                          `json:"name"`
	Id       string                                          `json:"id"`
	Duration *AnsiblePlaybookJSONResultsPlayTaskItemDuration `json:"duration"`
}

type AnsiblePlaybookJSONResultsPlayTaskItemDuration struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// AnsiblePlaybookJSONResultsPlayDuration
type AnsiblePlaybookJSONResultsPlayDuration struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// AnsiblePlaybookJSONResultsStats
type AnsiblePlaybookJSONResultsStats struct {
	Changed     int `json:"changed"`
	Failures    int `json:"failures"`
	Ignored     int `json:"ignored"`
	Ok          int `json:"ok"`
	Rescued     int `json:"rescued"`
	Skipped     int `json:"skipped"`
	Unreachable int `json:"unreachable"`
}

func (s *AnsiblePlaybookJSONResultsStats) String() string {
	str := fmt.Sprintf(" Changed: %s", strconv.Itoa(s.Changed))
	str = fmt.Sprintf("%s Failures: %s", str, strconv.Itoa(s.Failures))
	str = fmt.Sprintf("%s Ignored: %s", str, strconv.Itoa(s.Ignored))
	str = fmt.Sprintf("%s Ok: %s", str, strconv.Itoa(s.Ok))
	str = fmt.Sprintf("%s Rescued: %s", str, strconv.Itoa(s.Rescued))
	str = fmt.Sprintf("%s Skipped: %s", str, strconv.Itoa(s.Skipped))
	str = fmt.Sprintf("%s Unreachable: %s", str, strconv.Itoa(s.Unreachable))

	return str
}
