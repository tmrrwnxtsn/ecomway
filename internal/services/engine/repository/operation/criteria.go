package operation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (r *Repository) whereStmt(c model.OperationCriteria) (string, []any, error) {
	nonNilCriteria, nonNilCriterionArgs := r.nonNilCriteriaArgs(c)
	if nonNilCriteria == 0 {
		return "", nil, errors.New("criteria args not stated")
	}

	whereValues := make([]string, 0, nonNilCriteria)
	args := make([]any, 0, nonNilCriterionArgs)
	currArgID := 1

	if c.ID != nil {
		whereValues = append(whereValues, fmt.Sprintf("%s.id=$%d", operationTableAbbr, currArgID))
		args = append(args, *c.ID)
		currArgID++
	}

	if c.UserID != nil {
		whereValues = append(whereValues, fmt.Sprintf("%s.user_id=$%d", operationTableAbbr, currArgID))
		args = append(args, *c.UserID)
		currArgID++
	}

	if c.Types != nil {
		argsInRoundBrackets := make([]string, 0, len(*c.Types))
		for _, t := range *c.Types {
			argsInRoundBrackets = append(argsInRoundBrackets, fmt.Sprintf("$%d", currArgID))
			args = append(args, t)
			currArgID++
		}
		inRoundBrackets := strings.Join(argsInRoundBrackets, ", ")
		whereValues = append(whereValues, fmt.Sprintf("%s.type IN (%s)", operationTableAbbr, inRoundBrackets))
	}

	if c.Statuses != nil {
		argsInRoundBrackets := make([]string, 0, len(*c.Statuses))
		for _, s := range *c.Statuses {
			argsInRoundBrackets = append(argsInRoundBrackets, fmt.Sprintf("$%d", currArgID))
			args = append(args, s)
			currArgID++
		}
		inRoundBrackets := strings.Join(argsInRoundBrackets, ", ")
		whereValues = append(whereValues, fmt.Sprintf("%s.status IN (%s)", operationTableAbbr, inRoundBrackets))
	}

	if len(c.StatusesByType) > 0 && c.Types == nil && c.Statuses == nil {
		var filter string
		for operationType, operationStatuses := range c.StatusesByType {
			if filter == "" {
				filter = fmt.Sprintf("(%s.type=$%d AND ", operationTableAbbr, currArgID)
			} else {
				filter += fmt.Sprintf("OR (%s.type=$%d AND ", operationTableAbbr, currArgID)
			}
			args = append(args, operationType)
			currArgID++

			argsInRoundBrackets := make([]string, 0, len(operationStatuses))
			for _, s := range operationStatuses {
				argsInRoundBrackets = append(argsInRoundBrackets, fmt.Sprintf("$%d", currArgID))
				args = append(args, s)
				currArgID++
			}
			inRoundBrackets := strings.Join(argsInRoundBrackets, ", ")
			filter += fmt.Sprintf("%s.status IN (%s))", operationTableAbbr, inRoundBrackets)
		}
		whereValues = append(whereValues, fmt.Sprintf("(%v)", filter))
	}

	if c.ExternalSystems != nil {
		argsInRoundBrackets := make([]string, 0, len(*c.ExternalSystems))
		for _, es := range *c.ExternalSystems {
			argsInRoundBrackets = append(argsInRoundBrackets, fmt.Sprintf("$%d", currArgID))
			args = append(args, es)
			currArgID++
		}
		inRoundBrackets := strings.Join(argsInRoundBrackets, ", ")
		whereValues = append(whereValues, fmt.Sprintf("%s.external_system IN (%s)", operationTableAbbr, inRoundBrackets))
	}

	if c.ExternalID != nil {
		whereValues = append(whereValues, fmt.Sprintf("%s.external_id=$%d", operationTableAbbr, currArgID))
		args = append(args, *c.ExternalID)
		currArgID++
	}

	if !c.CreatedAtFrom.IsZero() {
		whereValues = append(whereValues, fmt.Sprintf("%s.created_at>=$%d", operationTableAbbr, currArgID))
		args = append(args, c.CreatedAtFrom)
		currArgID++
	}

	if !c.CreatedAtTo.IsZero() {
		whereValues = append(whereValues, fmt.Sprintf("%s.created_at<=$%d", operationTableAbbr, currArgID))
		args = append(args, c.CreatedAtTo)
		currArgID++
	}

	whereStmt := strings.Join(whereValues, " AND ")
	return "WHERE " + whereStmt, args, nil
}

func (r *Repository) nonNilCriteriaArgs(c model.OperationCriteria) (int, int) {
	nonNilCriteria, nonNilCriteriaArgs := 0, 0
	if c.ID != nil {
		nonNilCriteria++
		nonNilCriteriaArgs++
	}
	if c.UserID != nil {
		nonNilCriteria++
		nonNilCriteriaArgs++
	}
	if c.ExternalID != nil {
		nonNilCriteria++
		nonNilCriteriaArgs++
	}
	if c.Types != nil {
		nonNilCriteria++
		nonNilCriteriaArgs += len(*c.Types)
	}
	if c.Statuses != nil {
		nonNilCriteria++
		nonNilCriteriaArgs += len(*c.Statuses)
	}
	if c.ExternalSystems != nil {
		nonNilCriteria++
		nonNilCriteriaArgs += len(*c.ExternalSystems)
	}
	if len(c.StatusesByType) > 0 {
		nonNilCriteria++
		for _, statuses := range c.StatusesByType {
			nonNilCriteriaArgs += len(statuses) + 1
		}
	}
	if !c.CreatedAtFrom.IsZero() {
		nonNilCriteria++
		nonNilCriteriaArgs++
	}
	if !c.CreatedAtTo.IsZero() {
		nonNilCriteria++
		nonNilCriteriaArgs++
	}
	return nonNilCriteria, nonNilCriteriaArgs
}
