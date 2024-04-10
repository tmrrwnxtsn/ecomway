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
	return whereStmt, args, nil
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
	if c.Types != nil {
		nonNilCriteria++
		nonNilCriteriaArgs += len(*c.Types)
	}
	if c.Statuses != nil {
		nonNilCriteria++
		nonNilCriteriaArgs += len(*c.Statuses)
	}
	if c.ExternalID != nil {
		nonNilCriteria++
		nonNilCriteriaArgs++
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
