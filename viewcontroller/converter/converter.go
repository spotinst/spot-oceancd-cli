package converter

import (
	"bytes"
	"fmt"
	"spot-oceancd-cli/pkg/oceancd/model/phase"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	"spot-oceancd-cli/pkg/oceancd/model/verification"
	"strconv"
	"strings"
	"unicode"
)

const stubCell = "--"

func Weight(phase phase.Phase) string {
	if phase.TrafficPercentage < 1 {
		return stubCell
	} else {
		return strconv.Itoa(phase.TrafficPercentage)
	}
}

func PhaseIndex(index int) string {
	return fmt.Sprintf("%s %02d", "Phase", index)
}

func PhaseName(phase phase.Phase) string {
	if phase.Name == "" {
		return stubCell
	}
	return phase.Name
}

func PhaseStatus(phase phase.Phase) string {
	return fromCamelCaseToNormal(string(phase.Status))
}

func VerificationStatus(verification verification.Verification) string {
	return strings.ToUpper(string(verification.Status))
}

func RolloutStatus(status rollout.Status) string {
	return fromCamelCaseToNormal(string(status))
}

func fromCamelCaseToNormal(raw string) string {
	buf := &bytes.Buffer{}
	for i, char := range raw {
		if i == 0 {
			buf.WriteRune(unicode.ToUpper(char))
			continue
		}
		if unicode.IsUpper(char) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(unicode.ToLower(char))
	}
	return buf.String()
}
