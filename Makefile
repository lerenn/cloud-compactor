PROJECT_ROOT_PATH := .

include $(PROJECT_ROOT_PATH)/tools/make/docker.mk
include $(PROJECT_ROOT_PATH)/tools/make/help.mk

.PHONY: run
run: docker/run ## Run the application