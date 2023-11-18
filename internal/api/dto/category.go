package dto

import (
	"sort"
	"techno-store/internal/domain/bo"
)

type Category struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id,omitempty"`
	Sequence int64  `json:"sequence"`
	StatusID int64  `json:"status_id"`
}

func ToCategoryDTO(bo bo.Category) Category {
	return Category{
		ID:       bo.ID,
		Name:     bo.Name,
		ParentID: bo.ParentID,
		Sequence: bo.Sequence,
		StatusID: bo.StatusID,
	}
}

func (c Category) Model() bo.Category {
	return bo.Category{
		ID:       c.ID,
		Name:     c.Name,
		ParentID: c.ParentID,
		Sequence: c.Sequence,
		StatusID: c.StatusID,
	}
}

type CategoryUpdate struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name"`
	StatusID *int64  `json:"status_id"`
	Sequence *int64  `json:"sequence"`
}

func (c CategoryUpdate) Model() bo.CategoryUpdate {
	return bo.CategoryUpdate{
		ID:       c.ID,
		Name:     c.Name,
		StatusID: c.StatusID,
		Sequence: c.Sequence,
	}
}

type CategoryTree struct {
	ID       int64           `json:"id"`
	Name     string          `json:"category_name"`
	ParentID int64           `json:"-"`
	Depth    int64           `json:"-"`
	Sequence int64           `json:"-"`
	Children []*CategoryTree `json:"children,omitempty"`
}

type CategoriesTree struct {
	Data []*CategoryTree `json:"data"`
}

func categoryCollectionDto(categories bo.CategoryCollection) []CategoryTree {
	var categoryDTOs []CategoryTree
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, CategoryTree{
			ID:       category.ID,
			Name:     category.Name,
			ParentID: category.ParentID,
			Sequence: category.Sequence,
		})
	}
	return categoryDTOs
}

func CategoriesToTree(pCategories bo.PaginatedCategoryCollection) CategoriesTree {
	tree := buildTree(categoryCollectionDto(pCategories.Data))
	sortCategoryTree(tree)

	return CategoriesTree{
		Data: tree,
	}
}

func buildTree(categories []CategoryTree) []*CategoryTree {
	var tree []*CategoryTree
	childrenMap := make(map[int64][]*CategoryTree)

	for i, cat := range categories {
		cat.Depth = 1 // Root nodes start with depth 1

		if cat.ParentID == 0 {
			// It's a root node
			tree = append(tree, &categories[i])
		} else {
			// It's a child node
			childrenMap[cat.ParentID] = append(childrenMap[cat.ParentID], &categories[i])
		}
	}

	for _, branch := range tree {
		if children, ok := childrenMap[branch.ID]; ok {
			// Increase the depth for each child
			for _, child := range children {
				child.Depth = branch.Depth + 1
			}

			// Only add children if parent's depth is less than 3
			if branch.Depth < 3 {
				branch.Children = children
			}
		}
	}

	return tree
}

func sortCategoryTree(categories []*CategoryTree) {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Sequence < categories[j].Sequence
	})

	// Recursively sort the children of each category
	for _, cat := range categories {
		if len(cat.Children) > 0 {
			sortCategoryTree(cat.Children)
		}
	}
}
