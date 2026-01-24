package helpers

import (
	"context"
	"fmt"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	"gorm.io/gorm"
)

func CreateCashFlow(ctx context.Context, data schemas.CashFlowCreateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.createCashFlow")
	defer span.End()

	cashflow := schemas.CashFlow{
		Amount:          data.Amount,
		Month:           data.Month,
		Day:             data.Day,
		IncomeOrExpense: data.IncomeOrExpense,
		CategoryId:      data.CategoryId,
		AssociationId:   data.AssociationId,
	}

	result := gorm.WithResult()
	err = gorm.G[schemas.CashFlow](schemas.DB, result).Create(ctx, &cashflow)

	if err != nil {
		return nil, err
	}

	return schemas.ConvertStructToMap(cashflow)
}

func FilterCashFlow(ctx context.Context, queryParams *map[string]interface{}) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.FilterCashFlow")
	defer span.End()

	queries := *queryParams
	offset := queries["page"].(int) - 1*queries["items"].(int)

	resp, err := gorm.G[schemas.CashFlow](schemas.DB).Offset(offset).Limit(queries["items"].(int)).Where("").Find(ctx)

	if err != nil {
		return nil, err
	}

	response = map[string]interface{}{
		"items": queries["items"].(int),
		"page":  queries["page"].(int),
		"data":  resp,
	}

	return response, err
}

func UpdateCashFlow(ctx context.Context, cashflowId int, updateCashflowSchema schemas.CashFlowUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.updateCashFlow")
	defer span.End()

	cashflow := schemas.CashFlow{
		Amount:          updateCashflowSchema.Amount,
		Month:           updateCashflowSchema.Month,
		Day:             updateCashflowSchema.Day,
		IncomeOrExpense: updateCashflowSchema.IncomeOrExpense,
		CategoryId:      updateCashflowSchema.CategoryId,
	}

	var rowsAffected int

	rowsAffected, err = gorm.G[schemas.CashFlow](schemas.DB).Where("id = ? ", cashflowId).Updates(ctx, cashflow)
	if err != nil {
		groupUpdateError := schemas.ErrorList{
			ResponseCode: "CF001",
			Message:      "An error occured",
			Errors:       []string{fmt.Sprintf("An error occured updating group with id %d", cashflowId)},
		}

		span.SetAttributes(utils.MapToAttributes(
			map[string]interface{}{"errors": []string{err.Error()}})...,
		)

		return schemas.ConvertStructToMap(groupUpdateError)
	}

	if rowsAffected < 1 {
		groupUpdateError := schemas.ErrorList{
			ResponseCode: "GR000",
			Message:      "An error occured",
			Errors:       []string{fmt.Sprintf("Cashflow with ID %d not found", cashflowId)},
		}

		span.SetAttributes(utils.MapToAttributes(
			map[string]interface{}{"errors": groupUpdateError.Errors})...,
		)

		return schemas.ConvertStructToMap(groupUpdateError)

	}
	// If update happend it means the group exists
	updatedGroup, _ := gorm.G[schemas.Group](schemas.DB).Where("id = ? ", cashflowId).First(ctx)
	return schemas.ConvertStructToMap(updatedGroup)

}

func GUDCashFlow(ctx context.Context, categoryId int, method string, updateCashFlowSchema schemas.CashFlowUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.GUDCashFlow")
	defer span.End()

	var category schemas.Category

	if method == "GET" {
		category, err = gorm.G[schemas.Category](schemas.DB).Where("id = ? ", categoryId).First(ctx)
		return schemas.ConvertStructToMap(category)
	} else if method == "DELETE" {
		_, err = gorm.G[schemas.Category](schemas.DB).Where("id = ? ", categoryId).Delete(ctx)
		response := map[string]interface{}{
			"message": "successful",
		}
		return response, nil
	} else {
		return UpdateCashFlow(ctx, categoryId, updateCashFlowSchema)
	}
}

type CashFlowNode struct {
	Id                    int             `json:"id"`
	NodeType              string          `json:"type"` // cashflow/category/topic
	ChildNodes            []*CashFlowNode `json:"categories"`
	ParentCategoryId      int             `json:"-"` // `json:"parent_category_id"`
	Name                  string          `json:"name"`
	TotalIn               int             `json:"total_in"`
	TotalOut              int             `json:"total_out"`
	ParentNodeType        string          `json:"-"`
	ParentAlreadyAssigned bool            `json:"-"`
}

// This should be enough to differenciate categories and expenses
type CategoryKey struct {
	Id   int
	Type string
}

func FetchCashFlowTree(ctx context.Context, topicId int) (response map[string]interface{}, err error) {
	cashFlows, err := gorm.G[schemas.CashFlow](schemas.DB).Where("association_id = ? ", topicId).Find(ctx)

	if len(cashFlows) < 1 {
		response = map[string]interface{}{
			"data": cashFlows,
		}
		return response, err
	}

	topic, err := gorm.G[schemas.Category](schemas.DB).Where("id = ? ", topicId).First(ctx)

	categoryMap := make(map[CategoryKey]*CashFlowNode)

	rootNode := CashFlowNode{
		Id:       topicId,
		NodeType: "topic",
		Name:     topic.Name,
	}

	rootKey := CategoryKey{
		Id:   topicId,
		Type: "topic",
	}

	categoryMap[rootKey] = &rootNode

	genParentNodes(ctx, rootKey, &categoryMap)

	for _, val := range cashFlows {
		cfKey := CategoryKey{
			Id:   int(val.ID),
			Type: val.IncomeOrExpense,
		}

		cfNode := CashFlowNode{
			Id:               int(val.ID),
			NodeType:         val.IncomeOrExpense,
			ParentCategoryId: val.CategoryId,
			Name:             val.Description, // Use decriptions for the final node i.e the expense or income
		}

		if val.IncomeOrExpense == "income" {
			cfNode.TotalIn = int(val.Amount)
		} else {
			cfNode.TotalOut = int(val.Amount)
		}

		categoryMap[cfKey] = &cfNode

		categoryKey := CategoryKey{
			Id:   int(val.CategoryId),
			Type: "category",
		}

		cfNode.assignCflowNodes(&categoryMap, categoryKey, nil)
	}

	summarizedData := *(categoryMap[rootKey])

	summaryMap, err := schemas.ConvertStructToMap(summarizedData)

	response = map[string]interface{}{
		"data": summaryMap,
	}

	return response, nil
}

// Use this to generate all catgory nodes at once as opposed to having to do this on every iteration of the cashflow
func genParentNodes(ctx context.Context, catKey CategoryKey, parentIdMapsPtr *map[CategoryKey](*CashFlowNode)) (err error) {

	// Get the categories where the parent is the category ID entered
	categories, err := gorm.G[schemas.Category](schemas.DB).Where("parent_category = ? ", catKey.Id).Find(ctx)
	if len(categories) < 1 {
		return
	}

	for _, val := range categories {
		key := CategoryKey{
			Id:   int(val.ID),
			Type: val.CategoryType,
		}
		(*parentIdMapsPtr)[key] = &CashFlowNode{
			Id:               int(val.ID),
			NodeType:         val.CategoryType,
			ParentCategoryId: catKey.Id,
			ParentNodeType:   catKey.Type,
			Name:             val.Name,
		}

		genParentNodes(ctx, key, parentIdMapsPtr)
	}

	return
}

/*
	categoryId int

-- This is the Id of the parent category that the current cashflow or category belongs to

	parentIdMapsPtr *map[int]ParentCategoryKey

-- Pointer to the map containing the categories and sub categories tied to the context topic.
Keyed by their IDs with the values being the details of the parent category

currentCFlowNodePtr *CashFlowNode
-- current node being updated. Ideally the parent to whatever node is passed as the childNode

childNode *CashFlowNode
-- prev node in heirachy
*/
func (currentCFlowNodePtr *CashFlowNode) assignCflowNodes(parentIdMapsPtr *map[CategoryKey]*CashFlowNode, parentKey CategoryKey, childNodePtr *CashFlowNode) {
	fmt.Println("?????00????\n")
	fmt.Println("currentCFlowNodePtr ", currentCFlowNodePtr)
	fmt.Println("currentCFlowNodePtr.ChildNodes ", currentCFlowNodePtr.ChildNodes)
	fmt.Println("parentKey  ", parentKey)
	fmt.Println("childNodePtr ", childNodePtr)
	fmt.Println(" currentCFlowNodePtr.ParentCategoryId  ", currentCFlowNodePtr.ParentCategoryId)

	// Add the child nodes
	if currentCFlowNodePtr.NodeType != "income" && currentCFlowNodePtr.NodeType != "expense" {
		currentCFlowNodePtr.ChildNodes = append(currentCFlowNodePtr.ChildNodes, childNodePtr)
	}
	// Inside the function, after adding the child:
	if childNodePtr != nil {
		currentCFlowNodePtr.TotalOut += childNodePtr.TotalOut
		currentCFlowNodePtr.TotalIn += childNodePtr.TotalIn
	}

	if currentCFlowNodePtr.NodeType == "topic" {
		return
	}

	// fmt.Println("The parent Node created ------------ ", parentNode)
	// Parent to the parent of the current node
	// i.e the parent for the next loop
	keyData := (*parentIdMapsPtr)[parentKey]

	if currentCFlowNodePtr.ParentAlreadyAssigned {
		return
	}

	// Given a category cannot be a child in 2 arrays
	// i.e cannot have two parents.
	// Once you have appended its pointer to an array of the parent don't do it again
	// In the event of very nested summaries it prevents redundant looping
	currentCFlowNodePtr.ParentAlreadyAssigned = true

	newParentKey := CategoryKey{
		Id:   keyData.ParentCategoryId,
		Type: keyData.ParentNodeType,
	}

	fmt.Println("keyData  ------------ ", keyData)
	fmt.Println("newKey  ------------ ", newParentKey)
	fmt.Println("The currentCFlowNodePtr Node updated  ------------ ", currentCFlowNodePtr)

	(*parentIdMapsPtr)[parentKey].assignCflowNodes(parentIdMapsPtr, newParentKey, currentCFlowNodePtr)

}
