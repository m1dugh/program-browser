package utils


func CompareStrings(a string, b string) int {
    la := len(a)
    lb := len(b)
    var min int

    if la > lb {
        min = lb
    } else {
        min = la
    }
    var res int = 0

    for i := 0; i < min && res == 0; i++ {
        res = int(a[i]) - int(b[i])
    }

    if res == 0 {
        if la < lb {
            return -1
        } else if la == lb {
            return 0
        } else {
            return 1
        }
    }

    return res
}
