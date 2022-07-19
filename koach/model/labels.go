package model

import "strings"

type LabelsFilter struct {
	Filter map[string]string
}

func (l *LabelsFilter) FromString(labelsFilterString string) {
	l.Filter = map[string]string{}

	for _, labelFilter := range strings.Split(labelsFilterString, ",") {
		if strings.Contains(labelFilter, "=") {
			labelFilterSplit := strings.Split(labelFilter, "=")

			labelFilterKey := labelFilterSplit[0]
			labelFilterValue := labelFilterSplit[1]

			l.Filter[labelFilterKey] = labelFilterValue
		}
	}
}
