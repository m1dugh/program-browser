package utils

type StringSet []string;

func NewStringSet(values []string) *StringSet {
    
    var res *StringSet = &StringSet{}
    *res = make([]string, len(values))

    for _, v := range values {
        res.AddWord(v)
    }

    return res
}

func (set *StringSet) Length() int {
    return len(*set)
}

/// Returns the underlying array.
/// Caution: Since it does not return a copy, modifying the underlying array
/// can break the StringSet
func (set *StringSet) UnderlyingArray() []string {
    return *set
}

func (set *StringSet) _binsearch(value string) (int, bool) {

    start := 0
    end := len(*set)

    for start < end {
        middle := start + (end - start) / 2
        s := (*set)[middle]
        res := CompareStrings(value, s)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *StringSet) _insertAt(value string, pos int) {
    *set = append(*set, value)
    for i := len(*set) - 1; i > pos; i-- {
        (*set)[i] = (*set)[i - 1]
    }

    (*set)[pos] = value
}

func (set *StringSet) AddWord(value string) bool {
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *StringSet) ContainsWord(value string) bool {
    _, found := set._binsearch(value)
    return found
}

func (set *StringSet) AddAll(other *StringSet) int {
    var count int = 0
    for _, v := range *other {
        if set.AddWord(v) {
            count++
        }
    }

    return count
}

func (set *StringSet) ToArray() []string {
    res := make([]string, len(*set))
    copy(res, *set)
    return res
}

type Comparable interface {
    Compare(value interface{}) int
}

type ComparableSet[T Comparable] []T

func NewComparableSet[T Comparable](values []T) *ComparableSet[T] {

    var res *ComparableSet[T] = &ComparableSet[T]{}
    *res = make(ComparableSet[T], len(values))

    for _, v := range values {
        res.AddElement(v)
    }

    return res
}

func (set *ComparableSet[T]) _binsearch(value T) (int, bool) {

    start := 0
    end := len(*set)

    for start < end {
        middle := start + (end - start) / 2
        elem := (*set)[middle]
        res := value.Compare(elem)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *ComparableSet[T]) _insertAt(value T, pos int) {
    *set = append(*set, value)
    for i := len(*set) - 1; i > pos; i-- {
        (*set)[i] = (*set)[i - 1]
    }

    (*set)[pos] = value
}

func (set *ComparableSet[T]) AddElement(value T) bool {
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *ComparableSet[T]) Contains(value T) bool {
    _, found := set._binsearch(value)
    return found
}

func (set *ComparableSet[T]) AddAll(other *ComparableSet[T]) int {
    var count int = 0
    for _, v := range *other {
        if set.AddElement(v) {
            count++
        }
    }

    return count
}

func (set *ComparableSet[T]) Length() int {
    return len(*set)
}

/// Returns the underlying array.
/// Caution: Since it does not return a copy, modifying the underlying array
/// can break the ComparableSet
func (set *ComparableSet[T]) UnderlyingArray() []T {
    return *set
}

func (set *ComparableSet[T]) ToArray() []T {
    res := make([]T, len(*set))
    copy(res, *set)
    return res
}
