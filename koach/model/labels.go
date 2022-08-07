package model

import "strings"

func LabelsFromString(labelsFilterString string) map[string]string {
	labels := map[string]string{}

	for _, labelFilter := range strings.Split(labelsFilterString, ",") {
		if strings.Contains(labelFilter, "=") {
			labelFilterSplit := strings.Split(labelFilter, "=")

			labelFilterKey := labelFilterSplit[0]
			labelFilterValue := labelFilterSplit[1]

			labels[labelFilterKey] = labelFilterValue
		}
	}

	return labels
}
