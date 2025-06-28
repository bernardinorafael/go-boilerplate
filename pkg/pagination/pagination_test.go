package pagination

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPagination_New(t *testing.T) {
	tests := []struct {
		name         string
		data         []string
		totalItems   int
		currentPage  int
		itemsPerPage int
		expected     Paginated[string]
	}{
		{
			name:         "should return correct pagination for first page with data",
			data:         []string{"item-1", "item-2", "item-3"},
			totalItems:   10,
			currentPage:  1,
			itemsPerPage: 3,
			expected: Paginated[string]{
				Data: []string{"item-1", "item-2", "item-3"},
				Meta: Metadata{
					TotalItems:      10,
					CurrentPage:     1,
					ItemsPerPage:    3,
					TotalPages:      4,
					HasPreviousPage: false,
					HasNextPage:     true,
					IsFirstPage:     true,
					IsLastPage:      false,
				},
			},
		},
		{
			name:         "should return correct pagination for second page with data",
			data:         []string{"item-4", "item-5", "item-6"},
			totalItems:   10,
			currentPage:  2,
			itemsPerPage: 3,
			expected: Paginated[string]{
				Data: []string{"item-4", "item-5", "item-6"},
				Meta: Metadata{
					TotalItems:      10,
					CurrentPage:     2,
					ItemsPerPage:    3,
					TotalPages:      4,
					HasPreviousPage: true,
					HasNextPage:     true,
					IsFirstPage:     false,
					IsLastPage:      false,
				},
			},
		},
		{
			name:         "should return correct pagination for last page with data",
			data:         []string{"item-10"},
			totalItems:   10,
			currentPage:  4,
			itemsPerPage: 3,
			expected: Paginated[string]{
				Data: []string{"item-10"},
				Meta: Metadata{
					TotalItems:      10,
					CurrentPage:     4,
					ItemsPerPage:    3,
					TotalPages:      4,
					HasPreviousPage: true,
					HasNextPage:     false,
					IsFirstPage:     false,
					IsLastPage:      true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.data, tt.totalItems, tt.currentPage, tt.itemsPerPage)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("New() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_CalculatePages(t *testing.T) {
	baseTestName := "should return %d when totalItems=%d and itemsPerPage=%d"

	tests := []struct {
		totalItems   int
		itemsPerPage int
		expected     int
	}{
		{totalItems: 0, itemsPerPage: 5, expected: 0},
		{totalItems: 1, itemsPerPage: 5, expected: 1},
		{totalItems: 5, itemsPerPage: 5, expected: 1},
		{totalItems: 6, itemsPerPage: 5, expected: 2},
		{totalItems: 10, itemsPerPage: 5, expected: 2},
		{totalItems: 11, itemsPerPage: 5, expected: 3},
		{totalItems: 15, itemsPerPage: 5, expected: 3},
		{totalItems: 16, itemsPerPage: 5, expected: 4},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(baseTestName, tt.expected, tt.totalItems, tt.itemsPerPage), func(t *testing.T) {
			got := meta(tt.totalItems, 1, tt.itemsPerPage)
			if got.TotalPages != tt.expected {
				t.Errorf(
					"TotalPages calculation: totalItems=%d, itemsPerPage=%d, got=%d, expected=%d",
					tt.totalItems,
					tt.itemsPerPage,
					got.TotalPages,
					tt.expected,
				)
			}
		})
	}
}
