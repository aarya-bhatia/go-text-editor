package internal

func FlattenList[T any](array [][]T) []T {
  res := []T{}
  for _, element := range array {
    res = append(res, element...)
  }
  return res
}
