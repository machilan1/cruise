codegen-field-debug:
	@read -p "Enter the domain path: " DOMAIN_PATH; \
	read -p "Enter the model (empty if file name is model): " PARTIAL_MODEL; \
	read -p "Enter the nwe field name: " FIELD_NAME; \
	APP_MODEL_PATH="./internal/app/domain/"$$DOMAIN_PATH"api/model"$$PARTIAL_MODEL".go"; \
	BUSINESS_MODEL_PATH="./internal/business/domain/"$$DOMAIN_PATH"/model"$$PARTIAL_MODEL".go"; \
	STORAGE_MODEL_PATH="./internal/business/domain/"$$DOMAIN_PATH"/stores/"$$DOMAIN_PATH"db/model"$$PARTIAL_MODEL".go"; \

codegen-field-basic:
	@read -p "Enter the domain path: " DOMAIN_PATH; \
	read -p "Enter the model (empty if file name is model): " PARTIAL_MODEL; \
	read -p "Enter the new field name: " FIELD_NAME; \
	read -p "Enter the new field type (string, int, etc): " FIELD_TYPE; \
	read -p "Enter the app abbr for translation method: " abbr1; \
	read -p "Enter the db abbr for translation method: " abbr2; \
	APP_MODEL_PATH="./internal/app/domain/$$DOMAIN_PATH""api/model$$PARTIAL_MODEL.go"; \
	echo $$APP_MODEL_PATH; \
	BUSINESS_MODEL_PATH="./internal/business/domain/$$DOMAIN_PATH/model$$PARTIAL_MODEL.go"; \
	echo $$BUSINESS_MODEL_PATH; \
	STORAGE_MODEL_PATH="./internal/business/domain/$$DOMAIN_PATH/stores/$$DOMAIN_PATH""db/model$$PARTIAL_MODEL.go"; \
	echo $$STORAGE_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{AD\}/ {print "    " field_name " "ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{AN\}/ {print "    " field_name " "ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{AU\}/ {print "    " field_name " *"ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v abr="$$abbr1" '/\/\/ codegen:\{tAD\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v abr="app" '/\/\/ codegen:\{tBN\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v abr="app" '/\/\/ codegen:\{tBU\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{BD\}/ {print "    " field_name " "ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{BN\}/ {print "    " field_name " "ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{BU\}/ {print "    " field_name " *"ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v ft="$$FIELD_TYPE" '/\/\/ codegen:\{SD\}/ {print "    " field_name " "ft" ""`db:\"" tolower(field_name) "\"`"} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v abr="$$abbr1" '/\/\/ codegen:\{tSD\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \
	awk -v field_name="$$FIELD_NAME" -v abr="$$abbr2" '/\/\/ codegen:\{tBD\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \
