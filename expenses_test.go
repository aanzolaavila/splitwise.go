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

func Test_expensesParamsToUrlValues(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
	)

	const timeFormat = time.RFC3339
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	ps := ExpensesParams{
		// ints
		ExpensesGroupId:  1,
		ExpensesFriendId: "2",
		ExpensesLimit:    3,
		ExpensesOffset:   "4",
		// dates
		ExpensesDatedBefore:   now,
		ExpensesDatedAfter:    now.Format(timeFormat),
		ExpensesUpdatedBefore: now,
		ExpensesUpdatedAfter:  yesterday,
	}

	vals, err := expensesParamsToUrlValues(ps)
	require.NoError(err)
	require.NotNil(vals)
	assert.Len(vals, len(ps))

	test := func(field expensesParam, expected string) {
		if k := string(field); assert.Contains(vals, k) {
			vs := vals[k]
			require.Len(vs, 1)
			v := vs[0]
			assert.Equal(expected, v)
		}
	}

	test(ExpensesGroupId, "1")
	test(ExpensesFriendId, "2")
	test(ExpensesLimit, "3")
	test(ExpensesOffset, "4")
	test(ExpensesDatedBefore, now.Format(timeFormat))
	test(ExpensesDatedAfter, now.Format(timeFormat))
	test(ExpensesUpdatedBefore, now.Format(timeFormat))
	test(ExpensesUpdatedAfter, yesterday.Format(timeFormat))
}

func Test_expensesParamsToUrlValues_ErrorCases(t *testing.T) {
	ps := ExpensesParams{
		// ints
		ExpensesGroupId: 1.0,
	}

	vals, err := expensesParamsToUrlValues(ps)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
	assert.Zero(t, vals)

	ps = ExpensesParams{
		// dates
		ExpensesDatedBefore: 2.0,
	}

	vals, err = expensesParamsToUrlValues(ps)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
	assert.Zero(t, vals)
}
