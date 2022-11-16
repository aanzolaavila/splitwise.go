package splitwise

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getAndCheckIntExpensesParam(t *testing.T) {
	const (
		testField                = ExpensesGroupId
		expensesGroupIdAsStr     = "25"
		expensesGroupIdAsInt     = 25
		expensesGroupInvalidStr  = "invalid"
		expensesGroupInvalidType = 25.0
	)

	ps := ExpensesParams{
		testField: expensesGroupIdAsStr,
	}

	t.Run("StringRepresentingIntShouldPass", func(t *testing.T) {
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesGroupIdAsStr, s)
	})

	t.Run("IntShouldPass", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupIdAsInt
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesGroupIdAsStr, s)
	})

	t.Run("StringNOTRepresentingIntShouldFail", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupInvalidStr
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("InvalidTypeShouldFail", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupInvalidType
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("NonExistentParamShouldGiveZeroValueAndNotFail", func(t *testing.T) {
		s, err := getAndCheckIntExpensesParam(ps, ExpensesDatedAfter)
		assert.NoError(t, err)
		assert.Zero(t, s)
	})
}

func Test_getAndCheckDateExpensesParam(t *testing.T) {
	const (
		testField                             = ExpensesDatedBefore
		expensesDatedBeforeAsStr       string = "2022-01-01T12:00:00Z"
		expensesDatedBeforeInvalidStr         = "invalid"
		expensesDatedBeforeInvalidType        = 20.0
	)

	expensesDatedBeforeAsDate, err := time.Parse(time.RFC3339, expensesDatedBeforeAsStr)
	if err != nil {
		require.FailNowf(t, "failed to create date for testing: %s", err.Error())
	}

	ps := ExpensesParams{
		ExpensesDatedBefore: expensesDatedBeforeAsDate,
	}

	t.Run("DateTypeShouldPass", func(t *testing.T) {
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesDatedBeforeAsStr, s)
	})

	t.Run("StringAsValidDateShouldPass", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeAsStr
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesDatedBeforeAsStr, s)
	})

	t.Run("InvalidStringShouldFail", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeInvalidStr
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("NonExistentParamShouldGiveZeroValueAndNotFail", func(t *testing.T) {
		s, err := getAndCheckDateExpensesParam(ps, ExpensesDatedAfter)
		assert.NoError(t, err)
		assert.Zero(t, s)
	})

	t.Run("InvalidTypeShouldFail", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeInvalidType
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})
}
