
compiler := ../../tools/JackCompiler.sh
ext := jack
dir := Test

define compilation
	cp "$(1).$(ext)" "$(1)$(dir)"
	@$(compiler) "$(1)$(dir)"
endef

define test_compilation
	cp "$(1).$(ext)" "$(dir)"
	@$(compiler) "$(dir)"
endef

%:
	$(call compilation,$(@))

test-%:
	$(call test_compilation,$($*))

