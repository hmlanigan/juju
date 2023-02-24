test_resources() {
	if [ "$(skip 'test_resources')" ]; then
		echo "==> TEST SKIPPED: Resources tests"
		return
	fi

	set_verbosity

	echo "==> Checking for dependencies"
	check_dependencies juju

	file="${TEST_DIR}/test-resources.log"

	bootstrap "test-resources" "${file}"

	test_basic_resources
	test_refresh_resources
	test_attach_resources

	destroy_controller "test-resources"
}
