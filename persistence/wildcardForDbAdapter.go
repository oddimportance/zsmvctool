package persistence

import ()

type Wildcard string

type OrderBy string

type MathComparision string

const (
	notWildcard      Wildcard = ""
	wildcardPrefix            = "%#"
	wildcardTrailing          = "#%"
	wildcardBoth              = "%#%"
)

const (
	desc OrderBy = "DESC"
	asc          = "ASC"
)

const (
	lessThan           MathComparision = "<"
	greaterThan                        = ">"
	lessThanEqualTo                    = "<="
	greaterThanEqualTo                 = ">="
	equalTo                            = "="
	like                               = "LIKE"
	isNull                             = "IS NULL"
	notLike                            = "NOT LIKE"
)

// Not a wildcard search object
func (w Wildcard) NotWildcard() Wildcard {
	return notWildcard
}

// Use this to search wildcard search
// with LIKE '%test'
func (w Wildcard) PrefixWildcard() Wildcard {
	return wildcardPrefix
}

// Use this to search wildcard search
// with LIKE 'test%'
func (w Wildcard) TrailingWildcard() Wildcard {
	return wildcardTrailing
}

// Use this to search wildcard search
// with LIKE '%test%'
func (w Wildcard) Wildcard() Wildcard {
	return wildcardBoth
}

func (o OrderBy) Desc() OrderBy {
	return desc
}

func (o OrderBy) Asc() OrderBy {
	return asc
}

func (m MathComparision) LessThan() MathComparision {
	return lessThan
}

func (m MathComparision) GreaterThan() MathComparision {
	return greaterThan
}

func (m MathComparision) LessThanEqualTo() MathComparision {
	return lessThanEqualTo
}

func (m MathComparision) GreaterThanEqualTo() MathComparision {
	return greaterThanEqualTo
}

func (m MathComparision) EqalTo() MathComparision {
	return equalTo
}

func (m MathComparision) Like() MathComparision {
	return like
}

func (m MathComparision) NotLike() MathComparision {
	return notLike
}

func (m MathComparision) IsNull() MathComparision {
	return isNull
}
