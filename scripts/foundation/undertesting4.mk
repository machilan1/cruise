codegen-field-sub-struct:
	@read -p "Enter the domain path: " DOMAIN_PATH; \
	read -p "Enter the domain path with first letter capitalize: " DOMAINU; \
	read -p "Enter the model (empty if file name is model): " PARTIAL_MODEL; \
	read -p "Enter the new field name: " FIELD_NAME; \
	read -p "Enter the new field name plural name: " FIELD_PLURAL_NAME; \
	read -p "Enter the new field name plural name with first letter lower case: " FIELD_PLURAL_NAMEL; \
	read -p "Enter the app abbr for translation method: " abbr1; \
	read -p "Enter the db abbr for translation method: " abbr2; \
	APP_MODEL_PATH="./internal/app/domain/$$DOMAIN_PATH""api/model$$PARTIAL_MODEL.go"; \
	echo $$APP_MODEL_PATH; \
	BUSINESS_MODEL_PATH="./internal/business/domain/$$DOMAIN_PATH/model$$PARTIAL_MODEL.go"; \
	echo $$BUSINESS_MODEL_PATH; \
	STORAGE_MODEL_PATH="./internal/business/domain/$$DOMAIN_PATH/stores/$$DOMAIN_PATH""db/model$$PARTIAL_MODEL.go"; \
	echo $$STORAGE_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]App""$$FIELD_NAME" '/\/\/ codegen:\{AD\}/ {print "    " field_name " "ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]AppNew""$$FIELD_NAME" '/\/\/ codegen:\{AN\}/ {print "    " field_name " "ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]AppUpdate""$$FIELD_NAME" '/\/\/ codegen:\{AU\}/ {print "    " field_name " "ft" ""`json:\"" tolower(field_name) "\"`"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]""$$FIELD_NAME" '/\/\/ codegen:\{BD\}/ {print "    " field_name " "ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]New""$$FIELD_NAME" '/\/\/ codegen:\{BN\}/ {print "    " field_name " "ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="[]Update""$$FIELD_NAME" '/\/\/ codegen:\{BU\}/ {print "    " field_name " "ft" "} {print}' $$BUSINESS_MODEL_PATH > temp.go && mv temp.go $$BUSINESS_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v ft="dbjson.JSONColumn[[]db""$$FIELD_NAME""]" '/\/\/ codegen:\{SD\}/ {print "    " field_name " "ft" ""`db:\"" tolower(field_name) "\"`"} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="toApp""$$FIELD_PLURAL_NAME""(""$$abbr1" '/\/\/ codegen:\{tAD\}/ {print "    " field_name ": "abr"."field_name"),"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" -v field_nameL="$$FIELD_PLURAL_NAMEL" '/\/\/ codegen:\{tManyBN\}/ {print "    " field_nameL ",err := toCoreNew"field_name"("abr"."field_name")"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" '/\/\/ codegen:\{tManyBN\}/ {print "    if err!= nil \{"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v domain="$$DOMAIN_PATH" -v domainU="$$DOMAINU" '/\/\/ codegen:\{tManyBN\}/ {print "    return "domain".New"domainU"\{\}, err"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" '/\/\/ codegen:\{tManyBN\}/ {print "    \}"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" -v field_nameL="$$FIELD_PLURAL_NAMEL" '/\/\/ codegen:\{tManyBU\}/ {print "    " field_nameL ",err := toCoreUpdate"field_name"("abr"."field_name")"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" '/\/\/ codegen:\{tManyBU\}/ {print "    if err!= nil \{"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v domain="$$DOMAIN_PATH" -v domainU="$$DOMAINU" '/\/\/ codegen:\{tManyBU\}/ {print "    return "domain".Update"domainU"\{\}, err"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v abr="app" '/\/\/ codegen:\{tManyBU\}/ {print "    \}"} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v field_nameL="$$FIELD_PLURAL_NAMEL" '/\/\/ codegen:\{tBN\}/ {print "    " field_name ": "field_nameL","} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
	awk -v field_name="$$FIELD_PLURAL_NAME" -v field_nameL="$$FIELD_PLURAL_NAMEL" '/\/\/ codegen:\{tBU\}/ {print "    " field_name ": "field_nameL","} {print}' $$APP_MODEL_PATH > temp.go && mv temp.go $$APP_MODEL_PATH; \
#	awk -v field_name="$$FIELD_NAME" -v abr="$$abbr1" '/\/\/ codegen:\{tSD\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \
#	awk -v field_name="$$FIELD_NAME" -v abr="$$abbr2" '/\/\/ codegen:\{tBD\}/ {print "    " field_name ": "abr"."field_name","} {print}' $$STORAGE_MODEL_PATH > temp.go && mv temp.go $$STORAGE_MODEL_PATH; \



codegen-field-ss: codegen-field-sub-struct

now:
	@echo 'if fail to generate please run "brew install coreutils"...' ; \
	echo '' ;\
	S=$${S:-0} ;\
	gdate -d "$$S seconds" +"%Y-%m-%dT%H:%M:%S%:z"

later:
	@$(MAKE) now S=25

t-period:
	@$(MAKE) now S=30 && $(MAKE) now S=3630

tp: t-period