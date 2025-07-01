package main

import "testing"

func TestPsqlStyle(t *testing.T) {
	// Test cases for camelStyle function
	tests := []struct {
		input    string
		expected string
	}{
		{"PurchaseOrderID", "purchase_order_id"},
		{"ShopID", "shop_id"},
		//{"ProductSKUID", "product_sku_id"},
		{"ProductID", "product_id"},
		{"ProductName", "product_name"},
		{"ProductFeaturedImage", "product_featured_image"},
		{"ProductFirstOptionID", "product_first_option_id"},
		{"ProductFirstOptionName", "product_first_option_name"},
		{"ProductFirstOptionChoiceID", "product_first_option_choice_id"},
		{"ProductFirstOptionChoiceValue", "product_first_option_choice_value"},
	}

	for _, test := range tests {
		result := psqlStyle(test.input)
		if result != test.expected {
			t.Errorf("psqlStyle(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
