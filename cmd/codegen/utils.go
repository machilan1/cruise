package main

import (
	"strings"
	"unicode"
)

const (
	directAppNameKey = 0
	varTypeKey       = 1
	dummyKey         = 2
	dbPrefixKey      = 3
	directJSONName   = 4
	directDBName     = 5
)

const (
	noDBPrefix = ""
	typeString = "string"
	typeInt    = "int"
	typeBool   = "bool"
	typeTime   = "time.Time"

	dbPrefix    = "auction"
	dbTableName = "auctions"
)

func projectFieldsz() [][]string {
	// directAppName, varType, dummyKey, dbPrefix, directJSONName, directDBName
	data := [][]string{
		{"Date", typeTime, "", noDBPrefix, "", ""},
		{"Note", typeString, "", noDBPrefix, "", ""},
	}
	//data := [][]string{
	//	{"PurchaseOrderID", typeInt, "", noDBPrefix, "", ""},
	//	{"Etag", typeInt, "", noDBPrefix, "", ""},
	//
	//	{"GrossPaymentPriceMicros", typeInt, "", noDBPrefix, "", ""},
	//
	//	{"TransferredPriceMicros", typeInt, "", noDBPrefix, "", ""},
	//	{"CommissionPriceMicros", typeInt, "", noDBPrefix, "", ""},
	//	{"ShippingFeeMicros", typeInt, "", noDBPrefix, "", ""},
	//	{"PaymentFeeMicros", typeInt, "", noDBPrefix, "", ""},
	//}
	//data := [][]string{
	//	{"PurchaseOrderID", typeInt, "", noDBPrefix, "", "purchase_order_id"},
	//	{"ShopID", typeInt, "", noDBPrefix, "", "shop_id"},
	//	{"ProductSKUID", typeInt, "", noDBPrefix, "", "product_sku_id"},
	//	{"ProductID", typeInt, "", noDBPrefix, "", "product_id"},
	//	{"ProductName", typeString, "", noDBPrefix, "", "product_name"},
	//	{"ProductFeaturedImage", typeString, "", noDBPrefix, "", "product_featured_image"},
	//	{"ProductFirstOptionID", typeInt, "", noDBPrefix, "", "product_first_option_id"},
	//	{"ProductFirstOptionName", typeString, "", noDBPrefix, "", "product_first_option_name"},
	//	{"ProductFirstOptionChoiceID", typeInt, "", noDBPrefix, "", "product_first_option_choice_id"},
	//	{"ProductFirstOptionChoiceValue", typeString, "", noDBPrefix, "", "product_first_option_choice_value"},
	//	{"ProductSecondOptionID", typeInt, "", noDBPrefix, "", "product_second_option_id"},
	//	{"ProductSecondOptionName", typeString, "", noDBPrefix, "", "product_second_option_name"},
	//	{"ProductSecondOptionChoiceID", typeInt, "", noDBPrefix, "", "product_second_option_choice_id"},
	//	{"ProductSecondOptionChoiceValue", typeString, "", noDBPrefix, "", "product_second_option_choice_value"},
	//	{"Quantity", typeInt, "", noDBPrefix, "", "quantity"},
	//	{"UnitPriceMicros", typeInt, "", noDBPrefix, "", "unit_price_1000000x"},
	//	{"TotalPriceMicros", typeInt, "", noDBPrefix, "", "total_price_1000000x"},
	//	{"AddOnDiscountEventID", typeInt, "", noDBPrefix, "", "add_on_discount_event_id"},
	//	{"AddOnGiftEventID", typeInt, "", noDBPrefix, "", "add_on_gift_event_id"},
	//	{"BundleEventID", typeInt, "", noDBPrefix, "", "bundle_event_id"},
	//	{"DiscountEventID", typeInt, "", noDBPrefix, "", "discount_event_id"},
	//}
	//data := [][]string{
	//	{"ShopID", typeInt, "", noDBPrefix, camelStyle("shopID"), "shop_id"},
	//	{"CartID", typeInt, "", noDBPrefix, camelStyle("CartID"), "cart_id"},
	//	{"CustomerID", typeInt, "", noDBPrefix, camelStyle("CustomerID"), "customer_id"},
	//	{"PaymentRecordID", typeInt, "", noDBPrefix, camelStyle("PaymentRecordID"), "payment_record_id"},
	//
	//	{"TotalCompareAtPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalCompareAtPriceMicros"), "total_compare_at_price_1000000x"},
	//	{"TotalProductPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalProductPriceMicros"), "total_product_price_1000000x"},
	//	{"TotalDiscountPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalDiscountPriceMicros"), "total_discount_price_1000000x"},
	//	{"TotalBundleSavedPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalBundleSavedPriceMicros"), "total_bundle_saved_price_1000000x"},
	//
	//	{"TotalShippingFeeMicros", typeInt, "", noDBPrefix, camelStyle("TotalShippingFeeMicros"), "total_shipping_fee_1000000x"},
	//	{"TotalSavedShippingFeeMicros", typeInt, "", noDBPrefix, camelStyle("TotalSavedShippingFeeMicros"), "total_saved_shipping_fee_1000000x"},
	//	{"TotalPriceWithFeesMicros", typeInt, "", noDBPrefix, camelStyle("TotalPriceWithFeesMicros"), "total_price_with_fees_1000000x"},
	//	{"TotalPaymentPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalPaymentPriceMicros"), "total_payment_price_1000000x"},
	//
	//	{"TotalCouponSavedPriceMicros", typeInt, "", noDBPrefix, camelStyle("TotalCouponSavedPriceMicros"), "total_coupon_saved_price_1000000x"},
	//}
	//data := [][]string{
	//	{"AddOnEventID", typeInt, "", noDBPrefix, "addOnEventId", "add_on_event_id"},
	//	{"ShopID", typeTime, "", noDBPrefix, "shopId", "shop_id"},
	//	{"AddOnEventProductType", typeString, "", noDBPrefix, "addOnEventProductType", "add_on_event_product_type"},
	//	{"ProductID", typeInt, "", noDBPrefix, "productId", "product_id"},
	//	{"ProductSKUID", typeInt, "", noDBPrefix, "productSKUId", "product_sku_id"},
	//	{"IsEffective", typeBool, "", noDBPrefix, "isEffective", "is_effective"},
	//	{"EventPriceMicros", typeInt, "", noDBPrefix, "eventPriceMicros", "event_price_1000000x"},
	//	{"MaxQuantity", typeInt, "", noDBPrefix, "maxQuantity", "max_quantity"},
	//}

	for i := range data {
		if data[i][directJSONName] == "" {
			data[i][directJSONName] = camelStyle(data[i][directAppNameKey])
		}
		if data[i][directDBName] == "" {
			data[i][directDBName] = psqlStyle(data[i][directAppNameKey])
		}
	}

	return data
}

func patchPointerTypeFields(data [][]string) [][]string {
	for _, v := range data {
		v[varTypeKey] = "*" + v[varTypeKey]
	}
	return data
}

func coreFieldsz(data [][]string) []string {
	result := make([]string, 0)
	for _, v := range data {
		str := v[0:2]
		result = append(result, strings.Join(str, "  "))
	}

	return result
}

func dbFieldsz(data [][]string, otm bool) []string {
	result := make([]string, 0, len(data))
	for _, v := range data {
		dbResult := strings.ToLower(v[directAppNameKey])

		if v[dbPrefixKey] != "" {
			dbResult = v[dbPrefixKey] + "_" + dbResult
		}

		if v[directDBName] != "" {
			dbResult = v[directDBName]
		}
		if otm {
			dbResult = `db:"` + dbResult + `"` + ` json:"` + dbResult + `"`
		} else {
			dbResult = `db:"` + dbResult + `"`
		}
		dbResult = "`" + dbResult + "`"
		resul := make([]string, 3)
		resul[0] = v[directAppNameKey]
		resul[1] = v[varTypeKey]
		resul[2] = dbResult
		result = append(result, strings.Join(resul, "  "))
	}
	return result
}

func jsonFieldsz(data [][]string) []string {
	result := make([]string, 0, len(data))
	for _, v := range data {
		jsonResult := strings.ToLower(v[directAppNameKey][0:1]) + v[directAppNameKey][1:]
		if jsonResult[len(jsonResult)-2:] == "ItemID" {
			jsonResult = jsonResult[0:len(jsonResult)-2] + "Id"
		}
		if v[directJSONName] != "" {
			jsonResult = v[directJSONName]
		}
		jsonResult = `json:"` + jsonResult + `"`
		jsonResult = "`" + jsonResult + "`"
		resul := make([]string, 3)
		resul[0] = v[directAppNameKey]
		resul[1] = v[varTypeKey]
		resul[2] = jsonResult
		result = append(result, strings.Join(resul, "  "))
	}
	return result
}

func toStructFields(data [][]string, abbr string) []string {
	result := make([]string, 0, len(data))
	for _, v := range data {
		tempString := v[directAppNameKey] + ":" + " " + abbr + "." + v[directAppNameKey] + ","
		result = append(result, tempString)
	}
	return result
}

func toAppStructFieldsz(data [][]string, abbr string) []string {
	return toStructFields(data, abbr)
}

func toCoreNewStructFieldsz(data [][]string) []string {
	return toStructFields(data, "app")
}

func toCoreUpdateStructFieldsz(data [][]string) []string {
	return toStructFields(data, "app")
}

func toDBStructFieldsz(data [][]string, abbr string) []string {
	return toStructFields(data, abbr)
}

func toCoreStructFieldsz(data [][]string, abbr string) []string {
	abbr = "db" + strings.ToUpper(abbr[0:1]) + abbr[1:]
	return toStructFields(data, abbr)
}

func coreCreateFunctionz(data [][]string, abbr string) []string {
	abbr = "n" + strings.ToUpper(abbr[0:1]) + abbr[1:]
	return toStructFields(data, abbr)
}

func coreUpdateFunctionz(data [][]string, abbr string) []string {
	return toPointerUpdateIfNotNilFields(data, abbr)
}

func toPointerUpdateIfNotNilFields(data [][]string, abbr string) []string {
	result := make([]string, 0, len(data))
	abbrU := strings.ToUpper(abbr[0:1]) + abbr[1:]
	for _, v := range data {
		tempString := "if u" + abbrU + "." + v[directAppNameKey] + "!=nil{ " + abbr + "." + v[directAppNameKey] + "= *u" + abbrU + "." + v[directAppNameKey] + "}"
		result = append(result, tempString)
	}
	return result
}

func camelStyle(str string) string {
	if len(str) == 0 {
		return str
	}
	if len(str) > 2 && str[len(str)-2:] == "ID" {
		// 如果是以ID結尾的，且不只有ID兩個字，則轉換為小寫開頭的Id
		str = str[:1] + str[1:len(str)-2] + "Id"
	}

	// 從大寫開頭的CamelCase轉小寫開頭的CamelCase
	return strings.ToLower(str[:1]) + str[1:]
}

func psqlStyle(str string) string {
	if len(str) == 0 {
		return str
	}

	var builder strings.Builder

	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			prev := rune(str[i-1])

			// 加底線的條件：
			// - 前一個是小寫（像：helloWorld → hello_world）
			// - 前一個是大寫，但下一個是小寫（像：IDCard → id_card）
			if unicode.IsLower(prev) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1]))) {
				builder.WriteByte('_')
			}
		}
		builder.WriteRune(r)
	}

	return strings.ToLower(builder.String())
}
