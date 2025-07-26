package generator

import "ml-based-cache/internal/models"

func Converter(items []int) *[]models.CacheItem {
	converted := make([]models.CacheItem, 0, len(items))

	for _, v := range items {
		converted = append(converted, models.CacheItem{
			Val: v,
		})
	}

	return &converted
}

func ReConverter(items *[]models.CacheItem) []int {
	converted := make([]int, 0, len(*items))

	for _, v := range *items {
		converted = append(converted, v.Val)
	}

	return converted
}
