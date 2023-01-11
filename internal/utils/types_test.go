package utils

import (
    "testing"
    "reflect"
    "fmt"
)

type testVal struct {
    val int
    name string
}

func (v testVal) String() string {
    return fmt.Sprintf("%d: %s", v.val, v.name)
}

func (t testVal) Compare(other interface{}) int {
    t2, ok := other.(testVal)
    if ok {
        return t.val - t2.val
    } else {
        return t.val
    }
}

func TestComparableSet(t *testing.T) {
    var set *ComparableSet[testVal] = new(ComparableSet[testVal])
    set.AddElement(testVal{
        3,
        "test1",
    })

    set.AddElement(testVal{
        2,
        "test0",
    })

    set.AddElement(testVal{
        4,
        "test2",
    })
    expected := []testVal{
        testVal{
            2,
            "test0",
        },
        testVal{
            3,
            "test1",
        },
        testVal{
            4,
            "test2",
        },
    }
    res := set.ToArray()
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("value mismatch: \nreceived: %s\nexpected: %s", res, expected)
    }
}
